package workers_pool

type RequestHandler func(interface{})

type Request struct {
	Data    interface{}
	Handler RequestHandler
}
