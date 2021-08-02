// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils_test

import (
	"github.com/osgochina/donkeygo/internal/utils"
	"github.com/osgochina/donkeygo/test/dtest"
	"io/ioutil"
	"testing"
)

func Test_ReadCloser(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		var (
			n    int
			b    = make([]byte, 3)
			body = utils.NewReadCloser([]byte{1, 2, 3, 4}, false)
		)
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{1, 2, 3})
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{4})

		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{})
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{})
	})
	dtest.C(t, func(t *dtest.T) {
		var (
			r    []byte
			body = utils.NewReadCloser([]byte{1, 2, 3, 4}, false)
		)
		r, _ = ioutil.ReadAll(body)
		t.Assert(r, []byte{1, 2, 3, 4})
		r, _ = ioutil.ReadAll(body)
		t.Assert(r, []byte{})
	})
	dtest.C(t, func(t *dtest.T) {
		var (
			n    int
			r    []byte
			b    = make([]byte, 3)
			body = utils.NewReadCloser([]byte{1, 2, 3, 4}, true)
		)
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{1, 2, 3})
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{4})

		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{1, 2, 3})
		n, _ = body.Read(b)
		t.Assert(b[:n], []byte{4})

		r, _ = ioutil.ReadAll(body)
		t.Assert(r, []byte{1, 2, 3, 4})
		r, _ = ioutil.ReadAll(body)
		t.Assert(r, []byte{1, 2, 3, 4})
	})
}

func Test_RemoveSymbols(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(utils.RemoveSymbols(`-a-b._a c1!@#$%^&*()_+:";'.,'01`), `abac101`)
	})
}
