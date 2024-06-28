// Package streamdeck interfaces with the Stream Deck plugin API,
// allowing you to create go-based plugins for the platform.
package streamdeck

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"reflect"

	"github.com/tardisx/streamdeck-plugin/events"

	"github.com/gorilla/websocket"
)

// these are automatically populated by parseFlags
var flagPort int
var flagEvent, flagInfo string
var UUID string // the UUID this plugin is assigned

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
	Debug(string, ...any)
}

type nullLogger struct{}

func (nl nullLogger) Info(string, ...any)  {}
func (nl nullLogger) Error(string, ...any) {}
func (nl nullLogger) Debug(string, ...any) {}

type Connection struct {
	ws       *websocket.Conn
	logger   logger
	handlers map[reflect.Type]reflect.Value
	done     chan (bool)
}

// New creates a new struct for communication with the streamdeck
// plugin API. The websocket will not connect until Connect is called.
func New() Connection {
	return Connection{
		handlers: make(map[reflect.Type]reflect.Value),
		logger:   nullLogger{},
		done:     make(chan bool),
	}
}

// NewWithLogger is the same as New, but allows you to set a logger
// for debugging the websocket connection.
func NewWithLogger(l logger) Connection {
	c := New()
	c.logger = l
	return c
}

// parseFlags parses the command line flags to get the values provided
// by the Stream Deck plugin API.
func parseFlags() {
	flag.IntVar(&flagPort, "port", 0, "streamdeck sdk port")
	flag.StringVar(&flagEvent, "registerEvent", "", "streamdeck sdk register event")
	flag.StringVar(&flagInfo, "info", "", "streamdeck application info")
	flag.StringVar(&UUID, "pluginUUID", "", "uuid")
	flag.Parse()
}

// Connect connects the plugin to the Stream Deck API via the websocket.
// Once connected, events will be passed to handlers you have registered.
// Handlers should thus be registered via RegisterHandler before calling
// Connect.
// Connect returns immediately if the connection is successful, you should
// then call WaitForPluginExit to block until the connection is closed.
func (conn *Connection) Connect() error {

	parseFlags()
	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d", flagPort), nil)
	if err != nil {
		return err
	}

	conn.ws = c
	msg := events.ESOpenMessage{
		ESCommonNoContext: events.ESCommonNoContext{
			Event: flagEvent,
		},
		UUID: UUID,
	}
	conn.logger.Debug(fmt.Sprintf("writing openMessage: %+v", msg))
	err = c.WriteJSON(msg)
	if err != nil {
		conn.logger.Error(err.Error())
		panic(err)
	}

	// run the reader forever
	conn.logger.Info("starting reader")
	go conn.reader()

	return nil
}

// WaitForPluginExit waits until the Stream Deck API closes
// the websocket connection.
func (conn *Connection) WaitForPluginExit() {
	<-conn.done
}

// RegisterHandler registers a function to be called for a particular event. The
// event to be handled is determined by the functions single parameter type.
// This should be called before Connect to be sure your application is ready to
// receive events. You can register as many handlers as you like, but only one
// function per event type. This function will panic if the wrong kind of
// function is passed in, or if you try to register more than one for a single
// event type.
func (conn *Connection) RegisterHandler(handler any) {
	hType := reflect.TypeOf(handler)
	if hType.Kind() != reflect.Func {
		panic("handler must be a function")
	}
	// Assuming the function takes exactly one argument
	if hType.NumIn() != 1 {
		panic("handler func must take exactly one argument")
	}

	argType := hType.In(0)

	// check its a valid one (one that matches an event type)
	if !events.ValidEventType(argType) {
		panic("you cannot register a handler with this argument type")
	}

	_, alreadyExists := conn.handlers[argType]
	if alreadyExists {
		panic("handler for " + argType.Name() + " already exists")
	}

	conn.handlers[argType] = reflect.ValueOf(handler)
}

// Send sends a message to the API. It should be one of the
// events.ES* structs, such as events.ESOpenURL.
func (conn *Connection) Send(e any) error {
	b, _ := json.Marshal(e)
	conn.logger.Debug(fmt.Sprintf("sending: %s", string(b)))

	return conn.ws.WriteJSON(e)
}

func (conn *Connection) handle(event any) {
	// conn.logger.Debug(fmt.Sprintf("handle: incoming a %T", event))
	argType := reflect.TypeOf(event)
	handler, ok := conn.handlers[argType]
	if !ok {
		conn.logger.Debug(fmt.Sprintf("handle: no handler registered for type %s", argType))
		return
	} else {
		conn.logger.Debug(fmt.Sprintf("handle: found handler function for type %s", argType))

		v := []reflect.Value{reflect.ValueOf(event)}
		// conn.logger.Debug(fmt.Sprintf("handle: handler func: %+v", handler))
		// conn.logger.Debug(fmt.Sprintf("handle: handler var: %+v", v))

		// conn.logger.Debug("handle: calling handler function")

		handler.Call(v)
	}
}

func (conn *Connection) reader() {
	for {
		_, r, err := conn.ws.NextReader()
		if err != nil {
			conn.logger.Error(err.Error())
			break
		}

		b := bytes.Buffer{}
		r = io.TeeReader(r, &b)
		base := events.ERBase{}
		err = json.NewDecoder(r).Decode(&base)
		if err != nil {
			conn.logger.Error("cannot decode: " + err.Error())
			continue
		}

		t, ok := events.TypeForEvent(base.Event)
		if !ok {
			conn.logger.Error(fmt.Sprintf("no type registered for event '%s'", base.Event))
			continue
		}

		d, err := conn.unmarshalToConcrete(t, b.Bytes())
		if err != nil {
			conn.logger.Error("cannot unmarshal: " + err.Error())
			continue
		}
		conn.handle(d)
	}
	conn.logger.Info("websocket closed, shutting down reader")
	conn.done <- true
}

// unmarshalToConcrete takes a reflect.Type and a byte slice, creates a concrete
// instance of the type and unmarshals the JSON byte slice into it.
func (conn *Connection) unmarshalToConcrete(t reflect.Type, b []byte) (any, error) {
	// t is a reflect.Type of the thing we need to decode into
	d := reflect.New(t).Interface()
	// conn.logger.Info(fmt.Sprintf("instance is a %T", d))

	err := json.Unmarshal(b, &d)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal this:\n%s\ninto a %v (%v)\nbecause: %s", string(b), d, t, err.Error())
	}

	// get concrete instance of d into v
	v := reflect.ValueOf(d).Elem().Interface()
	// conn.logger.Debug(fmt.Sprintf("NOW instance is a %T", v))
	// conn.logger.Debug(fmt.Sprintf("reader: unmarshalled to: %+v", v))
	return v, nil
}
