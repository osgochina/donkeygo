package dutil

import (
	"github.com/osgochina/donkeygo/util/dconv"
	"strings"
)

// Comparator 对比方法
// 返回的数字表示的含义:
//    负数 , if a < b
//    0     , if a == b
//    正数 , if a > b
type Comparator func(a, b interface{}) int

func ComparatorString(a, b interface{}) int {
	return strings.Compare(dconv.String(a), dconv.String(b))
}

// ComparatorInt provides a basic comparison on int.
func ComparatorInt(a, b interface{}) int {
	return dconv.Int(a) - dconv.Int(b)
}

// ComparatorInt8 provides a basic comparison on int8.
func ComparatorInt8(a, b interface{}) int {
	return int(dconv.Int8(a) - dconv.Int8(b))
}

// ComparatorInt16 provides a basic comparison on int16.
func ComparatorInt16(a, b interface{}) int {
	return int(dconv.Int16(a) - dconv.Int16(b))
}

// ComparatorInt32 provides a basic comparison on int32.
func ComparatorInt32(a, b interface{}) int {
	return int(dconv.Int32(a) - dconv.Int32(b))
}

// ComparatorInt64 provides a basic comparison on int64.
func ComparatorInt64(a, b interface{}) int {
	return int(dconv.Int64(a) - dconv.Int64(b))
}

// ComparatorUint provides a basic comparison on uint.
func ComparatorUint(a, b interface{}) int {
	return int(dconv.Uint(a) - dconv.Uint(b))
}

// ComparatorUint8 provides a basic comparison on uint8.
func ComparatorUint8(a, b interface{}) int {
	return int(dconv.Uint8(a) - dconv.Uint8(b))
}

// ComparatorUint16 provides a basic comparison on uint16.
func ComparatorUint16(a, b interface{}) int {
	return int(dconv.Uint16(a) - dconv.Uint16(b))
}

// ComparatorUint32 provides a basic comparison on uint32.
func ComparatorUint32(a, b interface{}) int {
	return int(dconv.Uint32(a) - dconv.Uint32(b))
}

// ComparatorUint64 provides a basic comparison on uint64.
func ComparatorUint64(a, b interface{}) int {
	return int(dconv.Uint64(a) - dconv.Uint64(b))
}

// ComparatorFloat32 provides a basic comparison on float32.
func ComparatorFloat32(a, b interface{}) int {
	aFloat := dconv.Float32(a)
	bFloat := dconv.Float32(b)
	if aFloat == bFloat {
		return 0
	}
	if aFloat > bFloat {
		return 1
	}
	return -1
}

// ComparatorFloat64 provides a basic comparison on float64.
func ComparatorFloat64(a, b interface{}) int {
	aFloat := dconv.Float64(a)
	bFloat := dconv.Float64(b)
	if aFloat == bFloat {
		return 0
	}
	if aFloat > bFloat {
		return 1
	}
	return -1
}

// ComparatorByte provides a basic comparison on byte.
func ComparatorByte(a, b interface{}) int {
	return int(dconv.Byte(a) - dconv.Byte(b))
}

// ComparatorRune provides a basic comparison on rune.
func ComparatorRune(a, b interface{}) int {
	return int(dconv.Rune(a) - dconv.Rune(b))
}

// ComparatorTime provides a basic comparison on time.Time.
func ComparatorTime(a, b interface{}) int {
	aTime := dconv.Time(a)
	bTime := dconv.Time(b)
	switch {
	case aTime.After(bTime):
		return 1
	case aTime.Before(bTime):
		return -1
	default:
		return 0
	}
}
