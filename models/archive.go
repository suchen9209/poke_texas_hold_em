package models

import (
	"container/list"
)

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
	EVENT_LICENSING
	EVENT_PUBLIC_CARD
	EVENT_CLEAR_CARD
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

type CardInfo struct {
	Type      EventType
	User      string
	Position  int
	Timestamp int // Unix timestamp (secs)
	Card      Card
}

type ClientMessage struct {
	Message string
	Type    string
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
