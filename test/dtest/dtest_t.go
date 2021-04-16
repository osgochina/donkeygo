package dtest

import "testing"

//T 是测试单元的管理对象
type T struct {
	*testing.T
}

// Assert 判断<value>和<expect>两个对象值是否相等
func (that *T) Assert(value, expect interface{}) {
	Assert(value, expect)
}

// AssertEQ 判断<value>和<expect>两个对象是否相等,必须值与类型都相等
func (that *T) AssertEQ(value, expect interface{}) {
	AssertEQ(value, expect)
}

// AssertNE 判断<value>和<expect>的值是否相等，如果相等则报错
func (that *T) AssertNE(value, expect interface{}) {
	AssertNE(value, expect)
}

// AssertNQ  判断<value>和<expect>的值和类型是否相等，如果相等则报错
func (that *T) AssertNQ(value, expect interface{}) {
	AssertNQ(value, expect)
}

//AssertGT 判断<value>的值是否大于<expect>的值
func (that *T) AssertGT(value, expect interface{}) {
	AssertGT(value, expect)
}

//AssertGE 判断<value>的值是否大于等于<expect>的值
func (that *T) AssertGE(value, expect interface{}) {
	AssertGE(value, expect)
}

// AssertLT 判断<value>的值是否小于<expect>的值
func (that *T) AssertLT(value, expect interface{}) {
	AssertLT(value, expect)
}

// AssertLE 判断<value>的值是否小于等于<expect>的值
func (that *T) AssertLE(value, expect interface{}) {
	AssertLE(value, expect)
}

// AssertIN 判断<expect>的值是否包含<value>的值
func (that *T) AssertIN(value, expect interface{}) {
	AssertIN(value, expect)
}

//AssertNI not in 判断<expect>的值是否包含<value>的值,如果包含则失败
func (that *T) AssertNI(value, expect interface{}) {
	AssertNI(value, expect)
}

//Error 产生错误
func (that *T) Error(message ...interface{}) {
	Error(message...)
}

//Fatal 输出错误，并退出
func (that *T) Fatal(message ...interface{}) {
	Fatal(message...)
}
