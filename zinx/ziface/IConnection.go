package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	Conn() net.Conn
	ConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(uint32, []byte)
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}
