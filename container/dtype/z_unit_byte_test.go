// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtype_test

import (
	"github.com/osgochina/donkeygo/container/dtype"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/util/dconv"
	"sync"
	"testing"
)

func Test_Byte(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var wg sync.WaitGroup
		addTimes := 127
		i := dtype.NewByte(byte(0))
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(byte(1)), byte(0))
		t.AssertEQ(iClone.Val(), byte(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(byte(addTimes), i.Val())

		//空参测试
		i1 := dtype.NewByte()
		t.AssertEQ(i1.Val(), byte(0))
	})
}

func Test_Byte_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewByte(49)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)
	})
	// Unmarshal
	dtest.C(t, func(t *dtest.T) {
		var err error
		i := dtype.NewByte()
		err = json.UnmarshalUseNumber([]byte("49"), &i)
		t.Assert(err, nil)
		t.Assert(i.Val(), "49")
	})
}

func Test_Byte_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Byte
	}
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "2",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "2")
	})
}
