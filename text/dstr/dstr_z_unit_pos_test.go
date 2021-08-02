// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package dstr_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"testing"
)

func Test_Pos(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(dstr.Pos(s1, "ab"), 0)
		t.Assert(dstr.Pos(s1, "ab", 2), 7)
		t.Assert(dstr.Pos(s1, "abd", 0), -1)
		t.Assert(dstr.Pos(s1, "e", -4), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.Pos(s1, "爱"), 3)
		t.Assert(dstr.Pos(s1, "C"), 6)
		t.Assert(dstr.Pos(s1, "China"), 6)
	})
}

func Test_PosRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(dstr.PosRune(s1, "ab"), 0)
		t.Assert(dstr.PosRune(s1, "ab", 2), 7)
		t.Assert(dstr.PosRune(s1, "abd", 0), -1)
		t.Assert(dstr.PosRune(s1, "e", -4), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosRune(s1, "爱"), 1)
		t.Assert(dstr.PosRune(s1, "C"), 2)
		t.Assert(dstr.PosRune(s1, "China"), 2)
	})
}

func Test_PosI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(dstr.PosI(s1, "zz"), -1)
		t.Assert(dstr.PosI(s1, "ab"), 0)
		t.Assert(dstr.PosI(s1, "ef", 2), 4)
		t.Assert(dstr.PosI(s1, "abd", 0), -1)
		t.Assert(dstr.PosI(s1, "E", -4), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosI(s1, "爱"), 3)
		t.Assert(dstr.PosI(s1, "c"), 6)
		t.Assert(dstr.PosI(s1, "china"), 6)
	})
}

func Test_PosIRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		t.Assert(dstr.PosIRune(s1, "zz"), -1)
		t.Assert(dstr.PosIRune(s1, "ab"), 0)
		t.Assert(dstr.PosIRune(s1, "ef", 2), 4)
		t.Assert(dstr.PosIRune(s1, "abd", 0), -1)
		t.Assert(dstr.PosIRune(s1, "E", -4), 11)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosIRune(s1, "爱"), 1)
		t.Assert(dstr.PosIRune(s1, "c"), 2)
		t.Assert(dstr.PosIRune(s1, "china"), 2)
	})
}

func Test_PosR(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(dstr.PosR(s1, "zz"), -1)
		t.Assert(dstr.PosR(s1, "ab"), 7)
		t.Assert(dstr.PosR(s2, "ab", -2), 0)
		t.Assert(dstr.PosR(s1, "ef"), 11)
		t.Assert(dstr.PosR(s1, "abd", 0), -1)
		t.Assert(dstr.PosR(s1, "e", -4), -1)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosR(s1, "爱"), 3)
		t.Assert(dstr.PosR(s1, "C"), 6)
		t.Assert(dstr.PosR(s1, "China"), 6)
	})
}

func Test_PosRRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(dstr.PosRRune(s1, "zz"), -1)
		t.Assert(dstr.PosRRune(s1, "ab"), 7)
		t.Assert(dstr.PosRRune(s2, "ab", -2), 0)
		t.Assert(dstr.PosRRune(s1, "ef"), 11)
		t.Assert(dstr.PosRRune(s1, "abd", 0), -1)
		t.Assert(dstr.PosRRune(s1, "e", -4), -1)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosRRune(s1, "爱"), 1)
		t.Assert(dstr.PosRRune(s1, "C"), 2)
		t.Assert(dstr.PosRRune(s1, "China"), 2)
	})
}

func Test_PosRI(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(dstr.PosRI(s1, "zz"), -1)
		t.Assert(dstr.PosRI(s1, "AB"), 7)
		t.Assert(dstr.PosRI(s2, "AB", -2), 0)
		t.Assert(dstr.PosRI(s1, "EF"), 11)
		t.Assert(dstr.PosRI(s1, "abd", 0), -1)
		t.Assert(dstr.PosRI(s1, "e", -5), 4)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosRI(s1, "爱"), 3)
		t.Assert(dstr.PosRI(s1, "C"), 19)
		t.Assert(dstr.PosRI(s1, "China"), 6)
	})
}

func Test_PosRIRune(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s1 := "abcdEFGabcdefg"
		s2 := "abcdEFGz1cdeab"
		t.Assert(dstr.PosRIRune(s1, "zz"), -1)
		t.Assert(dstr.PosRIRune(s1, "AB"), 7)
		t.Assert(dstr.PosRIRune(s2, "AB", -2), 0)
		t.Assert(dstr.PosRIRune(s1, "EF"), 11)
		t.Assert(dstr.PosRIRune(s1, "abd", 0), -1)
		t.Assert(dstr.PosRIRune(s1, "e", -5), 4)
	})
	dtest.C(t, func(t *dtest.T) {
		s1 := "我爱China very much"
		t.Assert(dstr.PosRIRune(s1, "爱"), 1)
		t.Assert(dstr.PosRIRune(s1, "C"), 15)
		t.Assert(dstr.PosRIRune(s1, "China"), 2)
	})
}
