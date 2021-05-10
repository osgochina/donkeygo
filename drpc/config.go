package drpc

import (
	"net"
	"time"
)

type EndpointConfig struct {

	// 网络类型; tcp, tcp4, tcp6, unix, unixpacket, kcp or quic"
	Network string

	//作为服务器角色时候，本地监听地址
	listenAddr net.Addr

	// 默认的消息体编码格式
	DefaultBodyCodec string
	//默认会话生命周期
	DefaultSessionAge time.Duration
	//默认单次请求生命周期
	DefaultContextAge time.Duration
	//慢处理定义时间
	slowCometDuration time.Duration

	//是否打印会话中请求的 body或 metadata
	PrintDetail bool
	// 是否统计消耗时间
	CountTime bool

	// 作为客户端角色时，请求服务端的超时时间
	DialTimeout time.Duration
	//作为客户端角色时，请求服务端时候，本地使用的地址端口
	localAddr net.Addr
	// 在链接中断时候，试图链接服务端的最大重试次数。仅限客户端角色使用
	RedialTimes int
	//仅限客户端角色使用 试图链接服务端时候，重试的时间间隔
	RedialInterval time.Duration
}

func (that *EndpointConfig) check() (err error) {

	return nil
}
