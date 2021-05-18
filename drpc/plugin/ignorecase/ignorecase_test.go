package ignorecase_test

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/drpc"
	"github.com/osgochina/donkeygo/drpc/message"
	"github.com/osgochina/donkeygo/drpc/plugin/ignorecase"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

type Home struct {
	drpc.CallCtx
}

func (that *Home) Test(arg *map[string]string) (map[string]interface{}, *drpc.Status) {
	that.Session().Push("/push/test", map[string]string{
		"your_id": dconv.String(that.PeekMeta("peer_id")),
	})
	return map[string]interface{}{
		"arg": *arg,
	}, nil
}

func TestIgnoreCase(t *testing.T) {

	dtest.C(t, func(t *dtest.T) {
		// Server
		srv := drpc.NewEndpoint(drpc.EndpointConfig{ListenPort: 9090, Network: "tcp"}, ignorecase.NewIgnoreCase())
		srv.RouteCall(new(Home))
		go srv.ListenAndServe()
		time.Sleep(1e9)

		// Client
		cli := drpc.NewEndpoint(drpc.EndpointConfig{Network: "tcp"}, ignorecase.NewIgnoreCase())
		cli.RoutePush(new(Push))
		sess, stat := cli.Dial(":9090")
		if !stat.OK() {
			t.Fatal(stat)
		}
		var result interface{}
		stat = sess.Call("/home/TesT",
			map[string]string{
				"author": "henrylee2cn",
			},
			&result,
			message.WithSetMeta("peer_id", "110"),
		).Status()
		if !stat.OK() {
			t.Error(stat)
		}
		t.Logf("result:%v", result)
		time.Sleep(3e9)
	})

}

type Push struct {
	drpc.PushCtx
}

func (p *Push) Test(arg *map[string]string) *drpc.Status {
	glog.Infof("receive push(%s):\narg: %#v\n", p.IP(), arg)
	return nil
}
