package mysql

import (
	"github.com/rock-go/rock-chameleon-go/mysql/auth"
	"github.com/rock-go/rock-chameleon-go/mysql/server"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xreflect"
)

type config struct {
	Name     string    `lua:"name"     type:"string"`
	Bind     string    `lua:"bind"     type:"string"`
	Auth     auth.Auth `lua:"auth"     type:"object"`
	Database *EngineDB `lua:"database" type:"object"`

	CodeVM string
}

func newConfig(L *lua.LState) *config {
	tab := L.CheckTable(1)
	cfg := new(config)

	cfg.CodeVM = L.CodeVM()

	if e := xreflect.ToStruct(tab, cfg); e != nil {
		L.RaiseError("%v", e)
		return cfg
	}

	if e := cfg.verify(); e != nil {
		L.RaiseError("%v", e)
		return cfg
	}

	return cfg
}

func (cfg *config) verify() error {
	if e := auxlib.Name(cfg.Name); e != nil {
		return e
	}

	return nil
}

func (cfg *config) toSerCfg() server.Config {
	return server.Config{
		Protocol: "tcp",
		Address:  cfg.Bind,
		Auth:     cfg.Auth,
		CodeVM:   func() string { return cfg.CodeVM },
	}
}
