package main

import (
	"donkeygo/drpc"
	"donkeygo/drpc/message"
	"github.com/gogf/gf/os/glog"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{})
	defer cli.Close()

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial(":9090")
	if !stat.OK() {
		glog.Fatalf("%v", stat)
	}
	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		message.WithSetMeta("author", "henrylee2cn"),
	).Status()
	if !stat.OK() {
		glog.Fatalf("%v", stat)
	}
	glog.Printf("result: %d", result)
	glog.Printf("Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	glog.Printf("%s", *arg)
	return nil
}
