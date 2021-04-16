package drpc

import (
	"donkeygo/container/dmap"
	"donkeygo/drpc/socket"
	"net"
)

type CtxSession interface {

	// ID 获取id
	ID() string

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 返回交换区的内容
	Swap() *dmap.Map

	// CloseNotify 返回该链接被关闭时候的通知
	CloseNotify() <-chan struct{}

	// Health 检查该session是否健康
	Health() bool

	// SetID 设置session id
	SetID(newID string)

	// GetProtoFunc 获取当前session使用的通信协议
	GetProtoFunc() socket.ProtoFunc
}
