package openzwave

import (
	"time"
	"fmt"
	"sync"
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

type Listener chan <- Notification

type NodeAddress struct {
	HomeId uint32
	NodeId uint8
}

type ZwaveInterface struct {
	Manager Manager
	HomeIdList map[uint32]string

	listenerMutex sync.Mutex
	globalListeners map[Listener] Listener
	notifyDispatcher CallbackAdapter
}

func (zi *ZwaveInterface) AddListener(listener Listener) (){
	zi.listenerMutex.Lock()
	defer zi.listenerMutex.Unlock()

	zi.globalListeners[listener] = listener
}

func (zi *ZwaveInterface) RemoveListener(listener Listener) (){
	zi.listenerMutex.Lock()
	defer zi.listenerMutex.Unlock()

	delete(zi.globalListeners, listener)
}

func ensure(path string) string {
	if path != "" && path[len(path)-1] != '/' {
		return path + "/"
	}
	return path
}

func MakeZwaveInterface(options *ZwaveConfig) (*ZwaveInterface, error) {
	var erroring_out = true

	config := OptionsCreate(ensure(options.ZwaveLibDir), ensure(options.StateDir), "")
	defer func() {
		if erroring_out {
			OptionsDestroy()
		}
	}()

	config.AddOptionBool("SaveConfiguration", options.StateDir != "")
	config.AddOptionString("UserPath", ensure(options.StateDir), false)
	config.AddOptionString("ConfigPath", ensure(options.StateDir), false)

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
	zi.globalListeners = make(map[Listener]Listener)

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

	zi.notifyDispatcher = MakeDelegatingHandler(func(notify Notification) {
		zi.listenerMutex.Lock()
		defer zi.listenerMutex.Unlock()
		for k, _ := range zi.globalListeners {
			k <- notify
		}
	})
	AddNotifyHandler(manager, zi.notifyDispatcher)

	erroring_out = false
	return zi, nil
}

func (zwaveInterface *ZwaveInterface) StopZwaveInterface() {
	RemoveNotifyHandler(zwaveInterface.Manager, zwaveInterface.notifyDispatcher)
	DeleteCallbackAdapter(zwaveInterface.notifyDispatcher)

	zwaveInterface.Manager = nil
	zwaveInterface.HomeIdList = make(map[uint32]string)

	ManagerDestroy()
	OptionsDestroy()
}
