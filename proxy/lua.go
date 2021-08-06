package proxy

import "github.com/rock-go/rock/lua"

func newLuaProxyChameleon(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.Name, TProxy)
	if proc.IsNil() {
		proc.Set(newProxyGo(cfg))
	} else {
		proc.Value.(*proxyGo).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(uv lua.UserKV) {
	uv.Set("proxy", lua.NewFunction(newLuaProxyChameleon))
}
