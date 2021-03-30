package dmap

type (
	Map     = AnyAnyMap // Map is alias of AnyAnyMap.
	HashMap = AnyAnyMap // HashMap is alias of AnyAnyMap.
)

//创建一个map对象
func New(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

//创建一个map对象，并把参数data作为值
func NewFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}

//创建一个map对象
func NewHashMap(safe ...bool) *Map {
	return NewAnyAnyMap(safe...)
}

//创建一个map对象，并把参数data作为值
func NewHashMapFrom(data map[interface{}]interface{}, safe ...bool) *Map {
	return NewAnyAnyMapFrom(data, safe...)
}
