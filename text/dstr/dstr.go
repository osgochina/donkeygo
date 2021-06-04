package dstr

import (
	"bytes"
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/util/dconv"
	"strings"
)

func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// SnakeString converts the accepted string to a snake string (XxYy to xx_yy)
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	for _, d := range dconv.Bytes(s) {
		if d >= 'A' && d <= 'Z' {
			if j {
				data = append(data, '_')
				j = false
			}
		} else if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(dconv.String(data))
}

// IsNumeric 字符串是否是数字
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}

// QuoteMeta 返回一个带有反斜杠字符(\)的str版本
//在每个字符之前: .\+*?[^]($)
func QuoteMeta(str string, chars ...string) string {
	var buf bytes.Buffer
	for _, char := range str {
		if len(chars) > 0 {
			for _, c := range chars[0] {
				if c == char {
					buf.WriteRune('\\')
					break
				}
			}
		} else {
			switch char {
			case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
				buf.WriteRune('\\')
			}
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// Split 分割字符串
func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

// Replace 替换内容
// origin 原始内容
// search 要查找的内容
// replace 要替换的内容
// count 替换次数
func Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}
