// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dlog

import (
	"bytes"
	"fmt"
	"github.com/osgochina/donkeygo/os/dfile"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
	"time"
)

func Test_To(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		w := bytes.NewBuffer(nil)
		To(w).Error(1, 2, 3)
		To(w).Errorf("%d %d %d", 1, 2, 3)
		t.Assert(dstr.Count(w.String(), defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(w.String(), "1 2 3"), 2)
	})
}

func Test_Path(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Cat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		cat := "category"
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Cat(cat).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Cat(cat).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, cat, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(content, "1 2 3"), 2)
	})
}

func Test_Level(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Level(LevelProd).Stdout(false).Debug(1, 2, 3)
		Path(path).File(file).Level(LevelProd).Stdout(false).Debug("%d %d %d", 1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelDebug]), 0)
		t.Assert(dstr.Count(content, "1 2 3"), 0)
	})
}

func Test_Skip(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Skip(10).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(content, "1 2 3"), 2)
		t.Assert(dstr.Count(content, "Stack"), 1)
	})
}

func Test_Stack(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Stack(false).Stdout(false).Error(1, 2, 3)
		Path(path).File(file).Stdout(false).Errorf("%d %d %d", 1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 2)
		t.Assert(dstr.Count(content, "1 2 3"), 2)
		t.Assert(dstr.Count(content, "Stack"), 1)
	})
}

func Test_StackWithFilter(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).StackWithFilter("none").Stdout(false).Error(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
		t.Assert(dstr.Count(content, "Stack"), 1)
	})
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).StackWithFilter("donkeygo").Stdout(false).Error(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
		t.Assert(dstr.Count(content, "Stack"), 0)
	})
}

func Test_Header(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Header(true).Stdout(false).Error(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
	})
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Header(false).Stdout(false).Error(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelError]), 0)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
	})
}

func Test_Line(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Line(true).Stdout(false).Debug(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelDebug]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
		t.Assert(dstr.Count(content, ".go"), 1)
		t.Assert(dstr.Contains(content, dfile.Separator), true)
	})
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Line(false).Stdout(false).Debug(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelDebug]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
		t.Assert(dstr.Count(content, ".go"), 1)
		t.Assert(dstr.Contains(content, dfile.Separator), false)
	})
}

func Test_Async(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Async().Stdout(false).Debug(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(content, "")
		time.Sleep(200 * time.Millisecond)

		content = dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelDebug]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
	})

	dtest.C(t, func(t *dtest.T) {
		path := dfile.TempDir(dtime.TimestampNanoStr())
		file := fmt.Sprintf(`%d.log`, dtime.TimestampNano())

		err := dfile.Mkdir(path)
		t.Assert(err, nil)
		defer dfile.Remove(path)

		Path(path).File(file).Async(false).Stdout(false).Debug(1, 2, 3)
		content := dfile.GetContents(dfile.Join(path, file))
		t.Assert(dstr.Count(content, defaultLevelPrefixes[LevelDebug]), 1)
		t.Assert(dstr.Count(content, "1 2 3"), 1)
	})
}
