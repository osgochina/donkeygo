package dmutex_test

import (
	"github.com/osgochina/donkeygo/os/dmutex"
	"sync"
	"testing"
)

var (
	mu   = sync.Mutex{}
	rwmu = sync.RWMutex{}
	gmu  = dmutex.New()
)

func Benchmark_Mutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			mu.Unlock()
		}
	})
}

func Benchmark_RWMutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwmu.Lock()
			rwmu.Unlock()
		}
	})
}

func Benchmark_RWMutex_RLockRUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwmu.RLock()
			rwmu.RUnlock()
		}
	})
}

func Benchmark_DMutex_LockUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gmu.Lock()
			gmu.Unlock()
		}
	})
}

func Benchmark_DMutex_TryLock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if gmu.TryLock() {
				gmu.Unlock()
			}
		}
	})
}

func Benchmark_DMutex_RLockRUnlock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gmu.RLock()
			gmu.RUnlock()
		}
	})
}

func Benchmark_DMutex_TryRLock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if gmu.TryRLock() {
				gmu.RUnlock()
			}
		}
	})
}
