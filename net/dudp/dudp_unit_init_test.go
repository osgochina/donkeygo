package dudp_test

import (
	"github.com/osgochina/donkeygo/container/darray"
)

var (
	ports = darray.NewIntArray(true)
)

func init() {
	for i := 9000; i <= 10000; i++ {
		ports.Append(i)
	}
}
