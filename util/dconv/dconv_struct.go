package dconv

import (
	"encoding/json"
	"fmt"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/internal/empty"
	"github.com/osgochina/donkeygo/internal/structs"
	"github.com/osgochina/donkeygo/internal/utils"
	"reflect"
	"strings"
)

// Struct 将参数params 的key-value字典对象映射成pointer指向的struct的对象属性。
// mapping参数不是必须的，它的作用是映射转换的key名
func Struct(params interface{}, pointer interface{}, mapping ...map[string]string) (err error) {
	var keyToAttributeNameMapping map[string]string
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	return doStruct(params, pointer, keyToAttributeNameMapping, "")
}

func StructTag(params interface{}, pointer interface{}, priorityTag string) (err error) {
	return doStruct(params, pointer, nil, priorityTag)
}

func StructDeep(params interface{}, pointer interface{}, mapping ...map[string]string) error {
	var keyToAttributeNameMapping map[string]string
	if len(mapping) > 0 {
		keyToAttributeNameMapping = mapping[0]
	}
	return doStruct(params, pointer, keyToAttributeNameMapping, "")
}

func doStruct(params interface{}, pointer interface{}, mapping map[string]string, priorityTag string) (err error) {
	if params == nil {
		// If `params` is nil, no conversion.
		return nil
	}
	if pointer == nil {
		return derror.New("object pointer cannot be nil")
	}

	defer func() {
		// Catch the panic, especially the reflect operation panics.
		if exception := recover(); exception != nil {
			if e, ok := exception.(errorStack); ok {
				err = e
			} else {
				err = derror.NewSkipf(1, "%v", exception)
			}
		}
	}()

	// If given `params` is JSON, it then uses json.Unmarshal doing the converting.
	switch r := params.(type) {
	case []byte:
		if json.Valid(r) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.Unmarshal(r, rv.Interface())
				}
			} else {
				return json.Unmarshal(r, pointer)
			}
		}
	case string:
		if paramsBytes := []byte(r); json.Valid(paramsBytes) {
			if rv, ok := pointer.(reflect.Value); ok {
				if rv.Kind() == reflect.Ptr {
					return json.Unmarshal(paramsBytes, rv.Interface())
				}
			} else {
				return json.Unmarshal(paramsBytes, pointer)
			}
		}
	}

	var (
		paramsReflectValue      reflect.Value
		pointerReflectValue     reflect.Value
		pointerReflectKind      reflect.Kind
		pointerElemReflectValue reflect.Value // The pointed element.
	)
	if v, ok := params.(reflect.Value); ok {
		paramsReflectValue = v
	} else {
		paramsReflectValue = reflect.ValueOf(params)
	}
	if v, ok := pointer.(reflect.Value); ok {
		pointerReflectValue = v
		pointerElemReflectValue = v
	} else {
		pointerReflectValue = reflect.ValueOf(pointer)
		pointerReflectKind = pointerReflectValue.Kind()
		if pointerReflectKind != reflect.Ptr {
			return derror.Newf("object pointer should be type of '*struct', but got '%v'", pointerReflectKind)
		}
		// Using IsNil on reflect.Ptr variable is OK.
		if !pointerReflectValue.IsValid() || pointerReflectValue.IsNil() {
			return derror.New("object pointer cannot be nil")
		}
		pointerElemReflectValue = pointerReflectValue.Elem()
	}
	// If `params` and `pointer` are the same type, the do directly assignment.
	// For performance enhancement purpose.
	if pointerElemReflectValue.IsValid() && pointerElemReflectValue.Type() == paramsReflectValue.Type() {
		pointerElemReflectValue.Set(paramsReflectValue)
		return nil
	}

	// Normal unmarshalling interfaces checks.
	if err, ok := bindVarToReflectValueWithInterfaceCheck(pointerReflectValue, params); ok {
		return err
	}

	// It automatically creates struct object if necessary.
	// For example, if `pointer` is **User, then `elem` is *User, which is a pointer to User.
	if pointerElemReflectValue.Kind() == reflect.Ptr {
		if !pointerElemReflectValue.IsValid() || pointerElemReflectValue.IsNil() {
			e := reflect.New(pointerElemReflectValue.Type().Elem()).Elem()
			pointerElemReflectValue.Set(e.Addr())
		}
		//if v, ok := pointerElemReflectValue.Interface().(apiUnmarshalValue); ok {
		//	return v.UnmarshalValue(params)
		//}
		// Note that it's `pointerElemReflectValue` here not `pointerReflectValue`.
		if err, ok := bindVarToReflectValueWithInterfaceCheck(pointerElemReflectValue, params); ok {
			return err
		}
		// Retrieve its element, may be struct at last.
		pointerElemReflectValue = pointerElemReflectValue.Elem()
	}

	// paramsMap is the map[string]interface{} type variable for params.
	// DO NOT use MapDeep here.
	paramsMap := Map(params)
	if paramsMap == nil {
		return derror.Newf("convert params to map failed: %v", params)
	}

	// It only performs one converting to the same attribute.
	// doneMap is used to check repeated converting, its key is the real attribute name
	// of the struct.
	doneMap := make(map[string]struct{})

	// The key of the attrMap is the attribute name of the struct,
	// and the value is its replaced name for later comparison to improve performance.
	var (
		tempName       string
		elemFieldType  reflect.StructField
		elemFieldValue reflect.Value
		elemType       = pointerElemReflectValue.Type()
		attrMap        = make(map[string]string)
	)
	for i := 0; i < pointerElemReflectValue.NumField(); i++ {
		elemFieldType = elemType.Field(i)
		// Only do converting to public attributes.
		if !utils.IsLetterUpper(elemFieldType.Name[0]) {
			continue
		}
		// Maybe it's struct/*struct embedded.
		if elemFieldType.Anonymous {
			elemFieldValue = pointerElemReflectValue.Field(i)
			// Ignore the interface attribute if it's nil.
			if elemFieldValue.Kind() == reflect.Interface {
				elemFieldValue = elemFieldValue.Elem()
				if !elemFieldValue.IsValid() {
					continue
				}
			}
			if err = doStruct(paramsMap, elemFieldValue, mapping, priorityTag); err != nil {
				return err
			}
		} else {
			tempName = elemFieldType.Name
			attrMap[tempName] = utils.RemoveSymbols(tempName)
		}
	}
	if len(attrMap) == 0 {
		return nil
	}

	// The key of the tagMap is the attribute name of the struct,
	// and the value is its replaced tag name for later comparison to improve performance.
	var (
		tagMap           = make(map[string]string)
		priorityTagArray []string
	)
	if priorityTag != "" {
		priorityTagArray = append(utils.SplitAndTrim(priorityTag, ","), StructTagPriority...)
	} else {
		priorityTagArray = StructTagPriority
	}
	tagToNameMap, err := structs.TagMapName(pointerElemReflectValue, priorityTagArray)
	if err != nil {
		return err
	}
	for tagName, attributeName := range tagToNameMap {
		// If there's something else in the tag string,
		// it uses the first part which is split using char ','.
		// Eg:
		// orm:"id, priority"
		// orm:"name, with:uid=id"
		tagMap[attributeName] = utils.RemoveSymbols(strings.Split(tagName, ",")[0])
	}

	var (
		attrName  string
		checkName string
	)
	for mapK, mapV := range paramsMap {
		attrName = ""
		// It firstly checks the passed mapping rules.
		if len(mapping) > 0 {
			if passedAttrKey, ok := mapping[mapK]; ok {
				attrName = passedAttrKey
			}
		}
		// It secondly checks the predefined tags and matching rules.
		if attrName == "" {
			checkName = utils.RemoveSymbols(mapK)
			// Loop to find the matched attribute name with or without
			// string cases and chars like '-'/'_'/'.'/' '.

			// Matching the parameters to struct tag names.
			// The `tagV` is the attribute name of the struct.
			for attrKey, cmpKey := range tagMap {
				if strings.EqualFold(checkName, cmpKey) {
					attrName = attrKey
					break
				}
			}
			// Matching the parameters to struct attributes.
			if attrName == "" {
				for attrKey, cmpKey := range attrMap {
					// Eg:
					// UserName  eq user_name
					// User-Name eq username
					// username  eq userName
					// etc.
					if strings.EqualFold(checkName, cmpKey) {
						attrName = attrKey
						break
					}
				}
			}
		}

		// No matching, it gives up this attribute converting.
		if attrName == "" {
			continue
		}
		// If the attribute name is already checked converting, then skip it.
		if _, ok := doneMap[attrName]; ok {
			continue
		}
		// Mark it done.
		doneMap[attrName] = struct{}{}
		if err := bindVarToStructAttr(pointerElemReflectValue, attrName, mapV, mapping, priorityTag); err != nil {
			return err
		}
	}
	return nil
}

// bindVarToReflectValueWithInterfaceCheck does binding using common interfaces checks.
func bindVarToReflectValueWithInterfaceCheck(reflectValue reflect.Value, value interface{}) (err error, ok bool) {
	var pointer interface{}
	if reflectValue.Kind() != reflect.Ptr && reflectValue.CanAddr() {
		reflectValueAddr := reflectValue.Addr()
		if reflectValueAddr.IsNil() || !reflectValueAddr.IsValid() {
			return nil, false
		}
		// Not a pointer, but can token address, that makes it can be unmarshalled.
		pointer = reflectValue.Addr().Interface()
	} else {
		if reflectValue.IsNil() || !reflectValue.IsValid() {
			return nil, false
		}
		pointer = reflectValue.Interface()
	}
	if v, ok := pointer.(apiUnmarshalValue); ok {
		return v.UnmarshalValue(value), ok
	}
	if v, ok := pointer.(apiUnmarshalText); ok {
		if s, ok := value.(string); ok {
			return v.UnmarshalText([]byte(s)), ok
		}
		if b, ok := value.([]byte); ok {
			return v.UnmarshalText(b), ok
		}
	}
	if v, ok := pointer.(apiSet); ok {
		v.Set(value)
		return nil, ok
	}
	return nil, false
}

func bindVarToStructAttr(elem reflect.Value, name string, value interface{}, mapping map[string]string, priorityTag string) (err error) {
	structFieldValue := elem.FieldByName(name)
	if !structFieldValue.IsValid() {
		return nil
	}
	// CanSet checks whether attribute is public accessible.
	if !structFieldValue.CanSet() {
		return nil
	}
	defer func() {
		if e := recover(); e != nil {
			if err = bindVarToReflectValue(structFieldValue, value, mapping, priorityTag); err != nil {
				err = derror.Wrapf(err, `error binding value to attribute "%s"`, name)
			}
		}
	}()
	// Directly converting.
	if empty.IsNil(value) {
		structFieldValue.Set(reflect.Zero(structFieldValue.Type()))
	} else {
		structFieldValue.Set(reflect.ValueOf(Convert(value, structFieldValue.Type().String())))
	}
	return nil
}

// bindVarToReflectValue sets `value` to reflect value object `structFieldValue`.
func bindVarToReflectValue(structFieldValue reflect.Value, value interface{}, mapping map[string]string, priorityTag string) (err error) {
	if err, ok := bindVarToReflectValueWithInterfaceCheck(structFieldValue, value); ok {
		return err
	}
	kind := structFieldValue.Kind()
	// Converting using interface, for some kinds.
	switch kind {
	case reflect.Slice, reflect.Array, reflect.Ptr, reflect.Interface:
		if !structFieldValue.IsNil() {
			if v, ok := structFieldValue.Interface().(apiSet); ok {
				v.Set(value)
				return nil
			}
		}
	}

	// Converting by kind.
	switch kind {
	case reflect.Struct:
		// Recursively converting for struct attribute.
		if err := doStruct(value, structFieldValue, nil, ""); err != nil {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	// Note that the slice element might be type of struct,
	// so it uses Struct function doing the converting internally.
	case reflect.Slice, reflect.Array:
		a := reflect.Value{}
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			a = reflect.MakeSlice(structFieldValue.Type(), v.Len(), v.Len())
			if v.Len() > 0 {
				t := a.Index(0).Type()
				for i := 0; i < v.Len(); i++ {
					if t.Kind() == reflect.Ptr {
						e := reflect.New(t.Elem()).Elem()
						if err := doStruct(v.Index(i).Interface(), e, nil, ""); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e.Addr())
					} else {
						e := reflect.New(t).Elem()
						if err := doStruct(v.Index(i).Interface(), e, nil, ""); err != nil {
							// Note there's reflect conversion mechanism here.
							e.Set(reflect.ValueOf(v.Index(i).Interface()).Convert(t))
						}
						a.Index(i).Set(e)
					}
				}
			}
		} else {
			a = reflect.MakeSlice(structFieldValue.Type(), 1, 1)
			t := a.Index(0).Type()
			if t.Kind() == reflect.Ptr {
				e := reflect.New(t.Elem()).Elem()
				if err := doStruct(value, e, nil, ""); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e.Addr())
			} else {
				e := reflect.New(t).Elem()
				if err := doStruct(value, e, nil, ""); err != nil {
					// Note there's reflect conversion mechanism here.
					e.Set(reflect.ValueOf(value).Convert(t))
				}
				a.Index(0).Set(e)
			}
		}
		structFieldValue.Set(a)

	case reflect.Ptr:
		item := reflect.New(structFieldValue.Type().Elem())
		if err, ok := bindVarToReflectValueWithInterfaceCheck(item, value); ok {
			structFieldValue.Set(item)
			return err
		}
		elem := item.Elem()
		if err = bindVarToReflectValue(elem, value, mapping, priorityTag); err == nil {
			structFieldValue.Set(elem.Addr())
		}

	// It mainly and specially handles the interface of nil value.
	case reflect.Interface:
		if value == nil {
			// Specially.
			structFieldValue.Set(reflect.ValueOf((*interface{})(nil)))
		} else {
			// Note there's reflect conversion mechanism here.
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		}

	default:
		defer func() {
			if e := recover(); e != nil {
				err = derror.New(
					fmt.Sprintf(`cannot convert value "%+v" to type "%s"`,
						value,
						structFieldValue.Type().String(),
					),
				)
			}
		}()
		// It here uses reflect converting `value` to type of the attribute and assigns
		// the result value to the attribute. It might fail and panic if the usual Go
		// conversion rules do not allow conversion.
		structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
	}
	return nil
}
