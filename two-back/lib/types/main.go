package types

// Room is defintion of rooms table item
type Room struct {
	RoomID  string
	Status  string
	User1ID string
	User2ID string
}

// User is defintion of users table item
type User struct {
	ConnectionID string
	RoomID       string
}
