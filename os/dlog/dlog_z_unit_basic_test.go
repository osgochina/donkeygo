// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog

import (
	"bytes"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
)

func Test_Print(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Print(1, 2, 3)
		l.Println(1, 2, 3)
		l.Printf("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), "["), 0)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 3)
	})
}

func Test_Debug(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Debug(1, 2, 3)
		l.Debugf("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelDebug]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Info(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Info(1, 2, 3)
		l.Infof("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelInfo]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Notice(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Notice(1, 2, 3)
		l.Noticef("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelNotice]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Warning(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Warning(1, 2, 3)
		l.Warningf("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelWarning]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Error(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		l := NewWithWriter(w)
		l.Error(1, 2, 3)
		l.Errorf("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}
