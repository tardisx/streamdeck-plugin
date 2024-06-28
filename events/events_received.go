package events

import (
	"encoding/json"
	"reflect"
)

// ValidEventType returns a boolean indicating whether or not
// this is a valid event type
func ValidEventType(t reflect.Type) bool {
	for i := range receivedEventTypeMap {
		if receivedEventTypeMap[i] == t {
			return true
		}
	}
	return false
}

// TypeForEvent returns the type for a particular event type string
func TypeForEvent(e string) (reflect.Type, bool) {
	t, ok := receivedEventTypeMap[e]
	return t, ok
}

var receivedEventTypeMap = map[string]reflect.Type{}

func init() {
	receivedEventTypeMap["keyUp"] = reflect.TypeOf(ERKeyUp{})
	receivedEventTypeMap["didReceiveSettingsPayload"] = reflect.TypeOf(ERDidReceiveSettingsPayload{})
	receivedEventTypeMap["didReceiveSettings"] = reflect.TypeOf(ERDidReceiveSettings{})
	receivedEventTypeMap["globalSettings"] = reflect.TypeOf(ERDidReceiveGlobalSettings{})
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
	Action  string `json:"action"` // The action's unique identifier. If your plugin supports multiple actions, you should use this value to see which action was triggered
	Event   string `json:"event"`
	Context string `json:"context"` // A value identifying the instance's action. You will need to pass this opaque value to several APIs like the setTitle API.
	Device  string `json:"device"`  // A value to identify the device.
}

// ERDidReceiveSettings - The didReceiveSettings event is received after calling the getSettings API to retrieve the persistent data stored for the action.
// https://docs.elgato.com/sdk/plugins/events-received#didreceivesettings
type ERDidReceiveSettings struct {
	ERCommon
	Payload ERDidReceiveSettingsPayload `json:"payload"`
}

type ERDidReceiveSettingsPayload struct {
	Settings    json.RawMessage `json:"settings"` // This JSON object contains persistently stored data
	Coordinates struct {        //The coordinates of the action triggered
		Column int `json:"column"`
		Row    int `json:"row"`
	} `json:"coordinates"`
	State           *int `json:"state"`           // Only set when the action has multiple states defined in its manifest.json. The 0-based value contains the current state of the action
	IsInMultiAction bool `json:"isInMultiAction"` // Boolean indicating if the action is inside a Multi-Action
}

// ERDidReceiveGlobalSettings - The didReceiveGlobalSettings event is received after calling the getGlobalSettings API to retrieve the global persistent data stored for the plugin.
// https://docs.elgato.com/sdk/plugins/events-received#didreceiveglobalsettings
type ERDidReceiveGlobalSettings struct {
	Event   string                            `json:"event"`
	Payload ERDidReceiveGlobalSettingsPayload `json:"payload"`
}

type ERDidReceiveGlobalSettingsPayload struct {
	Settings json.RawMessage `json:"settings"` // This JSON object contains persistently stored data
}

// ERDidReceiveDeepLink - Occurs when Stream Deck receives a deep-link message intended for the plugin. The message is re-routed to the plugin, and provided as part of the payload. One-way deep-link message can be routed to the plugin using the URL format:
//
//	streamdeck://plugins/message/<PLUGIN_UUID>/{MESSAGE}
//
// https://docs.elgato.com/sdk/plugins/events-received#didreceivedeeplink
type ERDidReceiveDeepLink struct {
	Event   string `json:"event"`
	Payload struct {
		Url string `json:"url"` // The deep-link URL, with the prefix omitted. For example the URL streamdeck://plugins/message/com.elgato.test/hello-world would result in a url of hello-world
	} `json:"payload"`
}

// ERTouchTap - When the user touches the display, the plugin will receive the touchTap event
//
//	https://docs.elgato.com/sdk/plugins/events-received#touchtap-sd
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

// ERDialDown - When the user presses the encoder down, the plugin will receive the dialDown event (SD+).
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

// ERDialUp - When the user releases a pressed encoder, the plugin will receive the dialUp event (SD+).
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

// ERDialRotate - When the user rotates the encoder, the plugin will receive the dialRotate event.
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
		Ticks   int  `json:"ticks"`   // The integer which holds the number of "ticks" on encoder rotation. Positive values are for clockwise rotation, negative values are for counterclockwise rotation, zero value is never happen
		Pressed bool `json:"pressed"` // Boolean which is true on rotation when encoder pressed
	} `json:"payload"`
}

// ERKeyDown - When the user presses a key, the plugin will receive the keyDown event.
// https://docs.elgato.com/sdk/plugins/events-received#keydown
type ERKeyDown struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State            *int `json:"state"`            // Only set when the action has multiple states defined in its manifest.json. The 0-based value contains the current state of the action.
		UserDesiredState *int `json:"userDesiredState"` // Only set when the action is triggered with a specific value from a Multi-Action. For example, if the user sets the Game Capture Record action to be disabled in a Multi-Action, you would see the value 1. 0 and 1 are valid.
		IsInMultiAction  bool `json:"isInMultiAction"`  // Boolean indicating if the action is inside a Multi-Action.
	} `json:"payload"`
}

// ERKeyUp - When the user releases a key, the plugin will receive the keyUp event
// https://docs.elgato.com/sdk/plugins/events-received#keyup
type ERKeyUp struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State            *int `json:"state"`            // Only set when the action has multiple states defined in its manifest.json. The 0-based value contains the current state of the action
		UserDesiredState *int `json:"userDesiredState"` // Only set when the action is triggered with a specific value from a Multi-Action. For example, if the user sets the Game Capture Record action to be disabled in a Multi-Action, you would see the value 1. 0 and 1 are valid.
		IsInMultiAction  bool `json:"isInMultiAction"`  // Boolean indicating if the action is inside a Multi-Action
	} `json:"payload"`
}

// ERWillAppear - When an instance of an action is displayed on Stream Deck, for example, when the hardware is first plugged in or when a folder containing that action is entered, the plugin will receive a willAppear event. You will see such an event when:
//   - the Stream Deck application is started
//   - the user switches between profiles
//   - the user sets a key to use your action
//
// https://docs.elgato.com/sdk/plugins/events-received#willappear
type ERWillAppear struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		Controller      string `json:"controller"`      // Defines the controller type the action is applicable to. Keypad refers to a standard action on a Stream Deck device, e.g. 1 of the 15 buttons on the Stream Deck MK.2, or a pedal on the Stream Deck Pedal, etc., whereas an Encoder refers to a dial / touchscreen on the Stream Deck+.
		State           *int   `json:"state"`           // Only set when the action has multiple states defined in its manifest.json. The 0-based value contains the current state of the action.
		IsInMultiAction bool   `json:"isInMultiAction"` // Boolean indicating if the action is inside a Multi-Action.
	} `json:"payload"`
}

// ERWillDisappear - When an instance of an action ceases to be displayed on Stream Deck, for example, when switching profiles or folders, the plugin will receive a willDisappear event. You will see such an event when:
//   - the user switches between profiles
//   - the user deletes an action
//
// https://docs.elgato.com/sdk/plugins/events-received#willdisappear
type ERWillDisappear struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"`
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		Controller      string `json:"controller"`      // Defines the controller type the action is applicable to. Keypad refers to a standard action on a Stream Deck device, e.g. 1 of the 15 buttons on the Stream Deck MK.2, or a pedal on the Stream Deck Pedal, etc., whereas an Encoder refers to a dial / touchscreen on the Stream Deck+.
		State           *int   `json:"state"`           // Only set when the action has multiple states defined in its manifest.json. The 0-based value contains the current state of the action.
		IsInMultiAction bool   `json:"isInMultiAction"` //Boolean indicating if the action is inside a Multi-Action
	} `json:"payload"`
}

// ERTitleParametersDidChange - When the user changes the title or title parameters of the instance of an action, the plugin will receive a titleParametersDidChange event
// https://docs.elgato.com/sdk/plugins/events-received#titleparametersdidchange
type ERTitleParametersDidChange struct {
	ERCommon
	Payload struct {
		Settings    json.RawMessage `json:"settings"` // This JSON object contains data that you can set and is stored persistently.
		Coordinates struct {
			Column int `json:"column"`
			Row    int `json:"row"`
		} `json:"coordinates"`
		State           int    `json:"state"` // This value indicates which state of the action the title or title parameters have been changed.
		Title           string `json:"title"` //The new title.
		TitleParameters struct {
			FontFamily     string `json:"fontFamily"`     //The font family for the title.
			FontSize       int    `json:"fontSize"`       // The font size for the title.
			FontStyle      string `json:"fontStyle"`      // The font style for the title
			FontUnderline  bool   `json:"fontUnderline"`  //Boolean indicating an underline under the title
			ShowTitle      bool   `json:"showTitle"`      //Boolean indicating if the title is visible
			TitleAlignment string `json:"titleAlignment"` //Vertical alignment of the title. Possible values are "top", "bottom" and "middle".
			TitleColor     string `json:"titleColor"`     // Title color.
		} `json:"titleParameters"`
	} `json:"payload"`
}

// ERDeviceDidConnect - When a device is plugged into the computer, the plugin will receive a deviceDidConnect event
// https://docs.elgato.com/sdk/plugins/events-received#devicedidconnect
type ERDeviceDidConnect struct {
	Event string `json:"event"`

	Device     string `json:"device"`
	DeviceInfo struct {
		Name       string `json:"name"` // The name of the device set by the user.
		DeviceType int    `json:"type"` // Type of device. Possible values are kESDSDKDeviceType_StreamDeck (0), kESDSDKDeviceType_StreamDeckMini (1), kESDSDKDeviceType_StreamDeckXL (2), kESDSDKDeviceType_StreamDeckMobile (3) and kESDSDKDeviceType_CorsairGKeys (4)
		Size       struct {
			Columns int `json:"columns"`
			Rows    int `json:"rows"`
		} `json:"size"` // The number of columns and rows of keys that the device owns
	} `json:"deviceInfo"`
}

// ERDeviceDidDisconnect - When a device is unplugged from the computer, the plugin will receive a deviceDidDisconnect event
// https://docs.elgato.com/sdk/plugins/events-received#devicediddisconnect
type ERDeviceDidDisconnect struct {
	Event  string `json:"event"`
	Device string `json:"device"`
}

// ERApplicationDidLaunch - A plugin can request in its manifest.json to be notified when some applications are launched or terminated. The manifest.json should contain an ApplicationsToMonitor object specifying the list of application identifiers to monitor. On macOS, the application bundle identifier is used while the exe filename is used on Windows.
// https://docs.elgato.com/sdk/plugins/events-received#applicationdidlaunch
type ERApplicationDidLaunch struct {
	Event   string `json:"event"`
	Payload struct {
		Application string `json:"application"` // The identifier of the application that has been launched
	}
}

// ERApplicationDidTerminate - A plugin can request in its manifest.json to be notified when some applications are launched or terminated. The manifest.json should contain an ApplicationsToMonitor object specifying the list of application identifiers to monitor. On macOS, the application bundle identifier is used while the exe filename is used on Windows.
// https://docs.elgato.com/sdk/plugins/events-received#applicationdidterminate
type ERApplicationDidTerminate struct {
	Event   string `json:"event"`
	Payload struct {
		Application string `json:"application"` // The identifier of the application that has been launched
	}
}

// When the computer wakes up, the plugin will receive the systemDidWakeUp event
// Several important points to note:
//   - A plugin could get multiple systemDidWakeUp events when waking up the computer
//   - When the plugin receives the systemDidWakeUp event, there is no guarantee that the devices are available
//
// https://docs.elgato.com/sdk/plugins/events-received#systemdidwakeup
type ERApplicationSystemDidWakeUp struct {
	Event string `json:"event"`
}

// ERApplicationPropertyInspectorDidAppear - The plugin will receive a propertyInspectorDidAppear event when the Property Inspector appears
// https://docs.elgato.com/sdk/plugins/events-received#propertyinspectordidappear
type ERApplicationPropertyInspectorDidAppear struct {
	ERCommon
}

// ERApplicationPropertyInspectorDidDisappear - The plugin will receive a propertyInspectorDidDisappear event when the Property Inspector disappears
// https://docs.elgato.com/sdk/plugins/events-received#propertyinspectordiddisappear
type ERApplicationPropertyInspectorDidDisappear struct {
	ERCommon
}

// ERApplicationPropertySendToPlugin - The plugin will receive a sendToPlugin event when the Property Inspector sends a sendToPlugin event
// https://docs.elgato.com/sdk/plugins/events-received#sendtoplugin
type ERApplicationPropertySendToPlugin struct {
	Action  string          `json:"action"`
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}

// ERApplicationPropertySendToPropertyInspector - The Property Inspector will receive a sendToPropertyInspector event when the plugin sends a sendToPropertyInspector event
// https://docs.elgato.com/sdk/plugins/events-received#sendtopropertyinspector
type ERApplicationPropertySendToPropertyInspector struct {
	Action  string          `json:"action"`
	Event   string          `json:"event"`
	Context string          `json:"context"`
	Payload json.RawMessage `json:"payload"`
}
