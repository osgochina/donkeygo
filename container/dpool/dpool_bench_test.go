// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package dpool_test

import (
	"github.com/osgochina/donkeygo/container/dpool"
	"sync"
	"testing"
	"time"
)

var pool = dpool.New(time.Hour, nil)
var syncPool = dpool.NewSyncPool(func() interface{} { return "hello" }, func(i interface{}) {
	i = "aaa"
})
var syncp = sync.Pool{}

func BenchmarkGPoolPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool.Put(i)
	}
}

func BenchmarkGPoolGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool.Get()
	}
}

func BenchmarkSyncPoolPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		syncPool.Put(i)
	}
}

func BenchmarkSyncPoolGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		syncPool.Get()
	}
}

func BenchmarkSyncPPut(b *testing.B) {
	for i := 0; i < b.N; i++ {
		syncp.Put(i)
	}
}

func BenchmarkSyncPGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		syncp.Get()
	}
}
