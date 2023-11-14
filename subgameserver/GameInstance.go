package main

type GameState int32

const (
	Game_Loading GameState = 0
	Game_Playing GameState = 1
)

type GameInstance struct {
	Playerlist map[string]Player
	State      GameState
	RoomID     string
	StartTime  int64
}
