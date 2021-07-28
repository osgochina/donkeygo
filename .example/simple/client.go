package main

import (
	"github.com/osgochina/donkeygo/drpc"
	"github.com/osgochina/donkeygo/drpc/message"
	"github.com/osgochina/donkeygo/os/dlog"
	"time"
)

func main() {

	cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true})
	defer cli.Close()

	cli.RoutePush(new(Push))

	sess, stat := cli.Dial("127.0.0.1:9091")
	if !stat.OK() {
		dlog.Fatalf("%v", stat)
	}
	var result int
	stat = sess.Call("/math/add",
		[]int{1, 2, 3, 4, 5},
		&result,
		message.WithSetMeta("author", "henrylee2cn"),
	).Status()
	if !stat.OK() {
		dlog.Fatalf("%v", stat)
	}
	dlog.Printf("result: %d", result)
	dlog.Printf("Wait 10 seconds to receive the push...")
	time.Sleep(time.Second * 10)
}

type Push struct {
	drpc.PushCtx
}

func (that *Push) Status(arg *string) *drpc.Status {
	dlog.Printf("%s", *arg)
	return nil
}
