package main

import "sync"

type RoomState int32

const (
	Room_State_Waiting RoomState = 0
	Room_State_Loading RoomState = 1
	Room_State_InGame  RoomState = 2
)

type Room struct {
	roomState      RoomState
	ID             string
	MapID          string
	Capacity       int
	OwnerID        string
	UserState      map[string]UserState //to do 不应该把玩家状态放在这，应该专门放入玩家状态列表
	GameServerIP   string
	GameServerPort int
	GameInstance   string
}

var RoomIndexGenerater int
var RoomIndexGeneraterLock sync.Mutex

func GenerateRoomIndex() int {
	RoomIndexGeneraterLock.Lock()
	RoomIndexGenerater += 1
	RoomIndexGeneraterLock.Unlock()
	return RoomIndexGenerater
}

func ClearRoomIndex() {
	RoomIndexGeneraterLock.Lock()
	RoomIndexGenerater = 0
	RoomIndexGeneraterLock.Unlock()
}
