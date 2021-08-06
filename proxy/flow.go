package proxy

import (
	"context"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/utils"
	"net"
)

type flowGo struct {
	src net.Conn
	dst net.Conn
}

func newFlowGo(src net.Conn, dst net.Conn) *flowGo {
	return &flowGo{src, dst}
}

func (f *flowGo) close() {
	f.src.Close()
	f.dst.Close()
}

func (f *flowGo) start(ctx context.Context) error {
	defer f.close()

	go func() {
		if _, err := utils.Copy(ctx, f.src, f.dst); err != nil {
			logger.Printf("flow:%v in:%s out:%s ", err, f.src.RemoteAddr().String(), f.dst.RemoteAddr().String())
			return
		}
	}()

	_, err := utils.Copy(ctx, f.dst, f.src)
	return err
}
