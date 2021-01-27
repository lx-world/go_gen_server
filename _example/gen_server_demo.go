package main

import (
	"context"
	"fmt"

	gen "go_gen_server"
)

type Counter struct {
	gen.GenServer
}

func (t *Counter) Init(v interface{}) interface{} {
	//fmt.Printf("init: %+v\n", v)
	return v
}

func (t *Counter) HandlerCall(state, v interface{}) (interface{}, interface{}) {
	//fmt.Printf("call: %v\n", v)
	tmp := v.(int)
	state1 := state.(int)

	return tmp + state1, tmp + state1
}

func (t *Counter) HandlerCast(state, v interface{}) (interface{}, interface{}) {
	fmt.Printf("cast: %v\n", v)
	return nil, nil
}

func (t *Counter) HandlerInfo(state, v interface{}) (interface{}, interface{}) {
	fmt.Printf("HandlerInfo: %v\n", v)
	return nil, nil
}

func (t *Counter) Terminate(v interface{}) {
	fmt.Println("Terminate")
}

func main() {
	// 计数器使用例子
	r, _ := gen.Regist("gen_server", new(Counter), 0)
	p, _ := r.Spawn("counter")
	p.Call(context.TODO(), p.Self(), 1)
	p.Call(context.TODO(), p.Self(), 1)
	p.Call(context.TODO(), p.Self(), 1)
	res4, _ := p.Call(context.TODO(), p.Self(), 1)
	fmt.Printf("counter: %+v\n", res4)

	r.Stop()
}
