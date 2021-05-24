package dtime

import "time"

// Time is a wrapper for time.Time for additional features.
type Time struct {
	wrapper
}

// apiUnixNano is an interface definition commonly for custom time.Time wrapper.
type apiUnixNano interface {
	UnixNano() int64
}

func New(param ...interface{}) *Time {
	if len(param) > 0 {
		switch r := param[0].(type) {
		case time.Time:
			return NewFromTime(r)
		case *time.Time:
			return NewFromTime(*r)
		case Time:
			return &r
		case *Time:
			return r
		case string:
			if len(param) > 1 {
				switch t := param[1].(type) {
				case string:
					return NewFromStrFormat(r, t)
				case []byte:
					return NewFromStrFormat(r, string(t))
				}
			}
			return NewFromStr(r)
		case []byte:
			if len(param) > 1 {
				switch t := param[1].(type) {
				case string:
					return NewFromStrFormat(string(r), t)
				case []byte:
					return NewFromStrFormat(string(r), string(t))
				}
			}
			return NewFromStr(string(r))
		case int:
			return NewFromTimeStamp(int64(r))
		case int64:
			return NewFromTimeStamp(r)
		default:
			if v, ok := r.(apiUnixNano); ok {
				return NewFromTimeStamp(v.UnixNano())
			}
		}
	}
	return &Time{
		wrapper{time.Time{}},
	}
}

// Now 返回一个time对象，时间设置为当前时间
func Now() *Time {
	return &Time{
		wrapper{time.Now()},
	}
}

// NewFromTime 通过time.time对象创建一个dtime对象
func NewFromTime(t time.Time) *Time {
	return &Time{
		wrapper{t},
	}
}

// NewFromStr 通过字符串创建一个dtime.time对象
func NewFromStr(str string) *Time {
	if t, err := StrToTime(str); err == nil {
		return t
	}
	return nil
}

// NewFromTimeStamp 通过时间戳创建dtime.Time对象
func NewFromTimeStamp(timestamp int64) *Time {
	if timestamp == 0 {
		return &Time{}
	}
	var sec, nano int64
	if timestamp > 1e9 {
		for timestamp < 1e18 {
			timestamp *= 10
		}
		sec = timestamp / 1e9
		nano = timestamp % 1e9
	} else {
		sec = timestamp
	}
	return &Time{
		wrapper{time.Unix(sec, nano)},
	}
}

// NewFromStrFormat creates and returns a Time object with given string and
// custom format like: Y-m-d H:i:s.
// Note that it returns nil if there's error occurs.
func NewFromStrFormat(str string, format string) *Time {
	if t, err := StrToTimeFormat(str, format); err == nil {
		return t
	}
	return nil
}

// Timestamp 当前时间戳
func (that *Time) Timestamp() int64 {
	return that.UnixNano() / 1e9
}

// Millisecond 返回毫秒时间戳
func (that *Time) Millisecond() int {
	return that.Time.Nanosecond() / 1e6
}

// Microsecond 返回微秒的时间戳
func (that *Time) Microsecond() int {
	return that.Time.Nanosecond() / 1e3
}

// Nanosecond 返回纳秒时间戳
func (that *Time) Nanosecond() int {
	return that.Time.Nanosecond()
}
