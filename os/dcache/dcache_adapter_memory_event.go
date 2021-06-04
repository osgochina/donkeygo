package dcache

// 缓存写入事件
type adapterMemoryEvent struct {
	k interface{} // Key.
	e int64       // Expire time in milliseconds.
}
