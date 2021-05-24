package dutil

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"strings"
)

// Comparator 对比方法
// 返回的数字表示的含义:
//    负数 , if a < b
//    0     , if a == b
//    正数 , if a > b
type Comparator func(a, b interface{}) int

func ComparatorString(a, b interface{}) int {
	return strings.Compare(dconv.String(a), dconv.String(b))
}
