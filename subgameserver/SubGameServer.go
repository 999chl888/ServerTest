package main

import (
	"context"
	"errors"
	"fmt"
	protocol "gameserver/protocol"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubGameServer struct {
	gameInstanceMap map[string]GameInstance
	protocol.UnimplementedSubGameServer
}

func (s *SubGameServer) InitServer() error {
	s.gameInstanceMap = make(map[string]GameInstance)
	return nil
}

func (s *SubGameServer) CreateGameWithInstance(ctx context.Context, in *protocol.CreateGameWithInstanceRequest) (*protocol.CreateGameWithInstanceResponse, error) {
	response := new(protocol.CreateGameWithInstanceResponse)
	_, ok := s.gameInstanceMap[in.GameInstanceID]
	if ok {
		response.Ret = false
		response.Message = "Game already Created"
		return response, nil
	}

	s.gameInstanceMap[in.GameInstanceID] = GameInstance{
		Playerlist: nil,
		State:      Game_Loading,
		RoomID:     in.RoomID,
	}

	gameinfo, err := s.GetGameInfoFromMainServer(in.RoomID)
	if err != nil {
		response.Ret = false
		response.Message = "Fail to get Info From MainServer"
		return response, nil
	}

	s.gameInstanceMap[in.GameInstanceID] = *gameinfo
	response.Ret = true
	response.Message = "success"
	return response, nil
}

func (s *SubGameServer) GetGameInfoFromMainServer(roomid string) (*GameInstance, error) {
	var request protocol.GetRoomStateRequest
	//fmt.Printf("request.CIDs: %v\n", request.CIDs)
	request.RoomID = roomid
	ip := "192.168.0.199:8555"
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	call := protocol.NewMainGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.GetRoomState(ctx, &request)
	if err != nil {
		return nil, err
	}
	if !response.GetRet() {
		return nil, errors.New("Fail to get userinfo from mainserver")
	}
	fmt.Printf("len(response.User): %v\n", len(response.User))
	playerlist := make(map[string]Player)
	for i := 0; i < len(response.User); i++ {
		playerlist[response.User[i].ID] = Player{
			UserID:       response.User[i].ID,
			Health:       100,
			Professionid: "zhanshi",
			PositionX:    0,
			PositionY:    0,
			Action:       "StartAction",
			PlayerState:  Player_Alive,
			Itemsid:      []string{},
		}
	}

	gameinfo := &GameInstance{
		State:      Game_Playing,
		RoomID:     roomid,
		Playerlist: playerlist,
		StartTime:  time.Now().Unix(),
	}

	return gameinfo, nil
}

func (s *SubGameServer) GetGameState(ctx context.Context, in *protocol.GetGameStateRequest) (*protocol.GetGameStateResponse, error) {
	response := new(protocol.GetGameStateResponse)
	gameinstance := s.gameInstanceMap[in.GameInstanceID]
	timerecorder := (time.Now().Unix() - gameinstance.StartTime) % 40
	if timerecorder < 30 {
		response.Ret = false
		response.Message = "Present GetMsg is forbidden"
		return response, nil
	}

	response.Gamestate = s.ConvertGameState2Protocol(gameinstance.State)
	for _, player := range gameinstance.Playerlist {
		response.Playerlist = append(response.Playerlist, &protocol.PlayerInfo{
			Userid:       player.UserID,
			Health:       int64(player.Health),
			Itemsid:      player.Itemsid,
			Professionid: player.Professionid,
			PositionX:    int64(player.PositionX),
			PositionY:    int64(player.PositionY),
			Action:       player.Action,
			PlayerState:  s.ConvertePlayerState2Protocol(player.PlayerState),
		})
	}
	response.Ret = true
	return response, nil
}

func (s *SubGameServer) ConvertePlayerState2Protocol(state playerState) protocol.PlyerState {
	return protocol.PlyerState(state)
}

func (s *SubGameServer) ConverteProtocal2PlayerState(state protocol.PlyerState) playerState {
	return playerState(state)
}

func (s *SubGameServer) ConvertGameState2Protocol(state GameState) protocol.GameState {
	return protocol.GameState(state)
}

func (s *SubGameServer) PostGameState(ctx context.Context, in *protocol.PostGameStateRequest) (*protocol.PostGameStateResponse, error) {
	response := new(protocol.PostGameStateResponse)
	gameinstance := s.gameInstanceMap[in.GameInstanceID]
	timerecorder := (time.Now().Unix() - gameinstance.StartTime) % 40
	if timerecorder >= 30 {
		response.Ret = false
		response.Message = "Present PostMsg is forbidden"
		return response, nil
	}

	gameinstance.Playerlist[in.PlayerInfo.Userid] = Player{
		UserID:       in.PlayerInfo.Userid,
		Health:       int(in.PlayerInfo.Health),
		Itemsid:      in.PlayerInfo.Itemsid,
		Professionid: in.PlayerInfo.Professionid,
		PositionX:    int(in.PlayerInfo.PositionX),
		PositionY:    int(in.PlayerInfo.PositionY),
		Action:       in.PlayerInfo.Action,
		PlayerState:  s.ConverteProtocal2PlayerState(in.PlayerInfo.PlayerState),
	}

	return response, nil
}
