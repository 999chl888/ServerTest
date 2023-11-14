package main

import (
	"fmt"
	"log"
	"net"
	"os"

	protocol "gameserver/protocol"

	"google.golang.org/grpc"
)

func main() {
	var gameserver MainServer
	err := gameserver.InitServer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8555))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rpcServer := grpc.NewServer()
	protocol.RegisterMainGameServer(rpcServer, &gameserver)
	log.Println("Init success...")
	if err = rpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
