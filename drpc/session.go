package drpc

import (
	"context"
	"donkeygo/container/dmap"
	"donkeygo/drpc/message"
	"donkeygo/drpc/proto"
	"donkeygo/drpc/socket"
	"donkeygo/drpc/status"
	"github.com/gogf/gf/os/glog"
	"net"
	"sync"
	"sync/atomic"
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

	//// AsyncCall 发送消息，并异步接收响应
	//AsyncCall(serviceMethod string, args interface{}, result interface{}, callCmdChan chan<- CallCmd, setting ...message.MsgSetting) CallCmd
	//
	//// Call 发送消息并获得响应值
	//Call(serviceMethod string, args interface{}, result interface{}, setting ...message.MsgSetting) CallCmd

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

// NOTE: Do not change the order
const (
	statusPreparing int32 = iota
	statusOk
	statusActiveClosing
	statusActiveClosed
	statusPassiveClosing
	statusPassiveClosed
	statusRedialing
	statusRedialFailed
)

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

func newSession(e *endpoint, conn net.Conn, protoFunc []proto.ProtoFunc) *session {
	var s = &session{
		endpoint:       e,
		getCallHandler: e.router.subRouter.getCall,
		getPushHandler: e.router.subRouter.getPush,
		timeNow:        e.timeNow,
		protoFuncList:  protoFunc,
		status:         statusPreparing,
		socket:         socket.NewSocket(conn, protoFunc...),
		closeNotifyCh:  make(chan struct{}),
		callCmdMap:     dmap.New(true),
		sessionAge:     e.defaultSessionAge,
		contextAge:     e.defaultContextAge,
	}
	return s
}

//原子性修改session的状态
func (that *session) changeStatus(stat int32) {
	atomic.StoreInt32(&that.status, stat)
}

//原子性尝试修改session的状态
// 从fromList的多个状态修改成to的状态，修改成功则返回true，失败返回false
func (that *session) tryChangeStatus(to int32, fromList ...int32) (changed bool) {
	for _, from := range fromList {
		if atomic.CompareAndSwapInt32(&that.status, from, to) {
			return true
		}
	}
	return false
}

//判断session的状态是否在checkList中，如果在checkList中，则返回true，否则返回false
func (that *session) checkStatus(checkList ...int32) bool {
	stat := atomic.LoadInt32(&that.status)
	for _, v := range checkList {
		if v == stat {
			return true
		}
	}
	return false
}

//原子性的获取当前session的状态
func (that *session) getStatus() int32 {
	return atomic.LoadInt32(&that.status)
}

//判断session的状态是否是开始或者将要结束
func (that *session) goonRead() bool {
	return that.checkStatus(statusOk, statusActiveClosing)
}

// 通知session准备关闭
func (that *session) notifyClosed() {
	if atomic.CompareAndSwapInt32(&that.didCloseNotify, 0, 1) {
		close(that.closeNotifyCh)
	}
}

// CloseNotify session将要关闭的通知
func (that *session) CloseNotify() <-chan struct{} {
	return that.closeNotifyCh
}

// IsActiveClosed 判断链接是否处已经关闭，并且并且是主动关闭的
func (that *session) IsActiveClosed() bool {
	return that.checkStatus(statusActiveClosed)
}

// IsPassiveClosed 判断链接是否已经关闭，并且是被动关闭的
func (that *session) IsPassiveClosed() bool {
	return that.checkStatus(statusPassiveClosed)
}

// Health 判断session会话是否可用
func (that *session) Health() bool {
	s := that.getStatus()
	if s == statusOk {
		return true
	}
	if that.redialForClientLocked == nil {
		return false
	}
	if s == statusPassiveClosed {
		return true
	}
	return false
}

func (that *session) Endpoint() Endpoint {
	return that.endpoint
}

// ID 获取session的id
func (that *session) ID() string {
	return that.socket.ID()
}

// SetID 修改id
func (that *session) SetID(newID string) {
	oldID := that.ID()
	if oldID == newID {
		return
	}
	that.socket.SetID(newID)
	hub := that.endpoint.sessHub
	hub.set(that)
	hub.delete(oldID)
	glog.Info("session changes id: %s -> %s", oldID, newID)
}

// ControlFD 处理底层fd
func (that *session) ControlFD(f func(fd uintptr)) error {
	that.lock.RLock()
	defer that.lock.RUnlock()
	return that.socket.ControlFD(f)
}

//获取会话的原始链接
func (that *session) getConn() net.Conn {
	return that.socket.Raw()
}

func (that *session) ModifySocket(fn func(conn net.Conn) (modifiedConn net.Conn, newProtoFunc proto.ProtoFunc)) {
	//conn := s.getConn()
	//modifiedConn, newProtoFunc := fn(conn)
	//isModifiedConn := modifiedConn != nil
	//isNewProtoFunc := newProtoFunc != nil
	//if isNewProtoFunc {
	//	s.protoFuncs = s.protoFuncs[:0]
	//	s.protoFuncs = append(s.protoFuncs, newProtoFunc)
	//}
	//if !isModifiedConn && !isNewProtoFunc {
	//	return
	//}
	//var pub goutil.Map
	//if s.socket.SwapLen() > 0 {
	//	pub = s.socket.Swap()
	//}
	//id := s.ID()
	//s.socket.Reset(modifiedConn, s.protoFuncs...)
	//s.socket.Swap(pub) // set the old swap
	//s.socket.SetID(id)
}

// GetProtoFunc 获取协议方法
func (that *session) GetProtoFunc() proto.ProtoFunc {
	if len(that.protoFuncList) > 0 && that.protoFuncList[0] != nil {
		return that.protoFuncList[0]
	}
	return socket.DefaultProtoFunc()
}

// LocalAddr 获取本地监听地址
func (that *session) LocalAddr() net.Addr {
	return that.socket.LocalAddr()
}

// RemoteAddr 获取远程链接的地址
func (that *session) RemoteAddr() net.Addr {
	return that.socket.RemoteAddr()
}

// SessionAge 获取session的生存周期
func (that *session) SessionAge() time.Duration {
	that.sessionAgeLock.RLock()
	age := that.sessionAge
	that.sessionAgeLock.RUnlock()
	return age
}

// SetSessionAge 设置会话的最大生命周期
func (that *session) SetSessionAge(duration time.Duration) {
	that.sessionAgeLock.Lock()
	that.sessionAge = duration
	if duration > 0 {
		_ = that.socket.SetReadDeadline(time.Now().Add(duration))
	} else {
		_ = that.socket.SetReadDeadline(time.Time{})
	}
	that.sessionAgeLock.Unlock()
}

// ContextAge 获取CALL 或者PUSH上下文的最大生命周期
func (that *session) ContextAge() time.Duration {
	that.contextAgeLock.RLock()
	age := that.contextAge
	that.contextAgeLock.RUnlock()
	return age
}

// SetContextAge 设置CALL 或者PUSH上下文的最大生命周期
func (that *session) SetContextAge(duration time.Duration) {
	that.contextAgeLock.Lock()
	that.contextAge = duration
	that.contextAgeLock.Unlock()
}

func (that *session) Close() error {
	that.lock.Lock()
	defer that.lock.Unlock()
	//return that.closeLocked()
	return nil
}
