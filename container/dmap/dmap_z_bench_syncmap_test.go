// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dmap_test

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"sync"
	"testing"
)

var gm = dmap.NewIntIntMap(true)
var sm = sync.Map{}

func Benchmark_dmapSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Set(i, i)
			i++
		}
	})
}

func Benchmark_SyncMapSet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Store(i, i)
			i++
		}
	})
}

func Benchmark_dmapGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Get(i)
			i++
		}
	})
}

func Benchmark_SyncMapGet(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Load(i)
			i++
		}
	})
}

func Benchmark_dmapRemove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gm.Remove(i)
			i++
		}
	})
}

func Benchmark_SyncMapRmove(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sm.Delete(i)
			i++
		}
	})
}
