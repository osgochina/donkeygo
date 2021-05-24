package dtime

import (
	"bytes"
	"github.com/osgochina/donkeygo/text/dregex"
	"strconv"
	"strings"
)

var (
	// Refer: http://php.net/manual/en/function.date.php
	formats = map[byte]string{
		'd': "02",                        // Day:    Day of the month, 2 digits with leading zeros. Eg: 01 to 31.
		'D': "Mon",                       // Day:    A textual representation of a day, three letters. Eg: Mon through Sun.
		'w': "Monday",                    // Day:    Numeric representation of the day of the week. Eg: 0 (for Sunday) through 6 (for Saturday).
		'N': "Monday",                    // Day:    ISO-8601 numeric representation of the day of the week. Eg: 1 (for Monday) through 7 (for Sunday).
		'j': "=j=02",                     // Day:    Day of the month without leading zeros. Eg: 1 to 31.
		'S': "02",                        // Day:    English ordinal suffix for the day of the month, 2 characters. Eg: st, nd, rd or th. Works well with j.
		'l': "Monday",                    // Day:    A full textual representation of the day of the week. Eg: Sunday through Saturday.
		'z': "",                          // Day:    The day of the year (starting from 0). Eg: 0 through 365.
		'W': "",                          // Week:   ISO-8601 week number of year, weeks starting on Monday. Eg: 42 (the 42nd week in the year).
		'F': "January",                   // Month:  A full textual representation of a month, such as January or March. Eg: January through December.
		'm': "01",                        // Month:  Numeric representation of a month, with leading zeros. Eg: 01 through 12.
		'M': "Jan",                       // Month:  A short textual representation of a month, three letters. Eg: Jan through Dec.
		'n': "1",                         // Month:  Numeric representation of a month, without leading zeros. Eg: 1 through 12.
		't': "",                          // Month:  Number of days in the given month. Eg: 28 through 31.
		'Y': "2006",                      // Year:   A full numeric representation of a year, 4 digits. Eg: 1999 or 2003.
		'y': "06",                        // Year:   A two digit representation of a year. Eg: 99 or 03.
		'a': "pm",                        // Time:   Lowercase Ante meridiem and Post meridiem. Eg: am or pm.
		'A': "PM",                        // Time:   Uppercase Ante meridiem and Post meridiem. Eg: AM or PM.
		'g': "3",                         // Time:   12-hour format of an hour without leading zeros. Eg: 1 through 12.
		'G': "=G=15",                     // Time:   24-hour format of an hour without leading zeros. Eg: 0 through 23.
		'h': "03",                        // Time:   12-hour format of an hour with leading zeros. Eg: 01 through 12.
		'H': "15",                        // Time:   24-hour format of an hour with leading zeros. Eg: 00 through 23.
		'i': "04",                        // Time:   Minutes with leading zeros. Eg: 00 to 59.
		's': "05",                        // Time:   Seconds with leading zeros. Eg: 00 through 59.
		'u': "=u=.000",                   // Time:   Milliseconds. Eg: 234, 678.
		'U': "",                          // Time:   Seconds since the Unix Epoch (January 1 1970 00:00:00 GMT).
		'O': "-0700",                     // Zone:   Difference to Greenwich time (GMT) in hours. Eg: +0200.
		'P': "-07:00",                    // Zone:   Difference to Greenwich time (GMT) with colon between hours and minutes. Eg: +02:00.
		'T': "MST",                       // Zone:   Timezone abbreviation. Eg: UTC, EST, MDT ...
		'c': "2006-01-02T15:04:05-07:00", // Format: ISO 8601 date. Eg: 2004-02-12T15:19:21+00:00.
		'r': "Mon, 02 Jan 06 15:04 MST",  // Format: RFC 2822 formatted date. Eg: Thu, 21 Dec 2000 16:01:07 +0200.
	}

	// Week to number mapping.
	weekMap = map[string]string{
		"Sunday":    "0",
		"Monday":    "1",
		"Tuesday":   "2",
		"Wednesday": "3",
		"Thursday":  "4",
		"Friday":    "5",
		"Saturday":  "6",
	}

	// Day count of each month which is not in leap year.
	dayOfMonth = []int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
)

// Format 格式换Time格式
func (that *Time) Format(format string) string {
	if that == nil {
		return ""
	}
	runes := []rune(format)
	buffer := bytes.NewBuffer(nil)
	for i := 0; i < len(runes); {
		switch runes[i] {
		case '\\':
			if i < len(runes)-1 {
				buffer.WriteRune(runes[i+1])
				i += 2
				continue
			} else {
				return buffer.String()
			}
		case 'W':
			buffer.WriteString(strconv.Itoa(that.WeeksOfYear()))
		case 'z':
			buffer.WriteString(strconv.Itoa(that.DayOfYear()))
		case 't':
			buffer.WriteString(strconv.Itoa(that.DaysInMonth()))
		case 'U':
			buffer.WriteString(strconv.FormatInt(that.Unix(), 10))
		default:
			if runes[i] > 255 {
				buffer.WriteRune(runes[i])
				break
			}
			if f, ok := formats[byte(runes[i])]; ok {
				result := that.Time.Format(f)
				// Particular chars should be handled here.
				switch runes[i] {
				case 'j':
					for _, s := range []string{"=j=0", "=j="} {
						result = strings.Replace(result, s, "", -1)
					}
					buffer.WriteString(result)
				case 'G':
					for _, s := range []string{"=G=0", "=G="} {
						result = strings.Replace(result, s, "", -1)
					}
					buffer.WriteString(result)
				case 'u':
					buffer.WriteString(strings.Replace(result, "=u=.", "", -1))
				case 'w':
					buffer.WriteString(weekMap[result])
				case 'N':
					buffer.WriteString(strings.Replace(weekMap[result], "0", "7", -1))
				case 'S':
					buffer.WriteString(formatMonthDaySuffixMap(result))
				default:
					buffer.WriteString(result)
				}
			} else {
				buffer.WriteRune(runes[i])
			}
		}
		i++
	}
	return buffer.String()
}

// FormatNew 转换成指定格式，并返回一个新的对象
func (that *Time) FormatNew(format string) *Time {
	if that == nil {
		return nil
	}
	return NewFromStr(that.Format(format))
}

// FormatTo 要转换成的格式
func (that *Time) FormatTo(format string) *Time {
	if that == nil {
		return nil
	}
	that.Time = NewFromStr(that.Format(format)).Time
	return that
}

// Layout 设置时间个格式化
func (that *Time) Layout(layout string) string {
	if that == nil {
		return ""
	}
	return that.Time.Format(layout)
}

// LayoutNew 使用layout格式格式化时间并返回一个新的对象
func (that *Time) LayoutNew(layout string) *Time {
	if that == nil {
		return nil
	}
	return NewFromStr(that.Layout(layout))
}

// LayoutTo 格式化格式
func (that *Time) LayoutTo(layout string) *Time {
	if that == nil {
		return nil
	}
	that.Time = NewFromStr(that.Layout(layout)).Time
	return that
}

// IsLeapYear 当年是否是闰年
func (that *Time) IsLeapYear() bool {
	year := that.Year()
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

// DayOfYear 当天是当年的第几天
func (that *Time) DayOfYear() int {
	day := that.Day()
	month := int(that.Month())
	if that.IsLeapYear() {
		if month > 2 {
			return dayOfMonth[month-1] + day
		}
		return dayOfMonth[month-1] + day - 1
	}
	return dayOfMonth[month-1] + day - 1
}

// DaysInMonth 返回当前月份的天数
func (that *Time) DaysInMonth() int {
	switch that.Month() {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	}
	if that.IsLeapYear() {
		return 29
	}
	return 28
}

// WeeksOfYear 返回当前处于当年的第几周
func (that *Time) WeeksOfYear() int {
	_, week := that.ISOWeek()
	return week
}

// 转换成标准格式
func formatToStdLayout(format string) string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(format); {
		switch format[i] {
		case '\\':
			if i < len(format)-1 {
				b.WriteByte(format[i+1])
				i += 2
				continue
			} else {
				return b.String()
			}

		default:
			if f, ok := formats[format[i]]; ok {
				// Handle particular chars.
				switch format[i] {
				case 'j':
					b.WriteString("2")
				case 'G':
					b.WriteString("15")
				case 'u':
					if i > 0 && format[i-1] == '.' {
						b.WriteString("000")
					} else {
						b.WriteString(".000")
					}

				default:
					b.WriteString(f)
				}
			} else {
				b.WriteByte(format[i])
			}
			i++
		}
	}
	return b.String()
}

// 将自定义格式转换成对应的正则表达式
func formatToRegexPattern(format string) string {
	s := dregex.Quote(formatToStdLayout(format))
	s, _ = dregex.ReplaceString(`[0-9]`, `[0-9]`, s)
	s, _ = dregex.ReplaceString(`[A-Za-z]`, `[A-Za-z]`, s)
	s, _ = dregex.ReplaceString(`\s+`, `\s+`, s)
	return s
}

// 返回当前日期的短英文缩写
func formatMonthDaySuffixMap(day string) string {
	switch day {
	case "01", "21", "31":
		return "st"
	case "02", "22":
		return "nd"
	case "03", "23":
		return "rd"
	default:
		return "th"
	}
}
