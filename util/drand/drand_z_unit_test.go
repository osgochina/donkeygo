// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package drand_test

import (
	"github.com/osgochina/donkeygo/test/dtest"
	"github.com/osgochina/donkeygo/text/dstr"
	"github.com/osgochina/donkeygo/util/drand"
	"testing"
	"time"
)

func Test_Intn(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 1000000; i++ {
			n := drand.Intn(100)
			t.AssertLT(n, 100)
			t.AssertGE(n, 0)
		}
		for i := 0; i < 1000000; i++ {
			n := drand.Intn(-100)
			t.AssertLE(n, 0)
			t.Assert(n, -100)
		}
	})
}

func Test_Meet(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(drand.Meet(100, 100), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(drand.Meet(0, 100), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.Meet(50, 100), []bool{true, false})
		}
	})
}

func Test_MeetProb(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(drand.MeetProb(1), true)
		}
		for i := 0; i < 100; i++ {
			t.Assert(drand.MeetProb(0), false)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.MeetProb(0.5), []bool{true, false})
		}
	})
}

func Test_N(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(drand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(drand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.N(1, 2), []int{1, 2})
		}
	})
}

func Test_D(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(drand.D(time.Second, time.Second), time.Second)
		}
		for i := 0; i < 100; i++ {
			t.Assert(drand.D(0, 0), time.Duration(0))
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(
				drand.D(1*time.Second, 3*time.Second),
				[]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second},
			)
		}
	})
}

func Test_Rand(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(drand.N(1, 1), 1)
		}
		for i := 0; i < 100; i++ {
			t.Assert(drand.N(0, 0), 0)
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.N(1, 2), []int{1, 2})
		}
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.N(-1, 2), []int{-1, 0, 1, 2})
		}
	})
}

func Test_S(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.S(5)), 5)
		}
	})
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.S(5, true)), 5)
		}
	})
}

func Test_B(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			b := drand.B(5)
			t.Assert(len(b), 5)
			t.AssertNE(b, make([]byte, 5))
		}
	})
}

func Test_Str(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.S(5)), 5)
		}
	})
}

func Test_RandStr(t *testing.T) {
	str := "我爱GoFrame"
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 10; i++ {
			s := drand.Str(str, 100000)
			t.Assert(dstr.Contains(s, "我"), true)
			t.Assert(dstr.Contains(s, "爱"), true)
			t.Assert(dstr.Contains(s, "G"), true)
			t.Assert(dstr.Contains(s, "o"), true)
			t.Assert(dstr.Contains(s, "F"), true)
			t.Assert(dstr.Contains(s, "r"), true)
			t.Assert(dstr.Contains(s, "a"), true)
			t.Assert(dstr.Contains(s, "m"), true)
			t.Assert(dstr.Contains(s, "e"), true)
			t.Assert(dstr.Contains(s, "w"), false)
		}
	})
}

func Test_Digits(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.Digits(5)), 5)
		}
	})
}

func Test_RandDigits(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.Digits(5)), 5)
		}
	})
}

func Test_Letters(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.Letters(5)), 5)
		}
	})
}

func Test_RandLetters(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.Assert(len(drand.Letters(5)), 5)
		}
	})
}

func Test_Perm(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for i := 0; i < 100; i++ {
			t.AssertIN(drand.Perm(5), []int{0, 1, 2, 3, 4})
		}
	})
}
