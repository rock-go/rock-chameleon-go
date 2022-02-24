package ssh

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xbase"
	"reflect"
)

var (
	xEnv      *xbase.EnvT
	sshTypeOf = reflect.TypeOf((*sshGo)(nil)).String()
)

func (s *sshGo) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {

	case "version":
		s.serv.Version = val.String()

	case "root":
		s.auth.Set("root", val.String())
	}
}

func newLuaSSH(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, sshTypeOf)
	if proc.IsNil() {
		proc.Set(newSSH(cfg))
	} else {
		proc.Data.(*sshGo).cfg = cfg
	}

	L.Push(proc)
	return 1
}

func Inject(env *xbase.EnvT, uv lua.UserKV) {
	xEnv = env
	uv.Set("ssh", lua.NewFunction(newLuaSSH))
}
