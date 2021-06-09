package dudp

import (
	"io"
	"net"
	"time"
)

// Conn udp链接处理对象
type Conn struct {
	*net.UDPConn                 // 底层 UDP 链接.
	remoteAddr     *net.UDPAddr  // 远端的地址.
	recvDeadline   time.Time     // 读取数据超时时间.
	sendDeadline   time.Time     // 发送数据超时时间.
	recvBufferWait time.Duration // 读取缓冲区的间隔时间.
}

const (
	dDefaultRetryInterval  = 100 * time.Millisecond // 默认重试的间隔时间.
	dDefaultReadBufferSize = 1024                   // (Byte)Buffer size.
	dRecvAllWaitTimeout    = time.Millisecond       // 默认读取缓冲区的间隔时间.
)

type Retry struct {
	Count    int           // 最大重试次数
	Interval time.Duration // 重试间隔.
}

// NewConn 创建一个UDP链接，remoteAddress必须传入。localAddress可选
func NewConn(remoteAddress string, localAddress ...string) (*Conn, error) {
	if conn, err := NewNetConn(remoteAddress, localAddress...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnByNetConn 传入底层的udp链接，包装出我们自定义的udp链接
func NewConnByNetConn(udp *net.UDPConn) *Conn {
	return &Conn{
		UDPConn:        udp,
		recvDeadline:   time.Time{},
		sendDeadline:   time.Time{},
		recvBufferWait: dRecvAllWaitTimeout,
	}
}

// Send 发送数据
func (that *Conn) Send(data []byte, retry ...Retry) (err error) {
	for {
		if that.remoteAddr != nil {
			_, err = that.WriteToUDP(data, that.remoteAddr)
		} else {
			_, err = that.Write(data)
		}
		if err != nil {
			// 如果链接已经断开
			if err == io.EOF {
				return err
			}
			// 不要重试，或重试到了设置的最大次数，则返回失败
			if len(retry) == 0 || retry[0].Count == 0 {
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = dDefaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv 读取指定字节的数据
func (that *Conn) Recv(bufferSize int, retry ...Retry) ([]byte, error) {
	var err error               // Reading error.
	var size int                // Reading size.
	var data []byte             // Buffer object.
	var remoteAddr *net.UDPAddr // Current remote address for reading.
	if bufferSize > 0 {
		data = make([]byte, bufferSize)
	} else {
		data = make([]byte, dDefaultReadBufferSize)
	}
	for {
		size, remoteAddr, err = that.ReadFromUDP(data)
		if err == nil {
			that.remoteAddr = remoteAddr
		}
		if err != nil {
			// 如果链接已经断开
			if err == io.EOF {
				break
			}
			if len(retry) > 0 {
				// 重试到了最大值
				if retry[0].Count == 0 {
					break
				}
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = dDefaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			break
		}
		break
	}
	return data[:size], err
}

// SendRecv 发送数据到远端，并等待接收远端数据
func (that *Conn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := that.Send(data, retry...); err == nil {
		return that.Recv(receive, retry...)
	} else {
		return nil, err
	}
}

// SetDeadline 设置链接的读写超时时间
func (that *Conn) SetDeadline(t time.Time) error {
	err := that.UDPConn.SetDeadline(t)
	if err == nil {
		that.recvDeadline = t
		that.sendDeadline = t
	}
	return err
}

// SetRecvDeadline 单独设置链接的读超时时间
func (that *Conn) SetRecvDeadline(t time.Time) error {
	err := that.SetReadDeadline(t)
	if err == nil {
		that.recvDeadline = t
	}
	return err
}

// SetSendDeadline 单独设置链接的写超时时间
func (that *Conn) SetSendDeadline(t time.Time) error {
	err := that.SetWriteDeadline(t)
	if err == nil {
		that.sendDeadline = t
	}
	return err
}

// RecvWithTimeout 读取数据，如果指定时间内未读取到数据，则报错返回
func (that *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err := that.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer that.SetRecvDeadline(time.Time{})
	data, err = that.Recv(length, retry...)
	return
}

// SendWithTimeout 发送数据，在指定的时间内如果还为发送完毕，则报错返回
func (that *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err := that.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer that.SetSendDeadline(time.Time{})
	err = that.Send(data, retry...)
	return
}

// SendRecvWithTimeout 发送数据，并阻塞等待读取数据，在指定的时间内为读取成功，则报错返回
func (that *Conn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := that.Send(data, retry...); err == nil {
		return that.RecvWithTimeout(receive, timeout, retry...)
	} else {
		return nil, err
	}
}

// SetRecvBufferWait 设置从连接读取所有数据时的缓冲区等待超时。
//等待时间不能太长，否则可能会延迟从远程地址接收数据。
func (that *Conn) SetRecvBufferWait(d time.Duration) {
	that.recvBufferWait = d
}

// RemoteAddr 获取远端地址
func (that *Conn) RemoteAddr() net.Addr {
	//return c.conn.RemoteAddr()
	return that.remoteAddr
}
