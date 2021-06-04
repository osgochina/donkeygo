package dtime

import (
	"fmt"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/text/dregex"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 一般时间使用的短写
const (
	D  = 24 * time.Hour
	H  = time.Hour
	M  = time.Minute
	S  = time.Second
	MS = time.Millisecond
	US = time.Microsecond
	NS = time.Nanosecond

	// 匹配的时间格式正则
	// Regular expression1(datetime separator supports '-', '/', '.').
	// Eg:
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2017-12-14 04:51:34 +0805 LMT",
	// "2006-01-02T15:04:05Z07:00",
	// "2014-01-17T01:19:15+08:00",
	// "2018-02-09T20:46:17.897Z",
	// "2018-02-09 20:46:17.897",
	// "2018-02-09T20:46:17Z",
	// "2018-02-09 20:46:17",
	// "2018/10/31 - 16:38:46"
	// "2018-02-09",
	// "2018.02.09",
	timeRegexPattern1 = `(\d{4}[-/\.]\d{1,2}[-/\.]\d{1,2})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression2(datetime separator supports '-', '/', '.').
	// Eg:
	// 01-Nov-2018 11:50:28
	// 01/Nov/2018 11:50:28
	// 01.Nov.2018 11:50:28
	// 01.Nov.2018:11:50:28
	timeRegexPattern2 = `(\d{1,2}[-/\.][A-Za-z]{3,}[-/\.]\d{4})[:\sT-]*(\d{0,2}:{0,1}\d{0,2}:{0,1}\d{0,2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`

	// Regular expression3(time).
	// Eg:
	// 11:50:28
	// 11:50:28.897
	timeRegexPattern3 = `(\d{2}):(\d{2}):(\d{2})\.{0,1}(\d{0,9})`
)

var (
	// It's more high performance using regular expression
	// than time.ParseInLocation to parse the datetime string.
	timeRegex1, _ = regexp.Compile(timeRegexPattern1)
	timeRegex2, _ = regexp.Compile(timeRegexPattern2)
	timeRegex3, _ = regexp.Compile(timeRegexPattern3)

	// Month words to arabic numerals mapping.
	monthMap = map[string]int{
		"jan":       1,
		"feb":       2,
		"mar":       3,
		"apr":       4,
		"may":       5,
		"jun":       6,
		"jul":       7,
		"aug":       8,
		"sep":       9,
		"sept":      9,
		"oct":       10,
		"nov":       11,
		"dec":       12,
		"january":   1,
		"february":  2,
		"march":     3,
		"april":     4,
		"june":      6,
		"july":      7,
		"august":    8,
		"september": 9,
		"october":   10,
		"november":  11,
		"december":  12,
	}
)

// SetTimeZone 设置时区
func SetTimeZone(zone string) error {
	location, err := time.LoadLocation(zone)
	if err != nil {
		return err
	}
	return os.Setenv("TZ", location.String())
}

// Timestamp 当前时间戳
func Timestamp() int64 {
	return Now().Timestamp()
}

// TimestampMilli 获取微秒时间戳
func TimestampMilli() int64 {
	return Now().TimestampMilli()
}

// TimestampMicro 获取毫秒时间戳
func TimestampMicro() int64 {
	return Now().TimestampMicro()
}

// TimestampNano 获取纳秒时间戳
func TimestampNano() int64 {
	return Now().TimestampNano()
}

// 解析日期字符串
func parseDateStr(s string) (year, month, day int) {
	array := strings.Split(s, "-")
	if len(array) < 3 {
		array = strings.Split(s, "/")
	}
	if len(array) < 3 {
		array = strings.Split(s, ".")
	}
	// Parsing failed.
	if len(array) < 3 {
		return
	}
	// Checking the year in head or tail.
	if utils.IsNumeric(array[1]) {
		year, _ = strconv.Atoi(array[0])
		month, _ = strconv.Atoi(array[1])
		day, _ = strconv.Atoi(array[2])
	} else {
		if v, ok := monthMap[strings.ToLower(array[1])]; ok {
			month = v
		} else {
			return
		}
		year, _ = strconv.Atoi(array[2])
		day, _ = strconv.Atoi(array[0])
	}
	return
}

func ParseDuration(s string) (time.Duration, error) {
	if utils.IsNumeric(s) {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	}
	match, err := dregex.MatchString(`^([\-\d]+)[dD](.*)$`, s)
	if err != nil {
		return 0, err
	}
	if len(match) == 3 {
		v, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.ParseDuration(fmt.Sprintf(`%dh%s`, v*24, match[2]))
	}
	return time.ParseDuration(s)
}

// StrToTime 通过字符串转换成dtime.Time对象
func StrToTime(str string, format ...string) (*Time, error) {
	if len(format) > 0 {
		return StrToTimeFormat(str, format[0])
	}
	if isTimestampStr(str) {
		timestamp, _ := strconv.ParseInt(str, 10, 64)
		return NewFromTimeStamp(timestamp), nil
	}
	var (
		year, month, day     int
		hour, min, sec, nsec int
		match                []string
		local                = time.Local
	)
	if match = timeRegex1.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex2.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		year, month, day = parseDateStr(match[1])
	} else if match = timeRegex3.FindStringSubmatch(str); len(match) > 0 && match[1] != "" {
		//for k, v := range match {
		//	match[k] = strings.TrimSpace(v)
		//}
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(match[1])
		min, _ = strconv.Atoi(match[2])
		sec, _ = strconv.Atoi(match[3])
		nsec, _ = strconv.Atoi(match[4])
		for i := 0; i < 9-len(match[4]); i++ {
			nsec *= 10
		}
		return NewFromTime(time.Date(0, time.Month(1), 1, hour, min, sec, nsec, local)), nil
	} else {
		return nil, derror.New("unsupported time format")
	}

	// Time
	if len(match[2]) > 0 {
		s := strings.Replace(match[2], ":", "", -1)
		if len(s) < 6 {
			s += strings.Repeat("0", 6-len(s))
		}
		hour, _ = strconv.Atoi(s[0:2])
		min, _ = strconv.Atoi(s[2:4])
		sec, _ = strconv.Atoi(s[4:6])
	}
	// Nanoseconds, check and perform bit filling
	if len(match[3]) > 0 {
		nsec, _ = strconv.Atoi(match[3])
		for i := 0; i < 9-len(match[3]); i++ {
			nsec *= 10
		}
	}
	// If there's zone information in the string,
	// it then performs time zone conversion, which converts the time zone to UTC.
	if match[4] != "" && match[6] == "" {
		match[6] = "000000"
	}
	// If there's offset in the string, it then firstly processes the offset.
	if match[6] != "" {
		zone := strings.Replace(match[6], ":", "", -1)
		zone = strings.TrimLeft(zone, "+-")
		if len(zone) <= 6 {
			zone += strings.Repeat("0", 6-len(zone))
			h, _ := strconv.Atoi(zone[0:2])
			m, _ := strconv.Atoi(zone[2:4])
			s, _ := strconv.Atoi(zone[4:6])
			if h > 24 || m > 59 || s > 59 {
				return nil, derror.Newf("invalid zone string: %s", match[6])
			}
			// Comparing the given time zone whether equals to current time zone,
			// it converts it to UTC if they does not equal.
			_, localOffset := time.Now().Zone()
			// Comparing in seconds.
			if (h*3600 + m*60 + s) != localOffset {
				local = time.UTC
				// UTC conversion.
				operation := match[5]
				if operation != "+" && operation != "-" {
					operation = "-"
				}
				switch operation {
				case "+":
					if h > 0 {
						hour -= h
					}
					if m > 0 {
						min -= m
					}
					if s > 0 {
						sec -= s
					}
				case "-":
					if h > 0 {
						hour += h
					}
					if m > 0 {
						min += m
					}
					if s > 0 {
						sec += s
					}
				}
			}
		}
	}
	if month <= 0 || day <= 0 {
		return nil, derror.New("invalid time string:" + str)
	}
	return NewFromTime(time.Date(year, time.Month(month), day, hour, min, sec, nsec, local)), nil
}

// StrToTimeFormat 传入时间字符串，要转换的格式，生成dtime.Time对象
func StrToTimeFormat(str string, format string) (*Time, error) {
	return StrToTimeLayout(str, formatToStdLayout(format))
}

// StrToTimeLayout 把str时间转换成layout布局的时间格式
func StrToTimeLayout(str string, layout string) (*Time, error) {
	if t, err := time.ParseInLocation(layout, str, time.Local); err == nil {
		return NewFromTime(t), nil
	} else {
		return nil, err
	}
}

// 是不是时间戳格式
func isTimestampStr(s string) bool {
	length := len(s)
	if length == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
