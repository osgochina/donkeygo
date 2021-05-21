package dregex

import "regexp"

// Quote 对<s>中的特殊字符加引号，并返回副本
// Eg: Quote(`[foo]`) returns `\[foo\]`.
func Quote(s string) string {
	return regexp.QuoteMeta(s)
}

// Validate 验证正则表达式规则是否合法
func Validate(pattern string) error {
	_, err := getRegexp(pattern)
	return err
}

// IsMatch 是否匹配
func IsMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

// IsMatchString 字符串是否匹配
func IsMatchString(pattern string, src string) bool {
	return IsMatch(pattern, []byte(src))
}

//Match 返回匹配的byte数组
func Match(pattern string, src []byte) ([][]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindSubmatch(src), nil
	} else {
		return nil, err
	}
}

//MatchString 返回<pattern>匹配的字符串
func MatchString(pattern string, src string) ([]string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindStringSubmatch(src), nil
	} else {
		return nil, err
	}
}

// ReplaceFunc 使用自定义func按照规则替换数据
func ReplaceFunc(pattern string, src []byte, replaceFunc func(b []byte) []byte) ([]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.ReplaceAllFunc(src, replaceFunc), nil
	} else {
		return nil, err
	}
}

// ReplaceStringFunc replace all matched <pattern> in string <src>
// with custom replacement function <replaceFunc>.
func ReplaceStringFunc(pattern string, src string, replaceFunc func(s string) string) (string, error) {
	bytes, err := ReplaceFunc(pattern, []byte(src), func(bytes []byte) []byte {
		return []byte(replaceFunc(string(bytes)))
	})
	return string(bytes), err
}
