package drpc

import (
	"context"
	"donkeygo/container/dmap"
	"donkeygo/drpc/codec"
	"donkeygo/drpc/message"
	"donkeygo/drpc/status"
	"github.com/gogf/gf/util/gconv"
	"reflect"
	"sync"
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
	PeekMeta(key string) interface{}

	// VisitMeta 浏览消息的元数据
	VisitMeta(f func(key, value interface{}) bool)

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

	// SetMeta 设置指定key的值
	SetMeta(key, value string)

	// AddTFilterId 设置回复消息传输层的编码过滤方法id
	AddTFilterId(filterID ...byte)
}

var emptyValue = reflect.Value{}

// handlerCtx是 PushCtx 和 CallCtx 的底层公共实例
type handlerCtx struct {
	sess    *session
	input   message.Message
	output  message.Message
	handler *Handler
	arg     reflect.Value // 消息传入的参数
	//callCmd         *callCmd
	swap            *dmap.Map
	start           int64
	cost            time.Duration
	pluginContainer *PluginContainer
	stat            *status.Status
	context         context.Context
}

//newReadHandleCtx 创建一个给request/response或push使用的上下文
func newReadHandleCtx() *handlerCtx {
	c := new(handlerCtx)
	c.input = message.NewMessage()
	c.input.SetNewBody(c.binding)
	c.output = message.NewMessage()
	return c
}

//会话上下文生成池
var handlerCtxPool = sync.Pool{
	New: func() interface{} {
		return newReadHandleCtx()
	},
}

func (that *handlerCtx) reInit(s *session) {
	that.sess = s
	that.swap = s.socket.Swap().Clone(true)
}

func (that *handlerCtx) clean() {
	that.sess = nil
	that.handler = nil
	that.arg = emptyValue
	that.swap = nil
	//that.callCmd = nil
	that.cost = 0
	that.stat = nil
	that.context = nil
	that.input.Reset(message.WithNewBody(that.binding))
	that.output.Reset()
}

func (that *handlerCtx) Endpoint() Endpoint {
	return that.sess.Endpoint()
}

func (that *handlerCtx) Session() CtxSession {
	return that.sess
}

func (that *handlerCtx) Input() message.Message {
	return that.input
}

func (that *handlerCtx) Output() message.Message {
	return that.output
}

func (that *handlerCtx) Swap() *dmap.Map {
	return that.swap
}

func (that *handlerCtx) Seq() int32 {
	return that.input.Seq()
}

func (that *handlerCtx) ServiceMethod() string {
	return that.input.ServiceMethod()
}

func (that *handlerCtx) ResetServiceMethod(serviceMethod string) {
	that.input.SetServiceMethod(serviceMethod)
}

func (that *handlerCtx) PeekMeta(key string) interface{} {
	return that.input.Meta().Get(key)
}

func (that *handlerCtx) VisitMeta(f func(key, value interface{}) bool) {
	that.input.Meta().Iterator(f)
}

func (that *handlerCtx) CopyMeta() *dmap.Map {
	return that.input.Meta().Clone(true)
}

func (that *handlerCtx) SetMeta(key, value string) {
	that.output.Meta().Set(key, value)
}

func (that *handlerCtx) GetBodyCodec() byte {
	return that.input.BodyCodec()
}

func (that *handlerCtx) SetBodyCodec(bodyCodec byte) {
	that.output.SetBodyCodec(bodyCodec)
}

func (that *handlerCtx) AddTFilterId(filterID ...byte) {
	_ = that.output.PipeTFilter().Append(filterID...)
}

func (that *handlerCtx) IP() string {
	return that.sess.RemoteAddr().String()
}

func (that *handlerCtx) RealIP() string {
	realIP := gconv.String(that.PeekMeta(message.MetaRealIP))
	if len(realIP) > 0 {
		return realIP
	}
	return that.sess.RemoteAddr().String()
}

func (that *handlerCtx) Context() context.Context {
	if that.context == nil {
		return that.input.Context()
	}
	return that.context
}

func (that *handlerCtx) setContext(ctx context.Context) {
	that.context = ctx
}

func (that *handlerCtx) StatusOK() bool {
	return that.stat.OK()
}

func (that *handlerCtx) Status() *status.Status {
	return that.stat
}

func (that *handlerCtx) InputBodyBytes() []byte {
	b, ok := that.input.Body().(*[]byte)
	if !ok {
		return nil
	}
	return *b
}

func (that *handlerCtx) Bind(v interface{}) (byte, error) {
	b := that.InputBodyBytes()
	if b == nil {
		return codec.NilCodecID, nil
	}
	that.input.SetBody(v)
	err := that.input.UnmarshalBody(b)
	return that.input.BodyCodec(), err
}

func (that *handlerCtx) ReplyBodyCodec() byte {
	id := that.output.BodyCodec()
	if id != codec.NilCodecID {
		return id
	}
	id, ok := GetAcceptBodyCodec(that.input.Meta())
	if ok {
		if _, err := codec.Get(id); err == nil {
			that.output.SetBodyCodec(id)
			return id
		}
	}
	id = that.input.BodyCodec()
	that.output.SetBodyCodec(id)
	return id
}

//读取消息时候，执行该方法，作用是根据消息类型对消息进行不同的构造操作
func (that *handlerCtx) binding(header message.Header) (body interface{}) {
	that.start = that.sess.timeNow()
	that.pluginContainer = that.sess.endpoint.pluginContainer
	switch header.MType() {
	case message.TypeReply:
		return that.bindReply(header)
	case message.TypePush:
		return that.bindPush(header)
	case message.TypeCall:
		return that.bindCall(header)
	default:
		that.stat = statCodeMTypeNotAllowed
		return nil
	}
}

func (that *handlerCtx) handle() {

}
