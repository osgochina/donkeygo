package dipv4_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/net/dipv4"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestIP(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		ipUint := dipv4.Ip2long("192.168.1.1")
		t.Assert(ipUint, 3232235777)
		ipStr := dipv4.Long2ip(3232235777)
		t.Assert(ipStr, "192.168.1.1")
		t.Assert(dipv4.Validate("192.168.1.1"), true)
		t.Assert(dipv4.Validate("192.168.1.1111"), false)
		host, port := dipv4.ParseAddress("192.168.1.1:90")
		t.Assert(host, "192.168.1.1")
		t.Assert(port, 90)
		host, port = dipv4.ParseAddress("192.168.1.190:aa")
		t.Assert(host, "")
		t.Assert(port, 0)
	})
}

func TestIntranetIp(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		fmt.Println(dipv4.GetIpArray())
		fmt.Println(dipv4.GetIntranetIp())
		fmt.Println(dipv4.GetIntranetIpArray())
		fmt.Println(dipv4.IsIntranet("127.0.0.1"))
	})
}

func TestMac(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		fmt.Println(dipv4.GetMac())
		fmt.Println(dipv4.GetMacArray())
	})
}

func TestLookup(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		fmt.Println(dipv4.GetHostByName("www.baidu.com"))
		fmt.Println(dipv4.GetHostsByName("www.baidu.com"))
		fmt.Println(dipv4.GetNameByAddr("8.8.8.8"))
	})
}
