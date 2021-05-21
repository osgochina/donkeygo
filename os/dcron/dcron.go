package dcron

import (
	"github.com/gogf/gf/os/gtimer"
	"math"
)

const (
	StatusReady   = gtimer.StatusReady   //已就绪
	StatusRunning = gtimer.StatusRunning // 运行中
	StatusStopped = gtimer.StatusStopped //已停止
	StatusClosed  = gtimer.StatusClosed  // 已关闭
	defaultTimes  = math.MaxInt32        // 默认运行次数
)
