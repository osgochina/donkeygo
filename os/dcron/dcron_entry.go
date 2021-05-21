package dcron

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dtimer"
	"time"
)

type Entry struct {
	cron     *Cron         //entry属于那个cron
	entry    *dtimer.Entry // dtimer的定时器对象
	schedule *cronSchedule // cron规则
	jobName  string
	times    *dtype.Int
	Name     string
	Job      func()
	Time     time.Time
}

func (that *Entry) check() {

}
