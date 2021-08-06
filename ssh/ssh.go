package ssh

import (
	"context"
	"errors"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/audit/event"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"time"
)

type sshGo struct {
	lua.Super
	cfg  *config
	auth *auth

	ctx    context.Context
	cancel context.CancelFunc

	serv *Server
}

func newSSH(cfg *config) *sshGo {
	s := &sshGo{cfg: cfg}
	s.T = TSSH
	s.S = lua.INIT
	s.auth = &auth{data: make(map[string]string)}
	return s
}

func (s *sshGo) Name() string {
	return s.cfg.name
}

func (s *sshGo) event(ctx Context, pass string, err error) {
	audit.Put(event.New("honey_ssh_auth",
		event.ERR(err),
		event.User(ctx.User()),
		event.Addr(ctx.RemoteAddr().String()),
		event.Infof("pass: %s", pass)))
}

var (
	invalidU = errors.New("not found user")
	invalidP = errors.New("invalid pass")
)

func (s *sshGo) doAuth(ctx Context, pass string) bool {
	var err error
	defer s.event(ctx, pass, err)

	name := ctx.User()
	v, ok := s.auth.Get(name)
	if !ok {
		err = errors.New("not found user")
		goto ERR
	}

	if v != pass {
		err = errors.New("invalid pass")
		goto ERR
	}

ERR:
	return false
}

func (s *sshGo) handler(sess Session) {
}

func (s *sshGo) Start() error {
	s.serv = s.cfg.toSSH(s.handler, s.doAuth)

	var err error
	tk := time.NewTicker(2 * time.Second)
	defer tk.Stop()

	go func() {
		err = s.serv.ListenAndServe()
	}()

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.S = lua.RUNNING
	s.U = time.Now()
	logger.Errorf("%s %s start succeed", s.Name(), s.Type())
	return err
}

func (s *sshGo) Close() error {
	s.cancel()
	return s.serv.Close()
}
