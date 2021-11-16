package ssh

import (
	"context"
	"errors"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/thread"
)

type sshGo struct {
	lua.Super
	cfg  *config
	auth *auth

	ctx    context.Context
	cancel context.CancelFunc

	codeVM func() string

	serv *Server
}

func newSSH(cfg *config) *sshGo {
	s := &sshGo{cfg: cfg}
	s.V(lua.INIT , sshTypeOf)
	s.auth = &auth{data: make(map[string]string)}
	return s
}

func (s *sshGo) Name() string {
	return s.cfg.name
}

func (s *sshGo) event(ctx Context, pass string, err error) {
	audit.NewEvent("chameleon" ,
		audit.Subject("ssh auth fail") ,
		audit.From(s.codeVM()) ,
		audit.User(ctx.User()) ,
		audit.Remote(ctx.RemoteAddr().String()),
		audit.Msg("pass: %s" , pass ),
		audit.E(err)).Put()
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

	if s.cfg.version != "" {
		s.serv.Version = s.cfg.version
	}

	var err error
	thread.Spawn(1 , func() {
		s.serv.CodeVM = s.codeVM
		err = s.serv.ListenAndServe()
	})

	if err != nil {
		return err
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	logger.Errorf("%s %s start succeed", s.Name(), s.Type())
	return err
}

func (s *sshGo) Close() error {
	s.cancel()
	return s.serv.Close()
}
