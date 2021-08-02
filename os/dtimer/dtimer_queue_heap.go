package dtimer

// Len 实现 sort.Interface 中的 Len 接口
func (that *priorityQueueHeap) Len() int {
	return len(that.array)
}

// Less 实现 sort.Interface 中的 Less 接口
func (that *priorityQueueHeap) Less(i, j int) bool {
	return that.array[i].priority < that.array[j].priority
}

// Swap 实现 sort.Interface 中的 Swap 接口
func (that *priorityQueueHeap) Swap(i, j int) {
	if len(that.array) == 0 {
		return
	}
	that.array[i], that.array[j] = that.array[j], that.array[i]
}

// Push 实现heap的方法
func (that *priorityQueueHeap) Push(x interface{}) {
	that.array = append(that.array, x.(priorityQueueItem))
}

// Pop 实现heap的方法
func (that *priorityQueueHeap) Pop() interface{} {
	length := len(that.array)
	if length == 0 {
		return nil
	}
	item := that.array[length-1]
	that.array = that.array[0 : length-1]
	return item
}
