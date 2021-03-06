// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
	"time"
)

func Test_Rotate_Size(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		l := dlog.New()
		p := dfile.TempDir(dtime.TimestampNanoStr())
		err := l.SetConfigWithMap(d.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateSize":           10,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.Assert(err, nil)
		defer dfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			fmt.Println("logging content index:", i)
			l.Print(s)
		}

		time.Sleep(time.Second * 3)

		files, err := dfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 2)

		content := dfile.GetContents(dfile.Join(p, "access.log"))
		t.Assert(dstr.Count(content, s), 1)

		time.Sleep(time.Second * 5)
		files, err = dfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
}

func Test_Rotate_Expire(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		l := dlog.New()
		p := dfile.TempDir(dtime.TimestampNanoStr())
		err := l.SetConfigWithMap(d.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateExpire":         time.Second,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.Assert(err, nil)
		defer dfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			l.Print(s)
		}

		files, err := dfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)

		t.Assert(dstr.Count(dfile.GetContents(dfile.Join(p, "access.log")), s), 10)

		time.Sleep(time.Second * 3)

		files, err = dfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 1)

		t.Assert(dstr.Count(dfile.GetContents(dfile.Join(p, "access.log")), s), 0)

		time.Sleep(time.Second * 5)
		files, err = dfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
}
