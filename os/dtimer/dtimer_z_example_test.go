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
