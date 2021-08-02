// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package dstr_test

import (
	"github.com/osgochina/donkeygo/frame/d"
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
)

func Test_Replace(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(dstr.Replace(s1, "ab", "AB"), "ABcdEFG乱入的中文ABcdefg")
		t.Assert(dstr.Replace(s1, "EF", "ef"), "abcdefG乱入的中文abcdefg")
		t.Assert(dstr.Replace(s1, "MN", "mn"), s1)

		t.Assert(dstr.ReplaceByArray(s1, d.ArrayStr{
			"a", "A",
			"A", "-",
			"a",
		}), "-bcdEFG乱入的中文-bcdefg")

		t.Assert(dstr.ReplaceByMap(s1, d.MapStrStr{
			"a": "A",
			"G": "g",
		}), "AbcdEFg乱入的中文Abcdefg")
	})
}

func Test_ReplaceI_1(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcd乱入的中文ABCD"
		s2 := "a"
		t.Assert(dstr.ReplaceI(s1, "ab", "aa"), "aacd乱入的中文aaCD")
		t.Assert(dstr.ReplaceI(s1, "ab", "aa", 0), "abcd乱入的中文ABCD")
		t.Assert(dstr.ReplaceI(s1, "ab", "aa", 1), "aacd乱入的中文ABCD")

		t.Assert(dstr.ReplaceI(s1, "abcd", "-"), "-乱入的中文-")
		t.Assert(dstr.ReplaceI(s1, "abcd", "-", 1), "-乱入的中文ABCD")

		t.Assert(dstr.ReplaceI(s1, "abcd乱入的", ""), "中文ABCD")
		t.Assert(dstr.ReplaceI(s1, "ABCD乱入的", ""), "中文ABCD")

		t.Assert(dstr.ReplaceI(s2, "A", "-"), "-")
		t.Assert(dstr.ReplaceI(s2, "a", "-"), "-")

		t.Assert(dstr.ReplaceIByArray(s1, d.ArrayStr{
			"abcd乱入的", "-",
			"-", "=",
			"a",
		}), "=中文ABCD")

		t.Assert(dstr.ReplaceIByMap(s1, d.MapStrStr{
			"ab": "-",
			"CD": "=",
		}), "-=乱入的中文-=")
	})
}

func Test_ToLower(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "abcdefg乱入的中文abcdefg"
		t.Assert(dstr.ToLower(s1), e1)
	})
}

func Test_ToUpper(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "ABCDEFG乱入的中文ABCDEFG"
		t.Assert(dstr.ToUpper(s1), e1)
	})
}

func Test_UcFirst(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "AbcdEFG乱入的中文abcdefg"
		t.Assert(dstr.UcFirst(""), "")
		t.Assert(dstr.UcFirst(s1), e1)
		t.Assert(dstr.UcFirst(e1), e1)
	})
}

func Test_LcFirst(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "AbcdEFG乱入的中文abcdefg"
		e1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(dstr.LcFirst(""), "")
		t.Assert(dstr.LcFirst(s1), e1)
		t.Assert(dstr.LcFirst(e1), e1)
	})
}

func Test_UcWords(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱GF: i love go frame"
		e1 := "我爱GF: I Love Go Frame"
		t.Assert(dstr.UcWords(s1), e1)
	})
}

func Test_IsLetterLower(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.IsLetterLower('a'), true)
		t.Assert(dstr.IsLetterLower('A'), false)
		t.Assert(dstr.IsLetterLower('1'), false)
	})
}

func Test_IsLetterUpper(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.IsLetterUpper('a'), false)
		t.Assert(dstr.IsLetterUpper('A'), true)
		t.Assert(dstr.IsLetterUpper('1'), false)
	})
}

func Test_IsNumeric(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.IsNumeric("1a我"), false)
		t.Assert(dstr.IsNumeric("0123"), true)
		t.Assert(dstr.IsNumeric("我是中国人"), false)
	})
}

func Test_SubStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.SubStr("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(dstr.SubStr("我爱GoFrame", 6), "GoFrame")
		t.Assert(dstr.SubStr("我爱GoFrame", 6, 2), "Go")
		t.Assert(dstr.SubStr("我爱GoFrame", -1, 30), "我爱GoFrame")
		t.Assert(dstr.SubStr("我爱GoFrame", 30, 30), "")
	})
}

func Test_SubStrRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.SubStrRune("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(dstr.SubStrRune("我爱GoFrame", 2), "GoFrame")
		t.Assert(dstr.SubStrRune("我爱GoFrame", 2, 2), "Go")
		t.Assert(dstr.SubStrRune("我爱GoFrame", -1, 30), "我爱GoFrame")
		t.Assert(dstr.SubStrRune("我爱GoFrame", 30, 30), "")
	})
}

func Test_StrLimit(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StrLimit("我爱GoFrame", 6), "我爱...")
		t.Assert(dstr.StrLimit("我爱GoFrame", 6, ""), "我爱")
		t.Assert(dstr.StrLimit("我爱GoFrame", 6, "**"), "我爱**")
		t.Assert(dstr.StrLimit("我爱GoFrame", 8, ""), "我爱Go")
		t.Assert(dstr.StrLimit("*", 4, ""), "*")
	})
}

func Test_StrLimitRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StrLimitRune("我爱GoFrame", 2), "我爱...")
		t.Assert(dstr.StrLimitRune("我爱GoFrame", 2, ""), "我爱")
		t.Assert(dstr.StrLimitRune("我爱GoFrame", 2, "**"), "我爱**")
		t.Assert(dstr.StrLimitRune("我爱GoFrame", 4, ""), "我爱Go")
		t.Assert(dstr.StrLimitRune("*", 4, ""), "*")
	})
}

func Test_HasPrefix(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.HasPrefix("我爱GoFrame", "我爱"), true)
		t.Assert(dstr.HasPrefix("en我爱GoFrame", "我爱"), false)
		t.Assert(dstr.HasPrefix("en我爱GoFrame", "en"), true)
	})
}

func Test_HasSuffix(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.HasSuffix("我爱GoFrame", "GoFrame"), true)
		t.Assert(dstr.HasSuffix("en我爱GoFrame", "a"), false)
		t.Assert(dstr.HasSuffix("GoFrame很棒", "棒"), true)
	})
}

func Test_Reverse(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Reverse("我爱123"), "321爱我")
	})
}

func Test_NumberFormat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.NumberFormat(1234567.8910, 2, ".", ","), "1,234,567.89")
		t.Assert(dstr.NumberFormat(1234567.8910, 2, "#", "/"), "1/234/567#89")
		t.Assert(dstr.NumberFormat(-1234567.8910, 2, "#", "/"), "-1/234/567#89")
	})
}

func Test_ChunkSplit(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.ChunkSplit("1234", 1, "#"), "1#2#3#4#")
		t.Assert(dstr.ChunkSplit("我爱123", 1, "#"), "我#爱#1#2#3#")
		t.Assert(dstr.ChunkSplit("1234", 1, ""), "1\r\n2\r\n3\r\n4\r\n")
	})
}

func Test_SplitAndTrim(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := `

010    

020  

`
		a := dstr.SplitAndTrim(s, "\n", "0")
		t.Assert(len(a), 2)
		t.Assert(a[0], "1")
		t.Assert(a[1], "2")
	})
}

func Test_Fields(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Fields("我爱 Go Frame"), []string{
			"我爱", "Go", "Frame",
		})
	})
}

func Test_CountWords(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.CountWords("我爱 Go Go Go"), map[string]int{
			"Go": 3,
			"我爱": 1,
		})
	})
}

func Test_CountChars(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.CountChars("我爱 Go Go Go"), map[string]int{
			" ": 3,
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
		t.Assert(dstr.CountChars("我爱 Go Go Go", true), map[string]int{
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
	})
}

func Test_WordWrap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.WordWrap("12 34", 2, "<br>"), "12<br>34")
		t.Assert(dstr.WordWrap("12 34", 2, "\n"), "12\n34")
		t.Assert(dstr.WordWrap("我爱 GF", 2, "\n"), "我爱\nGF")
		t.Assert(dstr.WordWrap("A very long woooooooooooooooooord. and something", 7, "<br>"),
			"A very<br>long<br>woooooooooooooooooord.<br>and<br>something")
	})
}

func Test_RuneLen(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.RuneLen("1234"), 4)
		t.Assert(dstr.RuneLen("我爱GoFrame"), 9)
	})
}

func Test_Repeat(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Repeat("go", 3), "gogogo")
		t.Assert(dstr.Repeat("好的", 3), "好的好的好的")
	})
}

func Test_Str(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Str("name@example.com", "@"), "@example.com")
		t.Assert(dstr.Str("name@example.com", ""), "")
		t.Assert(dstr.Str("name@example.com", "z"), "")
	})
}

func Test_StrEx(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StrEx("name@example.com", "@"), "example.com")
		t.Assert(dstr.StrEx("name@example.com", ""), "")
		t.Assert(dstr.StrEx("name@example.com", "z"), "")
	})
}

func Test_StrTill(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StrTill("name@example.com", "@"), "name@")
		t.Assert(dstr.StrTill("name@example.com", ""), "")
		t.Assert(dstr.StrTill("name@example.com", "z"), "")
	})
}

func Test_StrTillEx(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StrTillEx("name@example.com", "@"), "name")
		t.Assert(dstr.StrTillEx("name@example.com", ""), "")
		t.Assert(dstr.StrTillEx("name@example.com", "z"), "")
	})
}

func Test_Shuffle(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(len(dstr.Shuffle("123456")), 6)
	})
}

func Test_Split(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Split("1.2", "."), []string{"1", "2"})
		t.Assert(dstr.Split("我爱 - GoFrame", " - "), []string{"我爱", "GoFrame"})
	})
}

func Test_Join(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Join([]string{"我爱", "GoFrame"}, " - "), "我爱 - GoFrame")
	})
}

func Test_Explode(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Explode(" - ", "我爱 - GoFrame"), []string{"我爱", "GoFrame"})
	})
}

func Test_Implode(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Implode(" - ", []string{"我爱", "GoFrame"}), "我爱 - GoFrame")
	})
}

func Test_Chr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Chr(65), "A")
	})
}

func Test_Ord(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Ord("A"), 65)
	})
}

func Test_HideStr(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.HideStr("15928008611", 40, "*"), "159****8611")
		t.Assert(dstr.HideStr("john@kohd.cn", 40, "*"), "jo*n@kohd.cn")
		t.Assert(dstr.HideStr("张三", 50, "*"), "张*")
		t.Assert(dstr.HideStr("张小三", 50, "*"), "张*三")
		t.Assert(dstr.HideStr("欧阳小三", 50, "*"), "欧**三")
	})
}

func Test_Nl2Br(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Nl2Br("1\n2"), "1<br>2")
		t.Assert(dstr.Nl2Br("1\r\n2"), "1<br>2")
		t.Assert(dstr.Nl2Br("1\r\n2", true), "1<br />2")
	})
}

func Test_AddSlashes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.AddSlashes(`1'2"3\`), `1\'2\"3\\`)
	})
}

func Test_StripSlashes(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.StripSlashes(`1\'2\"3\\`), `1'2"3\`)
	})
}

func Test_QuoteMeta(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.QuoteMeta(`.\+*?[^]($)`), `\.\\\+\*\?\[\^\]\(\$\)`)
		t.Assert(dstr.QuoteMeta(`.\+*中国?[^]($)`), `\.\\\+\*中国\?\[\^\]\(\$\)`)
		t.Assert(dstr.QuoteMeta(`.''`, `'`), `.\'\'`)
		t.Assert(dstr.QuoteMeta(`中国.''`, `'`), `中国.\'\'`)
	})
}

func Test_Count(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := "abcdaAD"
		t.Assert(dstr.Count(s, "0"), 0)
		t.Assert(dstr.Count(s, "a"), 2)
		t.Assert(dstr.Count(s, "b"), 1)
		t.Assert(dstr.Count(s, "d"), 1)
	})
}

func Test_CountI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := "abcdaAD"
		t.Assert(dstr.CountI(s, "0"), 0)
		t.Assert(dstr.CountI(s, "a"), 3)
		t.Assert(dstr.CountI(s, "b"), 1)
		t.Assert(dstr.CountI(s, "d"), 2)
	})
}

func Test_Compare(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Compare("a", "b"), -1)
		t.Assert(dstr.Compare("a", "a"), 0)
		t.Assert(dstr.Compare("b", "a"), 1)
	})
}

func Test_Equal(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Equal("a", "A"), true)
		t.Assert(dstr.Equal("a", "a"), true)
		t.Assert(dstr.Equal("b", "a"), false)
	})
}

func Test_Contains(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.Contains("abc", "a"), true)
		t.Assert(dstr.Contains("abc", "A"), false)
		t.Assert(dstr.Contains("abc", "ab"), true)
		t.Assert(dstr.Contains("abc", "abc"), true)
	})
}

func Test_ContainsI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.ContainsI("abc", "a"), true)
		t.Assert(dstr.ContainsI("abc", "A"), true)
		t.Assert(dstr.ContainsI("abc", "Ab"), true)
		t.Assert(dstr.ContainsI("abc", "ABC"), true)
		t.Assert(dstr.ContainsI("abc", "ABCD"), false)
		t.Assert(dstr.ContainsI("abc", "D"), false)
	})
}

func Test_ContainsAny(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(dstr.ContainsAny("abc", "a"), true)
		t.Assert(dstr.ContainsAny("abc", "cd"), true)
		t.Assert(dstr.ContainsAny("abc", "de"), false)
		t.Assert(dstr.ContainsAny("abc", "A"), false)
	})
}

func Test_SearchArray(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a := d.SliceStr{"a", "b", "c"}
		t.AssertEQ(dstr.SearchArray(a, "a"), 0)
		t.AssertEQ(dstr.SearchArray(a, "b"), 1)
		t.AssertEQ(dstr.SearchArray(a, "c"), 2)
		t.AssertEQ(dstr.SearchArray(a, "d"), -1)
	})
}

func Test_InArray(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		a := d.SliceStr{"a", "b", "c"}
		t.AssertEQ(dstr.InArray(a, "a"), true)
		t.AssertEQ(dstr.InArray(a, "b"), true)
		t.AssertEQ(dstr.InArray(a, "c"), true)
		t.AssertEQ(dstr.InArray(a, "d"), false)
	})
}
