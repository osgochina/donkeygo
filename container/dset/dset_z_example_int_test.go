// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package dset_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/container/dset"
)

func ExampleIntSet_Contains() {
	var set dset.IntSet
	set.Add(1)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))

	// Output:
	// true
	// false
}
