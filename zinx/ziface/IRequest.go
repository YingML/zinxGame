package ziface

type IRequest interface {
	GetConn() IConnection
	GetReqID() uint32
	GetData() []byte
	GetMsgID() uint32
}
