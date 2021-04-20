// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package dcmd_test

import (
	"donkeygo/os/dcmd"
	"donkeygo/os/denv"
	"donkeygo/test/dtest"
	"testing"
)

func Test_Default(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		dcmd.Init([]string{"gf", "--force", "remove", "-fq", "-p=www", "path", "-n", "root"}...)
		t.Assert(len(dcmd.GetArgAll()), 2)
		t.Assert(dcmd.GetArg(1), "path")
		t.Assert(dcmd.GetArg(100, "test"), "test")
		t.Assert(dcmd.GetOpt("force"), "remove")
		t.Assert(dcmd.GetOpt("n"), "root")
		t.Assert(dcmd.ContainsOpt("fq"), true)
		t.Assert(dcmd.ContainsOpt("p"), true)
		t.Assert(dcmd.ContainsOpt("none"), false)
		t.Assert(dcmd.GetOpt("none", "value"), "value")
	})
	dtest.C(t, func(t *dtest.T) {
		dcmd.Init([]string{"gf", "gen", "-h"}...)
		t.Assert(len(dcmd.GetArgAll()), 2)
		t.Assert(dcmd.GetOpt("h"), "")
		t.Assert(dcmd.ContainsOpt("h"), true)
	})
}

func Test_BuildOptions(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		s := dcmd.BuildOptions(map[string]string{
			"n": "john",
		})
		t.Assert(s, "-n=john")
	})

	dtest.C(t, func(t *dtest.T) {
		s := dcmd.BuildOptions(map[string]string{
			"n": "john",
		}, "-test")
		t.Assert(s, "-testn=john")
	})
}

func Test_GetWithEnv(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		denv.Set("TEST", "1")
		defer denv.Remove("TEST")
		t.Assert(dcmd.GetOptWithEnv("test"), 1)
	})
	dtest.C(t, func(t *dtest.T) {
		denv.Set("TEST", "1")
		defer denv.Remove("TEST")
		dcmd.Init("-test", "2")
		t.Assert(dcmd.GetOptWithEnv("test"), 2)
	})
}
