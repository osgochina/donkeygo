package drpc

import (
	"crypto/tls"
	"donkeygo/drpc/socket"
	"donkeygo/drpc/status"
	"net"
	"sync"
	"time"
)

type BaseEndpoint interface {

	// Close 关闭该端点
	Close() (err error)

	// CountSession 统计该端点上的session数量
	CountSession() int

	// GetSession 获取指定ID的 session
	GetSession(sessionID string) (Session, bool)

	// RangeSession 循环迭代session
	RangeSession(fn func(sess Session) bool)

	// SetTLSConfig 设置证书配置
	SetTLSConfig(tlsConfig *tls.Config)

	// SetTLSConfigFromFile 从文件中读取证书并设置证书配置
	SetTLSConfigFromFile(tlsCertFile, tlsKeyFile string, insecureSkipVerifyForClient ...bool) error

	// TLSConfig tls配置对象
	TLSConfig() *tls.Config

	// PluginContainer 插件容器对象
	PluginContainer() *PluginContainer
}

type EarlyEndpoint interface {
	BaseEndpoint

	// Router 获取路由对象
	Router() *Router

	// SubRoute 获取分组路由对象
	SubRoute(pathPrefix string, plugin ...Plugin) *SubRouter

	// RouteCall 通过struct注册CALL类型的处理程序，并且返回注册的路径列表
	RouteCall(ctrlStruct interface{}, plugin ...Plugin) []string
	// RouteCallFunc 通过func注册CALL类型的处理程序，并且返回单个注册路径
	RouteCallFunc(callHandleFunc interface{}, plugin ...Plugin) string
	// RoutePush 通过struct注册PUSH类型的处理程序，并且返回注册的路径列表
	RoutePush(ctrlStruct interface{}, plugin ...Plugin) []string
	// RoutePushFunc 通过func注册PUSH类型的处理程序，并且返回单个注册路径
	RoutePushFunc(pushHandleFunc interface{}, plugin ...Plugin) string
	// SetUnknownCall 设置默认处理程序，当没有找到CALL的处理程序时将调用该处理程序。
	SetUnknownCall(fn func(UnknownCallCtx) (interface{}, *status.Status), plugin ...Plugin)
	// SetUnknownPush 设置默认处理程序，当没有找到PUSH的处理程序时将调用该处理程序。
	SetUnknownPush(fn func(UnknownPushCtx) *status.Status, plugin ...Plugin)
}

type Endpoint interface {
	EarlyEndpoint

	// ListenAndServe 打开服务监听
	ListenAndServe(protoFunc ...socket.ProtoFunc) error

	// Dial 作为客户端链接到指定的服务
	Dial(addr string, protoFunc ...socket.ProtoFunc) (Session, *status.Status)

	// ServeConn 传入指定的conn，生成session
	// 提示：
	// 1. 不支持断开链接后自动重拨
	// 2. 不检查TLS
	// 3. 执行 PostAcceptPlugin 插件
	ServeConn(conn net.Conn, protoFunc ...socket.ProtoFunc) (Session, *status.Status)
}

type endpoint struct {
	router            *Router
	pluginContainer   *PluginContainer
	sessHub           *SessionHub
	closeCh           chan struct{}
	defaultSessionAge time.Duration
	defaultContextAge time.Duration
	tlsConfig         *tls.Config
	slowCometDuration time.Duration
	timeNow           func() int64
	mu                sync.Mutex
	network           string
	defaultBodyCodec  byte
	printDetail       bool
	countTime         bool

	//只有作为server角色时候才有该对象
	listerAddr net.Addr
	listeners  map[net.Listener]struct{}

	//只有作为client角色时候才有该对象
	dialer *Dialer
}

// PluginContainer 获取端点的插件容器
func (that *endpoint) PluginContainer() *PluginContainer {
	return that.pluginContainer
}

// TLSConfig 获取该端点的证书信息
func (that *endpoint) TLSConfig() *tls.Config {
	return that.tlsConfig
}

// SetTLSConfig 设置该端点的证书信息
func (that *endpoint) SetTLSConfig(tlsConfig *tls.Config) {
	that.tlsConfig = tlsConfig
	that.dialer.tlsConfig = tlsConfig
}

// SetTLSConfigFromFile 通过文件生成端点的证书信息
func (that *endpoint) SetTLSConfigFromFile(tlsCertFile, tlsKeyFile string, insecureSkipVerifyForClient ...bool) error {
	tlsConfig, err := NewTLSConfigFromFile(tlsCertFile, tlsKeyFile, insecureSkipVerifyForClient...)
	if err == nil {
		that.SetTLSConfig(tlsConfig)
	}
	return err
}
