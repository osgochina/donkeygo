package dtcp

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"time"
)

//读取缓冲区的间隔时间
const dRecvAllWaitTimeout = time.Millisecond

//tcp链接对象
type Conn struct {
	net.Conn                     //tcp 链接的对象
	reader         *bufio.Reader // 读数据的缓冲区
	recvDeadline   time.Time     // 读数据超时时间
	sendDeadline   time.Time     //发送数据超时时间
	recvBufferWait time.Duration //读取数据到缓冲区等待时间
}

//NewConn 创建一个tcp链接
func NewConn(addr string, timeout ...time.Duration) (*Conn, error) {
	conn, err := NewNetConn(addr, timeout...)
	if err != nil {
		return nil, err
	}
	return NewConnByNetConn(conn), nil
}

//NewConnTLS 创建一个TLS加密的tcp链接
func NewConnTLS(addr string, tlsConfig *tls.Config) (*Conn, error) {
	if conn, err := NewNetConnTLS(addr, tlsConfig); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

//通过已有的net.Conn 创建一个tcp链接
func NewConnByNetConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:           conn,
		reader:         bufio.NewReader(conn),
		recvDeadline:   time.Time{},
		sendDeadline:   time.Time{},
		recvBufferWait: dRecvAllWaitTimeout,
	}
}

//发送数据到远端，没有任何缓冲，直接调用tcp write写入到系统缓冲区
func (that *Conn) Send(data []byte, retry ...Retry) error {
	for {
		if _, err := that.Write(data); err != nil {
			//链接被关闭
			if err == io.EOF {
				return err
			}
			//如果没有重试
			if len(retry) == 0 || retry[0].Count == 0 {
				return err
			}
			//需要重试
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

//读取数据
func (that *Conn) Recv(length int, retry ...Retry) ([]byte, error) {
	var err error
	var size int
	var index int
	var buffer []byte
	var bufferWait bool

	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = make([]byte, dDefaultReadBufferSize)
	}

	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = that.SetReadDeadline(time.Now().Add(that.recvBufferWait)); err != nil {
				return nil, err
			}
		}
		//读取
		size, err = that.reader.Read(buffer[index:])
		if size > 0 {
			index += size
			if length > 0 {
				//读取了指定大小的数据，则结束
				if index == length {
					break
				}
			} else {
				//未指定读取长度的情况下，判断当前读取的长度是否超过了读取buffer的大小
				if index > dDefaultReadBufferSize {
					buffer = append(buffer, make([]byte, dDefaultReadBufferSize)...)
				} else {
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			//链接断了
			if err == io.EOF {
				break
			}
			//超时了,重新设置超时时间
			if bufferWait && isTimeout(err) {
				if err = that.SetReadDeadline(that.recvDeadline); err != nil {
					return nil, err
				}
				err = nil
				break
			}
			//重试
			if len(retry) > 0 {
				// It fails even it retried.
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
		//只从缓冲区读取一次默认长度的数据
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

//获取以 \n 结尾的一行，但是返回的数据不包含 \n
func (that *Conn) RecvLine(retry ...Retry) ([]byte, error) {
	var err error
	var buffer []byte
	data := make([]byte, 0)
	for {
		buffer, err = that.Recv(1, retry...)
		if len(buffer) > 0 {
			if buffer[0] == '\n' {
				data = append(data, buffer[:len(buffer)-1]...)
				break
			} else {
				data = append(data, buffer...)
			}
		}
		if err != nil {
			break
		}
	}
	return data, err
}

//在指定的时间内读取数据
func (that *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {

	if err = that.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer that.SetRecvDeadline(time.Time{})
	data, err = that.Recv(length, retry...)
	return data, err
}

//在指定的时间内发送数据
func (that *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err = that.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer that.SetSendDeadline(time.Time{})
	err = that.Send(data, retry...)
	return
}

//发送并且读取
func (that *Conn) SendRecv(data []byte, length int, retry ...Retry) ([]byte, error) {
	if err := that.Send(data, retry...); err == nil {
		return that.Recv(length, retry...)
	} else {
		return nil, err
	}
}

//设置超时时间
func (that *Conn) SetDeadline(t time.Time) error {
	err := that.Conn.SetDeadline(t)
	if err == nil {
		that.recvDeadline = t
		that.sendDeadline = t
	}
	return err
}

//设置读取超时时间
func (that *Conn) SetRecvDeadline(t time.Time) error {
	err := that.SetReadDeadline(t)
	if err == nil {
		that.recvDeadline = t
	}
	return err
}

//设置发送超时时间
func (that *Conn) SetSendDeadline(t time.Time) error {
	err := that.SetWriteDeadline(t)
	if err == nil {
		that.sendDeadline = t
	}
	return err
}
