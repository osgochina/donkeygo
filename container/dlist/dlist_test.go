package dlist_test

import (
	"fmt"
	"github.com/osgochina/donkeygo/container/dlist"
	"testing"
)

func ExampleNew() {
	n := 10
	l := dlist.New()
	for i := 0; i < n; i++ {
		l.PushBack(i)
	}
	fmt.Println(l.Len())
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())
	for i := 0; i < n; i++ {
		fmt.Print(l.PopFront())
	}
	l.Clear()
	fmt.Println()
	fmt.Println(l.Len())

	// Output:
	// 10
	// [0 1 2 3 4 5 6 7 8 9]
	// [9 8 7 6 5 4 3 2 1 0]
	// 0123456789
	// 0
}

func TestMove(t *testing.T) {
	l := dlist.New()
	e1 := l.PushBack(1)
	e2 := l.PushBack(2)
	e3 := l.PushBack(3)
	e4 := l.PushBack(4)
	e5 := l.PushBack(5)

	fmt.Println(l.FrontAll())
	l.MoveAfter(e4, e5)
	fmt.Println(l.FrontAll())
	l.MoveAfter(e1, e2)
	l.MoveAfter(e1, e3)
	l.MoveAfter(e1, e5)
}
