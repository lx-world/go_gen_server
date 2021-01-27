package go_gen_server

import (
	"context"
	"errors"
	"time"
)

const (
	defaultTimeOut     = 10 * time.Second
	defaultSendChanNum = 100
	defaultRspChanNum  = 100
)

type process struct {
	pname  string
	name   string
	chSend chan *pid
	chRsp  chan interface{}
	queue  *queue
	GenServerDefine
}

func newProcess(name, pname string, gen GenServerDefine) *process {
	p := &process{
		pname:           pname,
		name:            name,
		chSend:          make(chan *pid, defaultSendChanNum),
		chRsp:           make(chan interface{}, defaultRspChanNum),
		queue:           newQueue(),
		GenServerDefine: gen,
	}
	go p.loop(gen)
	return p
}

func (p *process) Call(ctx context.Context, pid *pid, v interface{}) (interface{}, error) {
	p.queue.Push(v)
	if pid == nil || pid.name == "" {
		return nil, errors.New("pid not exist")
	}
	pid.opt = _call
	p.chSend <- pid

	defaultCtx, cancel := context.WithTimeout(context.Background(), defaultTimeOut)
	defer cancel()

	for {
		select {
		case data := <-p.chRsp:
			return data, nil
		case <-ctx.Done():
			return nil, errors.New("gen_server call timeout")
		case <-defaultCtx.Done():
			return nil, errors.New("gen_server call default ctx timeout")
		}
	}
}

func (p *process) Cast(ctx context.Context, pid *pid, v interface{}) {
	p.queue.Push(v)
	if pid == nil || pid.name == "" {
		return
	}
	pid.opt = _cast
	p.chSend <- pid
}

func (p *process) Info(ctx context.Context, pid *pid, v interface{}) {
	p.queue.Push(v)
	if pid == nil || pid.name == "" {
		return
	}
	pid.opt = _info
	p.chSend <- pid
}

func (p *process) Self() *pid {
	return &pid{
		pname: p.pname,
		name:  p.name,
	}
}

func (p *process) loop(gen GenServerDefine) {
	defer func() {
		if err := recover(); err != nil {
			gen.Terminate(err)
		}
	}()

	for {
		select {
		case pd := <-p.chSend:
			if pd == nil {
				return
			}
			rstData, exist := registData(pd.pname)
			if !exist {
				continue
			}
			data, ok := rstData.Process(pd.name)
			if !ok {
				continue
			}
			res := p.queue.Pop()

			go func(proc *process, v interface{}) {
				var (
					rVal, upVal interface{}
				)
				switch pd.opt {
				case _call:
					rVal, upVal = proc.GenServerDefine.HandlerCall(rstData.state, v)
					p.chRsp <- rVal
				case _cast:
					proc.GenServerDefine.HandlerCast(rstData.state, v)
				default:
					proc.GenServerDefine.HandlerInfo(rstData.state, v)
				}
				rstData.updateStatus(upVal)
			}(data, res)
		}
	}
}
