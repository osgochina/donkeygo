package dconv

import (
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/os/dtime"
	"time"
)

// Duration 把参数转换成Duration格式
func Duration(any interface{}) time.Duration {
	// It's already this type.
	if v, ok := any.(time.Duration); ok {
		return v
	}
	s := String(any)
	if !utils.IsNumeric(s) {
		d, _ := dtime.ParseDuration(s)
		return d
	}
	return time.Duration(Int64(any))
}
