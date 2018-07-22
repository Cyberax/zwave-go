package openzwave

import (
	"time"
		"sync"
)

type ZwaveConfig struct {
	ZwaveLibDir  string
	StateDir     string
	NetworkKey   string
	LogFile      string
	LogToConsole bool
	InitTimeout  time.Duration
}

type Listener chan <- ZwaveNotify

type ZwaveInterface struct {
	Manager Manager
	notifyDispatcher CallbackAdapter

	listenerMutex sync.Mutex
	globalListeners map[Listener] Listener

	adapterMutex sync.Mutex
	adapters map[string] string
}

func (zi *ZwaveInterface) AddAdapter(path string) () {
	zi.adapterMutex.Lock()
	zi.adapters[path] = path
	zi.Manager.AddDriver(path)
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
	config.AddOptionString("ConfigPath", ensure(options.ZwaveLibDir), false)

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

	zi := new(ZwaveInterface)
	zi.Manager = manager
	zi.globalListeners = make(map[Listener]Listener)

	zi.notifyDispatcher = MakeDelegatingHandler(func(notify ZwaveNotify) {
		zi.listenerMutex.Lock()
		defer zi.listenerMutex.Unlock()
		for k := range zi.globalListeners {
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
	ManagerDestroy()
	OptionsDestroy()
}
