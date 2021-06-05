package dtime

import "database/sql/driver"

// Scan 从 database/sql包中扫描时间格式的数据到golang格式
func (that *Time) Scan(value interface{}) error {
	if that == nil {
		return nil
	}
	newTime := New(value)
	*that = *newTime
	return nil
}

// Value 把当前时间格式转换成database/sql/driver可以识别的数据格式
func (that *Time) Value() (driver.Value, error) {
	if that == nil {
		return nil, nil
	}
	if that.IsZero() {
		return nil, nil
	}
	return that.Time, nil
}
