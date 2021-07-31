package dmlock_test

import (
	"github.com/osgochina/donkeygo/os/dmlock"
	"testing"
)

var (
	lockKey = "This is the lock key for dmlock."
)

func Benchmark_DMLock_Lock_Unlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dmlock.Lock(lockKey)
		dmlock.Unlock(lockKey)
	}
}

func Benchmark_DMLock_RLock_RUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dmlock.RLock(lockKey)
		dmlock.RUnlock(lockKey)
	}
}

func Benchmark_DMLock_TryLock_Unlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if dmlock.TryLock(lockKey) {
			dmlock.Unlock(lockKey)
		}
	}
}

func Benchmark_DMLock_TryRLock_RUnlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if dmlock.TryRLock(lockKey) {
			dmlock.RUnlock(lockKey)
		}
	}
}
