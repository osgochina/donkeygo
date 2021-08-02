package dtimer

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"time"
)

// New 新建一个定时器
func New(options ...TimerOptions) *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: dtype.NewInt(StatusRunning),
		ticks:  dtype.NewInt64(),
	}
	if len(options) > 0 {
		t.options = options[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

// Add 添加定时任务
func (that *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return that.createEntry(interval, job, false, defaultTimes, StatusReady)
}

// AddEntry 添加定时任务
// interval 执行间隔
// job 要执行的任务
// singleton 是否并发限制
// times 要执行的次数
// status 默认状态
func (that *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return that.createEntry(interval, job, singleton, times, status)
}

// AddSingleton 添加并发限制任务
func (that *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return that.createEntry(interval, job, true, defaultTimes, StatusReady)
}

// AddOnce 添加只执行一次任务
func (that *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return that.createEntry(interval, job, true, 1, StatusReady)
}

// AddTimes 添加指定执行次数任务
func (that *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return that.createEntry(interval, job, true, times, StatusReady)
}

// DelayAdd 延迟添加任务
func (that *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.Add(interval, job)
	})
}

// DelayAddEntry 延迟添加更详细参数的任务
func (that *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	that.AddOnce(delay, func() {
		that.AddEntry(interval, job, singleton, times, status)
	})
}

// DelayAddSingleton 延迟添加并发限制任务
func (that *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddSingleton(interval, job)
	})
}

// DelayAddOnce 延迟添加执行一次任务
func (that *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddOnce(interval, job)
	})
}

// DelayAddTimes 延迟添加指定执行次数的任务
func (that *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddTimes(interval, times, job)
	})
}

// Start 启动定时器
func (that *Timer) Start() {
	that.status.Set(StatusRunning)
}

// Stop 停止定时器
func (that *Timer) Stop() {
	that.status.Set(StatusStopped)
}

// Close 关闭定时器
func (that *Timer) Close() {
	that.status.Set(StatusClosed)
}

// 创建并向定时器队列中推送一个任务
func (that *Timer) createEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	if times <= 0 {
		times = defaultTimes
	}
	// 获取当前任务的执行节拍，定时的时间/节拍间隔
	var (
		intervalTicksOfJob = int64(interval / that.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// 如果执行的节拍小于默认的定时节拍间隔，那么把它设置为最小运行节拍
		intervalTicksOfJob = 1
	}
	//下一次执行节拍
	nextTicks := that.ticks.Val() + intervalTicksOfJob
	entry := &Entry{
		job:       job,
		timer:     that,
		ticks:     intervalTicksOfJob,
		times:     dtype.NewInt(times),
		status:    dtype.NewInt(status),
		singleton: dtype.NewBool(singleton),
		nextTicks: dtype.NewInt64(nextTicks),
	}
	that.queue.Push(entry, nextTicks)
	return entry
}
