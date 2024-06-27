package streamdeck

import (
	"testing"
)

type testLogger struct {
	t *testing.T
}

func (tl testLogger) Info(s string, x ...any)  { tl.t.Log(s, x) }
func (tl testLogger) Debug(s string, x ...any) { tl.t.Log(s, x) }
func (tl testLogger) Error(s string, x ...any) { tl.t.Log(s, x) }

func TestReflection(t *testing.T) {

	c := NewWithLogger(testLogger{t: t})
	// incoming
	in := ERDidReceiveSettingsPayload{}

	ranHandler := false
	c.RegisterHandler(func(event ERDidReceiveSettingsPayload) {
		ranHandler = true
	})

	c.handle(in)

	if !ranHandler {
		t.Error("did not run handler")
	}

}

func TestUmmarshal(t *testing.T) {

	b := []byte(`
{
    "action": "com.elgato.example.action1",
    "event": "keyUp",
    "context": "ABC123",
    "device": "DEF456",
    "payload": {
        "settings": {},
        "coordinates": {
            "column": 3,
            "row": 1
        },
        "state": 0,
        "userDesiredState": 1,
        "isInMultiAction": false
    }
}`)

	c := NewWithLogger(testLogger{t: t})
	keyUp, err := c.unmarshalToConcrete(receivedEventTypeMap["keyUp"], b)

	if err != nil {
		t.Error(err)
	}

	realKeyUp, ok := keyUp.(ERKeyUp)
	if !ok {
		t.Errorf("wrong type (is %T)", keyUp)
	}
	if realKeyUp.Context != "ABC123" {
		t.Error("wrong value")
	}

}
