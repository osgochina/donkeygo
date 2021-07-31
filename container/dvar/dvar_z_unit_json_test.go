// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dvar_test

import (
	"github.com/osgochina/donkeygo/container/dvar"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"math"
	"testing"
)

func TestVar_Json(t *testing.T) {
	// Marshal
	dtest.C(t, func(t *dtest.T) {
		s := "i love gf"
		v := dvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})

	dtest.C(t, func(t *dtest.T) {
		s := int64(math.MaxInt64)
		v := dvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		t.Assert(err1, err2)
		t.Assert(b1, b2)
	})

	// Unmarshal
	dtest.C(t, func(t *dtest.T) {
		s := "i love gf"
		v := dvar.New(nil)
		b, err := json.Marshal(s)
		t.Assert(err, nil)

		err = json.UnmarshalUseNumber(b, v)
		t.Assert(err, nil)
		t.Assert(v.String(), s)
	})

	dtest.C(t, func(t *dtest.T) {
		var v dvar.Var
		s := "i love gf"
		b, err := json.Marshal(s)
		t.Assert(err, nil)

		err = json.UnmarshalUseNumber(b, &v)
		t.Assert(err, nil)
		t.Assert(v.String(), s)
	})
}
