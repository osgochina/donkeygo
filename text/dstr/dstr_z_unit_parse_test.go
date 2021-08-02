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
	"net/url"
	"testing"
)

func Test_Parse(t *testing.T) {
	// url
	dtest.C(t, func(t *dtest.T) {
		s := "goframe.org/index?name=john&score=100"
		u, err := url.Parse(s)
		t.Assert(err, nil)
		m, err := dstr.Parse(u.RawQuery)
		t.Assert(err, nil)
		t.Assert(m["name"], "john")
		t.Assert(m["score"], "100")

		// name overwrite
		m, err = dstr.Parse("a=1&a=2")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"a": 2,
		})
		// slice
		m, err = dstr.Parse("a[]=1&a[]=2")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"a": d.Slice{"1", "2"},
		})
		// map
		m, err = dstr.Parse("a=1&b=2&c=3")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"a": "1",
			"b": "2",
			"c": "3",
		})
		m, err = dstr.Parse("a=1&a=2&c=3")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"a": "2",
			"c": "3",
		})
		// map
		m, err = dstr.Parse("m[a]=1&m[b]=2&m[c]=3")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"m": d.Map{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		})
		m, err = dstr.Parse("m[a]=1&m[a]=2&m[b]=3")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"m": d.Map{
				"a": "2",
				"b": "3",
			},
		})
		// map - slice
		m, err = dstr.Parse("m[a][]=1&m[a][]=2")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"m": d.Map{
				"a": d.Slice{"1", "2"},
			},
		})
		m, err = dstr.Parse("m[a][b][]=1&m[a][b][]=2")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"m": d.Map{
				"a": d.Map{
					"b": d.Slice{"1", "2"},
				},
			},
		})
		// map - complicated
		m, err = dstr.Parse("m[a1][b1][c1][d1]=1&m[a2][b2]=2&m[a3][b3][c3]=3")
		t.Assert(err, nil)
		t.Assert(m, d.Map{
			"m": d.Map{
				"a1": d.Map{
					"b1": d.Map{
						"c1": d.Map{
							"d1": "1",
						},
					},
				},
				"a2": d.Map{
					"b2": "2",
				},
				"a3": d.Map{
					"b3": d.Map{
						"c3": "3",
					},
				},
			},
		})
	})
}
