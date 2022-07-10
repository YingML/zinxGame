package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"sync"
	"zinxGame/game/pb"
	"zinxGame/zinx/ziface"
)

//玩家对象
type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的连接(用于和客户端连接)
	X    float32
	Y    float32	//高度
	Z    float32
	V    float32	//旋转角度(0-360)
}

/*
	Player ID 生成器
 */
var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	//生成1个玩家ID
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	return &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}
}

// 提供一个发送给客户端消息的方法
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	msgData, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err : ", err)
		return
	}
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}
	p.Conn.SendMsg(msgID, msgData)
}

// 告知客户端玩家Pid, 同步已经上传的玩家ID给客户端
func (p *Player) SyncPid()  {
	// 创建MsgID: 1 的proto的数据
	protoMsg := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, protoMsg)
}

//向所有玩家广播自己地点
func (p *Player) SyncPosition()  {
	// 创建MsgID: 200 的proto的数据
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	players := p.GetSurroundPlayers()
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// 向周边玩家广播消息
func (p *Player) Talk(content string)  {
	protoMsg := &pb.BroadCast{
		Pid: p.Pid,
		Tp: 1, //tp -1 代表聊天广播
		Data: &pb.BroadCast_Content{Content: content},
	}
	players := WorldMgr.GetAllPlayers()
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// 将其他玩家位置广播给当前玩家客户端
func (p *Player) SyncPlayers() {
	players := p.GetSurroundPlayers()
	for _, player := range players {
		// 创建MsgID: 200 的proto的数据
		protoMsg := &pb.BroadCast{
			Pid: player.Pid,
			Tp: 2,
			Data: &pb.BroadCast_P{
				P: &pb.Position{
					X: player.X,
					Y: player.Y,
					Z: player.Z,
					V: player.V,
				},
			},
		}
		// 向当前玩家广播其他玩家位置信息
		p.SendMsg(200, protoMsg)
	}
}

// 获取当前玩家的周边玩家
func (p *Player) GetSurroundPlayers() []*Player {
	var players []*Player
	gid := WorldMgr.AoiMap.GetGidByPos(p.X, p.Z)
	grids := WorldMgr.AoiMap.GetSurroundsGridsByGid(gid)
	for _, grid := range grids {
		for _, pid := range grid.GetPlayerIDs() {
			players = append(players, WorldMgr.GetPlayerByPid(pid))
		}
	}
	return players
}

// 更新自己的位置并广播给周边玩家
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	p.SyncPosition()
}

// 向周边玩家广播当前玩家下线
func (p *Player) SyncOffLine() {
	protoMsg := &pb.SyncPid{Pid: p.Pid}

	players := p.GetSurroundPlayers()
	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}

	WorldMgr.RemovePlayerByPid(p.Pid)
	fmt.Printf("===> Player Pid: %d OffLine <===\n", p.Pid)
}