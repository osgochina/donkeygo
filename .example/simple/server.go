package main

import (
	"donkeygo/drpc"
	"fmt"
	"github.com/gogf/gf/os/glog"
	"time"
)

func main() {
	//开启信号监听
	go drpc.GraceSignal()

	svr := drpc.NewEndpoint(drpc.EndpointConfig{
		CountTime:   true,
		ListenPort:  9090,
		PrintDetail: true,
	})

	svr.RouteCall(new(Math))

	// broadcast per 5s
	go func() {
		for {
			time.Sleep(time.Second * 5)
			svr.RangeSession(func(sess drpc.Session) bool {
				sess.Push(
					"/push/status",
					fmt.Sprintf("this is a broadcast, server time: %v", time.Now()),
				)
				return true
			})
		}
	}()

	svr.ListenAndServe()
}

type Math struct {
	drpc.CallCtx
}

func (m *Math) Add(arg *[]int) (int, *drpc.Status) {
	// test meta
	glog.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
