package dtimer

import (
	"fmt"
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dtype"
	"time"
)

// Timer 是一个基于分层时间轮算法的定时任务管理器
type Timer struct {
	// 当前管理器状态
	status *dtype.Int
	//分层时间轮
	wheels []*wheel
	//时间轮层数
	length int
	// 每层时间轮刻度
	number int
	//定时器最小滴答声
	intervalMs int64
	// 获取当前时间的方法
	nowFunc func() time.Time
}

// NewTimer 创建定时器
func NewTimer(slot int, interval time.Duration, level ...int) *Timer {
	timer := doNewWithoutAutoStart(slot, interval, level...)
	timer.wheels[0].start()
	return timer
}

// 创建定时器并启动
func doNewWithoutAutoStart(slot int, interval time.Duration, level ...int) *Timer {
	if slot <= 0 {
		panic(fmt.Sprintf("invalid slot number: %d", slot))
	}
	length := defaultWheelLevel
	if len(level) > 0 {
		length = level[0]
	}

	timer := &Timer{
		status:     dtype.NewInt(StatusRunning),
		wheels:     make([]*wheel, length),
		length:     length,
		number:     slot,
		intervalMs: interval.Nanoseconds() / 1e6, //把滴答数的单位设置为毫秒，纳秒/1000000
		nowFunc: func() time.Time {
			return time.Now()
		},
	}

	//根据传入的时间轮长度，初始化时间轮
	for i := 0; i < length; i++ {
		//初始化第二轮=>第N轮
		if i > 0 {
			//上一个轮盘走完一圈所需要的时间，就是下一个轮盘的一次滴答声
			n := time.Duration(timer.wheels[i-1].totalMs) * time.Millisecond
			if n <= 0 {
				panic(fmt.Sprintf(`inteval is too large with level: %dms x %d`, interval, length))
			}
			//创建这一级的轮盘
			w := timer.newWheel(i, slot, n)
			timer.wheels[i] = w

			//给上一个轮盘添加一个定时任务，检查是否有任务到期
			timer.wheels[i-1].addEntry(n, w.proceed, false, defaultTimes, StatusReady)
			// 如果当前轮盘是最后一个轮盘，则为当前轮盘增加一个定时器，检查是否有任务到期
			if i == length-1 {
				timer.wheels[i].addEntry(n, w.proceed, false, defaultTimes, StatusReady)
			}
		} else {
			// 初始化第一个轮盘
			w := timer.newWheel(i, slot, interval)
			timer.wheels[i] = w
		}
	}

	return timer
}

//创建并返回一个轮盘
// level : 当前时间轮盘所在的层数
// slot ： 当前时间轮盘的刻度数量
// interval: 当前时间轮盘每次滴答需要消耗的时间
func (that *Timer) newWheel(level int, slot int, interval time.Duration) *wheel {
	w := &wheel{
		timer:      that,
		level:      level,                                      //表示当前轮盘是在第几层
		slots:      make([]*dlist.List, slot),                  //创建指定的刻度存放任务
		number:     int64(slot),                                //当前轮盘刻度数
		ticks:      dtype.NewInt64(),                           //记录当前轮盘走过的滴答数
		totalMs:    int64(slot) * interval.Nanoseconds() / 1e6, //当前轮盘要走完一圈需要的总毫秒数
		createMs:   time.Now().UnixNano() / 1e6,                //创建该轮盘的毫秒时间戳
		intervalMs: interval.Nanoseconds() / 1e6,               //当前时间盘的滴答间隔毫秒数
	}

	for i := int64(0); i < w.number; i++ {
		w.slots[i] = dlist.New(true)
	}
	return w
}

// Start 开始启动定时器
func (that *Timer) Start() {
	that.status.Set(StatusRunning)
}

// Stop 停止定时器
func (that *Timer) Stop() {
	that.status.Set(StatusStopped)
}

// Close 关闭定时器
func (that *Timer) Close() {
	that.status.Set(StatusClosed)
}

// Add 添加任务到定时作业系统
// interval 间隔时间
// job  要处理的方法
func (that *Timer) Add(interval time.Duration, job JobFunc) *Entry {
	return that.doAddEntry(interval, job, false, defaultTimes, StatusReady)
}

// AddEntry 更多参数的添加任务到定时作业系统
// interval 间隔时间
// job  要处理的方法
// singleton 并发限制，true 表示无论任务执行需要多久，同一时间只能执行一个方法，false表示不限制
// times 任务一共可以执行的次数
// status 任务首次添加时候的状态
func (that *Timer) AddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	return that.doAddEntry(interval, job, singleton, times, status)
}

// AddSingleton 添加一个并发限制任务
func (that *Timer) AddSingleton(interval time.Duration, job JobFunc) *Entry {
	return that.doAddEntry(interval, job, true, defaultTimes, StatusReady)
}

// AddOnce 添加一个只执行一次的定时任务
func (that *Timer) AddOnce(interval time.Duration, job JobFunc) *Entry {
	return that.doAddEntry(interval, job, true, 1, StatusReady)
}

// AddTimes 添加一个指定执行次数的定时任务
func (that *Timer) AddTimes(interval time.Duration, times int, job JobFunc) *Entry {
	return that.doAddEntry(interval, job, true, times, StatusReady)
}

// DelayAdd 添加一个延迟执行的定时任务
// delay 表示需要延迟的时间
func (that *Timer) DelayAdd(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.Add(interval, job)
	})
}

// DelayAddEntry 添加一个延迟执行的定时任务，并且支持更多参数
func (that *Timer) DelayAddEntry(delay time.Duration, interval time.Duration, job JobFunc, singleton bool, times int, status int) {
	that.AddOnce(delay, func() {
		that.AddEntry(interval, job, singleton, times, status)
	})
}

// DelayAddSingleton 添加一个延迟执行的并发限制任务
func (that *Timer) DelayAddSingleton(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddSingleton(interval, job)
	})
}

// DelayAddOnce 添加一个延迟执行的只执行一次的定时任务
func (that *Timer) DelayAddOnce(delay time.Duration, interval time.Duration, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddOnce(interval, job)
	})
}

// DelayAddTimes 添加一个延迟执行的确定执行次数的定时任务
func (that *Timer) DelayAddTimes(delay time.Duration, interval time.Duration, times int, job JobFunc) {
	that.AddOnce(delay, func() {
		that.AddTimes(interval, times, job)
	})
}

//添加任务到处理器
func (that *Timer) doAddEntry(interval time.Duration, job JobFunc, singleton bool, times int, status int) *Entry {
	// 通过任务的时间间隔获取它应该报错到那一层
	level := that.getLevelByIntervalMs(interval.Nanoseconds() / 1e6)
	return that.wheels[level].addEntry(interval, job, singleton, times, status)
}

// 添加一个带父条目的任务到处理器中
func (that *Timer) doAddEntryByParent(rollOn bool, nowMs, interval int64, parent *Entry) *Entry {
	level := that.getLevelByIntervalMs(interval)
	return that.wheels[level].addEntryByParent(rollOn, nowMs, interval, parent)
}

//通过传入的滴答时间间隔，判断需要把任务放到第几层轮盘
func (that *Timer) getLevelByIntervalMs(intervalMs int64) int {
	pos, cmp := that.binSearchIndex(intervalMs)

	switch cmp {
	//如果直接匹配上了，也不要直接使用
	case 0:
		fallthrough
	case -1:
		// -1 表示匹配到了最近的一个轮盘，但是精度不够，需要往下匹配，需要从这个轮盘往下，一个个的匹配，找到最适合的那个轮盘
		i := pos
		for ; i > 0; i-- {
			if intervalMs > that.wheels[i].intervalMs && intervalMs <= that.wheels[i].totalMs {
				return i
			}
		}
		return i
	case 1:
		// 1表示匹配到了最近的一个轮盘，精度有点小，需要循环往上匹配，找到最适合精度的那个轮盘
		i := pos
		for ; i < that.length-1; i++ {
			if intervalMs > that.wheels[i].intervalMs && intervalMs <= that.wheels[i].totalMs {
				return i
			}
		}
		return i
	}
	return 0
}

// 使用二分查找算法，查找指定n的时间间隔可能在第几层
func (that *Timer) binSearchIndex(n int64) (index int, result int) {
	//最小层在0层
	min := 0
	//最大层为设置层数-1
	max := that.length - 1

	mid := 0
	cmp := -2

	for min <= max {
		//轮盘中间层
		mid = min + (max-min)/2
		switch {
		case that.wheels[mid].intervalMs == n: //如果轮盘中间层的间隔等于要找到n，则表示匹配,直接返回命中的层
			cmp = 0
		case that.wheels[mid].intervalMs > n: //如果中间层的滴答毫秒数大于n，则表示当前层太大，精度不够，去更低层匹配
			cmp = -1
		case that.wheels[mid].intervalMs < n: // 如果中间层的滴答毫秒数小于n，则表示当前层太小，应该放到更高层去
			cmp = 1
		}
		switch cmp {
		case -1:
			max = mid - 1
		case 1:
			min = mid + 1
		case 0:
			return mid, cmp
		}
	}
	return mid, cmp
}
