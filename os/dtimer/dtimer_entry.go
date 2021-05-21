package dtimer

import "github.com/osgochina/donkeygo/container/dtype"

// Entry 定时器中的每个刻度中需要执行的任务
type Entry struct {
	wheel             *wheel      // 当前条目所属的时间轮盘
	job               JobFunc     // 将要执行的业务方法
	singleton         *dtype.Bool // 是否是单例模式
	status            *dtype.Int  // 当前业务方法执行的状态
	times             *dtype.Int  // 该条目能运行的最大次数，超过这个次数就会被销毁
	createMs          int64       //任务创建时候的毫秒时间戳
	createTicks       int64       // 任务创建时候所在轮盘的当前滴答声
	intervalMs        int64       // 任务的运行间隔毫秒数
	intervalTicks     int64       //任务运行时候间隔的滴答声
	installIntervalMs int64       //任务创建时候间隔的毫秒数

}

// 检查当前条目处理逻辑
// nowTicks: 传入当前时间轮盘的滴答声
// nowMs: 表示当前毫秒时间戳
// runnable: 是否需要运行这个条目
// addable: 是否需要把该条目重新添加到轮盘中，等带下一个到期时间
func (that *Entry) check(nowTicks int64, nowMs int64) (runnable, addable bool) {

	switch that.status.Val() {
	//如果当前条目是停止状态，说明该条目现在不需要现在执行，重新添加到轮盘中，等带下一个到期时间
	case StatusStopped:
		return false, true
	// 如果当前条目状态是关闭，则说明都不做
	case StatusClosed:
		return false, false
	//如果当前条目是重置，说明该条目不需要现在执行，重新添加到轮盘中，等带下一个到期时间
	case StatusReset:
		return false, true
	}

	//滴答声已经走动起来了，并且当前走到的位置已经命中了刻度
	diff := nowTicks - that.createTicks
	if diff > 0 && diff%that.intervalTicks == 0 {

		//如果当前任务所在的时间轮盘不是最小的那个，则需要把自己放到下一个轮盘中继续等待命中
		if that.wheel.level > 0 {

			diffMs := nowMs - that.createMs
			//如果间隔时间不足，则把它放到下一个刻度中执行
			if diffMs < that.wheel.timer.intervalMs {
				that.wheel.slots[(nowTicks+that.intervalTicks)%that.wheel.number].PushBack(that)
				return false, false
				// 正常
			} else if diffMs >= that.wheel.timer.intervalMs {
				//如果滴答一次的毫秒数 - 已等待的毫秒数，还要大于定时器最小滴答毫秒数，说明这个任务还要好久才能到期，把它重新放到轮盘中，等待下一次执行
				if leftMs := that.intervalMs - diffMs; leftMs > that.wheel.timer.intervalMs {
					that.wheel.timer.doAddEntryByParent(false, nowMs, leftMs, that)
					return false, false
				}
			}
		}
		//如果当前任务是单例任务，并且正在运行中,那么此次任务略过，把任务丢到处理器中，等待下次执行
		if that.IsSingleton() {
			if that.status.Set(StatusRunning) == StatusRunning {
				return false, true
			}
		}

		//如果任务的执行次数到了设置值，此次不需要执行，也不需要再次添加到任务处理器中
		times := that.times.Add(-1)
		if times <= 0 {
			if that.status.Set(StatusClosed) == StatusClosed || times < 0 {
				return false, false
			}
		}

		// 这里是一个快速判断逻辑，如果任务的最大执行次数小于20E，则让任务的执行次数变为最大
		// 针对不限制执行次数的任务使用
		if times < 2000000000 && times > 1000000000 {
			that.times.Set(defaultTimes)
		}
		return true, true
	}

	//默认情况下，如果未到执行时间，都把该条目放到轮盘中，等待到期执行
	return false, true
}

// Status 当前运行的状态
func (that *Entry) Status() int {
	return that.status.Val()
}

// SetStatus 设置条目状态
func (that *Entry) SetStatus(status int) int {
	return that.status.Set(status)
}

// Start 开始运行
func (that *Entry) Start() {
	that.status.Set(StatusReady)
}

// Stop 停止运行
func (that *Entry) Stop() {
	that.status.Set(StatusStopped)
}

// Reset 重置
func (that *Entry) Reset() {
	that.status.Set(StatusReset)
}

// Close 关闭
func (that *Entry) Close() {
	that.status.Set(StatusClosed)
}

// IsSingleton 判断该条目是否是单例模式
func (that *Entry) IsSingleton() bool {
	return that.singleton.Val()
}

// SetSingleton 设置该条目为单例模式
func (that *Entry) SetSingleton(enabled bool) {
	that.singleton.Set(enabled)
}

// SetTimes 设置当前条目运行时间
func (that *Entry) SetTimes(times int) {
	that.times.Set(times)
}

// Run 开始运行该条目
func (that *Entry) Run() {
	that.job()
}
