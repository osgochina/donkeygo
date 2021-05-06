package drpc

import (
	"crypto/tls"
	"net"
	"time"
)

type Dialer struct {
	network        string
	localAddr      net.Addr
	tlsConfig      *tls.Config
	dialTimeout    time.Duration
	redialInterval time.Duration
	redialTimes    int32
}
