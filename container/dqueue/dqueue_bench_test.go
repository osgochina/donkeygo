// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dqueue_test

import (
	"github.com/osgochina/donkeygo/container/dqueue"
	"testing"
)

var bn = 20000000
var length = 1000000
var qstatic = dqueue.New(length)
var qdynamic = dqueue.New()
var cany = make(chan interface{}, length)

func Benchmark_Dqueue_StaticPushAndPop(b *testing.B) {
	b.N = bn
	for i := 0; i < b.N; i++ {
		qstatic.Push(i)
		qstatic.Pop()
	}
}

func Benchmark_Dqueue_DynamicPushAndPop(b *testing.B) {
	b.N = bn
	for i := 0; i < b.N; i++ {
		qdynamic.Push(i)
		qdynamic.Pop()
	}
}

func Benchmark_Dqueue_DynamicPush(b *testing.B) {
	b.N = bn
	for i := 0; i < b.N; i++ {
		qdynamic.Push(i)
	}
}

func Benchmark_Dqueue_DynamicPop(b *testing.B) {
	b.N = bn
	for i := 0; i < b.N; i++ {
		qdynamic.Pop()
	}
}

func Benchmark_Channel_PushAndPop(b *testing.B) {
	b.N = bn
	for i := 0; i < b.N; i++ {
		cany <- i
		<-cany
	}
}
