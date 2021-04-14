package dtcp

import (
	"crypto/tls"
	"net"
	"time"
)

const (
	dDefaultConnTimeout    = 30 * time.Second       //默认链接超时时间
	dDefaultRetryInterval  = 100 * time.Millisecond // 读写重试间隔时间
	dDefaultReadBufferSize = 128                    //默认读取数据的buffer长度
)

//重试结构体
type Retry struct {
	Count    int           // Retry count.
	Interval time.Duration // Retry interval.
}

//NewNetConn 创建并返回一个 net.Conn 链接，addr的格式为："127.0.0.1:80",
//timeout为链接超时时间，如果不传则表示使用默认的30秒超时
func NewNetConn(addr string, timeout ...time.Duration) (net.Conn, error) {
	d := dDefaultConnTimeout
	if len(timeout) > 0 {
		d = timeout[0]
	}
	return net.DialTimeout("tcp", addr, d)
}

// NewNetConnTLS 创建并返回一个tls加密的 net.Conn 链接，addr的格式为："127.0.0.1:80",
func NewNetConnTLS(addr string, tslConfig *tls.Config, timeout ...time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: dDefaultConnTimeout,
	}
	if len(timeout) > 0 {
		dialer.Timeout = timeout[0]
	}
	return tls.DialWithDialer(dialer, "tcp", addr, tslConfig)
}

//判断链接返回的错误是否是超时错误
func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
