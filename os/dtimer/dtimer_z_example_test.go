// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dtimer_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/os/dtimer"
	"time"
)

func Example_add() {
	now := time.Now()
	interval := 1400 * time.Millisecond
	dtimer.Add(interval, func() {
		fmt.Println(time.Now(), time.Duration(time.Now().UnixNano()-now.UnixNano()))
		now = time.Now()
	})

	select {}
}
