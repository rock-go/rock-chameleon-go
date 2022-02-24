package stream

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/pipe"
	"github.com/rock-go/rock/xbase"
)

var xEnv *xbase.EnvT

/*
	chameleon.stream{
		name = "ssss",
		bind = "tcp://127.0.0.1:3390",
		remote = "tcp://172.31.61.67:3389",
	}
*/

func (s *stream) pipeL(L *lua.LState) int {
	pp := pipe.Check(L)
	if len(pp) > 0 {
		s.cfg.pipe = append(s.cfg.pipe, pp...)
	}

	return 0
}

func (s *stream) startL(L *lua.LState) int {
	if s.Code() != L.CodeVM() {
		L.RaiseError("%s %s proc start must be %s , got %s", s.Code(), s.Name(), L.CodeVM())
		return 0
	}

	xEnv.Start(s, func(err error) {
		L.RaiseError("%v", err)
	})
	return 0
}

func (s *stream) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "pipe":
		return L.NewFunction(s.pipeL)
	case "start":
		return L.NewFunction(s.startL)

	}
	return lua.LNil
}

func newLuaStreamChameleon(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, streamTypeOf)
	if proc.IsNil() {
		proc.Set(newStream(cfg))
	} else {
		proc.Data.(*stream).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(env *xbase.EnvT, uv lua.UserKV) {
	xEnv = env
	uv.Set("stream", lua.NewFunction(newLuaStreamChameleon))
}
