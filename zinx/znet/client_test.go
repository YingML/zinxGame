package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	fmt.Println("client start ...")
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的Message
		dp := NewDataPack()
		m := NewMessage(2,[]byte("Zinx v0.5 client Test Message"))
		data, err := dp.Pack(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := conn.Write(data); err != nil {
			fmt.Println(err)
			return
		}

		binaryHead := make([]byte, dp.headerSize)
		if _, err = io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("client recv data failed: ", err)
			return
		}

		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack data failed: ", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*Message)
			msg.data = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(conn, msg.data); err != nil {
				fmt.Println("client recv data failed: ", err)
				return
			}
			fmt.Printf("client recv data success. ID: %d, len: %d, data: %+v\n", msg.id, msg.dataLen, string(msg.data))
		}
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}