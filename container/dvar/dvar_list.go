package dvar

import (
	"github.com/osgochina/donkeygo/util/dutil"
)

// ListItemValues 检索并返回键值为<键值>的所有item struct/map的元素。
//注意，参数应该是包含map或struct元素的slice类型，否则返回一个空的slice。
func (that *Var) ListItemValues(key interface{}) (values []interface{}) {
	return dutil.ListItemValues(that.Val(), key)
}

// ListItemValuesUnique 检索并返回所有struct/map中key 的唯一元素。
//注意参数应该是包含map或struct元素的slice类型，否则返回一个空的slice。
func (that *Var) ListItemValuesUnique(key string) []interface{} {
	return dutil.ListItemValuesUnique(that.Val(), key)
}
