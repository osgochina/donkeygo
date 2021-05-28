package dcron

import (
	"github.com/gogf/gf/os/gtimer"
	"math"
	"time"
)

const (
	StatusReady   = gtimer.StatusReady   //已就绪
	StatusRunning = gtimer.StatusRunning // 运行中
	StatusStopped = gtimer.StatusStopped //已停止
	StatusClosed  = gtimer.StatusClosed  // 已关闭
	defaultTimes  = math.MaxInt32        // 默认运行次数
)

var (
	// 默认定时任务处理器
	defaultCron = NewCron()
)

// SetLogPath 设置定时任务处理详情的日志记录文件路径
func SetLogPath(path string) {
	defaultCron.SetLogPath(path)
}

// GetLogPath 获取日志路径
func GetLogPath() string {
	return defaultCron.GetLogPath()
}

// SetLogLevel 设置日志记录的级别
func SetLogLevel(level int) {
	defaultCron.SetLogLevel(level)
}

// GetLogLevel 获取日志记录的级别
func GetLogLevel() int {
	return defaultCron.GetLogLevel()
}

// Add 添加一个定时任务到处理器
func Add(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.Add(pattern, job, name...)
}

// AddSingleton 添加一个并发限制的定时任务到处理器
func AddSingleton(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddSingleton(pattern, job, name...)
}

// AddOnce 添加一个只执行一次的定时任务到处理器
func AddOnce(pattern string, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddOnce(pattern, job, name...)
}

// AddTimes 添加一个限制执行次数的定时任务到处理器
func AddTimes(pattern string, times int, job func(), name ...string) (*Entry, error) {
	return defaultCron.AddTimes(pattern, times, job, name...)
}

// DelayAdd 添加一个指定延迟时间，再解析规则的定时任务到处理器
func DelayAdd(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAdd(delay, pattern, job, name...)
}

// DelayAddSingleton 添加指定延迟时间，再解析规则，并且运行时候有并发限制的的任务
func DelayAddSingleton(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAddSingleton(delay, pattern, job, name...)
}

// DelayAddOnce 添加一个指定延迟时间，在解析规则，只运行一次的定时任务
func DelayAddOnce(delay time.Duration, pattern string, job func(), name ...string) {
	defaultCron.DelayAddOnce(delay, pattern, job, name...)
}

// DelayAddTimes 添加一个指定延迟时间，再解析定时规则，并且运行指定运行次数的定时任务。
func DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ...string) {
	defaultCron.DelayAddTimes(delay, pattern, times, job, name...)
}

// Search 通过名字搜索处理器中的任务
func Search(name string) *Entry {
	return defaultCron.Search(name)
}

// Remove 通过名字移除定时任务
func Remove(name string) {
	defaultCron.Remove(name)
}

// Size 获取处理器中一共有多少任务
func Size() int {
	return defaultCron.Size()
}

// Entries 返回处理器中的任务列表
func Entries() []*Entry {
	return defaultCron.Entries()
}

// Start 启动指定任务
func Start(name string) {
	defaultCron.Start(name)
}

// Stop 关闭指定任务
func Stop(name string) {
	defaultCron.Stop(name)
}

// CronNext 获取表达式指定次数的运行的时间列表
func CronNext(pattern string, t time.Time, queryTimes ...int) ([]time.Time, error) {
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}
	times := 1
	if len(queryTimes) > 0 {
		times = queryTimes[0]
	}
	var tArr []time.Time
	for ; times > 0; times-- {
		t = schedule.Next(t)
		tArr = append(tArr, t)
	}
	return tArr, nil
}
