package rabbitmqsubscriber

func (r *controller) RegisterRouter() {
	r.funcHandlers["sendMessage002"] = r.sendMessageHandler
}
