package dqueue

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dtype"
	"math"
)

// Queue 队列
type Queue struct {
	limit  int              // 队列长度限制
	list   *dlist.List      //存储数据的容器
	closed *dtype.Bool      // 队列是否关闭
	events chan struct{}    //队列数据写入时候触发该事件
	C      chan interface{} // 读取数据的channel
}

const (
	defaultQueueSize = 10000 // 默认的队列长度
	defaultBatchSize = 10    // 一次从队列存储容器中获取数据的数量
)

// New 创建一个新的队列
// 如果传入了limit则表示是一个固定大小的队列，底层直接使用标准库的channel驱动。
func New(limit ...int) *Queue {
	q := &Queue{
		closed: dtype.NewBool(),
	}
	// 如果是固定大小的队列，则直接使用底层channel通道处理
	if len(limit) > 0 && limit[0] > 0 {
		q.limit = limit[0]
		q.C = make(chan interface{}, q.limit)
	} else {
		q.list = dlist.New(true)
		q.events = make(chan struct{}, math.MaxInt32)  // 入队事件触发通道
		q.C = make(chan interface{}, defaultQueueSize) //默认的队列大小
		go q.asyncLoopFromListToChannel()
	}
	return q
}

//异步处理队列
func (that *Queue) asyncLoopFromListToChannel() {
	defer func() {
		if that.closed.Val() {
			_ = recover()
		}
	}()

	for !that.closed.Val() {
		//等待写入事件
		<-that.events
		for !that.closed.Val() {
			// 从队列中批量取出数据
			if length := that.list.Len(); length > 0 {
				if length > defaultBatchSize {
					length = defaultBatchSize
				}
				for _, v := range that.list.PopFronts(length) {
					// 这个地方有个特殊的用法，当that.C 关闭时，如果正在阻塞的写入到C中，则会发生错误。
					// 不过没关系，defer 已经注册了捕获操作，会捕获它，忽略它
					that.C <- v
				}
				//如果队列中已经没有数据了，则结束循环
			} else {
				break
			}
		}
		// 写入事件可能会同时存在很多次，但是这些写入，在上面批量从队列中获取数据已经处理了，所以需要把这些事件清空
		for i := 0; i < len(that.events)-1; i++ {
			<-that.events
		}

	}
	// 队列如果关闭了，则需要关闭底层channel通道
	close(that.C)
}

// Push 入队
func (that *Queue) Push(v interface{}) {
	if that.limit > 0 {
		that.C <- v
	} else {
		that.list.PushBack(v)
		//如果入队事件队列的channel大小小于默认队列大小，则继续发送入队事件，如果超出了，则不需要发送，说明消费端处理不过来
		if len(that.events) < defaultQueueSize {
			that.events <- struct{}{}
		}
	}
}

// Pop 出队消费
func (that *Queue) Pop() interface{} {
	return <-that.C
}

// Close 关闭队列
func (that *Queue) Close() {
	that.closed.Set(true)
	if that.events != nil {
		close(that.events)
	}
	if that.limit > 0 {
		close(that.C)
	}
	// 把通道中的最后数据都出队列，如果已关闭，返回页没问题
	for i := 0; i < defaultBatchSize; i++ {
		that.Pop()
	}
}

// Len 返回队列长度，队列长度等于 在list列表中的数据 + 在C通道中的数据
func (that *Queue) Len() (length int) {
	if that.list != nil {
		length += that.list.Len()
	}
	length += len(that.C)
	return
}

// Size len的别名
func (that *Queue) Size() int {
	return that.Len()
}
