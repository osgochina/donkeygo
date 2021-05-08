package drpc

import (
	"github.com/gogf/gf/os/glog"
	"net"
)

// Plugin 插件的基础对象
type Plugin interface {
	Name() string
}

// BeforeNewEndpointPlugin 创建Endpoint之前触发该事件
type BeforeNewEndpointPlugin interface {
	Plugin
	BeforeNewEndpoint(*EndpointConfig, *PluginContainer) error
}

// beforeNewEndpoint 在创建endpoint之前执行已定义的插件。
func (that *PluginContainer) beforeNewEndpoint(endpointConfig *EndpointConfig) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(BeforeNewEndpointPlugin); ok {
			if err = _plugin.BeforeNewEndpoint(endpointConfig, that); err != nil {
				glog.Fatalf("[BeforeNewEndpoint:%s] %s", plugin.Name(), err.Error())
				return
			}
		}
	}
}

// AfterNewEndpointPlugin 创建Endpoint之后触发该事件
type AfterNewEndpointPlugin interface {
	Plugin
	AfterNewEndpoint(EarlyEndpoint) error
}

// afterNewEndpoint 创建Endpoint之后执行已定义的插件
func (that *PluginContainer) afterNewEndpoint(e EarlyEndpoint) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterNewEndpointPlugin); ok {
			if err = _plugin.AfterNewEndpoint(e); err != nil {
				glog.Fatalf("[AfterNewEndpoint:%s] %s", plugin.Name(), err.Error())
				return
			}
		}
	}
}

// AfterRegRouterPlugin 路由注册成功触发该事件
type AfterRegRouterPlugin interface {
	Plugin
	AfterRegRouter(*Handler) error
}

// afterRegRouter 路由注册成功触发该事件
func (that *pluginSingleContainer) afterRegRouter(h *Handler) {
	var err error
	for _, plugin := range that.plugins {
		if _plugin, ok := plugin.(AfterRegRouterPlugin); ok {
			if err = _plugin.AfterRegRouter(h); err != nil {
				glog.Fatalf("[AfterRegRouter:%s] register handler:%s %s, error:%s", plugin.Name(), h.RouterTypeName(), h.Name(), err.Error())
				return
			}
		}
	}
}

// AfterListenPlugin 监听以后触发该事件
type AfterListenPlugin interface {
	Plugin
	AfterListen(net.Addr) error
}

// AfterDialPlugin 作为客户端链接到服务端成功以后触发该事件
type AfterDialPlugin interface {
	Plugin
	AfterDial(sess EarlySession, isRedial bool) *Status
}

// AfterAcceptPlugin 作为服务端，接收到客户端的链接后触发该事件
type AfterAcceptPlugin interface {
	Plugin
	AfterAccept(EarlySession) *Status
}

// BeforeWriteCallPlugin 写入CALL消息之前触发该事件
type BeforeWriteCallPlugin interface {
	Plugin
	BeforeWriteCall(WriteCtx) *Status
}

// AfterWriteCallPlugin 写入CALL消息成功之后触发该事件
type AfterWriteCallPlugin interface {
	Plugin
	AfterWriteCall(WriteCtx) *Status
}

// BeforeWriteReplyPlugin 写入Reply消息之前触发该事件
type BeforeWriteReplyPlugin interface {
	Plugin
	BeforeWriteReply(WriteCtx) *Status
}

// AfterWriteReplyPlugin 写入Reply消息成功之后触发该事件
type AfterWriteReplyPlugin interface {
	Plugin
	AfterWriteReply(WriteCtx) *Status
}

// BeforeWritePushPlugin 写入PUSH消息之前触发该事件
type BeforeWritePushPlugin interface {
	Plugin
	BeforeWritePush(WriteCtx) *Status
}

// AfterWritePushPlugin 写入PUSH消息成功之后触发该事件
type AfterWritePushPlugin interface {
	Plugin
	AfterWritePush(WriteCtx) *Status
}

// BeforeReadHeaderPlugin 执行读取Header之前触发该事件
type BeforeReadHeaderPlugin interface {
	Plugin
	BeforeReadHeader(EarlyCtx) error
}

// AfterReadCallHeaderPlugin 读取CALL消息的Header之后触发该事件
type AfterReadCallHeaderPlugin interface {
	Plugin
	AfterReadCallHeader(ReadCtx) *Status
}

// BeforeReadCallBodyPlugin 读取CALL 消息的body之前触发该事件
type BeforeReadCallBodyPlugin interface {
	Plugin
	BeforeReadCallBody(ReadCtx) *Status
}

// AfterReadCallBodyPlugin 读取CALL消息的body之后触发该事件
type AfterReadCallBodyPlugin interface {
	Plugin
	AfterReadCallBody(ReadCtx) *Status
}

// AfterReadPushHeaderPlugin 读取PUSH消息Header之后触发该事件
type AfterReadPushHeaderPlugin interface {
	Plugin
	AfterReadPushHeader(ReadCtx) *Status
}

// BeforeReadPushBodyPlugin 读取PUSH消息body之前触发该事件
type BeforeReadPushBodyPlugin interface {
	Plugin
	BeforeReadPushBody(ReadCtx) *Status
}

// AfterReadPushBodyPlugin 读取PUSH消息body之后触发该事件
type AfterReadPushBodyPlugin interface {
	Plugin
	AfterReadPushBody(ReadCtx) *Status
}

// AfterReadReplyHeaderPlugin 读取REPLY消息Header之前触发该事件
type AfterReadReplyHeaderPlugin interface {
	Plugin
	AfterReadReplyHeader(ReadCtx) *Status
}

// BeforeReadReplyBodyPlugin 读取REPLY消息body之前触发该事件
type BeforeReadReplyBodyPlugin interface {
	Plugin
	BeforeReadReplyBody(ReadCtx) *Status
}

// AfterReadReplyBodyPlugin 读取REPLY消息body之后触发该事件
type AfterReadReplyBodyPlugin interface {
	Plugin
	AfterReadReplyBody(ReadCtx) *Status
}

// AfterDisconnectPlugin 断开会话以后触发该事件
type AfterDisconnectPlugin interface {
	Plugin
	AfterDisconnect(BaseSession) *Status
}
