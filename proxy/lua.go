package proxy

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"github.com/rock-go/rock/xbase"
)

var xEnv *xbase.EnvT

func (p *proxyGo) pipeL(L *lua.LState) int {
	pp := pipe.Check(L)
	if len(pp) == 0 {
		return 0
	}
	p.cfg.pipe = append(p.cfg.pipe , pp...)
	return 0
}

func (p *proxyGo) startL(L *lua.LState) int {
	if p.Code() != L.CodeVM() {
		L.RaiseError("%s %s proc start must be %s , got %s", p.Code(), p.Name(), L.CodeVM())
		return 0
	}

	xEnv.Start(p, func(err error) {
		L.RaiseError("%v", err)
	})
	return 0
}

func (p *proxyGo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return L.NewFunction(p.pipeL)
	case "start":
		return L.NewFunction(p.startL)

	}

	return lua.LNil
}

/*
	chameleon.proxy{
		name = "xxxxx",
		bind = "tcp://127.0.0.1:3309",
		remote = "tcp://172.31.61.67:3389"
	}

*/
func newLuaProxyChameleon(L *lua.LState) int {
	cfg := newConfig(L)

	proc := L.NewProc(cfg.Name, proxyTypeOf)
	if proc.IsNil() {
		proc.Set(newProxyGo(cfg))
	} else {
		proc.Data.(*proxyGo).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(env *xbase.EnvT, uv lua.UserKV) {
	xEnv = env
	uv.Set("proxy", lua.NewFunction(newLuaProxyChameleon))
}
