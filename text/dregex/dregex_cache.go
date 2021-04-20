package dregex

import (
	"regexp"
	"sync"
)

var (
	regexMu = sync.RWMutex{}
	//*regexp.Regexp对象缓存
	//使用读写锁保证并发安全，但是没有过期逻辑
	regexMap = make(map[string]*regexp.Regexp)
)

//getRegexp 通过pattern获取*regexp.Regexp对象，并且使用缓存
func getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	regexMu.RLock()
	regex = regexMap[pattern]
	regexMu.RUnlock()
	if regex != nil {
		return regex, nil
	}
	regex, err = regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexMu.Lock()
	regexMap[pattern] = regex
	regexMu.Unlock()
	return regex, nil
}
