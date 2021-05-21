package dtimer_test

import (
	"github.com/osgochina/donkeygo/os/dtimer"
	"testing"
	"time"
)

var (
	timer = dtimer.NewTimer(5, 30*time.Millisecond)
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Add(time.Hour, func() {

		})
	}
}

func Benchmark_StartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Start()
		timer.Stop()
	}
}
