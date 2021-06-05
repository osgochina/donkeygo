// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtime_test

import (
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func Test_Json_Pointer(t *testing.T) {
	// Marshal
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time *dtime.Time
		}
		t1 := new(T)
		s := "2006-01-02 15:04:05"
		t1.Time = dtime.NewFromStr(s)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time *dtime.Time
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":null}`)
	})
	// Marshal nil omitempty
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time *dtime.Time `json:"time,omitempty"`
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{}`)
	})
	// Unmarshal
	dtest.C(t, func(t *dtest.T) {
		var t1 dtime.Time
		s := []byte(`"2006-01-02 15:04:05"`)
		err := json.UnmarshalUseNumber(s, &t1)
		t.Assert(err, nil)
		t.Assert(t1.String(), "2006-01-02 15:04:05")
	})
}

func Test_Json_Struct(t *testing.T) {
	// Marshal
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time dtime.Time
		}
		t1 := new(T)
		s := "2006-01-02 15:04:05"
		t1.Time = *dtime.NewFromStr(s)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":"2006-01-02 15:04:05"}`)
	})
	// Marshal nil
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time dtime.Time
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"Time":""}`)
	})
	// Marshal nil omitempty
	dtest.C(t, func(t *dtest.T) {
		type T struct {
			Time dtime.Time `json:"time,omitempty"`
		}
		t1 := new(T)
		j, err := json.Marshal(t1)
		t.Assert(err, nil)
		t.Assert(j, `{"time":""}`)
	})

}
