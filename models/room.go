package models

type Room struct {
	Id           int `form:"-"`
	CreateUserId int
	RoomName     string `form:"room_name"`
	RoomPassword string `form:"room_password"`
	CardType     string `form:"room_card_type"`
}

const RoomShortType = "short"
const RoomLongType = "long"

func CreateRoom(r *Room) int64 {
	insert, err := o.Insert(r)
	if err != nil {
		return 0
	}
	return insert
}

func GetOnlineRoom() []Room {
	s := o.QueryTable("room")
	var RoomList []Room
	_, _ = s.All(&RoomList)
	return RoomList
}

func CloseRoom(roomId int) {
	o.QueryTable("room").Filter("id", roomId).Delete()
}
