package dstr

import "strings"

// IsSubDomain 检查subDomain是否为主域名mainDomain的子域名
// It supports '*' in <mainDomain>.
func IsSubDomain(subDomain string, mainDomain string) bool {
	if p := strings.IndexByte(subDomain, ':'); p != -1 {
		subDomain = subDomain[0:p]
	}
	if p := strings.IndexByte(mainDomain, ':'); p != -1 {
		mainDomain = mainDomain[0:p]
	}
	subArray := strings.Split(subDomain, ".")
	mainArray := strings.Split(mainDomain, ".")
	subLength := len(subArray)
	mainLength := len(mainArray)
	// Eg:
	// "s.s.goframe.org" is not sub-domain of "*.goframe.org"
	// but
	// "s.s.goframe.org" is not sub-domain of "goframe.org"
	if mainLength > 2 && subLength > mainLength {
		return false
	}
	minLength := subLength
	if mainLength < minLength {
		minLength = mainLength
	}
	for i := minLength; i > 0; i-- {
		if mainArray[mainLength-i] == "*" {
			continue
		}
		if mainArray[mainLength-i] != subArray[subLength-i] {
			return false
		}
	}
	return true
}
