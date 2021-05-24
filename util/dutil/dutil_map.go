package dutil

import "github.com/osgochina/donkeygo/internal/utils"

// MapCopy 复制map数据到新的map对象
func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

// MapPossibleItemByKey 忽略map中key值中的符号，只对比字符和数字，如果能匹配，则返回。注意这个方法的性能会比较低
func MapPossibleItemByKey(data map[string]interface{}, key string) (foundKey string, foundValue interface{}) {
	if len(data) == 0 {
		return
	}
	if v, ok := data[key]; ok {
		return key, v
	}
	// Loop checking.
	for k, v := range data {
		if utils.EqualFoldWithoutChars(k, key) {
			return k, v
		}
	}
	return "", nil
}
