package events

import (
	"encoding/json"
)

type EventTarget int

const EventTargetBoth = EventTarget(0)
const EventTargetHardware = EventTarget(1)
const EventTargetSoftware = EventTarget(2)

type ESCommon struct {
	Event   string `json:"event"`   // name of this event type
	Context string `json:"context"` // A value to Identify the instance's action or Property Inspector. This value is received by the Property Inspector as a parameter of the connectElgatoStreamDeckSocket function.
}

type ESCommonNoContext struct {
	Event string `json:"event"` // name of this event type
}

type ESOpenMessage struct {
	ESCommonNoContext
	UUID string `json:"uuid"`
}

// https://docs.elgato.com/sdk/plugins/events-sent#setsettings

type ESSetSettings struct {
	ESCommon
	Payload json.RawMessage
}

func NewESSetSettings(context string, payload json.RawMessage) ESSetSettings {
	return ESSetSettings{
		ESCommon: ESCommon{
			Event:   "setSettings",
			Context: context,
		},
		Payload: payload,
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#getsettings

type ESGetSettings struct {
	ESCommon
}

func NewESGetSettings(context string) ESGetSettings {
	return ESGetSettings{
		ESCommon: ESCommon{
			Event:   "getSettings",
			Context: context,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#setglobalsettings

type ESSetGlobalSettings struct {
	ESCommon
	Payload json.RawMessage
}

func NewESSetGlobalSettings(context string, payload json.RawMessage) ESSetGlobalSettings {
	return ESSetGlobalSettings{
		ESCommon: ESCommon{
			Event:   "setGlobalSettings",
			Context: context,
		},
		Payload: payload,
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#getglobalsettings

type ESGetGlobalSettings struct {
	ESCommon
}

func NewESGetGlobalSettings(context string) ESGetGlobalSettings {
	return ESGetGlobalSettings{
		ESCommon: ESCommon{
			Event:   "getGlobalSettings",
			Context: context,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#openurl

type ESOpenURL struct {
	ESCommonNoContext
	Payload ESOpenURLPayload `json:"payload"`
}

type ESOpenURLPayload struct {
	URL string `json:"url"`
}

func NewESOpenURL(url string) ESOpenURL {
	return ESOpenURL{
		ESCommonNoContext: ESCommonNoContext{Event: "openUrl"},
		Payload: ESOpenURLPayload{
			URL: url,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#logmessage
type ESLogMessage struct {
	ESCommonNoContext
	Payload ESLogMessagePayload `json:"payload"`
}

type ESLogMessagePayload struct {
	Message string `json:"message"`
}

func NewESLogMessage(message string) ESLogMessage {
	return ESLogMessage{
		ESCommonNoContext: ESCommonNoContext{Event: "logMessage"},
		Payload:           ESLogMessagePayload{Message: message},
	}
}

// setTitle https://docs.elgato.com/sdk/plugins/events-sent#settitle
type ESSetTitle struct {
	ESCommon
	Payload ESSetTitlePayload `json:"payload"`
}

type ESSetTitlePayload struct {
	Title  string      `json:"title"`
	Target EventTarget `json:"target"`
	State  int         `json:"state"`
}

func NewESSetTitle(context string, title string, target EventTarget, state int) ESSetTitle {
	return ESSetTitle{
		ESCommon: ESCommon{
			Event:   "setTitle",
			Context: context,
		},
		Payload: ESSetTitlePayload{
			Title:  title,
			Target: target,
			State:  state,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#setimage
type ESSetImage struct {
	ESCommon
	Payload ESSetImagePayload `json:"payload"`
}

type ESSetImagePayload struct {
	Image  string      `json:"image"`
	Target EventTarget `json:"target"`
	State  int         `json:"state"`
}

func NewESSetImage(context string, imageBase64 string, target EventTarget, state int) ESSetImage {
	return ESSetImage{
		ESCommon: ESCommon{
			Event:   "setImage",
			Context: context,
		},
		Payload: ESSetImagePayload{
			Image:  imageBase64,
			Target: target,
			State:  state,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#setfeedback-sd
type ESSetFeedback struct {
	ESCommon
	Payload json.RawMessage `json:"payload"`
}

func NewESSetFeedback(context string, payload json.RawMessage) ESSetFeedback {
	return ESSetFeedback{
		ESCommon: ESCommon{
			Event:   "setFeedback",
			Context: context,
		},
		Payload: payload,
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#setfeedbacklayout-sd
type ESSetFeedbackLayout struct {
	ESCommon
	Payload ESSetFeedbackLayoutPayload `json:"payload"`
}

type ESSetFeedbackLayoutPayload struct {
	Layout string `json:"layout"`
}

func NewESSetFeedbackLayout(context string, layout string) ESSetFeedbackLayout {
	return ESSetFeedbackLayout{
		ESCommon: ESCommon{
			Event:   "setFeedbackLayout",
			Context: context,
		},
		Payload: ESSetFeedbackLayoutPayload{
			Layout: layout,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#settriggerdescription-sd
type ESSetTriggerDescription struct {
	ESCommon
	Payload ESSetTriggerDescriptionPayload `json:"payload"`
}

type ESSetTriggerDescriptionPayload struct {
	Rotate    string `json:"rotate,omitempty"`
	Push      string `json:"push,omitempty"`
	Touch     string `json:"touch,omitempty"`
	LongTouch string `json:"longTouch,omitempty"`
}

func NewESSetTriggerDescription(context string, rotate, push, touch, longTouch string) ESSetTriggerDescription {
	return ESSetTriggerDescription{
		ESCommon: ESCommon{
			Event:   "setTriggerDescription",
			Context: context,
		},
		Payload: ESSetTriggerDescriptionPayload{
			Rotate:    rotate,
			Push:      push,
			Touch:     touch,
			LongTouch: longTouch,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#showalert
type ESShowAlert struct {
	ESCommon
}

func NewESShowAlert(context string) ESShowAlert {
	return ESShowAlert{
		ESCommon: ESCommon{
			Event:   "showAlert",
			Context: context,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#showok
type ESShowOK struct {
	ESCommon
}

func NewESShowOK(context string) ESShowOK {
	return ESShowOK{
		ESCommon: ESCommon{
			Event:   "showOk",
			Context: context,
		},
	}
}

//https://docs.elgato.com/sdk/plugins/events-sent#setstate

type ESSetState struct {
	ESCommon
	Payload ESSetStatePayload `json:"payload"`
}

type ESSetStatePayload struct {
	State int `json:"state"`
}

func NewESSetState(context string, state int) ESSetState {
	return ESSetState{
		ESCommon: ESCommon{
			Event:   "setState",
			Context: context,
		},
		Payload: ESSetStatePayload{
			State: state,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#switchtoprofile

type ESSwitchToProfile struct {
	ESCommon
	Device  string                   `json:"device"`
	Payload ESSwitchToProfilePayload `json:"payload"`
}

type ESSwitchToProfilePayload struct {
	Profile string `json:"profile"`
	Page    int    `json:"page"`
}

func NewESSwitchToProfile(context string, device, profileName string, page int) ESSwitchToProfile {
	return ESSwitchToProfile{
		ESCommon: ESCommon{
			Event:   "switchToProfile",
			Context: context,
		},
		Device: device,
		Payload: ESSwitchToProfilePayload{
			Profile: profileName,
			Page:    page,
		},
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#sendtopropertyinspector

type ESSendToPropertyInspector struct {
	ESCommon
	Action  string `json:"action"`
	Payload json.RawMessage
}

func NewESSendToPropertyInspector(context string, action string, payload json.RawMessage) ESSendToPropertyInspector {
	return ESSendToPropertyInspector{
		ESCommon: ESCommon{
			Event:   "sendToPropertyInspector",
			Context: context,
		},
		Action:  action,
		Payload: payload,
	}
}

// https://docs.elgato.com/sdk/plugins/events-sent#sendtoplugin

type ESSendToPlugin struct {
	ESCommon
	Action  string `json:"action"`
	Payload json.RawMessage
}

func NewESSendToPlugin(context string, action string, payload json.RawMessage) ESSendToPlugin {
	return ESSendToPlugin{
		ESCommon: ESCommon{
			Event:   "sendToPlugin",
			Context: context,
		},
		Action:  action,
		Payload: payload,
	}
}
