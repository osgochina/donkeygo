package dtime

import (
	"bytes"
	"github.com/osgochina/donkeygo/errors/derror"
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

// New 通过参数创建时间对象
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

// IsZero 判断时间是否为0时间
func (that *Time) IsZero() bool {
	if that == nil {
		return true
	}
	return that.Time.IsZero()
}

// Clone 返回一个新的时间对象，值为当前时间
func (that *Time) Clone() *Time {
	return New(that.Time)
}

// Add 增加制定的时间
func (that *Time) Add(d time.Duration) *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.Add(d)
	return newTime
}

// AddStr 通过字符串格式的时间增加时间
func (that *Time) AddStr(duration string) (*Time, error) {
	if d, err := time.ParseDuration(duration); err != nil {
		return nil, err
	} else {
		return that.Add(d), nil
	}
}

// UTC 返回utc时间
func (that *Time) UTC() *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.UTC()
	return newTime
}

// ISO8601 formats the time as ISO8601 and returns it as string.
func (that *Time) ISO8601() string {
	return that.Layout("2006-01-02T15:04:05-07:00")
}

// RFC822 formats the time as RFC822 and returns it as string.
func (that *Time) RFC822() string {
	return that.Layout("Mon, 02 Jan 06 15:04 MST")
}

// AddDate adds year, month and day to the time.
func (that *Time) AddDate(years int, months int, days int) *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.AddDate(years, months, days)
	return newTime
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Round operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Round(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (that *Time) Round(d time.Duration) *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.Round(d)
	return newTime
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
//
// Truncate operates on the time as an absolute duration since the
// zero time; it does not operate on the presentation form of the
// time. Thus, Truncate(Hour) may return a time with a non-zero
// minute, depending on the time's Location.
func (that *Time) Truncate(d time.Duration) *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.Truncate(d)
	return newTime
}

// Equal reports whether t and u represent the same time instant.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
// See the documentation on the Time type for the pitfalls of using == with
// Time values; most code should use Equal instead.
func (that *Time) Equal(u *Time) bool {
	return that.Time.Equal(u.Time)
}

// Before reports whether the time instant t is before u.
func (that *Time) Before(u *Time) bool {
	return that.Time.Before(u.Time)
}

// After reports whether the time instant t is after u.
func (that *Time) After(u *Time) bool {
	return that.Time.After(u.Time)
}

// Sub returns the duration t-u. If the result exceeds the maximum (or minimum)
// value that can be stored in a Duration, the maximum (or minimum) duration
// will be returned.
// To compute t-d for a duration d, use t.Add(-d).
func (that *Time) Sub(u *Time) time.Duration {
	return that.Time.Sub(u.Time)
}

// StartOfMinute clones and returns a new time of which the seconds is set to 0.
func (that *Time) StartOfMinute() *Time {
	newTime := that.Clone()
	newTime.Time = newTime.Time.Truncate(time.Minute)
	return newTime
}

// StartOfHour clones and returns a new time of which the hour, minutes and seconds are set to 0.
func (that *Time) StartOfHour() *Time {
	y, m, d := that.Date()
	newTime := that.Clone()
	newTime.Time = time.Date(y, m, d, newTime.Time.Hour(), 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfDay clones and returns a new time which is the start of day, its time is set to 00:00:00.
// clone 一个新的时间，并返回这一天的开始时间
func (that *Time) StartOfDay() *Time {
	y, m, d := that.Date()
	newTime := that.Clone()
	newTime.Time = time.Date(y, m, d, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfWeek clones and returns a new time which is the first day of week and its time is set to
// 00:00:00.
func (that *Time) StartOfWeek() *Time {
	weekday := int(that.Weekday())
	return that.StartOfDay().AddDate(0, 0, -weekday)
}

// StartOfMonth clones and returns a new time which is the first day of the month and its is set to
// 00:00:00
func (that *Time) StartOfMonth() *Time {
	y, m, _ := that.Date()
	newTime := that.Clone()
	newTime.Time = time.Date(y, m, 1, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// StartOfQuarter clones and returns a new time which is the first day of the quarter and its time is set
// to 00:00:00.
func (that *Time) StartOfQuarter() *Time {
	month := that.StartOfMonth()
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// StartOfHalf clones and returns a new time which is the first day of the half year and its time is set
// to 00:00:00.
func (that *Time) StartOfHalf() *Time {
	month := that.StartOfMonth()
	offset := (int(month.Month()) - 1) % 6
	return month.AddDate(0, -offset, 0)
}

// StartOfYear clones and returns a new time which is the first day of the year and its time is set to
// 00:00:00.
func (that *Time) StartOfYear() *Time {
	y, _, _ := that.Date()
	newTime := that.Clone()
	newTime.Time = time.Date(y, time.January, 1, 0, 0, 0, 0, newTime.Time.Location())
	return newTime
}

// EndOfMinute clones and returns a new time of which the seconds is set to 59.
func (that *Time) EndOfMinute() *Time {
	return that.StartOfMinute().Add(time.Minute - time.Nanosecond)
}

// EndOfHour clones and returns a new time of which the minutes and seconds are both set to 59.
func (that *Time) EndOfHour() *Time {
	return that.StartOfHour().Add(time.Hour - time.Nanosecond)
}

// EndOfDay clones and returns a new time which is the end of day the and its time is set to 23:59:59.
func (that *Time) EndOfDay() *Time {
	y, m, d := that.Date()
	newTime := that.Clone()
	newTime.Time = time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), newTime.Time.Location())
	return newTime
}

// EndOfWeek clones and returns a new time which is the end of week and its time is set to 23:59:59.
func (that *Time) EndOfWeek() *Time {
	return that.StartOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

// EndOfMonth clones and returns a new time which is the end of the month and its time is set to 23:59:59.
func (that *Time) EndOfMonth() *Time {
	return that.StartOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// EndOfQuarter clones and returns a new time which is end of the quarter and its time is set to 23:59:59.
func (that *Time) EndOfQuarter() *Time {
	return that.StartOfQuarter().AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// EndOfHalf clones and returns a new time which is the end of the half year and its time is set to 23:59:59.
func (that *Time) EndOfHalf() *Time {
	return that.StartOfHalf().AddDate(0, 6, 0).Add(-time.Nanosecond)
}

// EndOfYear clones and returns a new time which is the end of the year and its time is set to 23:59:59.
func (that *Time) EndOfYear() *Time {
	return that.StartOfYear().AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (that *Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + that.String() + `"`), nil
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (that *Time) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		that.Time = time.Time{}
		return nil
	}
	newTime, err := StrToTime(string(bytes.Trim(b, `"`)))
	if err != nil {
		return err
	}
	that.Time = newTime.Time
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// Note that it overwrites the same implementer of `time.Time`.
func (that *Time) UnmarshalText(data []byte) error {
	vTime := New(data)
	if vTime != nil {
		*that = *vTime
		return nil
	}
	return derror.Newf(`invalid time value: %s`, data)
}

// NoValidation marks this struct object will not be validated by package gvalid.
func (that *Time) NoValidation() {

}
