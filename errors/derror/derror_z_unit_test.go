// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package derror_test

import (
	"errors"
	"fmt"
	"github.com/osgochina/donkeygo/errors/derror"
	"github.com/osgochina/donkeygo/internal/json"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func nilError() error {
	return nil
}

func Test_Nil(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		t.Assert(derror.New(""), nil)
		t.Assert(derror.Wrap(nilError(), "test"), nil)
	})
}

func Test_New(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.Newf("%d", 1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewSkip(1, "1")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewSkipf(1, "%d", 1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Wrap(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Wrapf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrapf(err, "%d", 2)
		err = derror.Wrapf(err, "%d", 3)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrapf(err, "%d", 2)
		err = derror.Wrapf(err, "%d", 3)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrapf(err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_WrapSkip(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.WrapSkip(1, err, "2")
		err = derror.WrapSkip(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.WrapSkip(1, err, "2")
		err = derror.WrapSkip(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.WrapSkip(1, err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_WrapSkipf(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.WrapSkipf(1, err, "2")
		err = derror.WrapSkipf(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.WrapSkipf(1, err, "2")
		err = derror.WrapSkipf(1, err, "3")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.WrapSkipf(1, err, "")
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "1")
	})
}

func Test_Cause(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		t.Assert(derror.Cause(err), err)
	})

	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.Assert(derror.Cause(err), "1")
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		t.Assert(derror.Cause(err), "1")
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.Assert(derror.Cause(err), "1")
	})
}

func Test_Format(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		t.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%s", err), "3: 2: 1")
		t.Assert(fmt.Sprintf("%v", err), "3: 2: 1")
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		t.Assert(fmt.Sprintf("%-s", err), "3")
		t.Assert(fmt.Sprintf("%-v", err), "3")
	})
}

func Test_Stack(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		t.Assert(fmt.Sprintf("%+v", err), "1")
	})

	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		t.AssertNE(fmt.Sprintf("%+v", err), "1")
		//fmt.Printf("%+v", err)
	})

	dtest.C(t, func(t *dtest.T) {
		err := derror.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.AssertNE(err, nil)
		//fmt.Printf("%+v", err)
	})
}

func Test_Current(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")
		t.Assert(derror.Current(err).Error(), "3")
	})
}

func Test_Next(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.Wrap(err, "3")
		t.Assert(err.Error(), "3: 2: 1")

		err = derror.Next(err)
		t.Assert(err.Error(), "2: 1")

		err = derror.Next(err)
		t.Assert(err.Error(), "1")

		err = derror.Next(err)
		t.Assert(err, nil)
	})
}

func Test_Code(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("123")
		t.Assert(derror.Code(err), -1)
		t.Assert(err.Error(), "123")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewCode(1, "123")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewCodef(1, "%s", "123")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewCodeSkip(1, 0, "123")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	dtest.C(t, func(t *dtest.T) {
		err := derror.NewCodeSkipf(1, 0, "%s", "123")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "123")
	})
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.WrapCode(1, err, "3")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.WrapCodef(1, err, "%s", "3")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.WrapCodeSkip(1, 100, err, "3")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
	dtest.C(t, func(t *dtest.T) {
		err := errors.New("1")
		err = derror.Wrap(err, "2")
		err = derror.WrapCodeSkipf(1, 100, err, "%s", "3")
		t.Assert(derror.Code(err), 1)
		t.Assert(err.Error(), "3: 2: 1")
	})
}

func Test_Json(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		err := derror.Wrap(derror.New("1"), "2")
		b, e := json.Marshal(err)
		t.Assert(e, nil)
		t.Assert(string(b), `"2: 1"`)
	})
}
