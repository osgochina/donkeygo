package dtimer

import (
	"fmt"
	"github.com/osgochina/donkeygo/os/dcmd"
	"math"
	"time"
)

const (
	StatusReady          = 0             // 定时器准备好了，随时可以运行
	StatusRunning        = 1             // 定时器运行中
	StatusStopped        = 2             // 定时器已经停止了
	StatusReset          = 3             // 定时器重置中
	StatusClosed         = -1            // 定时器已经关闭，等待被删除
	panicExit            = "exit"        // 内部使用的作业退出函数
	defaultTimes         = math.MaxInt32 // 默认限制运行次数，一个很大的数字。
	defaultSlotNumber    = 10            // 默认时间轮的刻度
	defaultWheelInterval = 60            // 默认时间轮的滴答间隔
	defaultWheelLevel    = 5             // 默认时间轮的层数
	cmdEnvKey            = "dk.dtimer"   //环境变量中的参数
)

var (
	defaultSlots    = dcmd.GetOptWithEnv(fmt.Sprintf("%s.slots", cmdEnvKey), defaultSlotNumber).Int()
	defaultLevel    = dcmd.GetOptWithEnv(fmt.Sprintf("%s.level", cmdEnvKey), defaultWheelLevel).Int()
	defaultInterval = dcmd.GetOptWithEnv(fmt.Sprintf("%s.interval", cmdEnvKey), defaultWheelInterval).Duration() * time.Millisecond
	defaultTimer    = NewTimer(defaultSlots, defaultInterval, defaultLevel)
)

// JobFunc is the job function.
type JobFunc = func()

// SetTimeout 设置一次性的任务，指定时间到了执行该任务
func SetTimeout(delay time.Duration, job JobFunc) {
	AddOnce(delay, job)
}

// SetInterval 设置循环任务，指定间隔时间内重复执行
func SetInterval(interval time.Duration, job JobFunc) {
	Add(interval, job)
}

// Add 添加任务到默认定时作业处理系统
func Add(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(interval, job)
}

// AddEntry 更多参数的添加任务到定时作业系统
// interval 间隔时间
// job  要处理的方法
// singleton 并发限制，true 表示无论任务执行需要多久，同一时间只能执行一个方法，false表示不限制
// times 任务一共可以执行的次数
// status 任务首次添加时候的状态
func AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return defaultTimer.AddEntry(interval, job, singleton, times, status)
}

// AddSingleton 添加一个并发限制任务
func AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddSingleton(interval, job)
}

// AddOnce 添加一个只执行一次的定时任务
func AddOnce(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddOnce(interval, job)
}

// AddTimes 添加一个指定执行次数的定时任务
func AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return defaultTimer.AddTimes(interval, times, job)
}

// DelayAdd 添加一个延迟执行的定时任务
// delay 表示需要延迟的时间
func DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAdd(delay, interval, job)
}

// DelayAddEntry 添加一个延迟执行的定时任务，并且支持更多参数
func DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	defaultTimer.DelayAddEntry(delay, interval, job, singleton, times, status)
}

// DelayAddSingleton 添加一个延迟执行的并发限制任务
func DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddSingleton(delay, interval, job)
}

// DelayAddOnce 添加一个延迟执行的只执行一次的定时任务
func DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddOnce(delay, interval, job)
}

// DelayAddTimes 添加一个延迟执行的确定执行次数的定时任务
func DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	defaultTimer.DelayAddTimes(delay, interval, times, job)
}

// Exit 退出任务
func Exit() {
	panic(panicExit)
}
