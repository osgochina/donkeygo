// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dmap_test

import (
	"github.com/osgochina/donkeygo/container/dmap"
	"github.com/osgochina/donkeygo/util/dutil"
	"testing"
)

var hashMap = dmap.New(true)
var listMap = dmap.NewListMap(true)
var treeMap = dmap.NewTreeMap(dutil.ComparatorInt, true)

func Benchmark_HashMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_ListMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			listMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_TreeMap_Set(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			treeMap.Set(i, i)
			i++
		}
	})
}

func Benchmark_HashMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hashMap.Get(i)
			i++
		}
	})
}

func Benchmark_ListMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			listMap.Get(i)
			i++
		}
	})
}

func Benchmark_TreeMap_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			treeMap.Get(i)
			i++
		}
	})
}
