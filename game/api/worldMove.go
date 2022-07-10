package api

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"zinxGame/game/core"
	"zinxGame/game/pb"
	"zinxGame/zinx/ziface"
	"zinxGame/zinx/znet"
)

type WorldMoveApi struct {
	znet.BaseRouter
}

func (m *WorldMoveApi) Handle(request ziface.IRequest)  {
	pos := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), pos)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error: ", err)
		return
	}

	pid, err := request.GetConn().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error, ", err)
		return
	}

	fmt.Printf("Player pid = %d, move(%f, %f, %f, %f)\n", pid, pos.X, pos.Y, pos.Z, pos.V)

	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.UpdatePos(pos.X, pos.Y, pos.Z, pos.V)
}