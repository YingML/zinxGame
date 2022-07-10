package core

import "fmt"

const (
	AoiMinX int = 0
	AoiMaxX int = 400
	AoiCntX int = 10
	AoiMinY int = 0
	AoiMaxY int = 400
	AoiCntY int = 10
)

type AOIMap struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CountX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向格子的数量
	CountY int

	//当前区域中有哪些格子map-key=格子的ID， value=格子对象
	grids map[int]*Grid
}

/*
	初始化一个AOI区域管理模块
*/
func NewAOIMap(minX int, minY int, maxX int, maxY int, countX int, countY int) *AOIMap {
	aoiMap := &AOIMap{
		MinX:   minX,
		MaxX:   maxX,
		CountX: countX,
		MinY:   minY,
		MaxY:   maxY,
		CountY: countY,
		grids:  make(map[int]*Grid),
	}

	width := aoiMap.gridWidth()
	length := aoiMap.gridLength()
	gid := 0
	for i := 0; i < countY; i++ {
		curY := length *i
		for j:=0; j<countX; j++ {
			curX := width*j
			grid := NewGrid(gid, curX, curX+width, curY, curY+length)
			aoiMap.grids[gid] = grid
			gid++
		}
	}

	return aoiMap
}

func (am *AOIMap) gridWidth() int {
	return (am.MaxX - am.MinX) / am.CountX
}

func (am *AOIMap) gridLength() int {
	return (am.MaxY - am.MinY) / am.CountY
}

func (am *AOIMap) String() string {
	return fmt.Sprintf("AOIManager: {minX: %d, minY: %d, maxX: %d, maxY: %d, countX: %d, countY: %d, grids: \n%+v}\n",am.MinX, am.MinY, am.MaxX, am.MaxY, am.CountX, am.MaxY, am.grids)
}

func (am *AOIMap) GetSurroundsGridsByGid(gid int) []*Grid {
	grid := am.grids[gid]
	var grids []*Grid

	grids = append(grids, grid)

	// gid所在行最左侧grid的id
	leftID := gid
	// 一行的个数
	lenX := 1

	// 右侧存在格子
	if grid.MaxX < am.MaxX {
		grids = append(grids, am.grids[gid+1])
		lenX++
	}
	// 左侧存在格子，设置leftID为最左侧grid的id
	if grid.MinX > am.MinX {
		grids = append(grids, am.grids[gid-1])
		leftID = gid-1
		lenX++
	}

	// 上面仍存在格子，则将上方的格子全部加入
	if grid.MinY > am.MinY {
		for i := 0; i < lenX; i++ {
			grids = append(grids, am.grids[leftID+i-am.CountX])
		}
	}

	// 下面仍存在格子，则将下方的格子全部加入
	if grid.MaxY < am.MaxY {
		for i := 0; i < lenX; i++ {
			grids = append(grids, am.grids[leftID+i+am.CountX])
		}
	}
	return grids
}

func (am *AOIMap) GetPidByPos(x, y float32) (playerIDs []int32)  {
	gridID := am.GetGidByPos(x, y)
	grids := am.GetSurroundsGridsByGid(gridID)
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		fmt.Printf("===> gridID: %d, playerIDs: %+v <===\n", grid.GID, grid.GetPlayerIDs())
	}
	return
}

func (am *AOIMap) GetGidByPos(x, y float32) int {
	idx := int(x) / am.gridWidth()
	idy := int(y) / am.gridLength()
	return idy * am.CountX + idx
}

func (am *AOIMap) AddPidToGrid(pid int32, gid int)  {
	am.grids[gid].Add(pid)
}

func (am *AOIMap) RemovePidFromGrid(pid int32, gid int) {
	am.grids[gid].Remove(pid)
}

func (am *AOIMap) GetPidsFromGid(gid int) (playerIDs []int32) {
	return am.grids[gid].GetPlayerIDs()
}

func (am *AOIMap) AddToGridByPos(pid  int32, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.AddPidToGrid(pid, gid)
}

func (am *AOIMap) RemoveFromGridByPos(pid  int32, x, y float32) {
	gid := am.GetGidByPos(x, y)
	am.RemovePidFromGrid(pid, gid)
}