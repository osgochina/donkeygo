package dstr

import (
	"bytes"
	"fmt"
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/util/dconv"
	"github.com/osgochina/donkeygo/util/drand"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// NotFoundIndex 字符串搜索的时候，没有找到字符串的位置，则返回-1
	NotFoundIndex = -1
)

// Replace 替换内容
// origin 原始内容
// search 要查找的内容
// replace 要替换的内容
// count 替换次数
func Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}

// ReplaceI 不区分大小写替换字符串
func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return origin
	}
	var (
		length      = len(search)
		searchLower = strings.ToLower(search)
	)
	for {
		originLower := strings.ToLower(origin)
		if pos := strings.Index(originLower, searchLower); pos != -1 {
			origin = origin[:pos] + replace + origin[pos+length:]
			if n--; n == 0 {
				break
			}
		} else {
			break
		}
	}
	return origin
}

// Count 统计substr在字符串s中出现的次数
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// CountI 不区分大小写统计substr在字符串s中出现的次数
func CountI(s, substr string) int {
	return strings.Count(ToLower(s), ToLower(substr))
}

// ReplaceByArray 使用数组提供要替换的字符串，替换原始字符串中的内容
func ReplaceByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = Replace(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceIByArray 不区分大小写使用数组提供要替换的字符串，替换原始字符串中的内容
func ReplaceIByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = ReplaceI(origin, array[i], array[i+1])
	}
	return origin
}

// ReplaceByMap 使用replaces的key去origin中查找对应的串，替换成key对应的value
func ReplaceByMap(origin string, replaces map[string]string) string {
	return utils.ReplaceByMap(origin, replaces)
}

// ReplaceIByMap 不区分大小写使用replaces的key去origin中查找对应的串，替换成key对应的value
func ReplaceIByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = ReplaceI(origin, k, v)
	}
	return origin
}

// ToLower 转换成大写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper 转换成小写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// UcFirst 首字母大写
func UcFirst(s string) string {
	return utils.UcFirst(s)
}

// LcFirst 首字母小写
func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if IsLetterUpper(s[0]) {
		return string(s[0]+32) + s[1:]
	}
	return s
}

// UcWords 字符串中的每个单词的首字母大写
func UcWords(str string) string {
	return strings.Title(str)
}

// IsLetterLower 是否是小写字符
func IsLetterLower(b byte) bool {
	return utils.IsLetterLower(b)
}

// IsLetterUpper 是否是大写字符
func IsLetterUpper(b byte) bool {
	return utils.IsLetterUpper(b)
}

// IsNumeric 字符串是否是数字
func IsNumeric(s string) bool {
	return utils.IsNumeric(s)
}

// SubStr 截取字符串，类似于php的substr函数
func SubStr(str string, start int, length ...int) (substr string) {
	lth := len(str)

	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= lth {
		start = lth
	}
	end := lth
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = lth
		}
	}
	if end > lth {
		end = lth
	}
	return str[start:end]
}

// SubStrRune 截取unicode 字符串
func SubStrRune(str string, start int, length ...int) (substr string) {
	// Converting to []rune to support unicode.
	rs := []rune(str)
	lth := len(rs)

	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= lth {
		start = lth
	}
	end := lth
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = lth
		}
	}
	if end > lth {
		end = lth
	}
	return string(rs[start:end])
}

// StrLimit 截取字符串指定长度，使用suffix填充
func StrLimit(str string, length int, suffix ...string) string {
	if len(str) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return str[0:length] + addStr
}

// StrLimitRune 截取unicode 字符串指定长度，使用suffix填充
func StrLimitRune(str string, length int, suffix ...string) string {
	rs := []rune(str)
	if len(rs) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return string(rs[0:length]) + addStr
}

// Reverse 反转字符串
func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// NumberFormat 以千位分隔符方式格式化一个数字
// <number>: 你要格式化的数字
// <decimals>: 要保留的小数位数
// <decPoint>: 指定小数点显示的字符
// <thousandsSep>: 指定千位分隔符显示的字符
// See http://php.net/manual/en/function.number-format.php.
func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(decimals)+"F", number)
	prefix, suffix := "", ""
	if decimals > 0 {
		prefix = str[:len(str)-(decimals+1)]
		suffix = str[len(str)-decimals:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if decimals > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}

// ChunkSplit 将字符串分割成更小的块
func ChunkSplit(body string, chunkLen int, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, endRunes := []rune(body), []rune(end)
	l := len(runes)
	if l <= 1 || l < chunkLen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(endRunes))
	for i := 0; i < l; i += chunkLen {
		if i+chunkLen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunkLen]...)
		}
		ns = append(ns, endRunes...)
	}
	return string(ns)
}

// Compare 返回一个整数，按字典顺序比较两个字符串。
// 如果返回 0 表示 a==b, -1 表示 a < b,  +1 表示 a > b.
func Compare(a, b string) int {
	return strings.Compare(a, b)
}

// Equal 使用 UTF-8 格式比较字符串是否相等
func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

// Fields 将字符串使用的单词作为切片数组返回
func Fields(str string) []string {
	return strings.Fields(str)
}

// HasPrefix 测试字符串s是否以prefix开头。
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix 测试字符串s是否以suffix结尾
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// CountWords 计算单词出现的次数，并以数组返回
func CountWords(str string) map[string]int {
	m := make(map[string]int)
	buffer := bytes.NewBuffer(nil)
	for _, r := range []rune(str) {
		if unicode.IsSpace(r) {
			if buffer.Len() > 0 {
				m[buffer.String()]++
				buffer.Reset()
			}
		} else {
			buffer.WriteRune(r)
		}
	}
	if buffer.Len() > 0 {
		m[buffer.String()]++
	}
	return m
}

// CountChars 计算字符出现的次数，并以数组返回
func CountChars(str string, noSpace ...bool) map[string]int {
	m := make(map[string]int)
	countSpace := true
	if len(noSpace) > 0 && noSpace[0] {
		countSpace = false
	}
	for _, r := range []rune(str) {
		if !countSpace && unicode.IsSpace(r) {
			continue
		}
		m[string(r)]++
	}
	return m
}

// WordWrap 打断字符串为指定数量的字串
// TODO: Enable cut parameter, see http://php.net/manual/en/function.wordwrap.php.
func WordWrap(str string, width int, br string) string {
	if br == "" {
		br = "\n"
	}
	var (
		current           int
		wordBuf, spaceBuf bytes.Buffer
		init              = make([]byte, 0, len(str))
		buf               = bytes.NewBuffer(init)
	)
	for _, char := range []rune(str) {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBuf.Len() > width {
					current = 0
				} else {
					current += spaceBuf.Len()
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += spaceBuf.Len() + wordBuf.Len()
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBuf.Len() + wordBuf.Len()
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			spaceBuf.WriteRune(char)
		} else {
			wordBuf.WriteRune(char)
			if current+spaceBuf.Len()+wordBuf.Len() > width && wordBuf.Len() < width {
				buf.WriteString(br)
				current = 0
				spaceBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBuf.Len() <= width {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}
	return buf.String()
}

// RuneLen 返回unicode的字符串长度。
func RuneLen(str string) int {
	return LenRune(str)
}

// LenRune 返回unicode的字符串长度。
func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}

// Repeat 返回一个新字符串，该字符串由输入字符串的乘数副本组成。
func Repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

// Str 查找字符串的首次出现的位置
// the first occurrence of <needle> to the end of <haystack>.
// See http://php.net/manual/en/function.strstr.php.
func Str(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	pos := strings.Index(haystack, needle)
	if pos == NotFoundIndex {
		return ""
	}
	return haystack[pos+len([]byte(needle))-1:]
}

// StrEx 返回字符串的一部分，从第一次出现到结束。
func StrEx(haystack string, needle string) string {
	if s := Str(haystack, needle); s != "" {
		return s[1:]
	}
	return ""
}

// StrTill returns part of <haystack> string ending to and including
// the first occurrence of <needle> from the start of <haystack>.
func StrTill(haystack string, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos == NotFoundIndex || pos == 0 {
		return ""
	}
	return haystack[:pos+1]
}

// StrTillEx returns part of <haystack> string ending to and excluding
// the first occurrence of <needle> from the start of <haystack>.
func StrTillEx(haystack string, needle string) string {
	pos := strings.Index(haystack, needle)
	if pos == NotFoundIndex || pos == 0 {
		return ""
	}
	return haystack[:pos]
}

// Shuffle randomly shuffles a string.
// It considers parameter <str> as unicode string.
func Shuffle(str string) string {
	runes := []rune(str)
	s := make([]rune, len(runes))
	for i, v := range drand.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

// Split 分割字符串
func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

// SplitAndTrim 分割字符串并去除分割后子串的头尾指定字符
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// SplitAndTrimSpace 分割字符串并去除分割后子串的头尾不可见字符
func SplitAndTrimSpace(str, delimiter string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = strings.TrimSpace(v)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// Join 连接字符串
func Join(array []string, sep string) string {
	return strings.Join(array, sep)
}

// JoinAny 连接任何东西
func JoinAny(array interface{}, sep string) string {
	return strings.Join(dconv.Strings(array), sep)
}

// Explode 分割字符串
func Explode(delimiter, str string) []string {
	return Split(str, delimiter)
}

// Implode 连接字符
func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

// Chr return the ascii string of a number(0-255).
func Chr(ascii int) string {
	return string([]byte{byte(ascii % 256)})
}

// Ord converts the first byte of a string to a value between 0 and 255.
func Ord(char string) int {
	return int(char[0])
}

// HideStr replaces part of the the string <str> to <hide> by <percentage> from the <middle>.
// It considers parameter <str> as unicode string.
func HideStr(str string, percent int, hide string) string {
	array := strings.Split(str, "@")
	if len(array) > 1 {
		str = array[0]
	}
	var (
		rs       = []rune(str)
		length   = len(rs)
		mid      = math.Floor(float64(length / 2))
		hideLen  = int(math.Floor(float64(length) * (float64(percent) / 100)))
		start    = int(mid - math.Floor(float64(hideLen)/2))
		hideStr  = []rune("")
		hideRune = []rune(hide)
	)
	for i := 0; i < hideLen; i++ {
		hideStr = append(hideStr, hideRune...)
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(string(rs[0:start]))
	buffer.WriteString(string(hideStr))
	buffer.WriteString(string(rs[start+hideLen:]))
	if len(array) > 1 {
		buffer.WriteString("@" + array[1])
	}
	return buffer.String()
}

// Nl2Br inserts HTML line breaks(<br>|<br />) before all newlines in a string:
// \n\r, \r\n, \r, \n.
// It considers parameter <str> as unicode string.
func Nl2Br(str string, isXhtml ...bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if len(isXhtml) > 0 && isXhtml[0] {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

// AddSlashes 给引号加反斜杠
func AddSlashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// StripSlashes 去除反斜杠
func StripSlashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// QuoteMeta 返回一个带有反斜杠字符(\)的str版本
//在每个字符之前: .\+*?[^]($)
func QuoteMeta(str string, chars ...string) string {
	var buf bytes.Buffer
	for _, char := range str {
		if len(chars) > 0 {
			for _, c := range chars[0] {
				if c == char {
					buf.WriteRune('\\')
					break
				}
			}
		} else {
			switch char {
			case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
				buf.WriteRune('\\')
			}
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

// SearchArray 在字符串数组中查找字符串s
func SearchArray(a []string, s string) int {
	for i, v := range a {
		if s == v {
			return i
		}
	}
	return NotFoundIndex
}

// InArray 判断字符串s是否在字符串数组a中
func InArray(a []string, s string) bool {
	return SearchArray(a, s) != NotFoundIndex
}

//SnakeString converts the accepted string to a snake string (XxYy to xx_yy)
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	for _, d := range dconv.Bytes(s) {
		if d >= 'A' && d <= 'Z' {
			if j {
				data = append(data, '_')
				j = false
			}
		} else if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(dconv.String(data))
}
