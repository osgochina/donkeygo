package dtimer

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dcmd"
	"math"
	"sync"
	"time"
)

// Timer 定时器对象
type Timer struct {
	mu      sync.RWMutex
	queue   *priorityQueue // 基于堆结构的优先队列
	status  *dtype.Int     // 当前定时器的状态
	ticks   *dtype.Int64   // 定时器进行的间隔数
	options TimerOptions   // 定时器选项
}

type TimerOptions struct {
	Interval time.Duration // Interval 定时器的触发间隔
}

const (
	StatusReady              = 0                    // 定时器准备好了，随时可以运行
	StatusRunning            = 1                    // 定时器运行中
	StatusStopped            = 2                    // 定时器已经停止了
	StatusClosed             = -1                   // 定时器已经关闭，等待被删除
	panicExit                = "exit"               // 内部使用的作业退出函数
	defaultTimes             = math.MaxInt32        // 默认限制运行次数，一个很大的数字。
	defaultTimerInterval     = 100                  // 默认触发间隔
	commandEnvKeyForInterval = "dk.dtimer.interval" //环境变量中的参数
)

var (
	defaultTimer    = New()
	defaultInterval = dcmd.GetOptWithEnv(commandEnvKeyForInterval, defaultTimerInterval).Duration() * time.Millisecond
)

// DefaultOptions 获取默认的配置选项
func DefaultOptions() TimerOptions {
	return TimerOptions{
		Interval: defaultInterval,
	}
}

// SetTimeout 设置延迟执行的一次性任务
func SetTimeout(delay time.Duration, job JobFunc) {
	AddOnce(delay, job)
}

// SetInterval 设置执行时间循环执行的任务
func SetInterval(interval time.Duration, job JobFunc) {
	Add(interval, job)
}

// Add 设置执行时间循环执行的任务
func Add(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(interval, job)
}

// AddEntry 通过参数设置任务
func AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return defaultTimer.AddEntry(interval, job, singleton, times, status)
}

// AddSingleton 添加并发限制任务
func AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddSingleton(interval, job)
}

// AddOnce 添加一次性任务
func AddOnce(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddOnce(interval, job)
}

// AddTimes 添加指定执行次数的任务
func AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return defaultTimer.AddTimes(interval, times, job)
}

// DelayAdd 延迟添加任务
func DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAdd(delay, interval, job)
}

// DelayAddEntry 延迟添加任务
func DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	defaultTimer.DelayAddEntry(delay, interval, job, singleton, times, status)
}

// DelayAddSingleton 延迟添加并发限制任务
func DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddSingleton(delay, interval, job)
}

// DelayAddOnce 延迟添加一次性任务
func DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddOnce(delay, interval, job)
}

// DelayAddTimes 延迟添加指定执行次数的任务
func DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	defaultTimer.DelayAddTimes(delay, interval, times, job)
}

// Exit 退出该任务，不再触发后续的执行
func Exit() {
	panic(panicExit)
}
