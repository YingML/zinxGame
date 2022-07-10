package znet

import (
	"fmt"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	dp := &DataPack{
		dataLenSize: 4,
		msgIDSize:   4,
		headerSize:  8,
	}
	msg := []byte("hello world")
	m1 := &Message{
		id:      1,
		dataLen: uint32(len(msg)),
		data: msg,
	}
	v, err := dp.Pack(m1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(v))
}

func TestDataPack_UnPack(t *testing.T) {
	dp := &DataPack{
		dataLenSize: 4,
		msgIDSize:   4,
		headerSize:  8,
	}
	msg := []byte("hello world")
	m1 := &Message{
		id:      uint32(1),
		dataLen: uint32(len(msg)),
		data: msg,
	}
	v, err := dp.Pack(m1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(v))
	m2, err := dp.UnPack(v)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(m2.GetData()))
}
