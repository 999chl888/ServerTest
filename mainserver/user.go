package main

type UserState int32

const (
	UnReady    UserState = 0
	Ready      UserState = 1
	Connecting UserState = 2
	InGame     UserState = 3
	Watching   UserState = 4
)

type User struct {
	UserID string
}
