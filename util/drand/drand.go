package drand

import (
	"encoding/binary"
	"time"
	"unsafe"
)

var (
	letters    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 52
	symbols    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"                   // 32
	digits     = "0123456789"                                           // 10
	characters = letters + digits + symbols                             // 94
)

// Intn 返回一个介于0-max之间的随机数
//1. ' max '只能大于0，否则直接返回' max ';
// 2。结果大于等于0，但小于' max ';
// 3。结果数为32位，小于math.MaxUint32。
func Intn(max int) int {
	if max <= 0 {
		return max
	}
	n := int(binary.LittleEndian.Uint32(<-bufferChan)) % max
	if (max > 0 && n < 0) || (max < 0 && n > 0) {
		return -n
	}
	return n
}

// B 检索并返回给定长度' n '的随机字节。
func B(n int) []byte {
	if n <= 0 {
		return nil
	}
	i := 0
	b := make([]byte, n)
	for {
		copy(b[i:], <-bufferChan)
		i += 4
		if i >= n {
			break
		}
	}
	return b
}

// N 返回在一个[min,max]之间的随机数，支持负数
func N(min, max int) int {
	if min >= max {
		return min
	}
	if min >= 0 {
		// 因为 Intn 不支持负数，
		// 我们应该先把值移到左边，
		// 然后调用 Intn 生成随机数，
		// 最后将结果移回右。
		return Intn(max-(min-0)+1) + (min - 0)
	}
	if min < 0 {
		// 因为 Intn 不支持负数，
		// 我们应该先把值移到左边，
		// 然后调用 Intn 生成随机数，
		// 最后将结果移回左边。
		return Intn(max+(0-min)+1) - (0 - min)
	}
	return 0
}

// S 返回长度为N的随机字符串 symbols 传入表示支持符号
func S(n int, symbols ...bool) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		if len(symbols) > 0 && symbols[0] {
			b[i] = characters[numberBytes[i]%94]
		} else {
			b[i] = characters[numberBytes[i]%62]
		}
	}
	return *(*string)(unsafe.Pointer(&b))
}

// D 返回随机时间
func D(min, max time.Duration) time.Duration {
	multiple := int64(1)
	if min != 0 {
		for min%10 == 0 {
			multiple *= 10
			min /= 10
			max /= 10
		}
	}
	n := int64(N(int(min), int(max)))
	return time.Duration(n * multiple)
}

// Str 返回给定字符串中的随机位数
func Str(s string, n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b     = make([]rune, n)
		runes = []rune(s)
	)
	if len(runes) <= 255 {
		numberBytes := B(n)
		for i := range b {
			b[i] = runes[int(numberBytes[i])%len(runes)]
		}
	} else {
		for i := range b {
			b[i] = runes[Intn(len(runes))]
		}
	}
	return string(b)
}

// Digits 返回随机数字
func Digits(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = digits[numberBytes[i]%10]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Letters 返回指定长度的随机字母
func Letters(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = letters[numberBytes[i]%52]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Symbols 返回长度为n的随机符号
func Symbols(n int) string {
	if n <= 0 {
		return ""
	}
	var (
		b           = make([]byte, n)
		numberBytes = B(n)
	)
	for i := range b {
		b[i] = symbols[numberBytes[i]%32]
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Perm 返回n数目的随机int切片
func Perm(n int) []int {
	m := make([]int, n)
	for i := 0; i < n; i++ {
		j := Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}

// Meet 随机对比概率
func Meet(num, total int) bool {
	return Intn(total) < num
}

// MeetProb 随机计算是否满足给定的概率。
func MeetProb(prob float32) bool {
	return Intn(1e7) < int(prob*1e7)
}
