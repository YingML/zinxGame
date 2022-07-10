package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinxGame/zinx/ziface"
)

const DefaultMaxPackageSize = 4096

type DataPack struct {
	dataLenSize    uint32
	msgIDSize      uint32
	headerSize     uint32
	maxPackageSize uint32
}

func NewDataPack() *DataPack {
	return &DataPack{
		dataLenSize:    4,
		msgIDSize:      4,
		headerSize:     4+4,
		maxPackageSize: DefaultMaxPackageSize,
	}
}

func (p *DataPack) Pack(msg ziface.IMessage) ([]byte, error)  {
	//创建1个存在bytes字节的缓冲
	dataBuff :=  bytes.NewBuffer([]byte{})
	//将dataLen写进缓冲中
	if err :=binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//将msgID写进缓冲中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	//将数据写进缓冲中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

func (p *DataPack) UnPack(bs []byte)(ziface.IMessage, error)  {
	m := &Message{}
	dataBuff := bytes.NewReader(bs)
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &m.dataLen); err != nil {
		return nil, err
	}
	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &m.id); err != nil {
		return nil, err
	}
	//判断dataLen是否超过定义的最大包大小
	if p.maxPackageSize > 0 && p.maxPackageSize < m.dataLen {
		return nil, errors.New("too large msg data recv")
	}
	return m, nil
}

func (p *DataPack) GetHeaderSize() uint32 {
	return p.headerSize
}