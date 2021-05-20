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

// Add 添加任务到默认定时作业处理系统
func Add(interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(interval, job)
}
