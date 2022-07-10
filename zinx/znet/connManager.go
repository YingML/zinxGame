package znet

import "C"
import (
	"errors"
	"fmt"
	"sync"
	"zinxGame/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	lock sync.RWMutex
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	//将Conn加入到ConnManager中
	c.connections[conn.ConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: cpnn num = ", c.Len())
}

func (c *ConnManager) GetConn(connID uint32) (ziface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if conn, ok :=  c.connections[connID]; ok {
		return conn, nil
	}else {
		return nil, errors.New("connection not found！")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()

	conn.Stop()
	delete(c.connections, conn.ConnID())
}

func (c *ConnManager) RemoveAll() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for connID, conn := range c.connections{
		conn.Stop()
		delete(c.connections, connID)
	}
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

var _ ziface.IConnManager = &ConnManager{}