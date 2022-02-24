package stream

import (
	"context"
	"fmt"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/transport"
	"github.com/rock-go/rock/transport/cli"
	"github.com/rock-go/rock/xlib"
	"net"
	"reflect"
)

var (
	streamTypeOf = reflect.TypeOf((*stream)(nil)).String()
)

type stream struct {
	lua.Super

	cfg *config
	cur config //保存当前启动 为了下次快速启动

	ln *xlib.Listener
}

func newStream(cfg *config) *stream {
	obj := &stream{cfg: cfg}
	obj.V(lua.INIT, streamTypeOf)
	return obj
}

func (st *stream) socket(conn net.Conn) (*cli.Stream, error) {
	host := st.cur.remote.Hostname()
	port := st.cur.remote.Port()
	if port == 0 {
		_, port = auxlib.ParseAddr(conn.LocalAddr())
	}

	if port == 0 {
		return nil, fmt.Errorf("invalid stream port")
	}

	return transport.Cli.Stream(map[string]interface{}{
		"type":    "forward",
		"network": st.cur.remote.Scheme(),
		"address": fmt.Sprintf("%s:%d", host, port),
	})
}

func (st *stream) pipe(ev *audit.Event) {
	n := len(st.cur.pipe)
	if n == 0 {
		ev.Put()
		return
	}

	for i := 0; i < n; i++ {
		pipe := st.cur.pipe[i]
		if e := pipe(ev, st.cur.co); e != nil {
			xEnv.Errorf("%s stream pipe fail %v", st.Name(), e)
		}
	}
}

func (st *stream) Code() string {
	return st.cfg.co.CodeVM()
}

func (st *stream) accept(ctx context.Context, conn net.Conn, stop context.CancelFunc) error {
	//toT nt
	ev := audit.NewEvent("chameleon").Alert().High().
		Subject("流式高交互代理蜜罐名命中").
		From(st.Code()).
		Remote(conn.RemoteAddr()).
		Msg("程序名称:%s %s://%s connect succeed",
			st.Name(), st.cur.bind.Scheme, conn.LocalAddr().String())
	st.pipe(ev)

	var toTn int64

	//接收的数据
	var rev int64

	//报错
	var err error

	//数据通道
	var socket *cli.Stream
	socket, err = st.socket(conn)
	ev = audit.NewEvent("chameleon").From(st.Code()).Remote(conn.RemoteAddr()).Alert().High()

	if err != nil {
		ev.Subject("流式高交互代理蜜罐恶意请求失败").
			Msg("程序名称:%s stream: %s", st.Name(), st.cur.remote.String()).E(err)
		st.pipe(ev)
		return err
	} else {
		ev.Subject("流式高交互代理蜜罐恶意请求成功").
			Msg("程序名称:%s stream: %s://%s", st.Name(), st.cur.remote.String()).E(err)
		st.pipe(ev)
	}

	xEnv.Spawn(0, func() {
		defer func() {
			stop()
			conn.Close()
		}()

		ev = audit.NewEvent("chameleon").From(st.Code()).Remote(conn.RemoteAddr()).Alert().High()

		toTn, err = auxlib.Copy(ctx, socket, conn)
		if err != nil {
			ev.Subject("流式高交互代理蜜罐上游请求关闭").
				Msg("程序名称: %s\n二段链接关闭 \n发送:%d", st.Name(), toTn).E(err)
			st.pipe(ev)

		} else {
			ev.Subject("流式高交互代理蜜罐上游请求结束").
				Msg("程序名称: %s 发送到远程 发送:%d", st.Name(), toTn)
			st.pipe(ev)
		}
	})

	xEnv.Spawn(0, func() {
		defer func() {
			stop()
			socket.Close()
		}()

		ev = audit.NewEvent("chameleon").From(st.Code()).Remote(conn).Alert().High()
		rev, err = auxlib.Copy(ctx, conn, socket)
		if err != nil {
			ev.Subject("流式高交互代理蜜罐请求关闭").
				Msg("程序名称: %s \n接收远程失败:%d", st.Name(), rev).E(err)
			st.pipe(ev)
		} else {
			ev.Subject("流式高交互代理蜜罐请求结束").
				Msg("程序名称: %s 接收远程结束 数量:%d", st.Name(), rev)
			st.pipe(ev)
		}
	})

	return err
}

func (st *stream) equal() bool {
	if st.cfg.remote.String() != st.cur.remote.String() {
		return false
	}

	if st.cfg.bind.String() != st.cur.bind.String() {
		return false
	}

	return true

}

func (st *stream) Listen() error {

	if st.ln == nil {
		goto conn
	}

conn:
	ln, err := xlib.Listen(st.cfg.bind)
	if err != nil {
		return err
	}
	st.ln = ln
	return nil
}

func (st *stream) start() (err error) {
	st.cur = *st.cfg
	xEnv.Spawn(100, func() {
		err = st.ln.OnAccept(st.accept)
	})
	return
}

func (st *stream) Start() error {

	if e := st.Listen(); e != nil {
		return e
	}

	return st.start()
}

//func (st *stream) Reload() (err error) {
//	if e := st.Listen(); e != nil {
//		return e
//	}
//
//	st.cur = *st.cfg
//	return nil
//}

func (st *stream) Close() error {
	return st.ln.Close()
}

func (st *stream) Name() string {
	return st.cur.name
}

func (st *stream) Type() string {
	return streamTypeOf
}
