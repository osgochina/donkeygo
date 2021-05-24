// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog

import (
	"bytes"
	"github.com/osgochina/donkeygo/test/dtest"
	"strings"
	"testing"
)

func Test_SetConfigWithMap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		l := New()
		m := map[string]interface{}{
			"path":     "/var/log",
			"level":    "all",
			"stdout":   false,
			"StStatus": 0,
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)
		t.Assert(l.config.Path, m["path"])
		t.Assert(l.config.Level, LevelAll)
		t.Assert(l.config.StdoutPrint, m["stdout"])
	})
}

func Test_SetConfigWithMap_LevelStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "all",
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)

		l.SetWriter(buffer)

		l.Debug("test")
		l.Warning("test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), true)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})

	dtest.C(t, func(t *dtest.T) {
		buffer := bytes.NewBuffer(nil)
		l := New()
		m := map[string]interface{}{
			"level": "warn",
		}
		err := l.SetConfigWithMap(m)
		t.Assert(err, nil)
		l.SetWriter(buffer)
		l.Debug("test")
		l.Warning("test")
		t.Assert(strings.Contains(buffer.String(), "DEBU"), false)
		t.Assert(strings.Contains(buffer.String(), "WARN"), true)
	})
}
