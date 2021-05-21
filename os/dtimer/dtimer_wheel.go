package dtimer

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dtype"
	"time"
)

//时间轮
type wheel struct {
	timer      *Timer        //该时间轮属于那个定时任务管理器
	level      int           //当前时间轮是第几层
	slots      []*dlist.List //刻度数组，每一个刻度代表一个任务槽
	number     int64         //该时间轮的刻度总数
	ticks      *dtype.Int64  //时间轮的滴答声，每次滴答表示走过了一个刻度
	totalMs    int64         //该时间轮走完一圈需要消耗多少毫秒 总毫秒数(totalMs)=刻度数(number)*每一个刻度间隔(intervalMs)
	createMs   int64         //该时间轮创建时的毫秒时间戳
	intervalMs int64         //该时间轮每一个刻度代表的毫秒数
}

// 时间轮开始走动，滴答声响起
func (that *wheel) start() {
	go func() {
		var tickDuration = time.Duration(that.intervalMs) * time.Millisecond
		var ticker = time.NewTicker(tickDuration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				switch that.timer.status.Val() {
				case StatusRunning:
					that.proceed()
				case StatusStopped:
					//什么也不做
				case StatusClosed:
					return
				}
			}
		}
	}()
}

//
func (that *wheel) proceed() {

	var (
		nowTicks = that.ticks.Add(1) //滴答声响起+1
		// 通过对当前滴答声取模，确定现在是处于哪个刻度，并且把这个刻度需要处理的任务拿出来
		list = that.slots[int(nowTicks%that.number)]
		// 当前刻度需要处理的任务数量
		length = list.Len()
		// 当前的毫秒时间戳
		nowMs = that.timer.nowFunc().UnixNano() / 1e6
	)
	//如果当前刻度需要处理的任务数大于0，则启动一个协程去处理这些任务，如果没有任务要处理，则什么都不做
	if length > 0 {
		go func(jobs *dlist.List, nowTicks int64) {

			var entry *Entry
			for i := length; i > 0; i-- {
				job := jobs.PopFront() //从队列头部取出任务
				if job == nil {
					break
				} else {
					entry = job.(*Entry)
				}
				//检查任务的运行时间是否正确
				runnable, addable := entry.check(nowTicks, nowMs)

				// 可以运行
				if runnable {
					go func(entry *Entry) {
						defer func() {
							//如果任务是主动退出的，则把任务关闭掉，否则报错出来
							if err := recover(); err != nil {
								if err != panicExit {
									panic(err)
								} else {
									entry.Close()
								}
							}
							//如果当前任务处理正常结束，则把任务设置为待处理，方便下次处理
							if entry.Status() == StatusRunning {
								entry.SetStatus(StatusReady)
							}
						}()
						entry.job()
					}(entry)
				}
				// 当前任务执行完毕，再次添加到处理系统，等待下一次执行
				if addable {
					if entry.Status() == StatusReset {
						entry.SetStatus(StatusReady)
					}
					entry.wheel.timer.doAddEntryByParent(!runnable, nowMs, entry.installIntervalMs, entry)
				}
			}

		}(list, nowTicks)
	}
}

//添加定时任务
// interval: 指定方法的执行的时间间隔
// jobFunc:  将要执行的方法
// singleton: 是否单例方法，也就是说同时只能有一个任务正在运行
// times:  当前任务可以运行的次数，超过这个次数，则会自动销毁这个任务
// status: 任务的当前状态
func (that *wheel) addEntry(interval time.Duration, jobFunc JobFunc, singleton bool, times int, status int) *Entry {
	// 可运行次数小于等于0，则赋值一个最大的值，表示不限制
	if times <= 0 {
		times = defaultTimes
	}
	var (
		//比如我传入运行间隔时间是10秒，那么intervalMs就是1万毫秒
		intervalMs = interval.Nanoseconds() / 1e6
		//要确定它是要多少个滴答声之后运行，则把运行间隔时间/当前时间轮盘的每次滴答需要多少毫秒，就可以得到intervalTicks是多少次滴答后运行
		intervalTicks = intervalMs / that.intervalMs
	)
	//如果间隔运行的滴答数为0，则表示按道理应该现在运行，但是现在正在添加，不能运行，那么就最近的一次滴答运行好了。设置intervalTicks=1
	if intervalTicks == 0 {
		intervalTicks = 1
	}

	nowMs := time.Now().UnixNano() / 1e6
	nowTicks := that.ticks.Val()
	entry := &Entry{
		wheel:             that,
		job:               jobFunc,
		singleton:         dtype.NewBool(singleton),
		status:            dtype.NewInt(status),
		times:             dtype.NewInt(times),
		createMs:          nowMs,
		createTicks:       nowTicks,
		intervalMs:        intervalMs,
		intervalTicks:     intervalTicks,
		installIntervalMs: intervalMs, //添加时候设置的间隔运行毫秒数
	}

	that.slots[(nowTicks+intervalTicks)%that.number].PushBack(entry)
	return entry
}

// 上一次任务命中以后，继承上一次任务，再次添加到任务处理器中，等待下一次条件达成，再次运行
func (that *wheel) addEntryByParent(rollOn bool, nowMs, interval int64, parent *Entry) *Entry {

	intervalTicks := interval / that.intervalMs
	if intervalTicks == 0 {
		intervalTicks = 1
	}

	nowTicks := that.ticks.Val()

	entry := &Entry{
		wheel:             that,
		job:               parent.job,
		singleton:         parent.singleton,
		status:            parent.status,
		times:             parent.times,
		createMs:          nowMs,
		createTicks:       nowTicks,
		intervalMs:        interval,
		intervalTicks:     intervalTicks,
		installIntervalMs: parent.installIntervalMs,
	}

	//如果是滚动处理
	if rollOn {
		entry.createMs = parent.createMs
		if parent.wheel.level == that.level {
			entry.createTicks = parent.createTicks
		}
	}
	that.slots[(nowTicks+intervalTicks)%that.number].PushBack(entry)
	return entry
}
