package models

type Room struct {
	Id           int    `form:"-"`
	RoomName     string `form:"room_name"`
	RoomPassword string `form:"room_password"`
}
