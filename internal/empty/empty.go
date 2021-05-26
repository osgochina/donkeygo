package empty

import "reflect"

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiInterfaces is used for type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// IsEmpty 判断任意一个值是否为空
//It returns true if <value> is in: 0, nil, false, "", len(slice/map/chan) == 0
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch val := value.(type) {
	case int:
		return val == 0
	case int8:
		return val == 0
	case int16:
		return val == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case uint:
		return val == 0
	case uint8:
		return val == 0
	case uint16:
		return val == 0
	case uint32:
		return val == 0
	case uint64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case bool:
		return val == false
	case string:
		return val == ""
	case []byte:
		return len(val) == 0
	case []rune:
		return len(val) == 0
	case []int:
		return len(val) == 0
	case []string:
		return len(val) == 0
	case []float32:
		return len(val) == 0
	case []float64:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	default:
		// Common interfaces checks.
		if f, ok := val.(apiString); ok {
			if f == nil {
				return true
			}
			return f.String() == ""
		}
		if f, ok := val.(apiInterfaces); ok {
			if f == nil {
				return true
			}
			return len(f.Interfaces()) == 0
		}
		if f, ok := value.(apiMapStrAny); ok {
			if f == nil {
				return true
			}
			return len(f.MapStrAny()) == 0
		}
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}
		switch rv.Kind() {
		case reflect.Bool:
			return !rv.Bool()
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			return rv.Int() == 0
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr:
			return rv.Uint() == 0
		case reflect.Float32,
			reflect.Float64:
			return rv.Float() == 0
		case reflect.String:
			return rv.Len() == 0
		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				if !IsEmpty(rv) {
					return false
				}
			}
			return true
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			return rv.Len() == 0
		case reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}

// IsNil 判断值是否为nil
func IsNil(value interface{}, traceSource ...bool) bool {
	if value == nil {
		return true
	}
	var rv reflect.Value
	if v, ok := value.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(value)
	}
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer:
		return !rv.IsValid() || rv.IsNil()

	case reflect.Ptr:
		if len(traceSource) > 0 && traceSource[0] {
			for rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if !rv.IsValid() {
				return true
			}
			if rv.Kind() == reflect.Ptr {
				return rv.IsNil()
			}
		} else {
			return !rv.IsValid() || rv.IsNil()
		}
	}
	return false
}
