package socket

import "donkeygo/container/dmap"

//socket 仓库
type HubSocket struct {
	// key: socket id (ip, name and so on)
	// value: Socket
	sockets *dmap.Map
}

//创建仓库
func NewHubSocket() *HubSocket {
	chub := &HubSocket{
		sockets: dmap.New(true),
	}
	return chub
}

//添加socket
func (that *HubSocket) Set(socket Socket) {
	_socket, ok := that.sockets.Search(socket.ID())
	if !ok {
		that.sockets.Set(socket.ID(), socket)
		return
	}
	if oldSocket := _socket.(Socket); socket != oldSocket {
		_ = oldSocket.Close()
	}
}

//获取socket
func (that *HubSocket) Get(id string) (Socket, bool) {
	_socket := that.sockets.Get(id)
	if _socket == nil {
		return nil, false
	}
	return _socket.(Socket), true
}

//遍历socket
func (that *HubSocket) Range(f func(Socket) bool) {
	that.sockets.Iterator(func(key, value interface{}) bool {
		return f(value.(Socket))
	})
}

//长度
func (that *HubSocket) Len() int {
	return that.sockets.Size()
}

//删除socket
func (that *HubSocket) Delete(id string) {
	that.sockets.Remove(id)
}

//修改socket的id
func (that *HubSocket) ChangeID(newID string, socket Socket) {
	oldID := socket.ID()
	socket.SetID(newID)
	that.Set(socket)
	if oldID != socket.RemoteAddr().String() {
		that.Delete(oldID)
	}
}
