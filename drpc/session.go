package drpc

import (
	"context"
	"donkeygo/container/dmap"
	"donkeygo/drpc/message"
	"donkeygo/drpc/proto"
	"donkeygo/drpc/socket"
	"donkeygo/drpc/status"
	"net"
	"sync"
	"time"
)

// EarlySession 尚未启动 goroutine 读取数据的链接会话
type EarlySession interface {
	Endpoint() Endpoint

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 临时存储区内容
	Swap() *dmap.Map

	// SetID 设置seaside
	SetID(newID string)

	// ControlFD 原始链接的fd
	ControlFD(f func(fd uintptr)) error

	// ModifySocket 修改session的底层socket
	ModifySocket(fn func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc))

	// GetProtoFunc 获取协议方法
	GetProtoFunc() proto.ProtoFunc

	// EarlySend 在会话刚建立的时候临时发送消息，不执行任何中间件
	EarlySend(mType byte, serviceMethod string, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)

	// EarlyReceive 在会话刚建立的时候临时接受信息，不执行任何中间件
	EarlyReceive(newArgs message.NewBodyFunc, ctx ...context.Context) (input message.Message)

	// EarlyCall 在会话刚建立的时候临时调用call发送和接受消息，不执行任何中间件
	EarlyCall(serviceMethod string, args, reply interface{}, callSetting ...message.MsgSetting) (opStat *status.Status)

	// EarlyReply 在会话刚建立的时候临时回复消息，不执行任何中间件
	EarlyReply(req message.Message, body interface{}, stat *status.Status, setting ...message.MsgSetting) (opStat *status.Status)

	// RawPush 发送原始push消息，不执行任何中间件
	RawPush(serviceMethod string, args interface{}, setting ...message.MsgSetting) (opStat *status.Status)

	// SessionAge 获取session最大的生存周期
	SessionAge() time.Duration

	// ContextAge 获取 CALL 和 PUSH 消息的最大生存周期
	ContextAge() time.Duration

	// SetSessionAge 设置session的最大生存周期
	SetSessionAge(duration time.Duration)

	// SetContextAge 设置单个 CALL 和 PUSH 消息的最大生存周期
	SetContextAge(duration time.Duration)
}

// BaseSession 基础的session
type BaseSession interface {
	Endpoint() Endpoint

	// ID 获取id
	ID() string

	// LocalAddr 本地地址
	LocalAddr() net.Addr

	// RemoteAddr 远端地址
	RemoteAddr() net.Addr

	// Swap 返回交换区的内容
	Swap() *dmap.Map
}

// CtxSession 在处理程序上下文中传递的会话对象
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

	// AsyncCall 发送消息，并异步接收响应
	AsyncCall(serviceMethod string, args interface{}, result interface{}, callCmdChan chan<- CallCmd, setting ...message.MsgSetting) CallCmd

	// Call 发送消息并获得响应值
	Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) CallCmd

	// Push 发送消息，不接收响应，只返回发送状态
	Push(serviceMethod string, args interface{}, setting ...message.MsgSetting) *status.Status

	// SessionAge 获取session最大的生存周期
	SessionAge() time.Duration

	// ContextAge 获取 CALL 和 PUSH 消息的最大生存周期
	ContextAge() time.Duration
}

type Session interface {
	Endpoint() Endpoint

	// SetID 设置session id
	SetID(newID string)

	// Close 关闭session
	Close() error

	CtxSession
}

type session struct {
	endpoint       *endpoint
	getCallHandler func(serviceMethodPath string) (*Handler, bool)
	getPushHandler func(serviceMethodPath string) (*Handler, bool)
	timeNow        func() int64
	callCmdMap     *dmap.Map
	protoFuncList  []proto.ProtoFunc
	socket         socket.Socket
	closeNotifyCh  chan struct{}
	writeLock      sync.Mutex
	sessionAge     time.Duration
	contextAge     time.Duration
	sessionAgeLock sync.RWMutex
	contextAgeLock sync.RWMutex
	lock           sync.RWMutex
	seq            int32
	status         int32
	didCloseNotify int32

	//链接如果断开，重新拨号，只有作为客户端角色的时候才有效果
	redialForClientLocked func() bool
}
