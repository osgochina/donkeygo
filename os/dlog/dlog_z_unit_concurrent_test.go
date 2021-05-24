// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog_test

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"sync"
	"testing"
)

func Test_Concurrent(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		c := 1000
		l := dlog.New()
		s := "@1234567890#"
		f := "test.log"
		p := dfile.TempDir(gtime.TimestampNanoStr())
		t.Assert(l.SetPath(p), nil)
		defer dfile.Remove(p)
		wg := sync.WaitGroup{}
		ch := make(chan struct{})
		for i := 0; i < c; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
				l.File(f).Stdout(false).Print(s)
			}()
		}
		close(ch)
		wg.Wait()
		content := dfile.GetContents(dfile.Join(p, f))
		t.Assert(dstr.Count(content, s), c)
	})
}
