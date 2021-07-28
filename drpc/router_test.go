package drpc

import (
	"github.com/osgochina/donkeygo/drpc/status"
	"github.com/osgochina/donkeygo/os/dlog"
	"github.com/osgochina/donkeygo/test/dtest"
	"testing"
)

//func TestHTTPServiceMethodMapper(t *testing.T) {
//	dtest.C(t, func(t *dtest.T) {
//		t.Assert(globalServiceMethodMapper("Abc", "Efg"), "/Abc/efg")
//		t.Assert(globalServiceMethodMapper("", ""), "/")
//	})
//}

func TestRouteCall(t *testing.T) {
	dtest.C(t, func(t *dtest.T) {
		root := newRouter(newPluginContainer())
		root.RouteCall(new(Math))
		h, found := root.subRouter.getCall("/math/add")
		if found {
			h.handleFunc(newReadHandleCtx(), h.NewArgValue())
		}
	})
}

type Math struct {
	name string
	CallCtx
}

// Add handles addition request
func (m *Math) Add(arg *[]int) (int, *status.Status) {
	// test meta
	dlog.Infof("author: %s", m.PeekMeta("author"))
	// add
	var r int
	for _, a := range *arg {
		r += a
	}
	// response
	return r, nil
}
