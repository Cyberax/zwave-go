package openzwave

type zwaveCallbackHandler struct {
	CallbackAdapter
	target func(notify ZwaveNotify)
}

func (ctx *zwaveCallbackHandler) OnManagerNotify(notify Notification) {
	ctx.target(ConvertNotification(notify))
}

func MakeDelegatingHandler(target func(notify ZwaveNotify)) CallbackAdapter {
	chandler := new(zwaveCallbackHandler)
    chandler.target = target
	adapter := NewDirectorCallbackAdapter(chandler)

	return adapter
}
