package ssh

import (
	"github.com/rock-go/rock/lua"
	"reflect"
	"strings"
)

var sshTypeOf = reflect.TypeOf((*sshGo)(nil)).String()

func (s *sshGo) NewIndex(L *lua.LState, key string, val lua.LValue) {
	if strings.HasPrefix(key, "auth_") {
		name := key[5:]
		pass := val.String()
		s.auth.Set(name, pass)
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

	L.Push(proc)
	return 1
}

func Inject(uv lua.UserKV) {
	uv.Set("ssh", lua.NewFunction(newLuaSSH))
}
