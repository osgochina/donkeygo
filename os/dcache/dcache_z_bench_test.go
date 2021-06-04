// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dcache_test

import (
	"github.com/osgochina/donkeygo/os/dcache"
	"testing"
)

var (
	localCache    = dcache.New()
	localCacheLru = dcache.New(10000)
)

func Benchmark_CacheSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Set(i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Get(i)
			i++
		}
	})
}

func Benchmark_CacheRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCache.Remove(i)
			i++
		}
	})
}

func Benchmark_CacheLruSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Set(i, i, 0)
			i++
		}
	})
}

func Benchmark_CacheLruGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Get(i)
			i++
		}
	})
}

func Benchmark_CacheLruRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			localCacheLru.Remove(i)
			i++
		}
	})
}
