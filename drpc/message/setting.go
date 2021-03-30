package message

import (
	"context"
	"donkeygo/drpc/status"
	"errors"
	"math"
)

var (
	messageSizeLimit uint32 = math.MaxUint32
	// ErrExceedMessageSizeLimit error
	ErrExceedMessageSizeLimit = errors.New("size of package exceeds limit")
)

type MsgSetting func(Message)

//什么也不做
func WithNothing() MsgSetting {
	return func(Message) {}
}

//设置消息的上下文对象
func WithContext(ctx context.Context) MsgSetting {
	return func(m Message) {
		m.(*message).ctx = ctx
	}
}

//设置消息的服务器接口名
func WithServiceMethod(serviceMethod string) MsgSetting {
	return func(m Message) {
		m.SetServiceMethod(serviceMethod)
	}
}

//设置消息的状态
func WithStatus(stat *status.Status) MsgSetting {
	return func(m Message) {
		m.SetStatus(stat)
	}
}

//添加消息的元数据
func WithSetMeta(key, value string) MsgSetting {
	return func(m Message) {
		m.Meta().Set(key, value)
	}
}

//删除消息元数据
func WithDelMeta(key string) MsgSetting {
	return func(m Message) {
		m.Meta().Remove(key)
	}
}

//设置消息的消息体编码格式
func WithBodyCodec(bodyCodec byte) MsgSetting {
	return func(m Message) {
		m.SetBodyCodec(bodyCodec)
	}
}

//设置消息体的内容
func WithBody(body interface{}) MsgSetting {
	return func(m Message) {
		m.SetBody(body)
	}
}

//设置创建消息体的函数
func WithNewBody(newBodyFunc NewBodyFunc) MsgSetting {
	return func(m Message) {
		m.SetNewBody(newBodyFunc)
	}
}

//设置消息的管道类型
func WithXFerPipe(filterID ...byte) MsgSetting {
	return func(m Message) {
		if err := m.PipeTFilter().Append(filterID...); err != nil {
			panic(err)
		}
	}
}

//获取消息的最大长度
func MsgSizeLimit() uint32 {
	return messageSizeLimit
}

//设置消息的最大长度
func SetMsgSizeLimit(maxMessageSize uint32) {
	if maxMessageSize <= 0 {
		messageSizeLimit = math.MaxUint32
	} else {
		messageSizeLimit = maxMessageSize
	}
}

//检查消息的最大长度
func checkMessageSize(messageSize uint32) error {
	if messageSize > messageSizeLimit {
		return ErrExceedMessageSizeLimit
	}
	return nil
}
