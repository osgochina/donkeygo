package dipv4

import (
	"encoding/binary"
	"fmt"
	"github.com/osgochina/donkeygo/text/dregex"
	"net"
	"strconv"
)

// Ip2long 把ip地址转换成一个uint32的整数
func Ip2long(ip string) uint32 {
	netIp := net.ParseIP(ip)
	if netIp == nil {
		return 0
	}
	return binary.BigEndian.Uint32(netIp.To4())
}

// Long2ip 把一个uint32的整数转换成一个ip地址
func Long2ip(long uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, long)
	return net.IP(ipByte).String()
}

// Validate 验证传入的参数是否是合法的ip地址
func Validate(ip string) bool {
	return dregex.IsMatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
}

// ParseAddress 把连起来的ip端口转换成分开的ip和端口
// Eg: 192.168.1.1:80 -> 192.168.1.1, 80
func ParseAddress(address string) (string, int) {
	match, err := dregex.MatchString(`^(.+):(\d+)$`, address)
	if err == nil && len(match) > 1 {
		i, _ := strconv.Atoi(match[2])
		return match[1], i
	}
	return "", 0
}

// GetSegment 返回传入ip地址的网段
// Eg: 192.168.2.102 -> 192.168.2
func GetSegment(ip string) string {
	match, err := dregex.MatchString(`^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$`, ip)
	if err != nil || len(match) < 4 {
		return ""
	}
	return fmt.Sprintf("%s.%s.%s", match[1], match[2], match[3])
}
