package dcron

import (
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/text/dregex"
	"strconv"
	"strings"
	"time"
)

// 定时任务计划表
type cronSchedule struct {
	create  int64  // 创建时间戳
	every   int64  // 运行间隔
	pattern string // 原始的cron 规则
	second  map[int]struct{}
	minute  map[int]struct{}
	hour    map[int]struct{}
	day     map[int]struct{}
	week    map[int]struct{}
	month   map[int]struct{}
}

const (
	// 匹配cron 规则，支持 秒 分 时 天 周 月
	dRegexForCron = `^([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,]+)\s+([\-/\d\*\?,A-Za-z]+)\s+([\-/\d\*\?,A-Za-z]+)$`
)

// 预先设置的特殊映射
var predefinedPatternMap = map[string]string{
	"@yearly":   "0 0 0 1 1 *",
	"@annually": "0 0 0 1 1 *",
	"@monthly":  "0 0 0 1 * *",
	"@weekly":   "0 0 0 * * 0",
	"@daily":    "0 0 0 * * *",
	"@midnight": "0 0 0 * * *",
	"@hourly":   "0 0 * * * *",
}

// 英文月份短标识与数字的对应关系
var monthMap = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,
}

// 英文周标识与数字的对应关系
var weekMap = map[string]int{
	"sun": 0,
	"mon": 1,
	"tue": 2,
	"wed": 3,
	"thu": 4,
	"fri": 5,
	"sat": 6,
}

// 创建定时规则
func newSchedule(pattern string) (*cronSchedule, error) {

	//匹配自定义时间表
	if match, _ := dregex.MatchString(`(@\w+)\s*(\w*)\s*`, pattern); len(match) > 0 {
		key := strings.ToLower(match[1])
		// 使用几个预定义的时间来代替cron表达式
		if v, found := predefinedPatternMap[key]; found {
			pattern = v
			//定义任务以固定的时间间隔执行，例如，@every 1h30m10s将表示添加任务之后每隔1小时30分10秒执行
		} else if strings.Compare(key, "@every") == 0 {
			if d, err := dtime.ParseDuration(match[2]); err != nil {
				return nil, err
			} else {
				return &cronSchedule{
					create:  time.Now().Unix(),
					every:   int64(d.Seconds()),
					pattern: pattern,
				}, nil
			}
		} else {
			return nil, errors.New(fmt.Sprintf(`invalid pattern: "%s"`, pattern))
		}
	}
	//匹配传统表达式，如：0 0 0 1 1 2
	if match, _ := dregex.MatchString(dRegexForCron, pattern); len(match) == 7 {
		schedule := &cronSchedule{
			create:  time.Now().Unix(),
			every:   0,
			pattern: pattern,
		}
		//秒
		if m, err := parseItem(match[1], 0, 59, false); err != nil {
			return nil, err
		} else {
			schedule.second = m
		}

		//分
		if m, err := parseItem(match[2], 0, 59, false); err != nil {
			return nil, err
		} else {
			schedule.minute = m
		}

		// 小时.
		if m, err := parseItem(match[3], 0, 23, false); err != nil {
			return nil, err
		} else {
			schedule.hour = m
		}

		// 天
		if m, err := parseItem(match[4], 1, 31, true); err != nil {
			return nil, err
		} else {
			schedule.day = m
		}
		//月
		if m, err := parseItem(match[5], 1, 12, false); err != nil {
			return nil, err
		} else {
			schedule.month = m
		}
		// 周
		if m, err := parseItem(match[6], 0, 6, true); err != nil {
			return nil, err
		} else {
			schedule.week = m
		}
		return schedule, nil
	} else {
		return nil, errors.New(fmt.Sprintf(`invalid pattern: "%s"`, pattern))
	}
}

// 解析规则中的每一项，并以map的形式返回
func parseItem(item string, min int, max int, allowQuestionMark bool) (map[int]struct{}, error) {
	m := make(map[int]struct{}, max-min+1)
	if item == "*" || (allowQuestionMark && item == "?") {
		for i := min; i <= max; i++ {
			m[i] = struct{}{}
		}
	} else {
		for _, item := range strings.Split(item, ",") {
			interval := 1
			intervalArray := strings.Split(item, "/")
			if len(intervalArray) == 2 {
				if i, err := strconv.Atoi(intervalArray[1]); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					interval = i
				}
			}
			var (
				rangeMin   = min
				rangeMax   = max
				fieldType  = byte(0)
				rangeArray = strings.Split(intervalArray[0], "-") // Like: 1-30, JAN-DEC
			)
			switch max {
			case 6:
				// It's checking week field.
				fieldType = 'w'
			case 12:
				// It's checking month field.
				fieldType = 'm'
			}
			// Eg: */5
			if rangeArray[0] != "*" {
				if i, err := parseItemValue(rangeArray[0], fieldType); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					rangeMin = i
					rangeMax = i
				}
			}
			if len(rangeArray) == 2 {
				if i, err := parseItemValue(rangeArray[1], fieldType); err != nil {
					return nil, errors.New(fmt.Sprintf(`invalid pattern item: "%s"`, item))
				} else {
					rangeMax = i
				}
			}
			for i := rangeMin; i <= rangeMax; i += interval {
				m[i] = struct{}{}
			}
		}
	}
	return m, nil
}

// 匹配每一项中的值
func parseItemValue(value string, fieldType byte) (int, error) {
	if dregex.IsMatchString(`^\d+$`, value) {
		// Pure number.
		if i, err := strconv.Atoi(value); err == nil {
			return i, nil
		}
	} else {
		// Check if contains letter,
		// it converts the value to number according to predefined map.
		switch fieldType {
		case 'm':
			if i, ok := monthMap[strings.ToLower(value)]; ok {
				return i, nil
			}
		case 'w':
			if i, ok := weekMap[strings.ToLower(value)]; ok {
				return i, nil
			}
		}
	}
	return 0, errors.New(fmt.Sprintf(`invalid pattern value: "%s"`, value))
}

// 是否命中规则
func (that *cronSchedule) meet(t time.Time) bool {
	if that.every != 0 {
		// It checks using interval.
		diff := t.Unix() - that.create
		if diff > 0 {
			return diff%that.every == 0
		}
		return false
	} else {
		// It checks using normal cron pattern.
		if _, ok := that.second[t.Second()]; !ok {
			return false
		}
		if _, ok := that.minute[t.Minute()]; !ok {
			return false
		}
		if _, ok := that.hour[t.Hour()]; !ok {
			return false
		}
		if _, ok := that.day[t.Day()]; !ok {
			return false
		}
		if _, ok := that.month[int(t.Month())]; !ok {
			return false
		}
		if _, ok := that.week[int(t.Weekday())]; !ok {
			return false
		}
		return true
	}
}
