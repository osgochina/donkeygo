package dconv

import "reflect"

// SliceAny 把任何对象转换成切片
func SliceAny(i interface{}) []interface{} {
	return Interfaces(i)
}

func Interfaces(i interface{}) []interface{} {
	if i == nil {
		return nil
	}
	if r, ok := i.([]interface{}); ok {
		return r
	} else if r, ok := i.(apiInterfaces); ok {
		return r.Interfaces()
	} else {
		var array []interface{}
		switch value := i.(type) {
		case []string:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int8:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int16:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int32:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []int64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint8:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint16:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []uint32:
			for _, v := range value {
				array = append(array, v)
			}
		case []uint64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []bool:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []float32:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		case []float64:
			array = make([]interface{}, len(value))
			for k, v := range value {
				array[k] = v
			}
		default:
			var (
				reflectValue = reflect.ValueOf(i)
				reflectKind  = reflectValue.Kind()
			)
			for reflectKind == reflect.Ptr {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
			switch reflectKind {
			case reflect.Slice, reflect.Array:
				array = make([]interface{}, reflectValue.Len())
				for i := 0; i < reflectValue.Len(); i++ {
					array[i] = reflectValue.Index(i).Interface()
				}
			default:
				return []interface{}{i}
			}
		}
		return array
	}
}
