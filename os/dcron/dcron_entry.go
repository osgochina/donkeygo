package dcron

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtimer"
	"time"
)

// Entry 定时任务条目
type Entry struct {
	cron     *Cron         //entry属于那个cron
	entry    *dtimer.Entry // dtimer的定时器对象
	schedule *cronSchedule // cron规则
	jobName  string        // 任务的方法名
	times    *dtype.Int    // 任务可以运行的次数
	Name     string        // 任务的自定义名字
	Job      func()        //任务处理方法地址
	Time     time.Time     // 任务建立时间
}

func (that *Entry) check() {
	if that.schedule.meet(time.Now()) {
		path := that.cron.GetLogPath()
		level := that.cron.GetLogLevel()
		switch that.cron.status.Val() {
		case StatusStopped:
			return
		case StatusClosed:
			that.Close()
			glog.Path(path).Level(level).Debugf("[dcron] %s(%s) %s removed", that.Name, that.schedule.pattern, that.jobName)
		case StatusReady:
			fallthrough
		case StatusRunning:
			times := that.times.Add(-1)
			if times <= 0 {
				if that.entry.SetStatus(StatusClosed) == StatusClosed || times < 0 {
					return
				}
			}
			if times < 2000000000 && times > 1000000000 {
				that.times.Set(defaultTimes)
			}
			glog.Path(path).Level(level).Debugf("[dcron] %s(%s) %s start", that.Name, that.schedule.pattern, that.jobName)
			defer func() {
				if err := recover(); err != nil {
					glog.Path(path).Level(level).Errorf("[dcron] %s(%s) %s end with error: %v", that.Name, that.schedule.pattern, that.jobName, err)
				} else {
					glog.Path(path).Level(level).Debugf("[dcron] %s(%s) %s end", that.Name, that.schedule.pattern, that.jobName)
				}
				if that.entry.Status() == StatusClosed {
					that.Close()
				}
			}()
			that.Job()
		}
	}
}

// IsSingleton 判断任务是否是单例模式
func (that *Entry) IsSingleton() bool {
	return that.entry.IsSingleton()
}

// SetSingleton 设置任务单例模式开关
func (that *Entry) SetSingleton(enabled bool) {
	that.entry.SetSingleton(enabled)
}

// SetTimes 设置任务运行次数
func (that *Entry) SetTimes(times int) {
	that.times.Set(times)
}

// Status 获取任务状态
func (that *Entry) Status() int {
	return that.entry.Status()
}

// SetStatus 设置任务状态
func (that *Entry) SetStatus(status int) int {
	return that.entry.SetStatus(status)
}

// Start 开始运行任务
func (that *Entry) Start() {
	that.entry.Start()
}

// Stop 停止任务
func (that *Entry) Stop() {
	that.entry.Stop()
}

// Close 关闭当前服务
func (that *Entry) Close() {
	that.cron.entries.Remove(that.Name)
	that.entry.Close()
}
