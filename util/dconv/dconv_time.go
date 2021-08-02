package dconv

import (
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/os/dtime"
	"time"
)

// Time 把任意类型any转换成 time.Time.
func Time(any interface{}, format ...string) time.Time {
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(time.Time); ok {
			return v
		}
	}
	if t := GTime(any, format...); t != nil {
		return t.Time
	}
	return time.Time{}
}

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

// GTime 把任意类型any转换成 *dtime.Time
func GTime(any interface{}, format ...string) *dtime.Time {
	if any == nil {
		return nil
	}
	if v, ok := any.(apiGTime); ok {
		return v.GTime(format...)
	}
	// It's already this type.
	if len(format) == 0 {
		if v, ok := any.(*dtime.Time); ok {
			return v
		}
		if t, ok := any.(time.Time); ok {
			return dtime.New(t)
		}
		if t, ok := any.(*time.Time); ok {
			return dtime.New(t)
		}
	}
	s := String(any)
	if len(s) == 0 {
		return dtime.New()
	}
	// Priority conversion using given format.
	if len(format) > 0 {
		t, _ := dtime.StrToTimeFormat(s, format[0])
		return t
	}
	if utils.IsNumeric(s) {
		return dtime.NewFromTimeStamp(Int64(s))
	} else {
		t, _ := dtime.StrToTime(s)
		return t
	}
}
