package dbinary

import (
	"encoding/binary"
	"math"
)

func LeDecodeToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

func LeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// 当b位数不够时，进行高位补0。
// 注意这里为了不影响原有输入参数，是采用的值复制设计。
func LeFillUpSize(b []byte, l int) []byte {
	if len(b) > l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
