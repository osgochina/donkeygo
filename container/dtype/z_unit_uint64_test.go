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
	"math"
	"sync"
	"testing"
)

type Temp struct {
	Name string
	Age  int
}

func Test_Uint64(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var wg sync.WaitGroup
		addTimes := 1000
		i := dtype.NewUint64(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(1), uint64(0))
		t.AssertEQ(iClone.Val(), uint64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(uint64(addTimes), i.Val())

		//空参测试
		i1 := dtype.NewUint64()
		t.AssertEQ(i1.Val(), uint64(0))
	})
}
func Test_Uint64_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		i := dtype.NewUint64(math.MaxUint64)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewUint64()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), i)
	})
}

func Test_Uint64_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Uint64
	}
	dtest.C(t, func(t *dtest.T) {
		var v *V
		err := dconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}
