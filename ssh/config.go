package ssh

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
)

type config struct {
	name   string
	bind   string
	prompt string
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}

	tab.ForEach(func(key lua.LValue, val lua.LValue) {
		switch key.String() {
		case "name":
			cfg.name = val.String()
		case "bind":
			cfg.bind = val.String()
		case "prompt":
			cfg.prompt = val.String()
		default:
			L.RaiseError("not found %s key", key.String())
		}
	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg

}

func (cfg *config) verify() error {
	if e := utils.Name(cfg.name); e != nil {
		return e
	}

	return nil
}

func (cfg *config) toSSH(h Handler, p PasswordHandler) *Server {
	return &Server{Addr: cfg.bind, Handler: h, PasswordHandler: p}
}
