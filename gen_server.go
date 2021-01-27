package go_gen_server

type GenServer struct{}

type GenServerDefine interface {
	Init(args interface{}) interface{}
	HandlerCall(interface{}, interface{}) (interface{}, interface{})
	HandlerCast(interface{}, interface{}) (interface{}, interface{})
	HandlerInfo(interface{}, interface{}) (interface{}, interface{})
	Terminate(interface{})
}
