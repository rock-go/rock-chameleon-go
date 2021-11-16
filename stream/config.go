package stream

import (
	"errors"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"strings"
)

type config struct {
	name                 string
	code                 string
	bind_network         string
	bind_address         string

	remote_network       string
	remote_address       string
}

func parseURL(L *lua.LState , vl lua.LValue) (network string , address string) {
	url := vl.String()

	if len(url) < 12 {
		L.RaiseError("invalid url")
		return
	}

	switch {
	case strings.HasPrefix(url , "tcp://"):
		return "tcp" , auxlib.CheckSocket(lua.S2L(url[6:]) , L)

	case strings.HasPrefix(url , "udp://"):
		return "udp" , auxlib.CheckSocket(lua.S2L(url[6:]) , L)

	default:

		L.RaiseError("invalid url , must be tcp://x.x.x.x:80 , got: %s" , url)
		return

	}
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}
	tab.Range(func(key string, lv lua.LValue) {
		switch key {
		case "name":
			cfg.name = auxlib.CheckProcName(lv, L)
		case "bind":
			cfg.bind_network , cfg.bind_address = parseURL(L , lv)
		case "remote":
			cfg.remote_network , cfg.remote_address = parseURL(L , lv)

		default:
			//todo
		}

	})

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}
	cfg.code = L.CodeVM()

	return cfg
}

func (cfg *config) verify() error {
	if  err := auxlib.Name(cfg.name); err != nil {
		return err
	}

	if cfg.remote_address == "" || cfg.remote_network == "" {
		return errors.New("invalid remote socket cfg")
	}

	if 	cfg.bind_network == "" || cfg.bind_address == "" {
		return errors.New("invalid bind socket cfg")
	}

	return nil
}