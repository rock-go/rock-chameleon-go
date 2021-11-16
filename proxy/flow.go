package proxy

import (
	"context"
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/auxlib"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/thread"
	"net"
)

type flowGo struct {
	ctx  context.Context
	stop context.CancelFunc

	src  net.Conn
	dst  net.Conn
}

func newFlowGo(src net.Conn, dst net.Conn) *flowGo {
	ctx , stop := context.WithCancel(context.Background())
	return &flowGo{ctx , stop , src, dst}
}

func (f *flowGo) close() {
	f.stop()
	f.src.Close()
	f.dst.Close()
}

func (f *flowGo) start() (err error) {
	defer f.close()

	thread.Spawn(0 , func(){
		if _, err := auxlib.Copy(f.ctx, f.src, f.dst); err != nil {
			audit.NewEvent("chameleon" ,
				audit.Subject("flow copy fail"),
				)
			logger.Printf("flow:%v in:%s out:%s ", err, f.src.RemoteAddr().String(), f.dst.RemoteAddr().String())
			return
		}
	})

	thread.Spawn(0 , func() {
		_, err = auxlib.Copy(f.ctx, f.dst, f.src)
	})

	return
}
