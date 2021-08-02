package dtimer

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"math"
)

// Entry 定时器要执行的任务
type Entry struct {
	job       JobFunc      // 要执行的func
	timer     *Timer       // 定时器对象
	ticks     int64        // 任务执行间隔，当前执行ticks+ticks等于下一次需要执行的ticks
	times     *dtype.Int   // job func要执行的次数
	status    *dtype.Int   // 任务状态
	singleton *dtype.Bool  // 是否是一次性任务
	nextTicks *dtype.Int64 // 该任务下一次执行的时间
}

type JobFunc = func()

// Status 获取当前任务的状态
func (that *Entry) Status() int {
	return that.status.Val()
}

// Run 运行任务
func (that *Entry) Run() {
	// 执行次数-1
	leftRunningTimes := that.times.Add(-1)
	// 检查执行次数，如果到了设置的值，则关闭该任务
	if leftRunningTimes < 0 {
		that.status.Set(StatusClosed)
		return
	}
	// 不限制执行次数
	if leftRunningTimes == math.MaxInt32-1 {
		that.times.Set(math.MaxInt32)
	}
	// 开启一个协程执行任务
	go func() {
		defer func() {
			// 如果任务执行报错，则截取该错误，看该任务是否是主动关闭，如果是主动关闭，则退出该任务的执行，
			// 以后也不会再次执行了，如果是运行报错，则继续网上panic，下次到了时间还会执行
			if err := recover(); err != nil {
				if err != panicExit {
					panic(err)
				} else {
					that.Close()
					return
				}
			}
			//把任务正在执行状态变更为准备执行
			if that.Status() == StatusRunning {
				that.SetStatus(StatusReady)
			}
		}()
		that.job()
	}()
}

// 检查任务是否可以执行
func (that *Entry) doCheckAndRunByTicks(currentTimerTicks int64) {
	// 当前ticks是否到了任务下一次需要执行的ticks
	if currentTimerTicks < that.nextTicks.Val() {
		return
	}
	// 设置下一次需要执行的ticks
	that.nextTicks.Set(currentTimerTicks + that.ticks)
	// 执行任务状态检查
	switch that.status.Val() {
	case StatusRunning:
		if that.IsSingleton() {
			return
		}
	case StatusReady:
		if !that.status.Cas(StatusReady, StatusRunning) {
			return
		}
	case StatusStopped:
		return
	case StatusClosed:
		return
	}
	// 执行任务
	that.Run()
}

// SetStatus 修改任务状态
func (that *Entry) SetStatus(status int) int {
	return that.status.Set(status)
}

// Start 启动任务，把任务状态修改为准备执行
func (that *Entry) Start() {
	that.status.Set(StatusReady)
}

// Stop 停止任务，把任务修正为停止状态，不同于close任务，停止状态可以从新开始
func (that *Entry) Stop() {
	that.status.Set(StatusStopped)
}

// Close 关闭该任务
func (that *Entry) Close() {
	that.status.Set(StatusClosed)
}

// Reset 重置任务运行节拍器
func (that *Entry) Reset() {
	that.nextTicks.Set(that.timer.ticks.Val() + that.ticks)
}

// IsSingleton 判断任务是否是并发限制任务
func (that *Entry) IsSingleton() bool {
	return that.singleton.Val()
}

// SetSingleton 设置任务为并发限制任务
func (that *Entry) SetSingleton(enabled bool) {
	that.singleton.Set(enabled)
}

// Job 当前定时任务需要执行的方法
func (that *Entry) Job() JobFunc {
	return that.job
}

// SetTimes 设置任务执行次数
func (that *Entry) SetTimes(times int) {
	that.times.Set(times)
}
