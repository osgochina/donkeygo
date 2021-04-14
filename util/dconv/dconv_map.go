package dconv

import (
	"donkeygo/internal/empty"
	"donkeygo/internal/utils"
	"encoding/json"
	"reflect"
	"strings"
)

// Map 把value值转换成map
func Map(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, false, tags...)
}

// MapDeep 把value值转换成map,并且支持递归转换
func MapDeep(value interface{}, tags ...string) map[string]interface{} {
	return doMapConvert(value, true, tags...)
}

// MapStrStr 转换成key和value都是字符串的map
func MapStrStr(value interface{}, tags ...string) map[string]string {
	if v, ok := value.(map[string]string); ok {
		return v
	}
	m := Map(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

// MapStrStrDeep 递归转换成key和value都是字符串的map
func MapStrStrDeep(value interface{}, tags ...string) map[string]string {
	if r, ok := value.(map[string]string); ok {
		return r
	}
	m := MapDeep(value, tags...)
	if len(m) > 0 {
		vMap := make(map[string]string, len(m))
		for k, v := range m {
			vMap[k] = String(v)
		}
		return vMap
	}
	return nil
}

//映射转换
func doMapConvert(value interface{}, recursive bool, tags ...string) map[string]interface{} {
	if value == nil {
		return nil
	}
	newTags := StructTagPriority
	switch len(tags) {
	case 0:
		//nothing
	case 1:
		newTags = append(strings.Split(tags[0], ","), StructTagPriority...)
	default:
		newTags = append(tags, StructTagPriority...)
	}
	dataMap := make(map[string]interface{})
	switch val := value.(type) {
	case string:
		if len(val) > 0 && val[0] == '{' && val[len(val)-1] == '}' {
			if err := json.Unmarshal([]byte(val), &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case []byte:
		if len(val) > 0 && val[0] == '{' && val[len(val)-1] == '}' {
			if err := json.Unmarshal(val, &dataMap); err != nil {
				return nil
			}
		} else {
			return nil
		}
	case map[interface{}]interface{}:
		for k, v := range val {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
		}
	case map[interface{}]string:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[interface{}]int:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[interface{}]uint:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[interface{}]float32:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[interface{}]float64:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[string]bool:
		for k, v := range val {
			dataMap[k] = v
		}
	case map[string]int:
		for k, v := range val {
			dataMap[k] = v
		}
	case map[string]uint:
		for k, v := range val {
			dataMap[k] = v
		}
	case map[string]float32:
		for k, v := range val {
			dataMap[k] = v
		}
	case map[string]float64:
		for k, v := range val {
			dataMap[k] = v
		}
	case map[string]interface{}:
		if recursive {
			for k, v := range val {
				dataMap[k] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
			}
		} else {
			return val
		}
	case map[int]interface{}:
		for k, v := range val {
			dataMap[String(k)] = doMapConvertForMapOrStructValue(false, v, recursive, newTags...)
		}
	case map[int]string:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	case map[uint]string:
		for k, v := range val {
			dataMap[String(k)] = v
		}
	default:
		var reflectValue reflect.Value
		if v, ok := value.(reflect.Value); ok {
			reflectValue = v
		} else {
			reflectValue = reflect.ValueOf(value)
		}
		reflectKind := reflectValue.Kind()
		for reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		// If <value> is type of array, it converts the value of even number index as its key and
		// the value of odd number index as its corresponding value, for example:
		// []string{"k1","v1","k2","v2"} => map[string]interface{}{"k1":"v1", "k2":"v2"}
		// []string{"k1","v1","k2"}      => map[string]interface{}{"k1":"v1", "k2":nil}
		case reflect.Slice, reflect.Array:
			length := reflectValue.Len()
			for i := 0; i < length; i += 2 {
				if i+1 < length {
					dataMap[String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
				} else {
					dataMap[String(reflectValue.Index(i).Interface())] = nil
				}
			}
		case reflect.Map, reflect.Struct:
			convertedValue := doMapConvertForMapOrStructValue(true, value, recursive, newTags...)
			if m, ok := convertedValue.(map[string]interface{}); ok {
				return m
			}
			return nil
		default:
			return nil
		}
	}
	return dataMap
}

func doMapConvertForMapOrStructValue(isRoot bool, value interface{}, recursive bool, tags ...string) interface{} {
	if isRoot == false && recursive == false {
		return value
	}
	var reflectValue reflect.Value
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		value = v.Interface()
	} else {
		reflectValue = reflect.ValueOf(value)
	}
	reflectKind := reflectValue.Kind()
	// If it is a pointer, we should find its real data type.
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		var (
			mapKeys = reflectValue.MapKeys()
			dataMap = make(map[string]interface{})
		)
		for _, k := range mapKeys {
			dataMap[String(k.Interface())] = doMapConvertForMapOrStructValue(
				false,
				reflectValue.MapIndex(k).Interface(),
				recursive,
				tags...,
			)
		}
		if len(dataMap) == 0 {
			return value
		}
		return dataMap
	case reflect.Struct:
		// Map converting interface check.
		if v, ok := value.(apiMapStrAny); ok {
			m := v.MapStrAny()
			if recursive {
				for k, v := range m {
					m[k] = doMapConvertForMapOrStructValue(false, v, recursive, tags...)
				}
			}
			return m
		}
		// Using reflect for converting.
		var (
			rtField     reflect.StructField
			rvField     reflect.Value
			dataMap     = make(map[string]interface{}) // result map.
			reflectType = reflectValue.Type()          // attribute value type.
			name        = ""                           // name may be the tag name or the struct attribute name.
		)
		for i := 0; i < reflectValue.NumField(); i++ {
			rtField = reflectType.Field(i)
			rvField = reflectValue.Field(i)
			// Only convert the public attributes.
			fieldName := rtField.Name
			if !utils.IsLetterUpper(fieldName[0]) {
				continue
			}
			name = ""
			fieldTag := rtField.Tag
			for _, tag := range tags {
				if name = fieldTag.Get(tag); name != "" {
					break
				}
			}
			if name == "" {
				name = fieldName
			} else {
				// Support json tag feature: -, omitempty
				name = strings.TrimSpace(name)
				if name == "-" {
					continue
				}
				array := strings.Split(name, ",")
				if len(array) > 1 {
					switch strings.TrimSpace(array[1]) {
					case "omitempty":
						if empty.IsEmpty(rvField.Interface()) {
							continue
						} else {
							name = strings.TrimSpace(array[0])
						}
					default:
						name = strings.TrimSpace(array[0])
					}
				}
			}
			if recursive || rtField.Anonymous {
				// Do map converting recursively.
				var (
					rvAttrField = rvField
					rvAttrKind  = rvField.Kind()
				)
				if rvAttrKind == reflect.Ptr {
					rvAttrField = rvField.Elem()
					rvAttrKind = rvAttrField.Kind()
				}
				switch rvAttrKind {
				case reflect.Struct:
					var (
						hasNoTag        = name == fieldName
						rvAttrInterface = rvAttrField.Interface()
					)
					if hasNoTag && rtField.Anonymous {
						// It means this attribute field has no tag.
						// Overwrite the attribute with sub-struct attribute fields.
						anonymousValue := doMapConvertForMapOrStructValue(false, rvAttrInterface, true, tags...)
						if m, ok := anonymousValue.(map[string]interface{}); ok {
							for k, v := range m {
								dataMap[k] = v
							}
						} else {
							dataMap[name] = rvAttrInterface
						}
					} else if !hasNoTag && rtField.Anonymous {
						// It means this attribute field has desired tag.
						dataMap[name] = doMapConvertForMapOrStructValue(false, rvAttrInterface, true, tags...)
					} else {
						dataMap[name] = doMapConvertForMapOrStructValue(false, rvAttrInterface, false, tags...)
					}

				// The struct attribute is type of slice.
				case reflect.Array, reflect.Slice:
					length := rvField.Len()
					if length == 0 {
						dataMap[name] = rvField.Interface()
						break
					}
					array := make([]interface{}, length)
					for i := 0; i < length; i++ {
						array[i] = doMapConvertForMapOrStructValue(false, rvField.Index(i), recursive, tags...)
					}
					dataMap[name] = array

				default:
					if rvField.IsValid() {
						dataMap[name] = reflectValue.Field(i).Interface()
					} else {
						dataMap[name] = nil
					}
				}
			} else {
				// No recursive map value converting
				if rvField.IsValid() {
					dataMap[name] = reflectValue.Field(i).Interface()
				} else {
					dataMap[name] = nil
				}
			}
		}
		if len(dataMap) == 0 {
			return value
		}
		return dataMap
		// The given value is type of slice.
	case reflect.Array, reflect.Slice:
		length := reflectValue.Len()
		if length == 0 {
			break
		}
		array := make([]interface{}, reflectValue.Len())
		for i := 0; i < length; i++ {
			array[i] = doMapConvertForMapOrStructValue(false, reflectValue.Index(i), recursive, tags...)
		}
		return array
	}
	return value
}
