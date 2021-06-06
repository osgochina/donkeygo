package dset_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/osgochina/donkeygo/container/dset"
)

func ExampleStrSet_Contains() {
	var set dset.StrSet
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// false
	// true
}

func ExampleStrSet_Walk() {
	var (
		set    dset.StrSet
		names  = g.SliceStr{"user", "user_detail"}
		prefix = "gf_"
	)
	set.Add(names...)
	// Add prefix for given table names.
	set.Walk(func(item string) string {
		return prefix + item
	})
	fmt.Println(set.Slice())

	// May Output:
	// [gf_user gf_user_detail]
}
