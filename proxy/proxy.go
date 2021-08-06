package proxy

import (
	"context"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/audit/event"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"net"
	"reflect"
	"time"
)

var TProxy = reflect.TypeOf((*proxyGo)(nil)).String()

type proxyGo struct {
	lua.Super
	cfg *config
	ln  net.Listener

	ctx    context.Context
	cancel context.CancelFunc
}

func newProxyGo(cfg *config) *proxyGo {
	p := &proxyGo{cfg: cfg}
	p.T = TProxy
	p.S = lua.INIT
	return p
}

func (p *proxyGo) Name() string {
	return p.cfg.Name
}

func (p *proxyGo) Start() error {
	ln, err := net.Listen(p.cfg.Protocol, p.cfg.Bind)
	if err != nil {
		return err
	}
	p.ln = ln
	p.ctx, p.cancel = context.WithCancel(context.Background())
	go p.accept()
	p.S = lua.RUNNING
	p.U = time.Now()

	return nil
}

func (p *proxyGo) Close() error {
	p.cancel()
	return p.ln.Close()
}

func (p *proxyGo) handle(src net.Conn) {
	var err error
	var dst net.Conn
	src_addr := src.RemoteAddr().String()

	//链接失败告警
	dst, err = net.Dial(p.cfg.Protocol, p.cfg.Remote)
	if err != nil {
		audit.Put(event.New(p.Name(), event.Addr(src_addr),
			event.Subject("%s 访问蜜罐失败", src_addr),
			event.ERR(err), event.Infof("服务端口:%s 后端地址:%s 链接失败", p.cfg.Bind, p.cfg.Remote)))
		return
	}
	//关闭
	defer dst.Close()

	//链接成功
	audit.Put(event.New(p.Name(), event.Addr(src.RemoteAddr().String()),
		event.Subject("%s 访问蜜罐成功", src_addr),
		event.Infof("服务端口:%s 后端地址:%s 会话端口:%s 链接成功",
			p.cfg.Bind, p.cfg.Remote, dst.RemoteAddr().String())))

	//结束告警
	defer func() {
		audit.Put(event.New(p.Name(), event.Addr(src.RemoteAddr().String()),
			event.Subject("%s 关闭蜜罐请求", src_addr),
			event.Infof("服务地址:%s 后端:%s 会话端口:%s 结束",
				p.cfg.Bind, p.cfg.Remote, dst.LocalAddr().String()),
			event.ERR(err)))
	}()

	flow := newFlowGo(src, dst)
	err = flow.start(p.ctx)
}

func (p *proxyGo) accept() {
	for {
		select {

		//控制退出
		case <-p.ctx.Done():
			return

		default:
			conn, err := p.ln.Accept()
			if err != nil {
				logger.Errorf("%s proxy accept %v", p.Name(), err)
				continue
			}
			go p.handle(conn)
		}

	}

}
