package dtimer

import (
	"container/heap"
	"github.com/osgochina/donkeygo/container/dtype"
	"math"
	"sync"
)

// 优先队列，使用go自带的head结构实现
type priorityQueue struct {
	mu             sync.RWMutex
	heap           *priorityQueueHeap //优先队列
	latestPriority *dtype.Int64       // 最后一个优先值
}

// 优先队列的堆
type priorityQueueHeap struct {
	array []priorityQueueItem
}

// 优先队列中的条目
type priorityQueueItem struct {
	value    interface{}
	priority int64 //优先级
}

// 创建一个优先队列
func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{
		heap: &priorityQueueHeap{
			array: make([]priorityQueueItem, 0),
		},
		latestPriority: dtype.NewInt64(math.MaxInt64), // 默认是int64的最大值，如果没有需要的任务添加进来，就不需要执行
	}
	heap.Init(queue.heap)
	return queue
}

func (that *priorityQueue) Len() int {
	that.mu.RLock()
	defer that.mu.RUnlock()
	return that.heap.Len()
}

func (that *priorityQueue) LatestPriority() int64 {
	return that.latestPriority.Val()
}

// Push 写入对象到优先队列，根据优先级priority确定其位置
func (that *priorityQueue) Push(value interface{}, priority int64) {
	that.mu.Lock()
	heap.Push(that.heap, priorityQueueItem{
		value:    value,
		priority: priority,
	})
	that.mu.Unlock()
	// 使用原子操作，更新最小优先级
	for {
		latestPriority := that.latestPriority.Val()
		if priority >= latestPriority {
			break
		}
		if that.latestPriority.Cas(latestPriority, priority) {
			break
		}
	}
}

// Pop 从优先队列中获取最优先的对象
func (that *priorityQueue) Pop() interface{} {
	that.mu.Lock()
	if v := heap.Pop(that.heap); v != nil {
		item := v.(priorityQueueItem)
		that.mu.Unlock()
		// 原子操作更新最新的优先值
		for {
			latestPriority := that.latestPriority.Val()
			if item.priority >= latestPriority {
				break
			}
			if that.latestPriority.Cas(latestPriority, item.priority) {
				break
			}
		}
		return item.value
	} else {
		that.mu.Unlock()
	}
	return nil
}
