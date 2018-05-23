package openzwave

type zwaveCallbackHandler struct {
	CallbackAdapter
	target func(notify Notification)
}

func (ctx *zwaveCallbackHandler) OnManagerNotify(notify Notification) {
	ctx.target(notify)
}

func MakeDelegatingHandler(target func(notify Notification)) CallbackAdapter {
	chandler := new(zwaveCallbackHandler)
    chandler.target = target
	adapter := NewDirectorCallbackAdapter(chandler)

	return adapter
}
