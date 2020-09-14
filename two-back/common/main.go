package common

// Room is defintion of the rooms table item
type Room struct {
	RoomID  string
	Status  string
	User1ID string
	User2ID string
}

// User is defintion of the users table item
type User struct {
	ConnectionID string
	RoomID       string
}

// MatchingQueueName is the queue name
var MatchingQueueName string = "matching"

// RoomTableName is name of the rooms table
var RoomTableName string = "rooms"

// UserTableName is name of the users table
var UserTableName string = "users"

// RoomStatusWaiting is status of the rooms table item
var RoomStatusWaiting string = "WAITING"

// RoomStatusPlaying is status of the rooms table item
var RoomStatusPlaying string = "PLAYING"
