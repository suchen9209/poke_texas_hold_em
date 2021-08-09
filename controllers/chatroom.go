// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"container/list"
	"encoding/json"
	"time"

	"poke/models"

	"github.com/beego/beego/v2/core/logs"

	"github.com/gorilla/websocket"
)

type Subscription struct {
	Archive []models.Event      // All the events from the archive.
	New     <-chan models.Event // New events coming in.
}

func newEvent(ep models.EventType, user, msg string) models.Event {
	return models.Event{ep, user, int(time.Now().Unix()), msg}
}

func sendCard(ep models.EventType, user string, position int, card models.Card) models.CardInfo {
	return models.CardInfo{
		Type:      ep,
		User:      user,
		Position:  position,
		Timestamp: int(time.Now().Unix()),
		Card:      card,
	}
}

func Join(name string, ws *websocket.Conn, ut models.UserType, user models.User) {
	subscribe <- Subscriber{Name: name, Conn: ws, UserType: ut, User: user}
}

func Leave(user models.User) {
	unsubscribe <- user
}

type Player struct {
	user     models.User
	Position int
	Conn     *websocket.Conn // Only for WebSocket users; otherwise nil.
}

type Subscriber struct {
	Name     string
	Conn     *websocket.Conn // Only for WebSocket users; otherwise nil.
	UserType models.UserType
	User     models.User
}

type UserInfoMsg struct {
	Type models.EventType
	Info []models.UserPointSeat
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 1000)
	// Channel for exit users.
	unsubscribe = make(chan models.User, 1000)
	// Send events here to publish them.
	publish         = make(chan models.Event, 1000)
	gameprocess     = make(chan models.CardInfo, 1000)
	userInfoChannel = make(chan models.Event, 1000)

	gameop = make(chan string)

	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()

	seat = list.New()
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case op := <-gameop:
			if op == "start" {
				userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
				models.InitCardMap()
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					ws := ss.Value.(Player).Conn
					if ws != nil {
						tmp_card := models.CardInfo{
							Type: models.EVENT_CLEAR_CARD,
						}
						msg, _ := json.Marshal(tmp_card)
						if ws.WriteMessage(websocket.TextMessage, msg) != nil {
							// User disconnected.
							unsubscribe <- ss.Value.(Player).user
						}
					}

				}
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					var tmp_poker = models.GetOneCard()
					gameprocess <- sendCard(models.EVENT_LICENSING, ss.Value.(Player).user.Name, ss.Value.(Player).Position, *tmp_poker)

				}
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					var tmp_poker = models.GetOneCard()
					gameprocess <- sendCard(models.EVENT_LICENSING, ss.Value.(Player).user.Name, ss.Value.(Player).Position, *tmp_poker)
				}

				for i := 1; i <= 5; i++ {
					var tmp_poker = models.GetOneCard()
					models.PublicCard[i] = *tmp_poker
				}

			}
			if op == "show_card3" {
				gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[1])
				gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[2])
				gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[3])
			}
			if op == "show_card4" {
				gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[4])
			}
			if op == "show_card5" {
				gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[5])
			}
		case process := <-gameprocess:
			if process.Type == models.EVENT_LICENSING {
				var other_process = process
				other_process.Card = models.Card{}
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					ws := ss.Value.(Player).Conn
					if ss.Value.(Player).user.Name == process.User {
						if ws != nil {
							msg, _ := json.Marshal(process)
							if ws.WriteMessage(websocket.TextMessage, msg) != nil {
								// User disconnected.
								unsubscribe <- ss.Value.(Player).user
							}
						}
					} else {
						if ws != nil {
							msg, _ := json.Marshal(other_process)
							if ws.WriteMessage(websocket.TextMessage, msg) != nil {
								// User disconnected.
								unsubscribe <- ss.Value.(Player).user
							}
						}
					}
				}
			}
			if process.Type == models.EVENT_PUBLIC_CARD {
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					ws := ss.Value.(Player).Conn
					if ws != nil {
						msg, _ := json.Marshal(process)
						if ws.WriteMessage(websocket.TextMessage, msg) != nil {
							// User disconnected.
							unsubscribe <- ss.Value.(Player).user
						}
					}

				}
			}
		case userInfo := <-userInfoChannel:
			logs.Info(userInfo)
			if userInfo.Type == models.EVENT_REFRESH_USER_INFO {
				var data UserInfoMsg
				data.Type = userInfo.Type
				data.Info = models.GetUserPointWithSeat(1)
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					ws := ss.Value.(Player).Conn
					if ws != nil {
						msg, _ := json.Marshal(data)
						if ws.WriteMessage(websocket.TextMessage, msg) != nil {
							// User disconnected.
							unsubscribe <- ss.Value.(Player).user
						}
					}

				}
			}
		case sub := <-subscribe:
			if !isUserExist(subscribers, sub.Name) {
				subscribers.PushBack(sub) // Add user to the end of list.
				if sub.UserType == models.POKER_PLAYER {
					pos := models.SetUserReturnPlayer(sub.User)
					var player = new(Player)
					player.user = sub.User
					player.Position = pos
					player.Conn = sub.Conn
					seat.PushBack(*player)
				}

				// Publish a JOIN event.
				publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
				logs.Info("New user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			} else {
				logs.Info("Old user:", sub.Name, ";WebSocket:", sub.Conn != nil)
			}
		case event := <-publish:
			// Notify waiting list.
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			broadcastWebSocket(event)
			models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				logs.Info("Message from", event.User, ";Content:", event.Content)
			}
		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).User == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
						logs.Error("WebSocket closed:", unsub)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub.Name, "") // Publish a LEAVE event.
					break
				}
			}
			for ss := seat.Front(); ss != nil; ss = ss.Next() {
				if ss.Value.(Player).user == unsub {
					seat.Remove(ss)
					models.RemoveGameUser(unsub.Id, 1)
					// Clone connection.
					ws := ss.Value.(Player).Conn
					if ws != nil {
						ws.Close()
						logs.Error("WebSocket closed:", unsub.Name)
					}
					publish <- newEvent(models.EVENT_LEAVE, unsub.Name, "") // Publish a LEAVE event.
					break
				}
			}
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(subscribers *list.List, user string) bool {
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if sub.Value.(Subscriber).Name == user {
			return true
		}
	}
	return false
}
