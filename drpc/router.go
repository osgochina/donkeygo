package drpc

import "reflect"

type Handler struct {
	// 绑定的handler名
	name string

	//参数类型
	argElem reflect.Type

	//返回值类型 注意：只有call消息才会有
	reply reflect.Type

	//处理该消息的方法
	handleFunc func(*handlerCtx, reflect.Type)

	// 不能匹配到绑定方法时，默认的处理方法
	unknownHandleFunc func(*handlerCtx)

	// 路由类型名字
	routerTypeName string

	//是否找到绑定方法
	isUnknown bool
}
