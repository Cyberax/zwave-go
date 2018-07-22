package openzwave


// The classification of a value to enable low level system or configuration parameters to be filtered by the application.
type ZwaveValueGenre int  
const (
	ZwaveValueGenre_Basic ZwaveValueGenre = 0 + iota // The 'level' as controlled by basic commands.  Usually duplicated by another command class.
	ZwaveValueGenre_User // Basic values an ordinary user would be interested in.
	ZwaveValueGenre_Config // Device-specific configuration parameters.  These cannot be automatically discovered via Z-Wave, and are usually described in the user manual instead.
	ZwaveValueGenre_System // Values of significance only to users who understand the Z-Wave protocol
	ZwaveValueGenre_Count // A count of the number of genres defined.  Not to be used as a genre itself.
)

// The type of data represented by the value object.
type ZwaveValueType int
const (
	ZwaveValueType_Bool ZwaveValueType = 0 + iota // Boolean, true or false
	ZwaveValueType_Byte	// 8-bit unsigned value
	ZwaveValueType_Decimal // Represents a non-integer value as a string to avoid floating point accuracy issues.
	ZwaveValueType_Int // 32-bit signed value
	ZwaveValueType_List	// List from which one item can be selected
	ZwaveValueType_Schedule	// Complex type used with the Climate Control Schedule command class
	ZwaveValueType_Short // 16-bit signed value
	ZwaveValueType_String // Text string
	ZwaveValueType_Button // A write-only value that is the equivalent of pressing a button to send a command to a device
	ZwaveValueType_Raw // A collection of bytes
	ZwaveValueType_Max = ZwaveValueType_Raw	// The highest-number type defined.  Not to be used as a type itself.
)

// A simple value object with ValueID structure details
type ZwaveValueID struct {
	// The Home ID of the driver that controls the node containing the value.
	HomeID uint32

	// The Zwave node ID, containing the value
	NodeID uint8

	// The genre of the value.  The genre classifies a value to enable
	// low-level system or configuration parameters to be filtered out by the application
	Genre ZwaveValueGenre

	// The Z-Wave command class that created and manages this value.  Knowledge of
	// command classes is not required to use OpenZWave, but this information is
	// exposed in case it is of interest.
	CommandClass uint8

	// Get the command class instance of this value.  It is possible for there to be
	// multiple instances of a command class, although currently it appears that
	// only the SensorMultilevel command class ever does this.  Knowledge of
	// instances and command classes is not required to use OpenZWave, but this
	// information is exposed in case it is of interest.
	Instance uint8

	// The value index.  The index is used to identify one of multiple
	// values created and managed by a command class.  In the case of configurable
	// parameters (handled by the configuration command class), the index is the
	// same as the parameter ID.  Knowledge of command classes is not required
	// to use OpenZWave, but this information is exposed in case it is of interest.
	Index uint8

	// The type of the value
	Type ZwaveValueType

	// An ID that specifies this ValueID within OpenZWave, not stable across restarts
	ComposedID uint64
}


// Notifications of various Z-Wave events sent to the watchers
type ZwaveNotificationType int

const (
	ZwaveType_ValueAdded ZwaveNotificationType = 0 + iota // A new node value has been added to OpenZWave's list. These notifications occur after a node has been discovered, and details of its command classes have been received.  Each command class may generate one or more values depending on the complexity of the item being represented.  
	ZwaveType_ValueRemoved // A node value has been removed from OpenZWave's list.  This only occurs when a node is removed. 
	ZwaveType_ValueChanged // A node value has been updated from the Z-Wave network and it is different from the previous value. 
	ZwaveType_ValueRefreshed // A node value has been updated from the Z-Wave network. 
	ZwaveType_Group // The associations for the node have changed. The application should rebuild any group information it holds about the node. 
	ZwaveType_NodeNew // A new node has been found (not already stored in zwcfg*.xml file) 
	ZwaveType_NodeAdded // A new node has been added to OpenZWave's list.  This may be due to a device being added to the Z-Wave network, or because the application is initializing itself. 
	ZwaveType_NodeRemoved // A node has been removed from OpenZWave's list.  This may be due to a device being removed from the Z-Wave network, or because the application is closing. 
	ZwaveType_NodeProtocolInfo // Basic node information has been received, such as whether the node is a listening device, a routing device and its baud rate and basic, generic and specific types. It is after this notification that you can call Manager::GetNodeType to obtain a label containing the device description. 
	ZwaveType_NodeNaming // One of the node names has changed (name, manufacturer, product). 
	ZwaveType_NodeEvent // A node has triggered an event.  This is commonly caused when a node sends a Basic_Set command to the controller.  The event value is stored in the notification. 
	ZwaveType_PollingDisabled // Polling of a node has been successfully turned off by a call to Manager::DisablePoll 
	ZwaveType_PollingEnabled // Polling of a node has been successfully turned on by a call to Manager::EnablePoll 
	ZwaveType_SceneEvent // Scene Activation Set received 
	ZwaveType_CreateButton // Handheld controller button event created 
	ZwaveType_DeleteButton // Handheld controller button event deleted 
	ZwaveType_ButtonOn // Handheld controller button on pressed event 
	ZwaveType_ButtonOff // Handheld controller button off pressed event 
	ZwaveType_DriverReady // A driver for a PC Z-Wave controller has been added and is ready to use.  The notification will contain the controller's Home ID, which is needed to call most of the Manager methods. 
	ZwaveType_DriverFailed // Driver failed to load 
	ZwaveType_DriverReset // All nodes and values for this driver have been removed.  This is sent instead of potentially hundreds of individual node and value notifications. 
	ZwaveType_EssentialNodeQueriesComplete // The queries on a node that are essential to its operation have been completed. The node can now handle incoming messages. 
	ZwaveType_NodeQueriesComplete // All the initialization queries on a node have been completed. 
	ZwaveType_AwakeNodesQueried // All awake nodes have been queried, so client application can expected complete data for these nodes. 
	ZwaveType_AllNodesQueriedSomeDead // All nodes have been queried but some dead nodes found. 
	ZwaveType_AllNodesQueried // All nodes have been queried, so client application can expected complete data. 
	ZwaveType_Notification // An error has occurred that we need to report. 
	ZwaveType_DriverRemoved // The Driver is being removed. (either due to Error or by request) Do Not Call Any Driver Related Methods after receiving this call 

	// When Controller Commands are executed, Notifications of Success/Failure etc are communicated via this Notification
	// Notification::GetEvent returns Driver::ControllerState and Notification::GetNotification returns Driver::ControllerError if there was a error
	ZwaveType_ControllerCommand

	ZwaveType_NodeReset // The Device has been reset and thus removed from the NodeList in OZW
)

// Notification codes.
type ZwaveNotificationCode int
const (
	ZwaveCode_MsgComplete ZwaveNotificationCode = 0 + iota // Completed messages 
	ZwaveCode_Timeout // Messages that timeout will send a Notification with this code. 
	ZwaveCode_NoOperation // Report on NoOperation message sent completion  
	ZwaveCode_Awake // Report when a sleeping node wakes up 
	ZwaveCode_Sleep // Report when a node goes to sleep 
	ZwaveCode_Dead // Report when a node is presumed dead 
	ZwaveCode_Alive// Report when a node is revived 
)

type ZwaveNotify struct {
	Type ZwaveNotificationType
	// Get the unique ValueID of any value involved in this notification.
	ValueID ZwaveValueID

	// The index of the association group that has been changed.
	// Only valid in Notification::Type_Group notifications.
	GroupIdx uint8

	// Get the event value of a notification.
	// Only valid in Notification::Type_NodeEvent and Notification::Type_ControllerCommand notifications.
	Event uint8

	// The button id of a notification.
	// Only valid in Notification::Type_CreateButton, Notification::Type_DeleteButton,
	// Notification::Type_ButtonOn and Notification::Type_ButtonOff notifications.
	ButtonId uint8

	// The scene Id of a notification.
	// Only valid in Notification::Type_SceneEvent notifications.
	SceneId uint8

	// The notification code from a notification.
	// Only valid for Notification::Type_Notification or
	// Notification::Type_ControllerCommand notifications.
	Notification uint8

	// The internal byte value of the notification. Should not normally need to be used.
	Byte uint8

	// Pre-rendered string description
	StringValue string
}

func ConvertValueID(valueId ValueID) ZwaveValueID {
	return ZwaveValueID{
		HomeID: valueId.GetHomeId(),
		NodeID: valueId.GetNodeId(),
		Genre: ZwaveValueGenre(valueId.GetGenre()),
		CommandClass: valueId.GetCommandClassId(),
		Instance: valueId.GetInstance(),
		Index: valueId.GetIndex(),
		Type: ZwaveValueType(valueId.GetType()),
		ComposedID: valueId.GetId(),
	}
}

func ConvertNotification(notify Notification) ZwaveNotify {
	var res = ZwaveNotify{
		Type: ZwaveNotificationType(notify.GetType()),
		ValueID: ConvertValueID(notify.GetValueID()),
		StringValue: notify.GetAsString(),
		Byte: notify.GetByte(),
	}

	switch res.Type {
	case ZwaveType_Group:
		res.GroupIdx = notify.GetGroupIdx()
	case ZwaveType_NodeEvent:
		res.Event = notify.GetEvent()
	case ZwaveType_CreateButton, ZwaveType_DeleteButton, ZwaveType_ButtonOn, ZwaveType_ButtonOff:
		res.ButtonId = notify.GetButtonId()
	case ZwaveType_SceneEvent:
		res.SceneId = notify.GetSceneId()
	case ZwaveType_Notification:
		res.Notification = notify.GetNotification()
	case ZwaveType_ControllerCommand:
		res.Notification = notify.GetNotification()
		res.Event = notify.GetEvent()
	}

	return res
}
