package redissubscriber

func (r *controller) RegisterRouter() {
	r.funcHandlers["sendMessage001"] = r.sendMessageHandler
}
