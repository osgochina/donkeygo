package dtimer_test

import (
	"github.com/osgochina/donkeygo/container/dlist"
	"github.com/osgochina/donkeygo/os/dtimer"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestSetTimeout(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.SetTimeout(200*time.Millisecond, func() {
			list.PushBack(1)
		})

		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 1)
	})
}
func TestSetInterval(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.SetInterval(300*time.Millisecond, func() {
			list.PushBack(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 3)
	})
}

func TestAddEntry(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.AddEntry(200*time.Millisecond, func() {
			list.PushBack(1)
		}, false, 2, dtimer.StatusReady)
		time.Sleep(1100 * time.Millisecond)
		t.Assert(list.Len(), 2)
	})
}

func TestAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.AddSingleton(200*time.Millisecond, func() {
			list.PushBack(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(1100 * time.Millisecond)
		t.Assert(list.Len(), 1)
	})
}

func TestAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.AddTimes(200*time.Millisecond, 2, func() {
			list.PushBack(1)
		})
		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 2)
	})
}

func TestDelayAdd(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.DelayAdd(500*time.Millisecond, 500*time.Millisecond, func() {
			list.PushBack(1)
		})
		time.Sleep(600 * time.Millisecond)
		t.Assert(list.Len(), 0)
		time.Sleep(600 * time.Millisecond)
		t.Assert(list.Len(), 1)
	})
}

func TestDelayAddEntry(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.DelayAddEntry(200*time.Millisecond, 200*time.Millisecond, func() {
			list.PushBack(1)
		}, false, 2, dtimer.StatusReady)
		time.Sleep(300 * time.Millisecond)
		t.Assert(list.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 2)
	})
}

func TestDelayAddSingleton(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.DelayAddSingleton(500*time.Millisecond, 500*time.Millisecond, func() {
			list.PushBack(1)
			time.Sleep(10000 * time.Millisecond)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(list.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 1)
	})
}

func TestDelayAddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.DelayAddOnce(200*time.Millisecond, 200*time.Millisecond, func() {
			list.PushBack(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(list.Len(), 0)
		time.Sleep(1000 * time.Millisecond)
		t.Assert(list.Len(), 1)
	})
}

func TestDelayAddTimes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		list := dlist.New(true)
		dtimer.DelayAddTimes(500*time.Millisecond, 500*time.Millisecond, 2, func() {
			list.PushBack(1)
		})
		time.Sleep(300 * time.Millisecond)
		t.Assert(list.Len(), 0)
		time.Sleep(1500 * time.Millisecond)
		t.Assert(list.Len(), 2)
	})
}
