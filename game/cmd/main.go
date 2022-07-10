package main

import (
	"flag"
	"fmt"
	"path"
	"runtime"
	"zinxGame/game/api"
	"zinxGame/game/core"
	"zinxGame/game/utils"
	"zinxGame/zinx/ziface"
	"zinxGame/zinx/znet"
)

func init() {
	var abPath string
	if _, filename, _, ok := runtime.Caller(0); ok {
		abPath = path.Dir(filename)
	}
	filePath := path.Join(path.Dir(abPath), "conf", "config.yaml")
	utils.ConfigPath = flag.String("config", filePath, "-c")
}

func OnConnectionAdd(conn ziface.IConnection)  {
	player := core.NewPlayer(conn)
	// 同步玩家的ID给客户端
	player.SyncPid()
	conn.SetProperty("pid", player.Pid)
	// 将玩家添加到地图中
	core.WorldMgr.AddPlayer(player)
	// 将玩家上线信息广播给所有玩家客户端
	player.SyncPosition()
	// 将其他玩家信息广播给当前玩家客户端
	player.SyncPlayers()

	fmt.Println(" ===> Player pid = ", player.Pid, "is arrived <===")
}

func OnConnectionRemove(conn ziface.IConnection) {
	pid, err := conn.GetProperty("pid")
	if err != nil {
		fmt.Println("Get Property pid err:", err)
		return
	}

	player := core.WorldMgr.GetPlayerByPid(pid.(int32))
	player.SyncOffLine()
}

func main()  {
	flag.Parse()
	utils.InitConfig()
	wp := znet.NewWorkPool(
		znet.WithTaskQueueSize(utils.Config.MaxTaskQueueSize),
		znet.WithWorkPoolSize(utils.Config.MaxWorkPoolSize),
	)

	mh := znet.NewMsgHandler(znet.WithWorkPool(wp))
	mh.AddRouter(2, &api.WorldChatApi{})
	mh.AddRouter(3, &api.WorldMoveApi{})

	s := znet.NewServer(mh,
		znet.WithServerName(utils.Config.Name),
		znet.WithIP(utils.Config.IP),
		znet.WithPort(utils.Config.Port),
		znet.WithMaxPackageSize(utils.Config.MaxPackageSize),
		znet.WithMaxConnection(utils.Config.MaxConnection),
	)
	//连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionRemove)
	s.Serve()
}
