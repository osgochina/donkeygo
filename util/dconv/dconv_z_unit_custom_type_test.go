// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dconv_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"testing"
	"time"
)

type Duration time.Duration

// UnmarshalText unmarshal text to duration.
func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}

func Test_Struct_CustomTimeDuration_Attribute(t *testing.T) {
	type A struct {
		Name    string
		Timeout Duration
	}
	dtest.C(t, func(t *dtest.T) {
		var a A
		err := dconv.Struct(d.Map{
			"name":    "john",
			"timeout": "1s",
		}, &a)
		t.Assert(err, nil)
	})
}
