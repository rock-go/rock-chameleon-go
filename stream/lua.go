package stream

import "github.com/rock-go/rock/lua"

func newLuaStreamChameleon(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name , streamTypeOf )
	if proc.IsNil() {
		proc.Set(newStream(cfg))
	} else {
		proc.Value.(*stream).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(uv lua.UserKV) {
	uv.Set("stream", lua.NewFunction(newLuaStreamChameleon))
}

