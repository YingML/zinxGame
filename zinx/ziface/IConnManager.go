package ziface

type IConnManager interface {
	Add(conn IConnection)
	GetConn(connID uint32) (IConnection, error)
	Len() int
	Remove(conn IConnection)
	RemoveAll()
}
