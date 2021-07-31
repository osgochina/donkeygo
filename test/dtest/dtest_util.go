package dtest

import (
	"fmt"
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/util/dconv"
	"os"
	"reflect"
	"testing"
)

const (
	pathFilterKey = "/test/dtest/dtest"
)

// C 创建一个测试单元用例
func C(t *testing.T, f func(t *T)) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n%s", err, ddebug.StackWithFilter(pathFilterKey))
			t.Fail()
		}
	}()

	f(&T{t})
}

// Case 创建一个单元测试用例
func Case(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n%s", err, ddebug.StackWithFilter(pathFilterKey))
			t.Fail()
		}
	}()
	f()
}

// Assert 判断<value>和<expect>两个对象值是否相等
func Assert(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}

	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	var (
		strValue  = dconv.String(value)
		strExpect = dconv.String(expect)
	)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
}

// AssertEQ 判断<value>和<expect>两个对象是否相等,必须值与类型都相等
func AssertEQ(value, expect interface{}) {
	// Value assert.
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err != nil {
			panic(err)
		}
		return
	}
	strValue := dconv.String(value)
	strExpect := dconv.String(expect)
	if strValue != strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v == %v`, strValue, strExpect))
	}
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 != t2 {
		panic(fmt.Sprintf(`[ASSERT] EXPECT TYPE %v[%v] == %v[%v]`, strValue, t1, strExpect, t2))
	}
}

// AssertNE 判断<value>和<expect>的值是否相等，如果相等则报错
func AssertNE(value, expect interface{}) {
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if err := compareMap(value, expect); err == nil {
			panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, value, expect))
		}
		return
	}
	var (
		strValue  = dconv.String(value)
		strExpect = dconv.String(expect)
	)
	if strValue == strExpect {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v != %v`, strValue, strExpect))
	}
}

// AssertNQ  判断<value>和<expect>的值和类型是否相等，如果相等则报错
func AssertNQ(value, expect interface{}) {
	// Type assert.
	t1 := reflect.TypeOf(value)
	t2 := reflect.TypeOf(expect)
	if t1 == t2 {
		panic(
			fmt.Sprintf(
				`[ASSERT] EXPECT TYPE %v[%v] != %v[%v]`,
				dconv.String(value), t1, dconv.String(expect), t2,
			),
		)
	}
	// Value assert.
	AssertNE(value, expect)
}

// AssertGT 判断<value>的值是否大于<expect>的值
func AssertGT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = dconv.String(value) > dconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = dconv.Int(value) > dconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = dconv.Uint(value) > dconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = dconv.Float64(value) > dconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v > %v`, value, expect))
	}
}

// AssertGE 判断<value>的值是否大于等于<expect>的值
func AssertGE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = dconv.String(value) >= dconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = dconv.Int64(value) >= dconv.Int64(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = dconv.Uint64(value) >= dconv.Uint64(expect)

	case reflect.Float32, reflect.Float64:
		passed = dconv.Float64(value) >= dconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(
			`[ASSERT] EXPECT %v(%v) >= %v(%v)`,
			value, reflect.ValueOf(value).Kind(),
			expect, reflect.ValueOf(expect).Kind(),
		))
	}
}

// AssertLT 判断<value>的值是否小于<expect>的值
func AssertLT(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = dconv.String(value) < dconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = dconv.Int(value) < dconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = dconv.Uint(value) < dconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = dconv.Float64(value) < dconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v < %v`, value, expect))
	}
}

// AssertLE 判断<value>的值是否小于等于<expect>的值
func AssertLE(value, expect interface{}) {
	passed := false
	switch reflect.ValueOf(expect).Kind() {
	case reflect.String:
		passed = dconv.String(value) <= dconv.String(expect)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		passed = dconv.Int(value) <= dconv.Int(expect)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		passed = dconv.Uint(value) <= dconv.Uint(expect)

	case reflect.Float32, reflect.Float64:
		passed = dconv.Float64(value) <= dconv.Float64(expect)
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v <= %v`, value, expect))
	}
}

// AssertIN 判断<expect>的值是否包含<value>的值
func AssertIN(value, expect interface{}) {
	passed := true
	expectKind := reflect.ValueOf(expect).Kind()
	switch expectKind {
	case reflect.Slice, reflect.Array:
		expectSlice := dconv.Strings(expect)
		for _, v1 := range dconv.Strings(value) {
			result := false
			for _, v2 := range expectSlice {
				if v1 == v2 {
					result = true
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v IN %v`, value, expect))
	}
}

// AssertNI not in 判断<expect>的值是否包含<value>的值,如果包含则失败
func AssertNI(value, expect interface{}) {
	passed := true
	expectKind := reflect.ValueOf(expect).Kind()
	switch expectKind {
	case reflect.Slice, reflect.Array:
		for _, v1 := range dconv.Strings(value) {
			result := true
			for _, v2 := range dconv.Strings(expect) {
				if v1 == v2 {
					result = false
					break
				}
			}
			if !result {
				passed = false
				break
			}
		}
	default:
		panic(fmt.Sprintf(`[ASSERT] INVALID EXPECT VALUE TYPE: %v`, expectKind))
	}
	if !passed {
		panic(fmt.Sprintf(`[ASSERT] EXPECT %v NOT IN %v`, value, expect))
	}
}

// Error 产生错误
func Error(message ...interface{}) {
	panic(fmt.Sprintf("[ERROR] %s", fmt.Sprint(message...)))
}

// Fatal 输出错误，并退出
func Fatal(message ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[FATAL] %s\n%s", fmt.Sprint(message...), ddebug.StackWithFilter(pathFilterKey))
	os.Exit(1)
}

//对比两个map对象是否相等
func compareMap(value, expect interface{}) error {
	rvValue := reflect.ValueOf(value)
	rvExpect := reflect.ValueOf(expect)
	if empty.IsNil(value) {
		value = nil
	}
	if rvExpect.Kind() == reflect.Map {
		if rvValue.Kind() == reflect.Map {
			if rvExpect.Len() == rvValue.Len() {
				// Turn two interface maps to the same type for comparison.
				// Direct use of rvValue.MapIndex(key).Interface() will panic
				// when the key types are inconsistent.
				mValue := make(map[string]string)
				mExpect := make(map[string]string)
				ksValue := rvValue.MapKeys()
				ksExpect := rvExpect.MapKeys()
				for _, key := range ksValue {
					mValue[dconv.String(key.Interface())] = dconv.String(rvValue.MapIndex(key).Interface())
				}
				for _, key := range ksExpect {
					mExpect[dconv.String(key.Interface())] = dconv.String(rvExpect.MapIndex(key).Interface())
				}
				for k, v := range mExpect {
					if v != mValue[k] {
						return fmt.Errorf(`[ASSERT] EXPECT VALUE map["%v"]:%v == map["%v"]:%v`+
							"\nGIVEN : %v\nEXPECT: %v", k, mValue[k], k, v, mValue, mExpect)
					}
				}
			} else {
				return fmt.Errorf(`[ASSERT] EXPECT MAP LENGTH %d == %d`, rvValue.Len(), rvExpect.Len())
			}
		} else {
			return fmt.Errorf(`[ASSERT] EXPECT VALUE TO BE A MAP`)
		}
	}
	return nil
}

//判断对象是否为nil
//func isNil(value interface{}) bool {
//	rv := reflect.ValueOf(value)
//	switch rv.Kind() {
//	case reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr, reflect.Func:
//		return rv.IsNil()
//	default:
//		return value == nil
//	}
//}

// AssertNil asserts `value` is nil.
func AssertNil(value interface{}) {
	if empty.IsNil(value) {
		return
	}
	if err, ok := value.(error); ok {
		panic(fmt.Sprintf(`%+v`, err))
	}
	AssertNE(value, nil)
}
