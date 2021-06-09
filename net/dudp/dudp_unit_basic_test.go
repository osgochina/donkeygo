package dudp_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/net/dudp"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

func Test_Basic(t *testing.T) {
	p, _ := ports.PopRand()
	s := dudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *dudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					dlog.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	// gudp.Conn.Send
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := dudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.Assert(err, nil)
			t.Assert(conn.Send([]byte(dconv.String(i))), nil)
			conn.Close()
		}
	})
	// gudp.Conn.SendRecv
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			conn, err := dudp.NewConn(fmt.Sprintf("127.0.0.1:%d", p))
			t.Assert(err, nil)
			result, err := conn.SendRecv([]byte(dconv.String(i)), -1)
			t.Assert(err, nil)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
			conn.Close()
		}
	})
	// gudp.Send
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			err := dudp.Send(fmt.Sprintf("127.0.0.1:%d", p), []byte(dconv.String(i)))
			t.Assert(err, nil)
		}
	})
	// gudp.SendRecv
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			result, err := dudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte(dconv.String(i)), -1)
			t.Assert(err, nil)
			t.Assert(string(result), fmt.Sprintf(`> %d`, i))
		}
	})
}

// If the read buffer size is less than the sent package size,
// the rest data would be dropped.
func Test_Buffer(t *testing.T) {
	p, _ := ports.PopRand()
	s := dudp.NewServer(fmt.Sprintf("127.0.0.1:%d", p), func(conn *dudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(1)
			if len(data) > 0 {
				if err := conn.Send(data); err != nil {
					dlog.Error(err)
				}
			}
			if err != nil {
				break
			}
		}
	})
	go s.Run()
	defer s.Close()
	time.Sleep(100 * time.Millisecond)
	dtest.C(t, func(t *dtest.T) {
		result, err := dudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("123"), -1)
		t.Assert(err, nil)
		t.Assert(string(result), "1")
	})
	dtest.C(t, func(t *dtest.T) {
		result, err := dudp.SendRecv(fmt.Sprintf("127.0.0.1:%d", p), []byte("456"), -1)
		t.Assert(err, nil)
		t.Assert(string(result), "4")
	})
}
