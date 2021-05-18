package drpc

import (
	"github.com/gogf/gf/os/glog"
	"github.com/osgochina/donkeygo/drpc/status"
	"github.com/osgochina/donkeygo/text/dstr"
	"github.com/pkg/errors"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// HTTPServiceMethodMapper like most RPC services service method mapper.
// Such as: user/get
// It is the default mapper.
// The mapping rule of struct(func) name to service methods:
//  `AaBb` -> `/aa_bb`
//  `ABcXYz` -> `/abc_xyz`
//  `Aa__Bb` -> `/aa_bb`
//  `aa__bb` -> `/aa_bb`
//  `ABC__XYZ` -> `/abc_xyz`
//  `Aa_Bb` -> `/aa/bb`
//  `aa_bb` -> `/aa/bb`
//  `ABC_XYZ` -> `/abc/xyz`
//
func HTTPServiceMethodMapper(prefix, name string) string {
	return path.Join("/", prefix, toServiceMethods(name, '/', true))
}

type ServiceMethodMapper func(prefix, name string) (serviceMethod string)

func SetServiceMethodMapper(mapper ServiceMethodMapper) {
	globalServiceMethodMapper = mapper
}

var globalServiceMethodMapper = HTTPServiceMethodMapper

// toServiceMethods maps struct(func) name to service methods.
func toServiceMethods(name string, sep rune, toSnake bool) string {
	var a = make([]rune, 0)
	var last rune
	for _, r := range name {
		if last == '_' {
			if r == '_' {
				last = '\x00'
				continue
			} else {
				a[len(a)-1] = sep
			}
		}
		if last == '\x00' && r == '_' {
			continue
		}
		a = append(a, r)
		last = r
	}
	name = string(a)
	if toSnake {
		name = dstr.SnakeString(name)
		name = strings.Replace(name, "__", "_", -1)
		name = strings.Replace(name, string(sep)+"_", string(sep), -1)
	}
	return name
}

const (
	pnPush        = "PUSH"
	pnCall        = "CALL"
	pnUnknownPush = "UNKNOWN_PUSH"
	pnUnknownCall = "UNKNOWN_CALL"
)

// Router 路由器
type Router struct {
	subRouter *SubRouter
}

type SubRouter struct {
	root            *Router
	callHandlers    map[string]*Handler
	pushHandlers    map[string]*Handler
	unknownCall     **Handler
	unknownPush     **Handler
	prefix          string
	pluginContainer *PluginContainer
}

//新建路由器
func newRouter(pluginContainer *PluginContainer) *Router {
	rootGroup := globalServiceMethodMapper("", "")
	root := &Router{
		subRouter: &SubRouter{
			callHandlers:    make(map[string]*Handler),
			pushHandlers:    make(map[string]*Handler),
			unknownCall:     new(*Handler),
			unknownPush:     new(*Handler),
			prefix:          rootGroup,
			pluginContainer: pluginContainer,
		},
	}
	root.subRouter.root = root
	return root
}

// SubRoute 添加处理程序组
func (that *Router) SubRoute(prefix string, plugin ...Plugin) *SubRouter {
	return that.subRouter.SubRoute(prefix, plugin...)
}

// RouteCall 注册struct对象到路由器
func (that *Router) RouteCall(callCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.subRouter.RouteCall(callCtrlStruct, plugin...)
}

// RouteCallFunc 注册func对象到路由器
func (that *Router) RouteCallFunc(callHandleFunc interface{}, plugin ...Plugin) string {
	return that.subRouter.RouteCallFunc(callHandleFunc, plugin...)
}

// RoutePush 注册 PUSH 类型的处理程序到路由器
func (that *Router) RoutePush(pushCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.subRouter.RoutePush(pushCtrlStruct, plugin...)
}

// RoutePushFunc 通过func注册PUSH类型的处理程序到路由器
func (that *Router) RoutePushFunc(pushHandleFunc interface{}, plugin ...Plugin) string {
	return that.subRouter.RoutePushFunc(pushHandleFunc, plugin...)
}

// SetUnknownCall 注册默认的未知CALL处理方法
func (that *Router) SetUnknownCall(fn func(UnknownCallCtx) (interface{}, *status.Status), plugin ...Plugin) {
	pluginContainer := that.subRouter.pluginContainer.cloneAndAppendMiddle(plugin...)
	warnInvalidHandlerHooks(plugin)

	var h = &Handler{
		name:            pnUnknownCall,
		isUnknown:       true,
		argElem:         reflect.TypeOf([]byte{}),
		pluginContainer: pluginContainer,
		unknownHandleFunc: func(ctx *handlerCtx) {
			body, stat := fn(ctx)
			if !stat.OK() {
				ctx.stat = stat
				ctx.output.SetStatus(stat)
			} else {
				ctx.output.SetBody(body)
			}
		},
	}
	if *that.subRouter.unknownCall == nil {
		glog.Printf("set %s handler", h.name)
	} else {
		glog.Warningf("covered %s handler", h.name)
	}
	that.subRouter.unknownCall = &h
}

// SetUnknownPush 注册未知PUSH处理方法
func (that *Router) SetUnknownPush(fn func(UnknownPushCtx) *status.Status, plugin ...Plugin) {
	pluginContainer := that.subRouter.pluginContainer.cloneAndAppendMiddle(plugin...)
	warnInvalidHandlerHooks(plugin)

	var h = &Handler{
		name:            pnUnknownPush,
		isUnknown:       true,
		argElem:         reflect.TypeOf([]byte{}),
		pluginContainer: pluginContainer,
		unknownHandleFunc: func(ctx *handlerCtx) {
			ctx.stat = fn(ctx)
		},
	}

	if *that.subRouter.unknownPush == nil {
		glog.Printf("set %s handler", h.name)
	} else {
		glog.Warningf("covered %s handler", h.name)
	}
	that.subRouter.unknownPush = &h
}

// Root 返回根路由器
func (that *SubRouter) Root() *Router {
	return that.root
}

// ToRouter 转换成添加了 SetUnknownCall 和 SetUnknownPush 方法的路由器
func (that *SubRouter) ToRouter() *Router {
	return &Router{subRouter: that}
}

// SubRoute 添加处理程序组
func (that *SubRouter) SubRoute(prefix string, plugin ...Plugin) *SubRouter {
	pluginContainer := that.pluginContainer.cloneAndAppendMiddle(plugin...)
	warnInvalidHandlerHooks(plugin)
	return &SubRouter{
		root:            that.root,
		callHandlers:    that.callHandlers,
		pushHandlers:    that.pushHandlers,
		unknownPush:     that.unknownPush,
		unknownCall:     that.unknownCall,
		prefix:          globalServiceMethodMapper(that.prefix, prefix),
		pluginContainer: pluginContainer,
	}
}

// RouteCall 通过struct注册多个 CALL 类型的处理程序，并返回它们的注册路径
func (that *SubRouter) RouteCall(callCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.reg(pnCall, makeCallHandlersFromStruct, callCtrlStruct, plugin)
}

// RouteCallFunc 通过func注册单个 CALL 类型的处理程序，并返回它的注册路径
func (that *SubRouter) RouteCallFunc(callHandleFunc interface{}, plugin ...Plugin) string {
	return that.reg(pnCall, makeCallHandlersFromFunc, callHandleFunc, plugin)[0]
}

// RoutePush 通过struct批量注册 PUSH 类型的处理程序，并返回它们的路径
func (that *SubRouter) RoutePush(pushCtrlStruct interface{}, plugin ...Plugin) []string {
	return that.reg(pnPush, makePushHandlersFromStruct, pushCtrlStruct, plugin)
}

// RoutePushFunc 通过func注册PUSH类型的处理程序，并返回它的路径
func (that *SubRouter) RoutePushFunc(pushHandleFunc interface{}, plugin ...Plugin) string {
	return that.reg(pnPush, makePushHandlersFromFunc, pushHandleFunc, plugin)[0]
}

//注册路由器
func (that *SubRouter) reg(
	routerTypeName string,
	handlerMaker func(string, interface{}, *PluginContainer) ([]*Handler, error),
	ctrlStruct interface{},
	plugins []Plugin,
) []string {
	pluginContainer := that.pluginContainer.cloneAndAppendMiddle(plugins...)
	warnInvalidHandlerHooks(plugins)
	handlers, err := handlerMaker(
		that.prefix,
		ctrlStruct,
		pluginContainer,
	)
	if err != nil {
		glog.Fatalf("%v", err)
	}
	var names []string
	var hadHandlers map[string]*Handler

	if routerTypeName == pnCall {
		hadHandlers = that.callHandlers
	} else {
		hadHandlers = that.pushHandlers
	}

	for _, h := range handlers {
		if _, ok := hadHandlers[h.name]; ok {
			glog.Fatalf("there is a handler conflict: %s", h.name)
		}
		h.routerTypeName = routerTypeName
		hadHandlers[h.name] = h
		//pluginContainer.postReg(h)
		glog.Printf("register %s handler: %s", routerTypeName, h.name)
		names = append(names, h.name)
	}
	return names
}

// 获取路由器中指定路径CALL的处理方法
func (that *SubRouter) getCall(uriPath string) (*Handler, bool) {
	t, ok := that.callHandlers[uriPath]
	if ok {
		return t, true
	}
	if unknown := *that.unknownCall; unknown != nil {
		return unknown, true
	}
	return nil, false
}

// 获取路由器中指定路径的PUSH处理方法，未找到则使用注册的默认方法
func (that *SubRouter) getPush(uriPath string) (*Handler, bool) {
	t, ok := that.pushHandlers[uriPath]
	if ok {
		return t, true
	}
	if unknown := *that.unknownPush; unknown != nil {
		return unknown, true
	}
	return nil, false
}

// callCtrlStruct 需要实现 CallCtx 接口
func makeCallHandlersFromStruct(prefix string, callCtrlStruct interface{}, pluginContainer *PluginContainer) ([]*Handler, error) {

	var (
		cType    = reflect.TypeOf(callCtrlStruct)
		handlers = make([]*Handler, 0, 1)
	)
	//判断是否是指针类型
	if cType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("call-handler: the type is not struct point: %s", cType.String())
	}
	var cTypeElem = cType.Elem()
	//判断是否是struct类型
	if cTypeElem.Kind() != reflect.Struct {
		return nil, errors.Errorf("call-handler: the type is not struct point: %s", cType.String())
	}
	//如果结构体没有实现 CallCtx 的方法，或者不是匿名结构体
	iType, ok := cTypeElem.FieldByName("CallCtx")
	if !ok || !iType.Anonymous {
		return nil, errors.Errorf("call-handler: the struct do not have anonymous field drpc.CallCtx: %s", cType.String())
	}

	var callCtxOffset = iType.Offset

	// 判断插件容器是否存在，不存在则创建一个
	if pluginContainer == nil {
		pluginContainer = newPluginContainer()
	}

	type CallCtrlValue struct {
		ctrl   reflect.Value
		ctxPtr *CallCtx
	}

	var pool = &sync.Pool{
		New: func() interface{} {
			ctrl := reflect.New(cTypeElem)
			callCtxPtr := ctrl.Pointer() + callCtxOffset
			ctxPtr := (*CallCtx)(unsafe.Pointer(callCtxPtr))
			return &CallCtrlValue{
				ctrl:   ctrl,
				ctxPtr: ctxPtr,
			}
		},
	}

	//循环判断结构体中的方法是否需要注册
	for m := 0; m < cType.NumMethod(); m++ {
		//获取方法信息
		method := cType.Method(m)
		mType := method.Type
		mName := method.Name
		//fmt.Println(mName,mType.NumIn())
		//如果结构体中的方法不存在包路径
		if method.PkgPath != "" {
			continue
		}
		//如果这个方法是 CallCtx 这个接口提供的，则跳过
		if isBelongToCallCtx(mName) {
			continue
		}
		//如果方法的入参不是两个参数
		if mType.NumIn() != 2 {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s needs one in argument, but have %d", cType.String(), mName, mType.NumIn())
		}
		//获取第0个参数，按理应该返回的是一个struct类型的指针
		structType := mType.In(0)
		if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s receiver need be a struct pointer: %s", cType.String(), mName, structType)
		}
		//第一个参数必须是可以被外部调用的方法或者是内置数据类型的指针
		argType := mType.In(1)
		if !isExportedOrBuiltinType(argType) {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s arg type not exported: %s", cType.String(), mName, argType)
		}
		//必须是指针类型
		if argType.Kind() != reflect.Ptr {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s arg type need be a pointer: %s", cType.String(), mName, argType)
		}
		// 返回值必须是两个
		if mType.NumOut() != 2 {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s needs two out arguments, but have %d", cType.String(), mName, mType.NumOut())
		}

		//返回值的第一个参数必须是一个可以外部使用的类型或者基础类型
		replyType := mType.Out(0)
		if !isExportedOrBuiltinType(replyType) {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s first reply type not exported: %s", cType.String(), mName, replyType)
		}
		//第二个返回值必须是*Status类型
		if returnType := mType.Out(1); !isStatusType(returnType.String()) {
			//if isBelongToCallCtx(mName) {
			//	continue
			//}
			return nil, errors.Errorf("call-handler: %s.%s second out argument %s is not *drpc.status.Status", cType.String(), mName, returnType)
		}

		var methodFunc = method.Func
		var handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			obj := pool.Get().(*CallCtrlValue)
			*obj.ctxPtr = ctx
			rets := methodFunc.Call([]reflect.Value{obj.ctrl, argValue})
			stat := (*status.Status)(unsafe.Pointer(rets[1].Pointer()))
			if !stat.OK() {
				ctx.stat = stat
				ctx.output.SetStatus(stat)
			} else {
				ctx.output.SetBody(rets[0].Interface())
			}
			pool.Put(obj)
		}
		handlers = append(handlers, &Handler{
			handleFunc:      handleFunc,
			argElem:         argType.Elem(),
			reply:           replyType,
			pluginContainer: pluginContainer,
			name: globalServiceMethodMapper(
				globalServiceMethodMapper(prefix, ctrlStructName(cType)),
				mName,
			),
		})
	}

	return handlers, nil
}

func makeCallHandlersFromFunc(prefix string, callHandleFunc interface{}, pluginContainer *PluginContainer) ([]*Handler, error) {

	var (
		cType      = reflect.TypeOf(callHandleFunc)
		cValue     = reflect.ValueOf(callHandleFunc)
		typeString = objectName(cValue)
	)
	//如果传入的不是Func 则报错
	if cType.Kind() != reflect.Func {
		return nil, errors.Errorf("call-handler: the type is not function: %s", typeString)
	}
	// 需要两个参数： CallCtx, *<T>.
	if cType.NumIn() != 2 {
		return nil, errors.Errorf("call-handler: %s needs two in argument, but have %d", typeString, cType.NumIn())
	}
	// 需要判断第二个参数是可以导出的或基础类型
	argType := cType.In(1)
	if !isExportedOrBuiltinType(argType) {
		return nil, errors.Errorf("call-handler: %s arg type not exported: %s", typeString, argType)
	}
	//第二个参数必须是指针类型
	if argType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("call-handler: %s arg type need be a pointer: %s", typeString, argType)
	}
	// 需要两个返回值: reply, *status.Status.
	if cType.NumOut() != 2 {
		return nil, errors.Errorf("call-handler: %s needs two out arguments, but have %d", typeString, cType.NumOut())
	}

	//该返回值必须是可以导出使用或者是基础类型
	replyType := cType.Out(0)
	if !isExportedOrBuiltinType(replyType) {
		return nil, errors.Errorf("call-handler: %s first reply type not exported: %s", typeString, replyType)
	}
	//第二个返回值必须是*Status类型
	if returnType := cType.Out(1); !isStatusType(returnType.String()) {
		return nil, errors.Errorf("call-handler: %s second out argument %s is not *drcp.status.Status", typeString, returnType)
	}
	//first agr need be a CallCtx (struct pointer or CallCtx).
	ctxType := cType.In(0)

	var handleFunc func(*handlerCtx, reflect.Value)

	switch ctxType.Kind() {
	default:
		return nil, errors.Errorf("call-handler: %s's first arg must be drpc.CallCtx type or struct pointer: %s", typeString, ctxType)
	case reflect.Interface:
		iFace := reflect.TypeOf((*CallCtx)(nil)).Elem()
		if !ctxType.Implements(iFace) ||
			!iFace.Implements(reflect.New(ctxType).Type().Elem()) {
			return nil, errors.Errorf("call-handler: %s's first arg must be drpc.CallCtx type or struct pointer: %s", typeString, ctxType)
		}
		handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			rets := cValue.Call([]reflect.Value{reflect.ValueOf(ctx), argValue})
			stat := (*status.Status)(unsafe.Pointer(rets[1].Pointer()))
			if !stat.OK() {
				ctx.stat = stat
				ctx.output.SetStatus(stat)
			} else {
				ctx.output.SetBody(rets[0].Interface())
			}
		}
	case reflect.Ptr:
		var ctxTypeElem = ctxType.Elem()
		//对象必须是struct类型
		if ctxTypeElem.Kind() != reflect.Struct {
			return nil, errors.Errorf("call-handler: %s's first arg must be drpc.CallCtx type or struct pointer: %s", typeString, ctxType)
		}
		//对象必须实现了CallCtx接口
		iType, ok := ctxTypeElem.FieldByName("CallCtx")
		if !ok || !iType.Anonymous {
			return nil, errors.Errorf("call-handler: %s's first arg do not have anonymous field drpc.CallCtx: %s", typeString, ctxType)
		}
		type CallCtrlValue struct {
			ctrl   reflect.Value
			ctxPtr *CallCtx
		}
		var callCtxOffset = iType.Offset
		var pool = &sync.Pool{
			New: func() interface{} {
				ctrl := reflect.New(ctxTypeElem)
				callCtxPtr := ctrl.Pointer() + callCtxOffset
				ctxPtr := (*CallCtx)(unsafe.Pointer(callCtxPtr))
				return &CallCtrlValue{
					ctrl:   ctrl,
					ctxPtr: ctxPtr,
				}
			},
		}
		handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			obj := pool.Get().(*CallCtrlValue)
			*obj.ctxPtr = ctx
			rets := cValue.Call([]reflect.Value{obj.ctrl, argValue})
			stat := (*status.Status)(unsafe.Pointer(rets[1].Pointer()))
			if !stat.OK() {
				ctx.stat = stat
				ctx.output.SetStatus(stat)
			} else {
				ctx.output.SetBody(rets[0].Interface())
			}
			pool.Put(obj)
		}
	}
	if pluginContainer == nil {
		pluginContainer = newPluginContainer()
	}
	return []*Handler{{
		name:            globalServiceMethodMapper(prefix, handlerFuncName(cValue)),
		handleFunc:      handleFunc,
		argElem:         argType.Elem(),
		reply:           replyType,
		pluginContainer: pluginContainer,
	}}, nil
}

func makePushHandlersFromStruct(prefix string, pushCtrlStruct interface{}, pluginContainer *PluginContainer) ([]*Handler, error) {

	var (
		cType    = reflect.TypeOf(pushCtrlStruct)
		handlers = make([]*Handler, 0, 1)
	)
	//判断传入的必须是指针类型
	if cType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("push-handler: the type is not struct point: %s", cType.String())
	}
	//必须是struct类型
	var cTypeElem = cType.Elem()
	if cTypeElem.Kind() != reflect.Struct {
		return nil, errors.Errorf("push-handler: the type is not struct point: %s", cType.String())
	}
	//必须实现了PushCtx接口
	iType, ok := cTypeElem.FieldByName("PushCtx")
	if !ok || !iType.Anonymous {
		return nil, errors.Errorf("push-handler: the struct do not have anonymous field drpc.PushCtx: %s", cType.String())
	}

	var pushCtxOffset = iType.Offset

	if pluginContainer == nil {
		pluginContainer = newPluginContainer()
	}

	type PushCtrlValue struct {
		ctrl   reflect.Value
		ctxPtr *PushCtx
	}
	var pool = &sync.Pool{
		New: func() interface{} {
			ctrl := reflect.New(cTypeElem)
			pushCtxPtr := ctrl.Pointer() + pushCtxOffset
			ctxPtr := (*PushCtx)(unsafe.Pointer(pushCtxPtr))
			return &PushCtrlValue{
				ctrl:   ctrl,
				ctxPtr: ctxPtr,
			}
		},
	}

	for m := 0; m < cType.NumMethod(); m++ {
		method := cType.Method(m)
		mType := method.Type
		mName := method.Name

		//方法必须是可以导出的
		if method.PkgPath != "" {
			continue
		}
		//如果是CallCtx接口的基础方法，则跳过
		if isBelongToPushCtx(mName) {
			continue
		}

		// Method needs two ins: receiver, *<T>.
		if mType.NumIn() != 2 {
			return nil, errors.Errorf("push-handler: %s.%s needs one in argument, but have %d", cType.String(), mName, mType.NumIn())
		}
		//第一个参数必须是struct类型
		structType := mType.In(0)
		if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
			return nil, errors.Errorf("push-handler: %s.%s receiver need be a struct pointer: %s", cType.String(), mName, structType)
		}
		//第二个参数必须是可以导出的或基础类型
		argType := mType.In(1)
		if !isExportedOrBuiltinType(argType) {
			return nil, errors.Errorf("push-handler: %s.%s arg type not exported: %s", cType.String(), mName, argType)
		}
		if argType.Kind() != reflect.Ptr {
			return nil, errors.Errorf("push-handler: %s.%s arg type need be a pointer: %s", cType.String(), mName, argType)
		}
		//返回参数如果不是一个
		if mType.NumOut() != 1 {
			return nil, errors.Errorf("push-handler: %s.%s needs one out arguments, but have %d", cType.String(), mName, mType.NumOut())
		}
		//返回参数必须是*Status类型
		if returnType := mType.Out(0); !isStatusType(returnType.String()) {
			return nil, errors.Errorf("push-handler: %s.%s out argument %s is not *drpc.Status", cType.String(), mName, returnType)
		}
		var methodFunc = method.Func
		var handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			obj := pool.Get().(*PushCtrlValue)
			*obj.ctxPtr = ctx
			rets := methodFunc.Call([]reflect.Value{obj.ctrl, argValue})
			ctx.stat = (*status.Status)(unsafe.Pointer(rets[0].Pointer()))
			pool.Put(obj)
		}
		handlers = append(handlers, &Handler{
			handleFunc:      handleFunc,
			argElem:         argType.Elem(),
			pluginContainer: pluginContainer,
			name: globalServiceMethodMapper(
				globalServiceMethodMapper(prefix, ctrlStructName(cType)),
				mName,
			),
		})
	}
	return handlers, nil
}

func makePushHandlersFromFunc(prefix string, pushHandleFunc interface{}, pluginContainer *PluginContainer) ([]*Handler, error) {

	var (
		cType      = reflect.TypeOf(pushHandleFunc)
		cValue     = reflect.ValueOf(pushHandleFunc)
		typeString = objectName(cValue)
	)
	if cType.Kind() != reflect.Func {
		return nil, errors.Errorf("push-handler: the type is not function: %s", typeString)
	}

	// needs one out: *Status.
	if cType.NumOut() != 1 {
		return nil, errors.Errorf("push-handler: %s needs one out arguments, but have %d", typeString, cType.NumOut())
	}
	if returnType := cType.Out(0); !isStatusType(returnType.String()) {
		return nil, errors.Errorf("push-handler: %s out argument %s is not *drpc.Status", typeString, returnType)
	}
	// needs two ins: PushCtx, *<T>.
	if cType.NumIn() != 2 {
		return nil, errors.Errorf("push-handler: %s needs two in argument, but have %d", typeString, cType.NumIn())
	}
	argType := cType.In(1)
	if !isExportedOrBuiltinType(argType) {
		return nil, errors.Errorf("push-handler: %s arg type not exported: %s", typeString, argType)
	}
	if argType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("push-handler: %s arg type need be a pointer: %s", typeString, argType)
	}
	// first agr need be a PushCtx (struct pointer or PushCtx).
	ctxType := cType.In(0)

	var handleFunc func(*handlerCtx, reflect.Value)
	switch ctxType.Kind() {
	default:
		return nil, errors.Errorf("push-handler: %s's first arg must be drpc.PushCtx type or struct pointer: %s", typeString, ctxType)
	case reflect.Interface:
		iFace := reflect.TypeOf((*PushCtx)(nil)).Elem()
		if !ctxType.Implements(iFace) ||
			!iFace.Implements(reflect.New(ctxType).Type().Elem()) {
			return nil, errors.Errorf("push-handler: %s's first arg need implement drpc.PushCtx: %s", typeString, ctxType)
		}

		handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			rets := cValue.Call([]reflect.Value{reflect.ValueOf(ctx), argValue})
			ctx.stat = (*status.Status)(unsafe.Pointer(rets[0].Pointer()))
		}
	case reflect.Ptr:
		var ctxTypeElem = ctxType.Elem()
		if ctxTypeElem.Kind() != reflect.Struct {
			return nil, errors.Errorf("push-handler: %s's first arg must be erpc.PushCtx type or struct pointer: %s", typeString, ctxType)
		}

		iType, ok := ctxTypeElem.FieldByName("PushCtx")
		if !ok || !iType.Anonymous {
			return nil, errors.Errorf("push-handler: %s's first arg do not have anonymous field drpc.PushCtx: %s", typeString, ctxType)
		}

		type PushCtrlValue struct {
			ctrl   reflect.Value
			ctxPtr *PushCtx
		}
		var pushCtxOffset = iType.Offset
		var pool = &sync.Pool{
			New: func() interface{} {
				ctrl := reflect.New(ctxTypeElem)
				pushCtxPtr := ctrl.Pointer() + pushCtxOffset
				ctxPtr := (*PushCtx)(unsafe.Pointer(pushCtxPtr))
				return &PushCtrlValue{
					ctrl:   ctrl,
					ctxPtr: ctxPtr,
				}
			},
		}

		handleFunc = func(ctx *handlerCtx, argValue reflect.Value) {
			obj := pool.Get().(*PushCtrlValue)
			*obj.ctxPtr = ctx
			rets := cValue.Call([]reflect.Value{obj.ctrl, argValue})
			ctx.stat = (*status.Status)(unsafe.Pointer(rets[0].Pointer()))
			pool.Put(obj)
		}
	}

	if pluginContainer == nil {
		pluginContainer = newPluginContainer()
	}
	return []*Handler{{
		name:            globalServiceMethodMapper(prefix, handlerFuncName(cValue)),
		handleFunc:      handleFunc,
		argElem:         argType.Elem(),
		pluginContainer: pluginContainer,
	}}, nil
}

var (
	typeOfCallCtx = reflect.TypeOf((*CallCtx)(nil)).Elem()
	typeOfPushCtx = reflect.TypeOf((*PushCtx)(nil)).Elem()
)

//判断方法是否属于 CallCtx
func isBelongToCallCtx(name string) bool {
	for m := 0; m < typeOfCallCtx.NumMethod(); m++ {
		if name == typeOfCallCtx.Method(m).Name {
			return true
		}
	}
	return false
}

func isBelongToPushCtx(name string) bool {
	for m := 0; m < typeOfPushCtx.NumMethod(); m++ {
		if name == typeOfPushCtx.Method(m).Name {
			return true
		}
	}
	return false
}

func isStatusType(s string) bool {
	return strings.HasPrefix(s, "*") && strings.HasSuffix(s, ".Status")
}

func ctrlStructName(cType reflect.Type) string {
	split := strings.Split(cType.String(), ".")
	return split[len(split)-1]
}

func handlerFuncName(v reflect.Value) string {
	str := objectName(v)
	split := strings.Split(str, ".")
	return split[len(split)-1]
}

func objectName(v reflect.Value) string {
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

// IsExportedOrBuiltinType 判断类型是否属于可以被外部调用的或者是基本类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExportedName(t.Name()) || t.PkgPath() == ""
}

// IsExportedName is this an exported - upper case - name?
func isExportedName(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}
