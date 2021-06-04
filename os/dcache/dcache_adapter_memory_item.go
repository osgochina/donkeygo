package dcache

import "github.com/osgochina/donkeygo/os/dtime"

// 内存缓存的值得数据结构
type adapterMemoryItem struct {
	value  interface{} // 真正的值
	expire int64       // 过期时间
}

// IsExpired 判断是否过期
func (that *adapterMemoryItem) IsExpired() bool {

	if that.expire >= dtime.TimestampMilli() {
		return false
	}
	return true
}
