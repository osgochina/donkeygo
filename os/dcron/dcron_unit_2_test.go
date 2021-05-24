// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dcron_test

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/os/dcron"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestCron_Entry_Operations(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			cron  = dcron.NewCron()
			array = darray.New(true)
		)
		cron.DelayAddTimes(500*time.Millisecond, "* * * * * *", 2, func() {
			glog.Println("add times")
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

	dtest.C(t, func(t *dtest.T) {
		var (
			cron  = dcron.NewCron()
			array = darray.New(true)
		)
		entry, err1 := cron.Add("* * * * * *", func() {
			glog.Println("add")
			array.Append(1)
		})
		t.Assert(err1, nil)
		t.Assert(array.Len(), 0)
		t.Assert(cron.Size(), 1)
		time.Sleep(1300 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Stop()
		time.Sleep(5000 * time.Millisecond)
		t.Assert(array.Len(), 1)
		t.Assert(cron.Size(), 1)
		entry.Start()
		glog.Println("start")
		time.Sleep(1000 * time.Millisecond)
		t.Assert(array.Len(), 2)
		t.Assert(cron.Size(), 1)
		entry.Close()
		time.Sleep(1200 * time.Millisecond)
		t.Assert(cron.Size(), 0)
	})
}
