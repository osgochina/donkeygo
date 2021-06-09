package dudp

import (
	"errors"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/util/dconv"
	"net"
)

const (
	dDefaultServer = "default"
)

// Server UDP server
type Server struct {
	conn    *Conn       // UDP server connection object.
	address string      // UDP server listening address.
	handler func(*Conn) // Handler for UDP connection.
}

var (
	// serverMapping 用于实例名到它的UDP服务器映射。
	serverMapping = dmap.NewStrAnyMap(true)
)

// GetServer 获取一个UDP服务器对象
func GetServer(name ...interface{}) *Server {
	serverName := dDefaultServer
	if len(name) > 0 && name[0] != "" {
		serverName = dconv.String(name[0])
	}
	if s := serverMapping.Get(serverName); s != nil {
		return s.(*Server)
	}
	s := NewServer("", nil)
	serverMapping.Set(serverName, s)
	return s
}

// NewServer 创建一个UDP服务器对象
func NewServer(address string, handler func(*Conn), name ...string) *Server {
	s := &Server{
		address: address,
		handler: handler,
	}
	if len(name) > 0 && name[0] != "" {
		serverMapping.Set(name[0], s)
	}
	return s
}

// SetAddress 设置要监听的本地地址
func (that *Server) SetAddress(address string) {
	that.address = address
}

// SetHandler 设置收到包后的处理函数
func (that *Server) SetHandler(handler func(*Conn)) {
	that.handler = handler
}

// Close 关闭该UDP服务
func (that *Server) Close() error {
	return that.conn.Close()
}

// Run 启动该UDP服务器
func (that *Server) Run() error {
	if that.handler == nil {
		err := errors.New("start running failed: socket handler not defined")
		dlog.Error(err)
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", that.address)
	if err != nil {
		dlog.Error(err)
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		dlog.Error(err)
		return err
	}
	that.conn = NewConnByNetConn(conn)
	that.handler(that.conn)
	return nil
}
