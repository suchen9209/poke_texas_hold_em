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
	publish         = make(chan models.SeatInfo, 1000)
	gameprocess     = make(chan models.CardInfo, 1000)
	userInfoChannel = make(chan models.Event, 1000)
	roundprocess    = make(chan models.RoundInfo, 1000)

	gameop = make(chan string)

	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()

	seat         = list.New()
	positionTurn int //记录当前轮到的玩家位置

	nowGameMatch models.GameMatch

	roundUserDetail = make(map[int]interface{})
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case op := <-gameop:
			if op == "start" {
				userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
				nowGameMatch = models.InitGameMatch(1)
				models.InitCardMap()
				nowGameMatch.GameStatus = "LICENSING"
				models.UpdateGameMatchStatus(nowGameMatch, "game_status")

				tmp_card := models.CardInfo{
					Type: models.EVENT_CLEAR_CARD,
				}
				sendMsgToSeat(tmp_card)

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
				nowGameMatch.GameStatus = "ROUND1"
				models.UpdateGameMatchStatus(nowGameMatch, "game_status")

				roundprocess <- models.RoundInfo{
					Type:            models.EVENT_ROUND_INFO,
					GM:              nowGameMatch,
					NowPosition:     nowGameMatch.SmallBindPosition,
					AllPointInRound: 0,
					MaxPoint:        0,
				}
				positionTurn = nowGameMatch.SmallBindPosition
				detailArr := models.GetRoundUserDetail(1)
				roundUserDetail = make(map[int]interface{})
				for _, v := range detailArr {
					roundUserDetail[v.Position] = v
				}
				//发送消息 通知小盲，大盲位置，已下注5 10
				//先写下注逻辑
				//所以这里先通知小盲的回合

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
			switch process.Type {
			case models.EVENT_LICENSING:
				var other_process = process
				other_process.Card = models.Card{}
				//这边需要特殊处理，所以暂时不走公共方法，后续如果有需要同样处理的位置再合
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
			case models.EVENT_PUBLIC_CARD:
				sendMsgToSeat(process)
			}

		case userInfo := <-userInfoChannel:
			logs.Info(userInfo)
			if userInfo.Type == models.EVENT_REFRESH_USER_INFO {
				var data UserInfoMsg
				data.Type = userInfo.Type
				data.Info = models.GetUserPointWithSeat(1)
				sendMsgToSeat(data)
			}
		case roundInfo := <-roundprocess:
			switch roundInfo.Type {
			case models.EVENT_ROUND_INFO:
				sendMsgToSeat(roundInfo)
			}

		case sub := <-subscribe:
			subscribers.PushBack(sub) // Add user to the end of list.
			if sub.UserType == models.POKER_PLAYER {
				gu := models.SetUserReturnPlayer(sub.User)

				var player = new(Player)
				player.user = sub.User
				player.Position = gu.Position
				player.Conn = sub.Conn
				seat.PushBack(*player)
				publish <- models.SeatInfo{
					Type:     models.EVENT_JOIN,
					GameUser: gu,
					User:     sub.User.Name,
				}
			}
			// Publish a JOIN event.
			// publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
			logs.Info("User:", sub.Name, ";WebSocket:", sub.Conn != nil)
		case event := <-publish:
			// Notify waiting list.
			for ch := waitingList.Back(); ch != nil; ch = ch.Prev() {
				ch.Value.(chan bool) <- true
				waitingList.Remove(ch)
			}

			if event.Type == models.EVENT_JOIN {
				for ss := seat.Front(); ss != nil; ss = ss.Next() {
					if ss.Value.(Player).user.Name == event.User {
						ws := ss.Value.(Player).Conn
						if ws != nil {
							msg, _ := json.Marshal(event)
							if ws.WriteMessage(websocket.TextMessage, msg) != nil {
								// User disconnected.
								unsubscribe <- ss.Value.(Player).user
							}
						}
						break
					}

				}
			}

			userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")

			// broadcastWebSocket(event)
			// models.NewArchive(event)

			if event.Type == models.EVENT_MESSAGE {
				logs.Info("Message from", event.User, ";Content:", event.GameUser)
			}
		case unsub := <-unsubscribe:
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
					publish <- models.SeatInfo{
						Type: models.EVENT_LEAVE,
						GameUser: models.GameUser{
							Position: ss.Value.(Player).Position,
						},
						User: unsub.Name,
					} // Publish a LEAVE event.
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

func sendMsgToSeat(data interface{}) {
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
