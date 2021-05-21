package dcron

import (
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtimer"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/util/dconv"
	"reflect"
	"runtime"
	"time"
)

type Cron struct {
	idGen    *dtype.Int64    // id生成器
	status   *dtype.Int      // 状态
	entries  *dmap.StrAnyMap // 任务列表
	logPath  *dtype.String   //日志存储路径
	logLevel *dtype.Int      //日志级别
}

// New 创建Cron对象
func New() *Cron {
	return &Cron{
		status:   dtype.NewInt(StatusRunning),
		entries:  dmap.NewStrAnyMap(true),
		logPath:  dtype.NewString(),
		logLevel: dtype.NewInt(glog.LEVEL_PROD),
	}
}

func (that *Cron) addEntry(pattern string, job func(), singleton bool, name ...string) (*Entry, error) {
	schedule, err := newSchedule(pattern)
	if err != nil {
		return nil, err
	}
	entry := &Entry{
		cron:     that,
		schedule: schedule,
		jobName:  runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name(),
		times:    dtype.NewInt(defaultTimes),
		Job:      job,
		Time:     time.Now(),
	}
	if len(name) > 0 {
		entry.Name = name[0]
	} else {
		entry.Name = "dcron-" + dconv.String(that.idGen.Add(1))
	}

	// 定时任务每秒检查一次，看看是否要执行
	entry.entry = dtimer.AddEntry(time.Second, entry.check, singleton, -1, gtimer.StatusStopped)
	that.entries.Set(entry.Name, entry)
	entry.entry.Start()
	return entry, nil
}
