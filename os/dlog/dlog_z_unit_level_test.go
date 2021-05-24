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

func Test_LevelPrefix(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		l := New()
		t.Assert(l.GetLevelPrefix(LevelDebug), defaultLevelPrefixes[LevelDebug])
		t.Assert(l.GetLevelPrefix(LevelInfo), defaultLevelPrefixes[LevelInfo])
		t.Assert(l.GetLevelPrefix(LevelNotice), defaultLevelPrefixes[LevelNotice])
		t.Assert(l.GetLevelPrefix(LevelWarning), defaultLevelPrefixes[LevelWarning])
		t.Assert(l.GetLevelPrefix(LevelProd), defaultLevelPrefixes[LevelProd])
		t.Assert(l.GetLevelPrefix(LevelCritical), defaultLevelPrefixes[LevelCritical])
		l.SetLevelPrefix(LevelDebug, "debug")
		t.Assert(l.GetLevelPrefix(LevelDebug), "debug")
		l.SetLevelPrefixes(map[int]string{
			LevelCritical: "critical",
		})
		t.Assert(l.GetLevelPrefix(LevelDebug), "debug")
		t.Assert(l.GetLevelPrefix(LevelInfo), defaultLevelPrefixes[LevelInfo])
		t.Assert(l.GetLevelPrefix(LevelNotice), defaultLevelPrefixes[LevelNotice])
		t.Assert(l.GetLevelPrefix(LevelWarning), defaultLevelPrefixes[LevelWarning])
		t.Assert(l.GetLevelPrefix(LevelError), defaultLevelPrefixes[LevelError])
		t.Assert(l.GetLevelPrefix(LevelCritical), "critical")
	})
	dtest.C(t, func(t *dtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		l.SetWriter(buffer)
		l.Debug("test1")
		t.Assert(dstr.Contains(buffer.String(), defaultLevelPrefixes[LevelDebug]), true)

		buffer.Reset()

		l.SetLevelPrefix(LevelDebug, "debug")
		l.Debug("test2")
		t.Assert(dstr.Contains(buffer.String(), defaultLevelPrefixes[LevelDebug]), false)
		t.Assert(dstr.Contains(buffer.String(), "debug"), true)

		buffer.Reset()
		l.SetLevelPrefixes(map[int]string{
			LevelError: "error",
		})
		l.Error("test3")
		t.Assert(dstr.Contains(buffer.String(), defaultLevelPrefixes[LevelError]), false)
		t.Assert(dstr.Contains(buffer.String(), "error"), true)
	})
}
