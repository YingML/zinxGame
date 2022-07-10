package ziface

type IDataPack interface {
	Pack(msg IMessage) ([]byte, error)
	UnPack(bs []byte)(IMessage, error)
	GetHeaderSize() uint32
}
