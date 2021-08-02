package dtimer

import (
	"time"
)

func (that *Timer) loop() {
	go func() {
		var (
			currentTimerTicks int64
			// 创建底层的最小定时器
			timerIntervalTicker = time.NewTicker(that.options.Interval)
		)
		defer timerIntervalTicker.Stop()
		for {
			select {
			case <-timerIntervalTicker.C:
				// 检查定时器的状态
				switch that.status.Val() {
				case StatusRunning:
					// 定时器进行中，定时器触发成功，则刻度前进1
					currentTimerTicks = that.ticks.Add(1)
					// 如果当前的刻度满足执行条件，则执行
					if currentTimerTicks >= that.queue.LatestPriority() {
						that.proceed(currentTimerTicks)
					}

				case StatusStopped:
					// Do nothing.

				case StatusClosed:
					// Timer exits.
					return
				}
			}
		}
	}()
}

// 执行任务
func (that *Timer) proceed(currentTimerTicks int64) {
	var (
		value interface{}
	)
	for {
		value = that.queue.Pop()
		if value == nil {
			break
		}
		entry := value.(*Entry)
		// 检查当前任务的节拍器，判断当前节拍是否到了任务需要执行的节拍，如果没到，则把任务重新丢入队列
		if jobNextTicks := entry.nextTicks.Val(); currentTimerTicks < jobNextTicks {
			// 不满足条件，把任务重新丢入队列
			that.queue.Push(entry, entry.nextTicks.Val())
			break
		}
		// 当前任务的节拍已经满足执行条件，但是还需要再次判断当前任务的节拍和状态是否可以执行，如果满足条件，则新开一个协程执行
		entry.doCheckAndRunByTicks(currentTimerTicks)
		// 再次判断任务状态，如果不是关闭状态，则放入队列，等待下次执行
		if entry.Status() != StatusClosed {
			that.queue.Push(entry, entry.nextTicks.Val())
		}
	}
}
