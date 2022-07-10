package znet

import "C"
import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"zinxGame/zinx/ziface"
)

type Conn struct {
	server   ziface.IServer
	conn     net.Conn
	id       uint32
	isClosed bool
	exitChan chan bool
	handler  ziface.IMsgHandler
	msgChan  chan ziface.IMessage

	// 连接属性集合
	property map[string]interface{}
	// 保护连接属性的锁
	propertyLock sync.RWMutex
}

func (c *Conn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Conn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if v, ok := c.property[key]; ok {
		return v, nil
	}else {
		return nil, errors.New("no property found")
	}
}

func (c *Conn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func NewConn(s ziface.IServer,c net.Conn, id uint32, h ziface.IMsgHandler) *Conn {
	return &Conn{
		server:   s,
		conn:     c,
		id:       id,
		isClosed: false,
		exitChan: make(chan bool, 1),
		handler:  h,
		msgChan:  make(chan ziface.IMessage, 10),
		property: make(map[string]interface{}),
	}
}

func (c *Conn) StartReader()  {
	defer c.server.GetConnManager().Remove(c)
	for {
		// 创建1个拆解包对象
		dp := c.server.GetDataPack()
		// 读取客户端的Msg Head 二进制流8个字节

		headData := make([]byte, dp.GetHeaderSize())
		if _, err := io.ReadFull(c.Conn(), headData); err != nil {
			fmt.Println("read msg head error", err)
			c.exitChan <- true
			break
		}
		//拆包， 得到MsgID和msgData, 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据dataLen，再次读取Data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn(), data); err != nil {
				fmt.Println("read msg data error", err)
			}
		}
		msg.SetData(data)
		reqID := rand.Uint32()
		req := NewRequest(reqID, c, msg)
		// Reader只负责从连接中读取数据并构建出IRequest供开发人员在实现Handler时处理
		if c.handler.GetWorkPool() != nil {
			c.handler.GetWorkPool().SendToTask(req)
		}else {
			c.handler.Serve(req)
		}
	}
}

func (c *Conn) StartWriter()  {
	for {
		select {
		case msg := <- c.msgChan:
			if msg == nil {
				if c.isClosed {
					fmt.Printf("[WARN] msg is nil! connection %d closed.", c.id)
					return
				}
				fmt.Printf("[ERROR] msg is nil, but connection %d still active.", c.id)
				continue
			}
			dp := c.server.GetDataPack()
			data, err := dp.Pack(msg)
			if err != nil {
				fmt.Println("Pack message failed: ", err.Error())
				continue
			}
			err = binary.Write(c.conn, binary.LittleEndian, data)
			if err != nil {
				fmt.Println("Write message failed: ", err.Error())
				continue
			}
		case exist := <- c.exitChan:
			if exist {
				return
			}
		}
	}
}

func (c *Conn) Start() {
	fmt.Println("Conn Start().. ConnID=", c.id)
	c.server.GetConnManager().Add(c)
	c.server.CallOnConnStart(c)
	go c.StartReader()

	if c.handler.GetWorkPool() == nil {
		go c.StartWriter()
	}
}

func (c *Conn) Stop() {
	c.server.CallOnConnStop(c)
	fmt.Println("Conn Stop().. ConnID=", c.id)

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	c.exitChan <- true
	// 关闭连接
	c.conn.Close()

	// 回收channel
	close(c.exitChan)
	close(c.msgChan)
}

func (c *Conn) Conn() net.Conn {
	return c.conn
}

func (c *Conn) ConnID() uint32 {
	return c.id
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) SendMsg(id uint32, data []byte) {
	msg := NewMessage(id, data)
	if c.handler.GetWorkPool() == nil {
		c.msgChan <- msg
		return
	}

	dp := c.server.GetDataPack()
	data, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack message failed: ", err.Error())
		return
	}
	err = binary.Write(c.conn, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Write message failed: ", err.Error())
		return
	}
}

type HandleFunc func(net.Conn, []byte, int) error

var _ ziface.IConnection = &Conn{}