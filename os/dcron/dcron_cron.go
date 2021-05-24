package dcron

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gtimer"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dlog"
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

// NewCron 创建Cron对象
func NewCron() *Cron {
	return &Cron{
		idGen:    dtype.NewInt64(),
		status:   dtype.NewInt(StatusRunning),
		entries:  dmap.NewStrAnyMap(true),
		logPath:  dtype.NewString(),
		logLevel: dtype.NewInt(dlog.LevelProd),
	}
}

//添加定时任务
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

// SetLogPath 设置日志路径
func (that *Cron) SetLogPath(path string) {
	that.logPath.Set(path)
}

// GetLogPath 获取日志路径
func (that *Cron) GetLogPath() string {
	return that.logPath.Val()
}

// SetLogLevel 设置日志等级
func (that *Cron) SetLogLevel(level int) {
	that.logLevel.Set(level)
}

// GetLogLevel 获取日志等级
func (that *Cron) GetLogLevel() int {
	return that.logLevel.Val()
}

// Add 添加定时任务
func (that *Cron) Add(pattern string, job func(), name ...string) (*Entry, error) {
	if len(name) > 0 {
		if that.Search(name[0]) != nil {
			return nil, errors.New(fmt.Sprintf(`cron job "%s" already exists`, name[0]))
		}
	}
	return that.addEntry(pattern, job, false, name...)
}

// AddSingleton 添加单例模式的定时任务
func (that *Cron) AddSingleton(pattern string, job func(), name ...string) (*Entry, error) {
	if entry, err := that.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetSingleton(true)
		return entry, nil
	}
}

// AddOnce 添加只执行一次的定时任务
func (that *Cron) AddOnce(pattern string, job func(), name ...string) (*Entry, error) {
	if entry, err := that.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetTimes(1)
		return entry, nil
	}
}

// AddTimes 添加指定执行次数的定时任务
func (that *Cron) AddTimes(pattern string, times int, job func(), name ...string) (*Entry, error) {
	if entry, err := that.Add(pattern, job, name...); err != nil {
		return nil, err
	} else {
		entry.SetTimes(times)
		return entry, nil
	}
}

// DelayAdd 延迟添加定时任务
func (that *Cron) DelayAdd(delay time.Duration, pattern string, job func(), name ...string) {
	dtimer.AddOnce(delay, func() {
		if _, err := that.Add(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddSingleton 延迟添加单例模式的定时任务
func (that *Cron) DelayAddSingleton(delay time.Duration, pattern string, job func(), name ...string) {
	dtimer.AddOnce(delay, func() {
		if _, err := that.AddSingleton(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddOnce 延迟添加执行一次的定时任务
func (that *Cron) DelayAddOnce(delay time.Duration, pattern string, job func(), name ...string) {
	dtimer.AddOnce(delay, func() {
		if _, err := that.AddOnce(pattern, job, name...); err != nil {
			panic(err)
		}
	})
}

// DelayAddTimes 延迟添加指定执行次数的定时任务
func (that *Cron) DelayAddTimes(delay time.Duration, pattern string, times int, job func(), name ...string) {
	dtimer.AddOnce(delay, func() {
		if _, err := that.AddTimes(pattern, times, job, name...); err != nil {
			panic(err)
		}
	})
}

// Search 通过任务名查找任务
func (that *Cron) Search(name string) *Entry {
	if v := that.entries.Get(name); v != nil {
		return v.(*Entry)
	}
	return nil
}

// Start 启动定时器
func (that *Cron) Start(name ...string) {
	if len(name) > 0 {
		for _, v := range name {
			if entry := that.Search(v); entry != nil {
				entry.Start()
			}
		}
	} else {
		that.status.Set(StatusReady)
	}
}

// Stop 关闭定时任务
func (that *Cron) Stop(name ...string) {
	if len(name) > 0 {
		for _, v := range name {
			if entry := that.Search(v); entry != nil {
				entry.Stop()
			}
		}
	} else {
		that.status.Set(StatusStopped)
	}
}

// Remove 移除定时任务
func (that *Cron) Remove(name string) {
	if v := that.entries.Get(name); v != nil {
		v.(*Entry).Close()
	}
}

// Close 关闭定时任务器
func (that *Cron) Close() {
	that.status.Set(StatusClosed)
}

// Size 有多少个定时任务在处理器中
func (that *Cron) Size() int {
	return that.entries.Size()
}

// Entries 导出任务切片列表
func (that *Cron) Entries() []*Entry {
	array := darray.NewSortedArraySize(that.entries.Size(), func(v1, v2 interface{}) int {
		entry1 := v1.(*Entry)
		entry2 := v2.(*Entry)
		if entry1.Time.Nanosecond() > entry2.Time.Nanosecond() {
			return 1
		}
		return -1
	}, true)
	that.entries.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			array.Add(v.(*Entry))
		}
	})
	entries := make([]*Entry, array.Len())
	array.RLockFunc(func(array []interface{}) {
		for k, v := range array {
			entries[k] = v.(*Entry)
		}
	})
	return entries
}
