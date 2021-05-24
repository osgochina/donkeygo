package dcron_test

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dcron"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestCron_Add_Close(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		_, err1 := cron.Add("* * * * * *", func() {
			glog.Println("cron1")
			array.Append(1)
		})
		_, err2 := cron.Add("* * * * * *", func() {
			glog.Println("cron2")
			array.Append(1)
		}, "test")
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(cron.Size(), 2)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 2)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 4)
		cron.Close()
		time.Sleep(1300 * time.Millisecond)
		fixedLength := array.Len()
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), fixedLength)
	})
}
func TestCron_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		cron.Add("* * * * * *", func() {}, "add")
		//fmt.Println("start", time.Now())
		cron.DelayAdd(time.Second, "* * * * * *", func() {}, "delay_add")
		t.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 2)

		cron.Remove("delay_add")
		t.Assert(cron.Size(), 1)

		entry1 := cron.Search("add")
		entry2 := cron.Search("test-none")
		t.AssertNE(entry1, nil)
		t.Assert(entry2, nil)
	})
}

func TestCron_Remove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.New(true)
		cron.Add("* * * * * *", func() {
			array.Append(1)
		}, "add")
		t.Assert(array.Len(), 0)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)

		cron.Remove("add")
		t.Assert(array.Len(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})
}

func TestCron_AddSingleton(t *testing.T) {
	// un used, can be removed
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		cron.Add("* * * * * *", func() {}, "add")
		cron.DelayAdd(time.Second, "* * * * * *", func() {}, "delay_add")
		t.Assert(cron.Size(), 1)
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 2)

		cron.Remove("delay_add")
		t.Assert(cron.Size(), 1)

		entry1 := cron.Search("add")
		entry2 := cron.Search("test-none")
		t.AssertNE(entry1, nil)
		t.Assert(entry2, nil)
	})
	// keep this
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.AddSingleton("* * * * * *", func() {
			array.Append(1)
			time.Sleep(50 * time.Second)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 1)
	})

}

func TestCron_AddOnce1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.AddOnce("* * * * * *", func() {
			array.Append(1)
		})
		cron.AddOnce("* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 2)
		time.Sleep(2500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_AddOnce2(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.AddOnce("@every 2s", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_AddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.AddTimes("* * * * * *", 2, func() {
			array.Append(1)
		})
		time.Sleep(3500 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_DelayAdd(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.DelayAdd(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
	})
}

func TestCron_DelayAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.DelayAddSingleton(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
			time.Sleep(10 * time.Second)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(2200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
	})
}

func TestCron_DelayAddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.DelayAddOnce(500*time.Millisecond, "* * * * * *", func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(2200 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 0)
	})
}

func TestCron_DelayAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cron := dcron.NewCron()
		array := darray.NewArray(true)
		cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
			array.Append(1)
		})
		t.Assert(cron.Size(), 0)
		time.Sleep(800 * time.Millisecond)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(3000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 0)
	})
}
