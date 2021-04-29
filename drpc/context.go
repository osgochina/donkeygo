package drpc

import (
	"context"
	"donkeygo/container/dmap"
	"donkeygo/drpc/message"
	"donkeygo/drpc/status"
	"reflect"
	"time"
)

// BaseCtx 基础上下文
type BaseCtx interface {

	// Endpoint 获取当前Endpoint
	Endpoint() Endpoint

	// Session 返回当前的session
	Session() CtxSession

	// IP 返回远端ip
	IP() string

	// RealIP 返回远端真实ip
	RealIP() string

	// Swap 返回自定义交换区数据
	Swap() *dmap.Map

	// Context 获取上下文
	Context() context.Context
}

// WriteCtx 写消息时使用的上下文方法
type WriteCtx interface {
	BaseCtx

	// Output 将要发送的消息对象
	Output() message.Message

	// StatusOK 状态是否ok
	StatusOK() bool

	// Status 当前步骤的状态
	Status() *status.Status
}

// InputCtx 该上下文是一个公共上下文
type inputCtx interface {
	BaseCtx

	// Seq 获取消息的序列号
	Seq() int32

	// PeekMeta 窥视消息的元数据
	PeekMeta(key string) []byte

	// VisitMeta 浏览消息的元数据
	VisitMeta(f func(key, value []byte))

	// CopyMeta 赋值消息的元数据
	CopyMeta() *dmap.Map

	// ServiceMethod 该消息需要访问的服务名
	ServiceMethod() string

	// ResetServiceMethod 重置该消息将要访问的服务名
	ResetServiceMethod(string)
}

// ReadCtx 读取消息使用的上下文
type ReadCtx interface {
	inputCtx

	// Input 获取传入的消息
	Input() message.Message

	// StatusOK 状态是否ok
	StatusOK() bool

	// Status 当前步骤的状态
	Status() *status.Status
}

// PushCtx push消息使用的上下文
type PushCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte
}

// CallCtx call消息使用的上下文
type CallCtx interface {
	inputCtx

	// Input 获取传入的消息
	Input() message.Message

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// Output 将要发送的消息对象
	Output() message.Message

	// ReplyBodyCodec 获取响应消息的编码格式
	ReplyBodyCodec() byte

	// SetBodyCodec 设置响应消息的编码格式
	SetBodyCodec(byte)

	// AddMeta 添加元数据
	AddMeta(key, value string)

	// SetMeta 设置指定key的值
	SetMeta(key, value string)

	// AddTFilterId 设置回复消息传输层的编码过滤方法id
	AddTFilterId(filterID ...byte)
}

// UnknownPushCtx 未知push消息的上下文
type UnknownPushCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// InputBodyBytes 传入消息体
	InputBodyBytes() []byte

	// Bind 如果push消息是未知的消息，则使用v对象解析消息内容
	Bind(v interface{}) (bodyCodec byte, err error)
}

// UnknownCallCtx 未知call消息的上下文
type UnknownCallCtx interface {
	inputCtx

	// GetBodyCodec 获取当前消息的编码格式
	GetBodyCodec() byte

	// InputBodyBytes 传入消息体
	InputBodyBytes() []byte

	// Bind 如果push消息是未知的消息，则使用v对象解析消息内容
	Bind(v interface{}) (bodyCodec byte, err error)

	// SetBodyCodec 设置回复消息的编码格式
	SetBodyCodec(byte)

	// AddMeta 添加元数据
	AddMeta(key, value string)

	// SetMeta 设置指定key的值
	SetMeta(key, value string)

	// AddTFilterId 设置回复消息传输层的编码过滤方法id
	AddTFilterId(filterID ...byte)
}

// handlerCtx是 PushCtx 和 CallCtx 的底层公共实例
type handlerCtx struct {
	sess    *session
	input   message.Message
	output  message.Message
	handler *Handler
	arg     reflect.Value // 消息传入的参数
	//callCmd         *callCmd
	swap    *dmap.Map
	start   int64
	cost    time.Duration
	stat    *status.Status
	context context.Context
}

func (that *handlerCtx) reInit(s *session) {
	that.sess = s
	that.swap = s.socket.Swap().Clone(true)
}
