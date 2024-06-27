package streamdeck

import (
	"encoding/json"
	"reflect"
)

var receivedEventTypeMap = map[string]reflect.Type{}

func init() {
	receivedEventTypeMap["keyUp"] = reflect.TypeOf(ERKeyUp{})
	receivedEventTypeMap["didReceiveSettingsPayload"] = reflect.TypeOf(ERDidReceiveSettingsPayload{})
	receivedEventTypeMap["didReceiveSettings"] = reflect.TypeOf(ERDidReceiveSettings{})
	receivedEventTypeMap["globalSettings"] = reflect.TypeOf(ERGlobalSettings{})
	receivedEventTypeMap["didReceiveDeepLink"] = reflect.TypeOf(ERDidReceiveDeepLink{})
	receivedEventTypeMap["touchTap"] = reflect.TypeOf(ERTouchTap{})
	receivedEventTypeMap["dialDown"] = reflect.TypeOf(ERDialDown{})
	receivedEventTypeMap["dialUp"] = reflect.TypeOf(ERDialUp{})
	receivedEventTypeMap["dialRotate"] = reflect.TypeOf(ERDialRotate{})
	receivedEventTypeMap["keyDown"] = reflect.TypeOf(ERKeyDown{})
	receivedEventTypeMap["willAppear"] = reflect.TypeOf(ERWillAppear{})
	receivedEventTypeMap["willDisappear"] = reflect.TypeOf(ERWillDisappear{})
	receivedEventTypeMap["titleParametersDidChange"] = reflect.TypeOf(ERTitleParametersDidChange{})
	receivedEventTypeMap["deviceDidConnect"] = reflect.TypeOf(ERDeviceDidConnect{})
	receivedEventTypeMap["deviceDidDisconnect"] = reflect.TypeOf(ERDeviceDidDisconnect{})
	receivedEventTypeMap["applicationDidLaunch"] = reflect.TypeOf(ERApplicationDidLaunch{})
	receivedEventTypeMap["applicationDidTerminate"] = reflect.TypeOf(ERApplicationDidTerminate{})
	receivedEventTypeMap["applicationSystemDidWakeUp"] = reflect.TypeOf(ERApplicationSystemDidWakeUp{})
	receivedEventTypeMap["applicationPropertyInspectorDidAppear"] = reflect.TypeOf(ERApplicationPropertyInspectorDidAppear{})
	receivedEventTypeMap["applicationPropertyInspectorDidDisappear"] = reflect.TypeOf(ERApplicationPropertyInspectorDidDisappear{})
	receivedEventTypeMap["applicationPropertySendToPlugin"] = reflect.TypeOf(ERApplicationPropertySendToPlugin{})
	receivedEventTypeMap["applicationPropertySendToPropertyInspector"] = reflect.TypeOf(ERApplicationPropertySendToPropertyInspector{})
}

type ERBase struct {
	Event string `json:"event"`
}

//{Action: Event:deviceDidConnect Context: Device:A1DA463F033AD2616E05636CD16F064F Payload:[]}"}

type ERCommon struct {
	Action  string `json:"action"`
	Event   string `json:"event"`
	Context string `json:"context"`
	Device  string `json:"device"`
}

// https://docs.elgato.com/sdk/plugins/events-received#didreceivesettings

type ERDidReceiveSettingsPayload struct {
	Settings    json.RawMessage `json:"settings"`
	Coordinates struct {
		Column int `json:"column"`
		Row    int `json:"row"`
	} `json:"coordinates"`
	IsInMultiAction bool `json:"isInMultiAction"`
}

type ERDidReceiveSettings struct {
	ERCommon
	Payload ERDidReceiveSettingsPayload `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#didreceiveglobalsettings

type ERGlobalSettings struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#didreceivedeeplink

type ERDidReceiveDeepLink struct {
	Event   string `json:"event"`
	Payload struct {
		Url string `json:"url"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#touchtap-sd

type ERTouchTap struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Controller  string          `json:"controller"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		TapPosition []int `json:"tapPos"`
		Hold        bool  `json:"hold"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#dialdown-sd

type ERDialDown struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Controller  string          `json:"controller"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#dialup-sd

type ERDialUp struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Controller  string          `json:"controller"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#dialrotate-sd

type ERDialRotate struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Controller  string          `json:"controller"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		Ticks   int  `json:"ticks"`
		Pressed bool `json:"pressed"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#keydown

type ERKeyDown struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State            int  `json:"state"`
		UserDesiredState int  `json:"userDesiredState"`
		IsInMultiAction  bool `json:"isInMultiAction"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#keyup
type ERKeyUp struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State            int  `json:"state"`
		UserDesiredState int  `json:"userDesiredState"`
		IsInMultiAction  bool `json:"isInMultiAction"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#willappear
type ERWillAppear struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		Controller      string `json:"controller"`
		State           int    `json:"state"`
		IsInMultiAction bool   `json:"isInMultiAction"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#willdisappear

const ERWillDisappearAction = "willDisappear"

type ERWillDisappear struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		Controller      string `json:"controller"`
		State           int    `json:"state"`
		IsInMultiAction bool   `json:"isInMultiAction"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#titleparametersdidchange

type ERTitleParametersDidChange struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State           int    `json:"state"`
		Title           string `json:"title"`
		TitleParameters struct {
			FontFamily     string `json:"fontFamily"`
			FontSize       int    `json:"fontSize"`
			FontStyle      string `json:"fontStyle"`
			FontUnderline  bool   `json:"fontUnderline"`
			ShowTitle      bool   `json:"showTitle"`
			TitleAlignment string `json:"titleAlignment"`
			TitleColor     string `json:"titleColor"`
		} `json:"titleParameters"`
	} `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#devicedidconnect

type ERDeviceDidConnect struct {
	Event string `json:"event"`

	Device     string `json:"device"`
	DeviceInfo struct {
		Name       string `json:"name"`
		DeviceType int    `json:"type"`
		Size       struct {
			Columns int `json:"columns"`
			Rows    int `json:"rows"`
		} `json:"size"`
	} `json:"deviceInfo"`
}

// https://docs.elgato.com/sdk/plugins/events-received#devicediddisconnect

type ERDeviceDidDisconnect struct {
	Event  string `json:"event"`
	Device string `json:"device"`
}

// https://docs.elgato.com/sdk/plugins/events-received#applicationdidlaunch

type ERApplicationDidLaunch struct {
	Event   string `json:"event"`
	Payload struct {
		Application string `json:"application"`
	}
}

// https://docs.elgato.com/sdk/plugins/events-received#applicationdidterminate

type ERApplicationDidTerminate struct {
	Event   string `json:"event"`
	Payload struct {
		Application string `json:"application"`
	}
}

// https://docs.elgato.com/sdk/plugins/events-received#systemdidwakeup

type ERApplicationSystemDidWakeUp struct {
	Event string `json:"event"`
}

// https://docs.elgato.com/sdk/plugins/events-received#propertyinspectordidappear
type ERApplicationPropertyInspectorDidAppear struct {
	ERCommon
}

// https://docs.elgato.com/sdk/plugins/events-received#propertyinspectordiddisappear

type ERApplicationPropertyInspectorDidDisappear struct {
	ERCommon
}

// https://docs.elgato.com/sdk/plugins/events-received#sendtoplugin

type ERApplicationPropertySendToPlugin struct {
	Action  string          `json:"action"`
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

// https://docs.elgato.com/sdk/plugins/events-received#sendtopropertyinspector
type ERApplicationPropertySendToPropertyInspector struct {
	Action  string          `json:"action"`
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}
