package ziface

type IServer interface {
	Start()
	Serve()
	Stop()
	GetConnManager() IConnManager
	SetOnConnStart(func(conn IConnection))
	SetOnConnStop(func(conn IConnection))
	CallOnConnStart(conn IConnection)
	CallOnConnStop(conn IConnection)
	GetDataPack() IDataPack
}
