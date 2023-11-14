package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gameserver/protocol"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	option := flag.String("mode", "", "Mode to test. It can be [createroom, joinroom, getroomstate, startgame]")
	userID := flag.String("userid", "testid", "Test userid")
	ip := flag.String("ip", "127.0.0.1:8555", "Main server IP")
	roomID := flag.String("roomid", "1", "Test room ID")
	instanceID := flag.String("instanceid", "1_GID", "Test Game instance ID")
	flag.Parse()
	switch *option {
	case "createroom":
		err := CreateRoom(*userID, *ip)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	case "joinroom":
		err := JoinRoom(*userID, *ip, *roomID)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	case "getroomstate":
		err := GetRoomState(*ip, *roomID)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	case "startgame":
		err := StartGame(*userID, *roomID, *ip)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	case "creategame":
		err := StartGameWithInstance(*instanceID, *ip, *roomID)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	case "Getgamestate":
		err := GetGameStateWithGameinstance(*instanceID, *ip)
		if err != nil {
			print(err.Error())
			os.Exit(1)
		}
	}

}

func CreateRoom(userid string, ip string) error {
	//var taskinfo TaskInfo
	var request protocol.CreateRoomRequest
	//fmt.Printf("request.CIDs: %v\n", request.CIDs)
	request.User = &protocol.User{
		ID: userid,
	}
	request.MapID = "testmap0"
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewMainGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.CreateRoom(ctx, &request)
	if err != nil {
		return err
	}
	if !response.Ret {
		return errors.New("Error occured while creating Room")
	}
	fmt.Printf("Create room success!,roomID is %v\n", response.RoomID)
	return nil
}

func JoinRoom(userid string, ip string, roomid string) error {
	//var taskinfo TaskInfo
	var request protocol.JoinRoomRequest
	//fmt.Printf("request.CIDs: %v\n", request.CIDs)
	request.RoomID = roomid
	request.User = &protocol.User{
		ID: userid,
	}
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewMainGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.JoinRoom(ctx, &request)
	if err != nil {
		return err
	}
	if !response.GetRet() {
		return errors.New("err")
	}
	fmt.Printf("Success to join the room %v\n", roomid)
	return nil
}

func GetRoomState(ip string, roomid string) error {
	//var taskinfo TaskInfo
	var request protocol.GetRoomStateRequest
	//fmt.Printf("request.CIDs: %v\n", request.CIDs)
	request.RoomID = roomid
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewMainGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.GetRoomState(ctx, &request)
	if err != nil {
		return err
	}
	if !response.GetRet() {
		return errors.New("err")
	}

	fmt.Printf("The map of the room is %v\n", response.Mapid)
	fmt.Printf("The room's Capacity is %v\n", response.Capacity)
	fmt.Printf("The room's game instance would be on %v:%v\n", response.GameserverIP, response.GameserverPort)
	fmt.Printf("The room's game instanceID is %v\n", response.GameInstanceID)
	fmt.Println("User list ======================")

	for user, state := range response.User {
		fmt.Printf("===user: %v state %v===\n", user, state)
	}
	return nil
}

func StartGame(userid string, roomid string, ip string) error {
	var request protocol.StartGameRequest
	request.RoomID = roomid
	request.UserID = userid
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewMainGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.StartGame(ctx, &request)
	if err != nil {
		return err
	}
	if !response.GetRet() {
		return errors.New(response.Message)
	}
	fmt.Printf("Create game from room success! Game Server is %v:%v", response.GameserverIP, response.GameserverPort)
	return nil
}

func StartGameWithInstance(instanceid string, ip string, roomid string) error {
	var request protocol.CreateGameWithInstanceRequest
	request.RoomID = roomid
	request.GameInstanceID = instanceid
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewSubGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.CreateGameWithInstance(ctx, &request)
	if err != nil {
		return err
	}
	if !response.GetRet() {
		return errors.New(response.Message)
	}

	fmt.Printf("Create game with game instance success!")

	return nil
}

func GetGameStateWithGameinstance(instanceid string, ip string) error {
	var request protocol.GetGameStateRequest
	request.GameInstanceID = instanceid
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	defer conn.Close()
	call := protocol.NewSubGameClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(10)*time.Second)
	defer cancel()

	response, err := call.GetGameState(ctx, &request)
	if err != nil {
		return err
	}
	if !response.GetRet() {
		fmt.Println("go here")
		return errors.New(response.Message)
	}

	fmt.Printf("Gamestate of the instance is %v\n", response.Gamestate)
	fmt.Println("====player state list====")
	for _, player := range response.Playerlist {
		fmt.Printf("Userid: %v\n", player.Userid)
		fmt.Printf("player: %v\n", player)
		fmt.Println("===========================")
	}

	return nil
}
