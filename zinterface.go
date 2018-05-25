package openzwave

import (
	"time"
	"fmt"
)

type ZwaveConfig struct {
	ZwaveLibDir  string
	StateDir     string
	Adapters     []string
	NetworkKey   string
	LogFile      string
	LogToConsole bool
	InitTimeout  time.Duration
}

type ZwaveInterface struct {
	Manager Manager
	HomeIdList map[uint32]string
}

func MakeZwaveInterface(options *ZwaveConfig) (*ZwaveInterface, error) {
	var erroring_out = true

	config := OptionsCreate(options.ZwaveLibDir, options.StateDir, "")
	defer func() {
		if erroring_out {
			OptionsDestroy()
		}
	}()

	config.AddOptionBool("SaveConfiguration", options.StateDir != "")
	config.AddOptionString("UserPath", options.StateDir, false)
	config.AddOptionString("ConfigPath", options.StateDir, false)

	config.AddOptionBool("ConsoleOutput", options.LogToConsole)
	config.AddOptionString("LogFileName", options.LogFile, false)
	config.AddOptionBool("AppendLogFile", true)

	config.AddOptionString("NetworkKey", options.NetworkKey, false)
	config.Lock()

	manager := ManagerCreate()
	defer func() {
		if erroring_out {
			ManagerDestroy()
		}
	}()

	homeIdChan := make(chan uint32, len(options.Adapters))
	initAdapter := MakeDelegatingHandler(func(notify Notification){
		if notify.GetType() == NotificationType_DriverReady {
			select {
			case homeIdChan <- notify.GetHomeId():
			default:
			}
		}
	})
	defer DeleteCallbackAdapter(initAdapter)

	AddNotifyHandler(manager, initAdapter)
	defer RemoveNotifyHandler(manager, initAdapter)

	for _, serialAdapter := range options.Adapters {
		manager.AddDriver(serialAdapter)
	}

	zi := new(ZwaveInterface)
	zi.Manager = manager
	zi.HomeIdList = make(map[uint32]string)

	timeout := time.After(options.InitTimeout)
	for len(zi.HomeIdList) < len(options.Adapters) {
		select {
		case <-timeout:
			return nil, fmt.Errorf("timed out initializing devices")
		case homeId := <- homeIdChan:
			fmt.Printf("Discovered device with homeid 0x%X on %s\n", homeId,
				manager.GetControllerPath(homeId))
			zi.HomeIdList[homeId] = manager.GetControllerPath(homeId)
		}
	}

	erroring_out = false
	return zi, nil
}

func (zwaveInterface *ZwaveInterface) StopZwaveInterface() {
	zwaveInterface.Manager = nil
	zwaveInterface.HomeIdList = make(map[uint32]string)
	ManagerDestroy()
	OptionsDestroy()
}

