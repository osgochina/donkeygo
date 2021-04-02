package socket

import (
	"errors"
	"net"
	"time"
)

const (
	normal      int32 = 0 //链接正常
	activeClose int32 = 1 //链接已关闭
)

var ErrProactivelyCloseSocket = errors.New("socket is closed proactively")

var (
	writeBuffer     int           = -1
	readBuffer      int           = -1
	changeKeepAlive bool          = false
	keepAlive       bool          = true
	keepAlivePeriod time.Duration = -1
	noDelay         bool          = true
)

func TryOptimize(conn net.Conn) {
	//if c, ok := conn.(ifaceSetKeepAlive); ok {
	//	if changeKeepAlive {
	//		c.SetKeepAlive(keepAlive)
	//	}
	//	if keepAlivePeriod >= 0 && keepAlive {
	//		c.SetKeepAlivePeriod(keepAlivePeriod)
	//	}
	//}
	//if c, ok := conn.(ifaceSetBuffer); ok {
	//	if readBuffer >= 0 {
	//		c.SetReadBuffer(readBuffer)
	//	}
	//	if writeBuffer >= 0 {
	//		c.SetWriteBuffer(writeBuffer)
	//	}
	//}
	//if c, ok := conn.(ifaceSetNoDelay); ok {
	//	if !noDelay {
	//		c.SetNoDelay(noDelay)
	//	}
	//}
}
