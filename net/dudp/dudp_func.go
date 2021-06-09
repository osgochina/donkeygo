package dudp

import "net"

// NewNetConn 创建一个udp链接
// 远端的地址端口，必传
// 本地使用的地址端口，可选
func NewNetConn(remoteAddress string, localAddress ...string) (*net.UDPConn, error) {
	var err error
	var remoteAddr, localAddr *net.UDPAddr
	remoteAddr, err = net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		return nil, err
	}
	if len(localAddress) > 0 {
		localAddr, err = net.ResolveUDPAddr("udp", localAddress[0])
		if err != nil {
			return nil, err
		}
	}
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Send 对远端的一个地址端口，发送udp包
// address: 远端地址
// data : 要发送的数据
// retry : 重试配置
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendRecv 对远端的一个地址端口，发送udp包，并等待接收
// address: 远端地址
// data : 要发送的数据
// receive : 要接收的字节数
// retry : 重试配置
func SendRecv(address string, data []byte, receive int, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecv(data, receive, retry...)
}
