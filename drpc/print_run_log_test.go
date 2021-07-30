package drpc

import (
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

func TestEnablePrintRunLog(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		_ = dlog.SetLevelStr("ALL")
		t.Assert(enablePrintRunLog(), true)
		_ = dlog.SetLevelStr("DEV")
		t.Assert(enablePrintRunLog(), true)
		_ = dlog.SetLevelStr("DEVELOP")
		t.Assert(enablePrintRunLog(), true)
		_ = dlog.SetLevelStr("PROD")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("PRODUCT")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("DEBU")
		t.Assert(enablePrintRunLog(), true)
		_ = dlog.SetLevelStr("DEBUG")
		t.Assert(enablePrintRunLog(), true)
		_ = dlog.SetLevelStr("INFO")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("NOTI")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("NOTICE")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("WARN")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("WARNING")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("ERRO")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("ERROR")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("CRIT")
		t.Assert(enablePrintRunLog(), false)
		_ = dlog.SetLevelStr("CRITICAL")
		t.Assert(enablePrintRunLog(), false)
	})
}
