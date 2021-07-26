// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dfsnotify_test

import (
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dfsnotify"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

func TestWatcher_AddOnce(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		value := dtype.New()
		path := dfile.TempDir(dconv.String(dtime.TimestampNano()))
		err := dfile.PutContents(path, "init")
		t.Assert(err, nil)
		defer dfile.Remove(path)

		time.Sleep(100 * time.Millisecond)
		callback1, err := dfsnotify.AddOnce("mywatch", path, func(event *dfsnotify.Event) {
			value.Set(1)
		})
		t.Assert(err, nil)
		callback2, err := dfsnotify.AddOnce("mywatch", path, func(event *dfsnotify.Event) {
			value.Set(2)
		})
		t.Assert(err, nil)
		t.Assert(callback2, nil)

		err = dfile.PutContents(path, "1")
		t.Assert(err, nil)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)

		err = dfsnotify.RemoveCallback(callback1.Id)
		t.Assert(err, nil)

		err = dfile.PutContents(path, "2")
		t.Assert(err, nil)

		time.Sleep(100 * time.Millisecond)
		t.Assert(value, 1)
	})
}

func TestWatcher_AddRemove(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path1 := dfile.TempDir() + dfile.Separator + dconv.String(dtime.TimestampNano())
		path2 := dfile.TempDir() + dfile.Separator + dconv.String(dtime.TimestampNano()) + "2"
		dfile.PutContents(path1, "1")
		defer func() {
			dfile.Remove(path1)
			dfile.Remove(path2)
		}()
		v := dtype.NewInt(1)
		callback, err := dfsnotify.Add(path1, func(event *dfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRename() {
				v.Set(3)
				dfsnotify.Exit()
				return
			}
		})
		t.Assert(err, nil)
		t.AssertNE(callback, nil)

		dfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		dfile.Rename(path1, path2)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})

	dtest.C(t, func(t *dtest.T) {
		path1 := dfile.TempDir() + dfile.Separator + dconv.String(dtime.TimestampNano())
		dfile.PutContents(path1, "1")
		defer func() {
			dfile.Remove(path1)
		}()
		v := dtype.NewInt(1)
		callback, err := dfsnotify.Add(path1, func(event *dfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
			if event.IsRemove() {
				v.Set(4)
				return
			}
		})
		t.Assert(err, nil)
		t.AssertNE(callback, nil)

		dfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		dfile.Remove(path1)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)

		dfile.PutContents(path1, "1")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 4)
	})
}

func TestWatcher_Callback1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path1 := dfile.TempDir(dtime.TimestampNanoStr())
		dfile.PutContents(path1, "1")
		defer func() {
			dfile.Remove(path1)
		}()
		v := dtype.NewInt(1)
		callback, err := dfsnotify.Add(path1, func(event *dfsnotify.Event) {
			if event.IsWrite() {
				v.Set(2)
				return
			}
		})
		t.Assert(err, nil)
		t.AssertNE(callback, nil)

		dfile.PutContents(path1, "2")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 2)

		v.Set(3)
		dfsnotify.RemoveCallback(callback.Id)
		dfile.PutContents(path1, "3")
		time.Sleep(100 * time.Millisecond)
		t.Assert(v.Val(), 3)
	})
}

func TestWatcher_Callback2(t *testing.T) {
	// multiple callbacks
	dtest.C(t, func(t *dtest.T) {
		path1 := dfile.TempDir(dtime.TimestampNanoStr())
		t.Assert(dfile.PutContents(path1, "1"), nil)
		defer func() {
			dfile.Remove(path1)
		}()
		v1 := dtype.NewInt(1)
		v2 := dtype.NewInt(1)
		callback1, err1 := dfsnotify.Add(path1, func(event *dfsnotify.Event) {
			if event.IsWrite() {
				v1.Set(2)
				return
			}
		})
		callback2, err2 := dfsnotify.Add(path1, func(event *dfsnotify.Event) {
			if event.IsWrite() {
				v2.Set(2)
				return
			}
		})
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.AssertNE(callback1, nil)
		t.AssertNE(callback2, nil)

		t.Assert(dfile.PutContents(path1, "2"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 2)
		t.Assert(v2.Val(), 2)

		v1.Set(3)
		v2.Set(3)
		dfsnotify.RemoveCallback(callback1.Id)
		t.Assert(dfile.PutContents(path1, "3"), nil)
		time.Sleep(100 * time.Millisecond)
		t.Assert(v1.Val(), 3)
		t.Assert(v2.Val(), 2)
	})
}

func TestWatcher_WatchFolderWithoutRecursively(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			err     error
			array   = darray.New(true)
			dirPath = dfile.TempDir(dtime.TimestampNanoStr())
		)
		err = dfile.Mkdir(dirPath)
		t.AssertNil(err)

		_, err = dfsnotify.Add(dirPath, func(event *dfsnotify.Event) {
			//fmt.Println(event.String())
			array.Append(1)
		}, false)
		t.AssertNil(err)
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 0)

		f, err := dfile.Create(dfile.Join(dirPath, "1"))
		t.AssertNil(err)
		t.AssertNil(f.Close())
		time.Sleep(time.Millisecond * 100)
		t.Assert(array.Len(), 1)
	})
}
