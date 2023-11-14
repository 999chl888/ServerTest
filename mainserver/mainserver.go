package main

import (
	"context"
	"errors"
	"fmt"
	protocol "gameserver/protocol"
	"gameserver/util"
	"log"
	"strconv"
)

type MainServer struct {
	serverIP string
	protocol.UnimplementedMainGameServer
	roomlist map[string]Room
}

func (s *MainServer) InitServer() error {
	ClearRoomIndex()
	s.roomlist = make(map[string]Room)
	ip, err := util.GetLocalIP()
	fmt.Printf("Server ip: %v\n", ip)
	if err != nil {
		log.Fatalf("get local ip failed")
	}
	s.serverIP = ip
	return nil
}

func (s *MainServer) CreateRoom(ctx context.Context, in *protocol.CreateRoomRequest) (*protocol.CreateRoomResponse, error) {
	response := new(protocol.CreateRoomResponse)
	room_id := GenerateRoomIndex()
	if room_id >= 100000 {
		ClearRoomIndex()
		room_id = 0
	}

	userstatemp := make(map[string]UserState)
	userstatemp[in.User.ID] = UnReady
	response.RoomID = strconv.Itoa(room_id)
	s.roomlist[response.RoomID] = Room{
		ID:        response.RoomID,
		MapID:     "000",
		Capacity:  8,
		OwnerID:   in.User.ID,
		UserState: userstatemp,
		roomState: RoomState(Room_State_Waiting),
	}

	response.Ret = true
	fmt.Printf("Room %v created successfully", room_id)
	return response, nil
}

func (s *MainServer) JoinRoom(ctx context.Context, in *protocol.JoinRoomRequest) (*protocol.JoinRoomResponse, error) {
	response := new(protocol.JoinRoomResponse)
	room, ok := s.roomlist[in.RoomID]
	if !ok {
		response.Ret = false
		response.Message = fmt.Sprintf("Undefined RoomID %v", in.RoomID)
		return response, nil
	}
	if room.roomState != RoomState(Room_State_Waiting) {
		response.Ret = false
		response.Message = "The room has already started game"
		return response, nil
	}
	_, ok = s.roomlist[in.RoomID].UserState[in.User.ID]
	if ok {
		response.Ret = false
		response.Message = fmt.Sprintf("Player %v is already in the room %v", in.User.ID, in.RoomID)
		return response, nil
	}
	room.UserState[in.User.ID] = UnReady
	response.Ret = true
	response.Message = ""
	return response, nil
}

func (s *MainServer) GetRoomState(ctx context.Context, in *protocol.GetRoomStateRequest) (*protocol.GetRoomStateResponse, error) {
	response := new(protocol.GetRoomStateResponse)
	room, ok := s.roomlist[in.RoomID]
	if !ok {
		response.Ret = false
		return response, errors.New(fmt.Sprintf("Undefined RoomID %v", in.RoomID))
	}
	response.Ret = true
	response.Mapid = room.MapID
	response.Capacity = int64(room.Capacity)
	response.GameInstanceID = room.GameInstance
	response.GameserverIP = room.GameServerIP
	response.GameserverPort = int64(room.GameServerPort)
	for userid, state := range room.UserState {
		response.User = append(response.User, &protocol.User{ID: userid, State: ConverteProtocalState2UserState(state)})
	}
	return response, nil
}

func ConverteProtocalState2UserState(userstate UserState) protocol.UserState {
	//这里没有直接调用转换，避免后续protol与userstate映射不一致产生bug，这里在增加用户状态后需要改动
	return protocol.UserState(userstate)
}

func (s *MainServer) StartGame(ctx context.Context, in *protocol.StartGameRequest) (*protocol.StartGameResponse, error) {
	response := new(protocol.StartGameResponse)
	preroom, ok := s.roomlist[in.RoomID]
	if !ok {
		response.Ret = false
		response.Message = "Undifined Room ID, Please recreate your room"
	}
	if in.UserID != preroom.OwnerID {
		response.Ret = false
		response.Message = "The user requesting does not have START-GAME permission"
	}

	serverip, serverport := s.GetSubServerDir()

	s.roomlist[in.RoomID] = Room{
		UserState:      preroom.UserState,
		Capacity:       preroom.Capacity,
		OwnerID:        preroom.OwnerID,
		roomState:      RoomState(Room_State_Loading),
		GameInstance:   in.RoomID + "_GID",
		GameServerIP:   serverip,
		GameServerPort: serverport,
		MapID:          preroom.MapID,
	}

	response.GameinstanceID = in.RoomID + "_GID"
	response.Ret = true
	response.GameserverIP = serverip
	response.GameserverPort = int64(serverport)
	return response, nil
}

func (s *MainServer) GetSubServerDir() (string, int) {
	return "192.168.0.199", 8333
}
