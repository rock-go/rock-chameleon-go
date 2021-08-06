package mysql

import (
	"context"
	"github.com/rock-go/rock-chameleon-go/mysql/engine"
	"github.com/rock-go/rock-chameleon-go/mysql/server"
	"github.com/rock-go/rock-chameleon-go/mysql/sql/information_schema"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"reflect"
	"time"
)

var TGoMySQL = reflect.TypeOf((*GoMysql)(nil)).String()

type GoMysql struct {
	lua.Super
	cfg *config
	ser *server.Server

	ctx    context.Context
	cancel context.CancelFunc
}

func newGoMysql(cfg *config) *GoMysql {
	m := &GoMysql{cfg: cfg}
	m.T = TGoMySQL
	m.S = lua.INIT
	return m
}

func (m *GoMysql) Name() string {
	return m.cfg.Name
}

func (m *GoMysql) Start() error {
	eg := engine.NewDefault()
	eg.AddDatabase(m.cfg.Database.obj)
	eg.AddDatabase(information_schema.NewInformationSchemaDatabase(eg.Catalog))

	s, err := server.NewDefaultServer(m.cfg.toSerCfg(), eg)
	if err != nil {
		return err
	}

	m.ser = s
	go func() {
		err = s.Start()
	}()

	m.U = time.Now()
	m.S = lua.RUNNING
	m.ctx, m.cancel = context.WithCancel(context.Background())
	logger.Errorf("%s %s start succeed", m.Name(), m.Type())
	return nil
}

func (m *GoMysql) Close() error {
	m.cancel()
	return m.ser.Close()
}
