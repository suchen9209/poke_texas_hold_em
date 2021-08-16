package models

import (
	"container/list"
)

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_LICENSING           //发牌
	EVENT_PUBLIC_CARD         //公共牌
	EVENT_CLEAR_CARD          //清理牌桌
	EVENT_REFRESH_USER_INFO   //更新用户信息
	EVENT_ROUND_INFO          //回合信息
	EVENT_USER_OPERATION_INFO //用户操作信息
	EVENT_GAME_END            //游戏结束
	EVENT_GAME_END_SHOW_CARD  //游戏结束判定卡牌信息
)

type UserType int

const (
	POKER_PLAYER = iota
	VIEWER
)

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	Timestamp int // Unix timestamp (secs)
	Content   string
}

type SeatInfo struct {
	Type      EventType // JOIN, LEAVE
	User      string
	Timestamp int // Unix timestamp (secs)
	GameUser  GameUser
}

type CardInfo struct {
	Type      EventType
	User      string
	Position  int
	Timestamp int // Unix timestamp (secs)
	Card      Card
}

type RoundInfo struct {
	Type            EventType
	GM              GameMatch
	NowPosition     int
	AllPointInRound int
	MaxPoint        int
	Detail          interface{}
}

type ClientMessage struct {
	Message   string
	Type      string
	Position  int
	UserId    int
	Point     int
	Operation string
	Name      string
}

const archiveSize = 20

// Event archives.
var archive = list.New()

// NewArchive saves new event to archive list.
func NewArchive(event Event) {
	if archive.Len() >= archiveSize {
		archive.Remove(archive.Front())
	}
	archive.PushBack(event)
}

// GetEvents returns all events after lastReceived.
func GetEvents(lastReceived int) []Event {
	events := make([]Event, 0, archive.Len())
	for event := archive.Front(); event != nil; event = event.Next() {
		e := event.Value.(Event)
		if e.Timestamp > int(lastReceived) {
			events = append(events, e)
		}
	}
	return events
}
