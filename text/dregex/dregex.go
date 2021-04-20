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
