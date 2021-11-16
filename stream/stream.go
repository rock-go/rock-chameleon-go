package stream

import (
	"context"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/thread"
	"github.com/rock-go/rock/transport"
	"github.com/rock-go/rock/transport/cli"
	"net"
	"reflect"
)

var (
	streamTypeOf = reflect.TypeOf((*stream)(nil)).String()
)

type stream struct {
	lua.Super

	cfg        *config
	cur        config //保存当前启动 为了下次快速启动

	ln         *auxlib.Listener
}

func newStream(cfg *config) *stream {
	obj := &stream{cfg: cfg}
	obj.V(lua.INIT , streamTypeOf)
	return obj
}

func (st *stream) accept( ctx context.Context , conn net.Conn , stop context.CancelFunc ) error {
	//toT nt

	audit.NewEvent("chameleon" ,
		audit.Subject("proxy honey stream hit"),
		audit.From(st.cfg.code),
		audit.Remote(conn.RemoteAddr()),
		audit.Msg("程序名称:%s 本地端口:%s://%s connect succeed" ,
			st.Name() , st.cur.bind_network , st.cur.bind_address),
	).Log().Put()

	var toTn int64

	//接收的数据
	var rev int64

	//报错
	var err error

	//数据通道
	var socket *cli.Stream
	socket , err = st.Transport()
	ev := audit.NewEvent("chameleon" , audit.From(st.cur.code) ,
		audit.Remote(conn.RemoteAddr()))

	if err != nil {
		ev.Subject("proxy honey upstream conn fail")
		ev.Msg("程序名称:%s stream: %s://%s" , st.Name() , st.cur.remote_network , st.cur.remote_address )
		ev.E(err).Log().Put()
	} else {
		ev.Subject("proxy honey upstream conn succeed")
		ev.Msg("程序名称:%s stream: %s://%s" , st.Name() , st.cur.remote_network , st.cur.remote_address )
		ev.E(err).Log().Put()
	}

	cancel := func() {
		stop()
		socket.Close()
	}

	thread.Spawn(0 , func() {
		defer cancel()

		ev = audit.NewEvent("chameleon" , audit.From(st.cur.code) , audit.Remote(conn.RemoteAddr()))

		toTn , err = auxlib.Copy(ctx, socket , conn)
		if err != nil {
			ev.Subject("proxy honey conn close")
			ev.Msg("程序名称: %s\n二段链接关闭 \n发送:%d 原因:%v" , st.Name() ,toTn , err )
			ev.E(err).Log()
		} else {
			ev.Subject("proxy honey conn over")
			ev.Msg("程序名称: %s 发送到远程 发送:%d 原因:%v" , st.Name() ,toTn , err)
			ev.Log().Put()
		}
	})

	thread.Spawn(0 , func() {
		defer cancel()
		ev = audit.NewEvent("chameleon" , audit.From(st.cur.code) , audit.Remote(conn.RemoteAddr()))
		rev , err = auxlib.Copy(ctx, conn , socket)
		if err != nil {
			ev.Subject("proxy honey conn failure")
			ev.Msg("程序名称: %s 接收远程失败:%d 原因:%v" , st.Name() , rev , err )
			ev.E(err).Log().Put()
		} else {
			ev.Subject("proxy honey conn over")
			ev.Msg("程序名称: %s 接收远程结束 数量:%d 原因:%v" , st.Name() , rev , err )
			ev.Log().Put()
		}
	})
	return err
}

func (st *stream) Listen() error {

	if st.ln == nil {
		goto conn
	}

	if  st.cfg.bind_address == st.cur.bind_address &&
		st.cfg.bind_network == st.cur.bind_network {

		st.ln.CloseActiveConn()
		return nil
	}

conn:
	ln , err := auxlib.Listen(st.cfg.bind_network , st.cfg.bind_address)
	if err != nil {
		return err
	}
	st.ln = ln
	return nil
}

func (st *stream) Transport() (*cli.Stream , error) {
	return transport.Cli.Stream(map[string]interface{}{
		"type":"forward",
		"network": st.cfg.remote_network,
		"address": st.cfg.remote_address,
	})

}

func (st *stream) start() (err error) {
	st.cur = *st.cfg
	thread.Spawn(1 , func() {
		err = st.ln.OnAccept(st.accept)
	})
	return
}

func (st *stream) Start() error {

	if e := st.Listen() ; e != nil {
		return e
	}

	return st.start()
}


func (st *stream) Reload() (err error) {
	if e := st.Listen(); e != nil {
		return e
	}

	st.cur = *st.cfg
	return nil
}

func (st *stream) Close() error {
	return st.ln.Close()
}

func (st *stream) Name() string {
	return st.cur.name
}

func (st *stream) Type() string {
	return streamTypeOf
}