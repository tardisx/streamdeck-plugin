package streamdeck

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"reflect"

	"github.com/gorilla/websocket"
)

var flagPort int
var flagEvent, flagInfo string
var UUID string // the UUID this plugin is assigned

func init() {
	flag.IntVar(&flagPort, "port", 0, "streamdeck sdk port")
	flag.StringVar(&flagEvent, "registerEvent", "", "streamdeck sdk register event")
	flag.StringVar(&flagInfo, "info", "", "streamdeck application info")
	flag.StringVar(&UUID, "pluginUUID", "", "uuid")
	flag.Parse()
}

type logger interface {
	Info(string, ...any)
	Error(string, ...any)
	Debug(string, ...any)
}

type nullLogger struct{}

func (nl nullLogger) Info(string, ...any)  {}
func (nl nullLogger) Error(string, ...any) {}
func (nl nullLogger) Debug(string, ...any) {}

type EventHandler struct {
	MsgType string
	Handler func()
}

type Connection struct {
	ws       *websocket.Conn
	logger   logger
	handlers map[reflect.Type]reflect.Value
	done     chan (bool)
}

func New() Connection {
	return Connection{
		handlers: make(map[reflect.Type]reflect.Value),
		logger:   nullLogger{},
		done:     make(chan bool),
	}
}

func NewWithLogger(l logger) Connection {
	c := New()
	c.logger = l
	return c
}

func (conn *Connection) Connect() error {

	c, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d", flagPort), nil)
	if err != nil {
		return err
	}

	conn.ws = c
	msg := ESOpenMessage{
		ESCommonNoContext: ESCommonNoContext{
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

func (c *Connection) WaitForPluginExit() {
	<-c.done
}

// RegisterHandler registers a function to be called for a particular event. The
// event to be handled is determined by the functions single parameter type.
// This should be called before Connect to be sure your application is ready to
// receive events. You can register as many handlers as you like, but only one
// function per event type. This function will panic if the wrong kind of
// function is passed in, or if you try to register more than one for a single
// event type.
func (r *Connection) RegisterHandler(handler any) {
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
	valid := false
	for i := range receivedEventTypeMap {
		if receivedEventTypeMap[i] == argType {
			valid = true
			break
		}
	}
	if !valid {
		panic("you cannot register a handler with this argument type")
	}

	_, alreadyExists := r.handlers[argType]
	if alreadyExists {
		panic("handler for " + argType.Name() + " already exists")
	}

	r.handlers[argType] = reflect.ValueOf(handler)
}

func (conn *Connection) Send(e any) error {
	b, _ := json.Marshal(e)
	conn.logger.Debug(fmt.Sprintf("sending: %s", string(b)))

	return conn.ws.WriteJSON(e)
}

func (conn *Connection) handle(event any) {
	conn.logger.Debug(fmt.Sprintf("handle: incoming a %T", event))
	argType := reflect.TypeOf(event)
	handler, ok := conn.handlers[argType]
	if !ok {
		conn.logger.Debug(fmt.Sprintf("handle: no handler registered for type %s", argType))
		return
	} else {
		conn.logger.Debug(fmt.Sprintf("handle: found handler function for type %s", argType))

		v := []reflect.Value{reflect.ValueOf(event)}
		conn.logger.Debug(fmt.Sprintf("handle: handler func: %+v", handler))
		conn.logger.Debug(fmt.Sprintf("handle: handler var: %+v", v))

		conn.logger.Debug("handle: calling handler function")

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
		base := ERBase{}
		err = json.NewDecoder(r).Decode(&base)
		if err != nil {
			conn.logger.Error("cannot decode: " + err.Error())
			continue
		}

		t, ok := receivedEventTypeMap[base.Event]
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
	conn.logger.Info(fmt.Sprintf("instance is a %T", d))

	err := json.Unmarshal(b, &d)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal this:\n%s\ninto a %v (%v)\nbecause: %s", string(b), d, t, err.Error())
	}

	// get concrete instance of d into v
	v := reflect.ValueOf(d).Elem().Interface()
	conn.logger.Info(fmt.Sprintf("NOW instance is a %T", v))

	conn.logger.Debug(fmt.Sprintf("reader: unmarshalled to: %+v", v))
	return v, nil
}
