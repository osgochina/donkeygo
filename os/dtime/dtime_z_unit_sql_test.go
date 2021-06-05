package dtime_test

import (
	"github.com/osgochina/donkeygo/os/dtime"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestTime_Scan(t1 *testing.T) {
	dtest.C(t1, func(t *dtest.T) {
		tt := dtime.Time{}
		//test string
		s := dtime.Now().String()
		t.Assert(tt.Scan(s), nil)
		t.Assert(tt.String(), s)
		//test nano
		n := dtime.TimestampNano()
		t.Assert(tt.Scan(n), nil)
		t.Assert(tt.TimestampNano(), n)
		//test nil
		none := (*dtime.Time)(nil)
		t.Assert(none.Scan(nil), nil)
		t.Assert(none, nil)
	})

}

func TestTime_Value(t1 *testing.T) {
	dtest.C(t1, func(t *dtest.T) {
		tt := dtime.Now()
		s, err := tt.Value()
		t.Assert(err, nil)
		t.Assert(s, tt.Time)
		//test nil
		none := (*dtime.Time)(nil)
		s, err = none.Value()
		t.Assert(err, nil)
		t.Assert(s, nil)

	})
}
