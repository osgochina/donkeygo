package dtimer

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/container/dtype"
	"time"
)

// Timer 是一个基于分层时间轮算法的定时任务管理器
type Timer struct {
	// 当前管理器状态
	status *dtype.Int
	//时间轮
	wheels []*wheel
	//时间轮层数
	length int
	// 每一个时间轮层数有多少个任务插槽
	number int
	//定时器最小触发单位
	intervalMs int64
	// 获取当前时间的方法
	nowFunc func() time.Time
}

//时间轮
type wheel struct {
	timer      *Timer       //该时间轮属于那个定时任务管理器
	level      int          //当前时间轮是第几层
	slots      []dlist.List //任务槽，用于分散存放任务
	number     int          //时间轮的槽数
	ticks      *dtype.Int64 //时间轮的滴答声
	totalMs    int64        //该时间轮走完一圈需要消耗多少毫秒
	createMs   int64        //该时间轮创建时的毫秒数
	intervalMs int64        //该时间轮每一个刻度代表的毫秒数
}
