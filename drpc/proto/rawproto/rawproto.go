package rawproto

import (
	"donkeygo/drpc/proto"
	"io"
	"sync"
)

/**
rpc协议的原始格式 使用网络字节序，大端
{4 bytes 表示整个消息的长度}
{1 byte  表示协议的版本}
{1 byte  表示传输管道过滤器id的长度}
{传输管道过滤器id序列化后的内容}
# 以下的内容都是经过传输管道过滤器处理过的数据
{1 bytes 表示消息序列号长度}
{序列号 16进制的32位int32}
{1 byte 表示消息类型} # CALL:1;REPLY:2;PUSH:3
{1 byte 表示请求服务的方法长度 service method length}
{service method}
{2 bytes status length}
{status(json)}
{2 bytes metadata length}
{metadata(json)}
{1 byte bode codec id}
{body}
**/

var RawProtoFunc = func(rw proto.IOWithReadBuffer) proto.Proto {
	return &rawProto{
		id:   6,
		name: "raw",
		r:    rw,
		w:    rw,
	}
}

var _ proto.Proto = new(rawProto)

type rawProto struct {
	r    io.Reader
	w    io.Writer
	rMu  sync.Mutex
	name string
	id   byte
}

func (that *rawProto) Version() (byte, string) {
	return that.id, that.name
}

func (that *rawProto) Pack(m proto.Message) error {
	return nil
}

func (that *rawProto) Unpack(m proto.Message) error {
	return nil
}
