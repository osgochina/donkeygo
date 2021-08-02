// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package darray_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/container/darray"
	"github.com/osgochina/donkeygo/frame/d"
)

func ExampleStrArray_Walk() {
	var array darray.StrArray
	tables := d.SliceStr{"user", "user_detail"}
	prefix := "gf_"
	array.Append(tables...)
	// Add prefix for given table names.
	array.Walk(func(value string) string {
		return prefix + value
	})
	fmt.Println(array.Slice())

	// Output:
	// [gf_user gf_user_detail]
}
