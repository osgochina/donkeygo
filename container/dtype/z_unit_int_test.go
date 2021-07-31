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

func Test_Int(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var wg sync.WaitGroup
		addTimes := 1000
		i := dtype.NewInt(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(1), 0)
		t.AssertEQ(iClone.Val(), 1)
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(addTimes, i.Val())

		//空参测试
		i1 := dtype.NewInt()
		t.AssertEQ(i1.Val(), 0)
	})
}

func Test_Int_JSON(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		v := 666
		i := dtype.NewInt(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := dtype.NewInt()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Int_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *dtype.Int
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
