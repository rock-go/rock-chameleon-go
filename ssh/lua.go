package ssh

import (
	"github.com/rock-go/rock/lua"
	"reflect"
)

var sshTypeOf = reflect.TypeOf((*sshGo)(nil)).String()

func (s *sshGo) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {

	case "version":
		s.serv.Version = val.String()

	case "root":
		s.auth.Set("root" , val.String())
	}
}

func newLuaSSH(L *lua.LState) int {
	cfg := newConfig(L)
	proc := L.NewProc(cfg.name, sshTypeOf)
	if proc.IsNil() {
		proc.Set(newSSH(cfg))
	} else {
		proc.Value.(*sshGo).cfg = cfg
	}

	proc.Value.(*sshGo).codeVM = L.CodeVM
	L.Push(proc)
	return 1
}

func Inject(uv lua.UserKV) {
	uv.Set("ssh", lua.NewFunction(newLuaSSH))
}
