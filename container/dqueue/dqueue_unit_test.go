package dqueue_test

import (
	"github.com/osgochina/donkeygo/container/dqueue"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
	"time"
)

func TestQueue_Len(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		max := 100
		for n := 10; n < max; n++ {
			q1 := dqueue.New(max)
			for i := 0; i < max; i++ {
				q1.Push(i)
			}
			t.Assert(q1.Len(), max)
			t.Assert(q1.Size(), max)
		}
	})
}

func TestQueue_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		q := dqueue.New()
		for i := 0; i < 100; i++ {
			q.Push(i)
		}
		t.Assert(q.Pop(), 0)
		t.Assert(q.Pop(), 1)
	})
}

func TestQueue_Pop(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		q1 := dqueue.New()
		q1.Push(1)
		q1.Push(2)
		q1.Push(3)
		q1.Push(4)
		i1 := q1.Pop()
		t.Assert(i1, 1)
	})
}

func TestQueue_Close(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		q1 := dqueue.New()
		q1.Push(1)
		q1.Push(2)
		time.Sleep(time.Millisecond)
		t.Assert(q1.Len(), 2)
		q1.Close()
	})
}
