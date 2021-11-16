package proxy

import (
	"context"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/thread"
	"net"
	"reflect"
)

var proxyTypeOf = reflect.TypeOf((*proxyGo)(nil)).String()

type proxyGo struct {
	lua.Super

	cfg *config
	cur config

	ln  *auxlib.Listener
}

func newProxyGo(cfg *config) *proxyGo {
	p := &proxyGo{cfg: cfg}
	p.V(lua.INIT , proxyTypeOf)
	return p
}

func (p *proxyGo) Name() string {
	return p.cfg.Name
}

func (p *proxyGo) Listen() error {

	if p.ln == nil {
		goto conn
	}

	if  p.cfg.Bind == p.cur.Bind &&
		p.cfg.Protocol == p.cur.Protocol {

		p.cur = *p.cfg
		p.ln.CloseActiveConn()
		return nil
	}

conn:
	ln , err := auxlib.Listen(p.cfg.Protocol, p.cfg.Bind)
	if err != nil {
		return err
	}
	p.ln = ln
	p.cur = *p.cfg

	return nil
}

func (p *proxyGo) Start() error {

	if e := p.Listen(); e != nil {
		return e
	}

	var err error
	thread.Spawn(3000 , func(){ err = p.ln.OnAccept(p.accept)})
	return err
}

func (p *proxyGo) Reload() error {
	return p.Listen()
}

func (p *proxyGo) Close() error {
	e := p.ln.Close()
	return e
}

func (p *proxyGo) accept(ctx context.Context , conn net.Conn , stop context.CancelFunc) error {
	ev := audit.NewEvent("chameleon",
		audit.Subject("proxy honey conn hit"),
		audit.From(p.cur.code),
		audit.Remote(conn.RemoteAddr()))

	dst , err := net.Dial(p.cur.Protocol , p.cur.Remote)
	if err != nil {
		ev.Msg("%s 服务端口:%s 后端地址:%s 链接失败",p.Name(), p.cfg.Bind, p.cfg.Remote).E(err).Log().Put()
	} else {
		ev.Msg("%s 服务端口:%s 后端地址:%s 链接成功",p.Name(), p.cfg.Bind, p.cfg.Remote).Log().Put()
	}

	cancel := func() {
		stop()
		dst.Close()
	}

	thread.Spawn(100 , func() {
		defer cancel()
		var toTn int64

		ev = audit.NewEvent("chameleon" , audit.From(p.cur.code) , audit.Remote(conn.RemoteAddr()))

		toTn , err = auxlib.Copy(ctx, dst , conn)
		ev.Msg("%s ")
		if err != nil {
			ev.Subject("proxy honey conn close")
			ev.Msg("程序名称: %s\n  代理关闭 \n发送:%d 原因:%v" , p.Name() ,toTn , err )
			ev.E(err).Log()
		} else {
			ev.Subject("proxy honey conn over")
			ev.Msg("程序名称: %s 发送到远程 发送:%d 原因:%v" , p.Name() ,toTn , err)
			ev.Log().Put()
		}
	})

	thread.Spawn(100 , func() {
		defer cancel()
		var rev int64

		ev = audit.NewEvent("chameleon" , audit.From(p.cur.code) , audit.Remote(conn.RemoteAddr()))
		rev , err = auxlib.Copy(ctx, conn , dst)
		if err != nil {
			ev.Subject("proxy honey conn close")
			ev.Msg("程序名称: %s 接收远程失败:%d 原因:%v" , p.Name() , rev , err )
			ev.E(err).Log().Put()
		} else {
			ev.Subject("proxy honey conn over")
			ev.Msg("程序名称: %s 接收远程结束 数量:%d 原因:%v" , p.Name() , rev , err )
			ev.Log().Put()
		}
	})

	return err
}