package dtime

import (
	"strconv"
	"time"
)

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

// NewFromStrFormat creates and returns a Time object with given string and
// custom format like: Y-m-d H:i:s.
// Note that it returns nil if there's error occurs.
func NewFromStrFormat(str string, format string) *Time {
	if t, err := StrToTimeFormat(str, format); err == nil {
		return t
	}
	return nil
}

// NewFromStrLayout creates and returns a Time object with given string and
// stdlib layout like: 2006-01-02 15:04:05.
// Note that it returns nil if there's error occurs.
func NewFromStrLayout(str string, layout string) *Time {
	if t, err := StrToTimeLayout(str, layout); err == nil {
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

// Timestamp 当前时间戳
func (that *Time) Timestamp() int64 {
	return that.UnixNano() / 1e9
}

// TimestampMilli 获取微秒时间戳
func (that *Time) TimestampMilli() int64 {
	return that.UnixNano() / 1e6
}

// TimestampMicro 获取毫秒时间戳
func (that *Time) TimestampMicro() int64 {
	return that.UnixNano() / 1e3
}

// TimestampNano 获取纳秒时间戳
func (that *Time) TimestampNano() int64 {
	return that.UnixNano()
}

// TimestampStr 获取时间戳的字符串表示
func (that *Time) TimestampStr() string {
	return strconv.FormatInt(that.Timestamp(), 10)
}

// TimestampMilliStr 获取微妙时间戳的字符串表示
func (that *Time) TimestampMilliStr() string {
	return strconv.FormatInt(that.TimestampMilli(), 10)
}

// TimestampMicroStr 获取毫秒时间戳的字符串表示
func (that *Time) TimestampMicroStr() string {
	return strconv.FormatInt(that.TimestampMicro(), 10)
}

// TimestampNanoStr 获取纳秒时间戳的字符串表示
func (that *Time) TimestampNanoStr() string {
	return strconv.FormatInt(that.TimestampNano(), 10)
}

// Month 月份
func (that *Time) Month() int {
	return int(that.Time.Month())
}

// Second 秒长度
func (that *Time) Second() int {
	return that.Time.Second()
}

// Millisecond 返回毫秒长度
func (that *Time) Millisecond() int {
	return that.Time.Nanosecond() / 1e6
}

// Microsecond 返回微秒长度
func (that *Time) Microsecond() int {
	return that.Time.Nanosecond() / 1e3
}

// Nanosecond 返回纳秒长度
func (that *Time) Nanosecond() int {
	return that.Time.Nanosecond()
}

// 返回时间的字符串格式
func (that *Time) String() string {
	if that == nil {
		return ""
	}
	if that.IsZero() {
		return ""
	}
	return that.Format("Y-m-d H:i:s")
}

// Clone 返回一个新的时间对象，值为当前时间
func (that *Time) Clone() *Time {
	return New(that.Time)
}
