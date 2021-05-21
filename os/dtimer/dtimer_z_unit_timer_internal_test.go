// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtimer

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestTimer_Proceed(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		index := dtype.NewInt()
		array := darray.New(true)
		timer := doNewWithoutAutoStart(10, 60*time.Millisecond, 6)
		timer.nowFunc = func() time.Time {
			return time.Now().Add(time.Duration(index.Add(1)) * time.Millisecond * 60)
		}
		timer.AddOnce(2*time.Second, func() {
			array.Append(1)
		})
		timer.AddOnce(1*time.Minute, func() {
			array.Append(2)
		})
		timer.AddOnce(5*time.Minute, func() {
			array.Append(3)
		})
		timer.AddOnce(1*time.Hour, func() {
			array.Append(4)
		})
		timer.AddOnce(100*time.Minute, func() {
			array.Append(5)
		})
		timer.AddOnce(2*time.Hour, func() {
			array.Append(6)
		})
		for i := 0; i < 500000; i++ {
			timer.wheels[0].proceed()
			time.Sleep(10 * time.Microsecond)
		}
		time.Sleep(time.Second)
		t.Assert(array.Slice(), []int{1, 2, 3, 4, 5, 6})
	})
}
