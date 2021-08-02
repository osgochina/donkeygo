// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dbase64_test

import (
	"github.com/osgochina/donkeygo/debug/ddebug"
	"github.com/osgochina/donkeygo/encoding/dbase64"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

type testPair struct {
	decoded, encoded string
}

var pairs = []testPair{
	// RFC 3548 examples
	{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+"},
	{"\x14\xfb\x9c\x03\xd9", "FPucA9k="},
	{"\x14\xfb\x9c\x03", "FPucAw=="},

	// RFC 4648 examples
	{"", ""},
	{"f", "Zg=="},
	{"fo", "Zm8="},
	{"foo", "Zm9v"},
	{"foob", "Zm9vYg=="},
	{"fooba", "Zm9vYmE="},
	{"foobar", "Zm9vYmFy"},

	// Wikipedia examples
	{"sure.", "c3VyZS4="},
	{"sure", "c3VyZQ=="},
	{"sur", "c3Vy"},
	{"su", "c3U="},
	{"leasure.", "bGVhc3VyZS4="},
	{"easure.", "ZWFzdXJlLg=="},
	{"asure.", "YXN1cmUu"},
	{"sure.", "c3VyZS4="},
}

func Test_Basic(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		for k := range pairs {
			// Encode
			t.Assert(dbase64.Encode([]byte(pairs[k].decoded)), []byte(pairs[k].encoded))
			t.Assert(dbase64.EncodeToString([]byte(pairs[k].decoded)), pairs[k].encoded)
			t.Assert(dbase64.EncodeString(pairs[k].decoded), pairs[k].encoded)

			// Decode
			r1, _ := dbase64.Decode([]byte(pairs[k].encoded))
			t.Assert(r1, []byte(pairs[k].decoded))

			r2, _ := dbase64.DecodeString(pairs[k].encoded)
			t.Assert(r2, []byte(pairs[k].decoded))

			r3, _ := dbase64.DecodeToString(pairs[k].encoded)
			t.Assert(r3, pairs[k].decoded)
		}
	})
}

func Test_File(t *testing.T) {
	path := ddebug.TestDataPath("test")
	expect := "dGVzdA=="
	dtest.C(t, func(t *dtest.T) {
		b, err := dbase64.EncodeFile(path)
		t.Assert(err, nil)
		t.Assert(string(b), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		s, err := dbase64.EncodeFileToString(path)
		t.Assert(err, nil)
		t.Assert(s, expect)
	})
}

func Test_File_Error(t *testing.T) {
	path := "none-exist-file"
	expect := ""
	dtest.C(t, func(t *dtest.T) {
		b, err := dbase64.EncodeFile(path)
		t.AssertNE(err, nil)
		t.Assert(string(b), expect)
	})
	dtest.C(t, func(t *dtest.T) {
		s, err := dbase64.EncodeFileToString(path)
		t.AssertNE(err, nil)
		t.Assert(s, expect)
	})
}
