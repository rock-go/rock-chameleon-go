package proxy

import (
	"errors"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
	"github.com/rock-go/rock/xreflect"
)

type config struct {
	Name     string `lua:"name"     type:"string"`
	Protocol string `lua:"protocol" type:"string"`
	Bind     string `lua:"bind"     type:"string"`
	Remote   string `lua:"remote"   type:"string"`
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := &config{}
	if e := xreflect.ToStruct(tab, cfg); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := utils.Name(cfg.Name); e != nil {
		return e
	}

	switch cfg.Protocol {
	case "tcp", "udp":
		//todo

	default:
		return errors.New("invalid protocol")
	}

	return nil
}
