package go_gen_server

import (
	"errors"
	"sync"
)

var (
	rst          sync.Once
	lck          *sync.RWMutex
	registGlobal map[string]*regist
)

func registData(name string) (*regist, bool) {
	lck.RLock()
	defer lck.RUnlock()

	data, ok := registGlobal[name]
	return data, ok
}

type regist struct {
	l     *sync.RWMutex
	m     map[string]*process
	state interface{}
	name  string
	GenServerDefine
}

func Regist(name string, gSrv interface{}, opts interface{}) (*regist, error) {
	rst.Do(func() {
		lck = new(sync.RWMutex)
		registGlobal = make(map[string]*regist)
	})

	lck.Lock()
	defer lck.Unlock()

	r, ok := registGlobal[name]
	if !ok {
		gen, ok := gSrv.(GenServerDefine)
		if !ok {
			return nil, errors.New("gen_server param error")
		}
		state := gen.Init(opts)
		r = &regist{
			l:               new(sync.RWMutex),
			m:               make(map[string]*process),
			GenServerDefine: gen,
			state:           state,
			name:            name,
		}
		registGlobal[name] = r
	}
	return r, nil
}

func (r *regist) Stop() {
	lck.Lock()
	defer lck.Unlock()
	delete(registGlobal, r.name)
}

func (r *regist) IsExist(name string) bool {
	r.l.RLock()
	_, ok := r.m[name]
	r.l.RUnlock()

	return ok
}

func (r *regist) set(name string, p *process) {
	r.l.Lock()
	r.m[name] = p
	r.l.Unlock()
}

func (r *regist) updateStatus(v interface{}) {
	r.l.Lock()
	r.state = v
	r.l.Unlock()
}

// Process get Process .
func (r *regist) Process(name string) (*process, bool) {
	r.l.RLock()
	res, ok := r.m[name]
	r.l.RUnlock()

	return res, ok
}

// Spawn start new thread, name is not same .
func (r *regist) Spawn(name string) (*process, error) {
	if r.IsExist(name) {
		return nil, errors.New("gen server name exist")
	}
	p := newProcess(name, r.name, r.GenServerDefine)
	r.set(name, p)
	return p, nil
}
