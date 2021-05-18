package heartbeat_test

import (
	"github.com/osgochina/donkeygo/drpc"
	"github.com/osgochina/donkeygo/drpc/plugin/heartbeat"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestHeartbeatCALl(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		svr := drpc.NewEndpoint(drpc.EndpointConfig{
			ListenPort:  9090,
			PrintDetail: true,
		}, heartbeat.NewPong())
		go svr.ListenAndServe()

		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true}, heartbeat.NewPing(3, true))
		cli.Dial(":9090")
		time.Sleep(time.Second * 20)

	})
}

func TestHeartbeatCALl2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		svr := drpc.NewEndpoint(drpc.EndpointConfig{
			ListenPort:  9090,
			PrintDetail: true,
		}, heartbeat.NewPong())
		go svr.ListenAndServe()

		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(drpc.EndpointConfig{PrintDetail: true}, heartbeat.NewPing(3, true))
		sess, _ := cli.Dial(":9090")
		for i := 0; i < 8; i++ {
			sess.Call("/", nil, nil).Status()
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 10)
	})
}

func TestHeartbeatPush1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9090, PrintDetail: true},
			heartbeat.NewPing(3, false),
		)
		go srv.ListenAndServe()
		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(
			drpc.EndpointConfig{PrintDetail: true},
			heartbeat.NewPong(),
		)
		cli.Dial(":9090")
		time.Sleep(time.Second * 10)
	})

}

func TestHeartbeatPush2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		srv := drpc.NewEndpoint(
			drpc.EndpointConfig{ListenPort: 9090, PrintDetail: true},
			heartbeat.NewPing(3, false),
		)
		go srv.ListenAndServe()
		time.Sleep(time.Second)

		cli := drpc.NewEndpoint(
			drpc.EndpointConfig{PrintDetail: true},
			heartbeat.NewPong(),
		)
		sess, _ := cli.Dial(":9090")
		for i := 0; i < 8; i++ {
			sess.Push("/", nil)
			time.Sleep(time.Second)
		}
		time.Sleep(time.Second * 5)
	})

}
