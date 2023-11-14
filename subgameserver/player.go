package main

type playerState int32

const (
	Player_Alive playerState = 0
	Player_Dead  playerState = 1
)

type Player struct {
	UserID       string
	Health       int
	Itemsid      []string
	Professionid string
	PositionX    int
	PositionY    int
	Action       string
	PlayerState  playerState
}
