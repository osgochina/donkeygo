package dconv

import (
	"donkeygo/encoding/dbinary"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func String(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	//case gtime.Time:
	//	if value.IsZero() {
	//		return ""
	//	}
	//	return value.String()
	//case *gtime.Time:
	//	if value == nil {
	//		return ""
	//	}
	//	return value.String()
	default:
		//// Empty checks.
		//if value == nil {
		//	return ""
		//}
		//if f, ok := value.(apiString); ok {
		//	// If the variable implements the String() interface,
		//	// then use that interface to perform the conversion
		//	return f.String()
		//}
		//if f, ok := value.(apiError); ok {
		//	// If the variable implements the Error() interface,
		//	// then use that interface to perform the conversion
		//	return f.Error()
		//}
		//// Reflect checks.
		//var (
		//	rv   = reflect.ValueOf(value)
		//	kind = rv.Kind()
		//)
		//switch kind {
		//case reflect.Chan,
		//	reflect.Map,
		//	reflect.Slice,
		//	reflect.Func,
		//	reflect.Ptr,
		//	reflect.Interface,
		//	reflect.UnsafePointer:
		//	if rv.IsNil() {
		//		return ""
		//	}
		//case reflect.String:
		//	return rv.String()
		//}
		//if kind == reflect.Ptr {
		//	return String(rv.Elem().Interface())
		//}
		//// Finally we use json.Marshal to convert.
		//if jsonContent, err := json.Marshal(value); err != nil {
		//	return fmt.Sprint(value)
		//} else {
		//	return string(jsonContent)
		//}
		return ""
	}
}

var emptyStringMap = map[string]struct{}{
	"":      {},
	"0":     {},
	"no":    {},
	"off":   {},
	"false": {},
}

func Bool(i interface{}) bool {
	if i == nil {
		return false
	}
	switch value := i.(type) {
	case bool:
		return value
	case []byte:
		if _, ok := emptyStringMap[strings.ToLower(string(value))]; ok {
			return false
		}
		return true
	case string:
		if _, ok := emptyStringMap[strings.ToLower(value)]; ok {
			return false
		}
		return true
	default:
		rv := reflect.ValueOf(i)
		switch rv.Kind() {
		case reflect.Ptr:
			return !rv.IsNil()
		case reflect.Map:
			fallthrough
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			return rv.Len() != 0
		case reflect.Struct:
			return true
		default:
			s := strings.ToLower(String(i))
			if _, ok := emptyStringMap[s]; ok {
				return false
			}
			return true
		}
	}

}

func Int(i interface{}) int {
	if i == nil {
		return 0
	}
	if v, ok := i.(int); ok {
		return v
	}
	return int(Int64(i))
}

func Int64(i interface{}) int64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case int64:
		return value
	case bool:
		if value {
			return 1
		}
		return 0
	case []byte:
		return dbinary.DecodeToInt64(value)
	default:
		s := String(i)
		isMinus := false

		//判断是否有符号
		if len(s) > 0 {
			if s[0] == '-' {
				isMinus = true
				s = s[1:]
			} else if s[0] == '+' {
				s = s[1:]
			}
		}
		//转换16进制
		if len(s) > 2 && s[0] == '0' && (s[0] == 'x' || s[0] == 'X') {
			if v, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
				if isMinus {
					return -v
				}
				return v
			}
		}
		//转换8进制
		if len(s) > 1 && s[0] == '0' {
			if v, e := strconv.ParseInt(s[1:], 8, 64); e == nil {
				if isMinus {
					return -v
				}
				return v
			}
		}

		//转换10进制
		if v, e := strconv.ParseInt(s, 10, 64); e == nil {
			if isMinus {
				return -v
			}
			return v
		}

		//实在无法处理，看样子是浮点数，先转换下，再强转
		return int64(Float64(value))
	}
}

func Float64(i interface{}) float64 {
	if i == nil {
		return 0
	}
	switch value := i.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return dbinary.DecodeToFloat64(value)
	default:
		v, _ := strconv.ParseFloat(String(i), 64)
		return v
	}
}
