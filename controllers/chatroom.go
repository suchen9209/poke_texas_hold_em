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
	publish              = make(chan models.SeatInfo, 1000)
	gameprocess          = make(chan models.CardInfo, 1000)
	userInfoChannel      = make(chan models.Event, 1000)
	roundprocess         = make(chan models.RoundInfo, 1000)
	userOperationProcess = make(chan models.UserOperationMsg, 1000)

	gameop = make(chan string)

	// Long polling waiting list.
	waitingList = list.New()
	subscribers = list.New()

	seat         = list.New()
	positionTurn int //记录当前轮到的玩家位置
	LimitPoint   int //当前轮最小点数值

	nowGameMatch models.GameMatch

	roundUserDetail = make(map[int]models.InRoundUserDetail)

	emptySend = 0 //游戏结束时清零 仅用于记录大小盲
	detailArr []models.InRoundUserDetail

	roundCheckNumber = 0
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case op := <-gameop:
			if op == "start" {
				startGame()

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
				roundInfo.Detail = roundUserDetail[roundInfo.NowPosition]
				roundInfo.AllPointInRound = getRoundPoint(nowGameMatch.GameStatus, false)
				sendMsgToSeat(roundInfo)
			}
		case uop := <-userOperationProcess:
			fromCheck := false
			uop.GameMatchLog.GameMatchId = nowGameMatch.Id
			switch uop.GameMatchLog.Operation {
			case models.GAME_OP_RAISE:
				a := roundUserDetail[uop.Position]
				a.RoundPoint = a.RoundPoint + uop.GameMatchLog.PointNumber
				a.Point = a.Point - uop.GameMatchLog.PointNumber
				roundUserDetail[uop.Position] = a
				models.ChangeUserPoint(a.UserId, -uop.GameMatchLog.PointNumber)
				LimitPoint = a.RoundPoint
			case models.GAME_OP_CALL:
				a := roundUserDetail[uop.Position]
				uop.GameMatchLog.PointNumber = LimitPoint - a.RoundPoint
				a.RoundPoint = LimitPoint
				a.Point = a.Point - uop.GameMatchLog.PointNumber
				roundUserDetail[uop.Position] = a
				models.ChangeUserPoint(a.UserId, -uop.GameMatchLog.PointNumber)
			case models.GAME_OP_CHECK:
				roundCheckNumber++
				fromCheck = true
				//do nothing
			case models.GAME_OP_ALLIN:
				a := roundUserDetail[uop.Position]
				userRePoint := models.GetUserPoint(a.UserId)
				a.RoundPoint = userRePoint
				a.Point = 0
				roundUserDetail[uop.Position] = a
				models.ChangeUserPoint(a.UserId, -userRePoint)
				LimitPoint = a.RoundPoint
			case models.GAME_OP_FOLD:
				switch nowGameMatch.GameStatus {
				case "ROUND1":
					nowGameMatch.Pot1st = nowGameMatch.Pot1st + roundUserDetail[uop.Position].RoundPoint
				case "ROUND2":
					nowGameMatch.Pot2nd = nowGameMatch.Pot1st + roundUserDetail[uop.Position].RoundPoint
				case "ROUND3":
					nowGameMatch.Pot3rd = nowGameMatch.Pot1st + roundUserDetail[uop.Position].RoundPoint
				case "ROUND4":
					nowGameMatch.Pot4th = nowGameMatch.Pot1st + roundUserDetail[uop.Position].RoundPoint
				}
				delete(roundUserDetail, uop.Position)
			}
			models.AddGameMatchLog(uop.GameMatchLog)
			sendMsgToSeat(uop)
			logs.Info("after user operation")
			logs.Info(roundUserDetail)
			if len(roundUserDetail) <= 1 {
				//game end
				GameEnd()
			} else {
				if fromCheck {
					if roundCheckNumber < len(roundUserDetail) {
						positionTurn = getNextPosition(roundUserDetail, positionTurn)
					} else {
						endRoundPoint()
					}
					nextRoundInfo()
				} else {
					var have_not_fill_point bool
					have_not_fill_point = false
					for _, v := range roundUserDetail {
						if v.RoundPoint != LimitPoint {
							have_not_fill_point = true
						}
					}

					if have_not_fill_point {
						//next user
						emptySend++
						if emptySend > 2 {
							positionTurn = getNextPosition(roundUserDetail, positionTurn)
							nextRoundInfo()
						}
					} else {
						//end this round
						endRoundPoint()
						nextRoundInfo()
					}
				}
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

//用于查找下一个位置的人，有则返回座位号，无则返回false
func getNextPosition(m map[int]models.InRoundUserDetail, initpostition int) int {
	user_number := len(detailArr)
	for i := initpostition + 1; i <= initpostition+user_number; i++ {
		key := i % user_number
		if key == 0 {
			key = user_number
		}
		if _, ok := m[key]; ok {
			return key
		}
	}
	return 0
}

func nextRoundInfo() {
	roundprocess <- models.RoundInfo{
		Type:        models.EVENT_ROUND_INFO,
		GM:          nowGameMatch,
		NowPosition: positionTurn,
		MaxPoint:    LimitPoint,
	}
}

/**
标志着一轮的结束
**/
func endRoundPoint() {

	roundCheckNumber = 0
	if _, ok := roundUserDetail[nowGameMatch.SmallBindPosition]; ok {
		positionTurn = nowGameMatch.SmallBindPosition
	} else {
		positionTurn = getNextPosition(roundUserDetail, nowGameMatch.SmallBindPosition)
	}
	LimitPoint = 0
	var nextRound string
	switch nowGameMatch.GameStatus {
	case "ROUND1":
		nextRound = "ROUND2"
		nowGameMatch.Pot1st = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[1])
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[2])
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[3])
	case "ROUND2":
		nextRound = "ROUND3"
		nowGameMatch.Pot2nd = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[4])
	case "ROUND3":
		nextRound = "ROUND4"
		nowGameMatch.Pot3rd = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[5])
	case "ROUND4":
		nextRound = "END"
		nowGameMatch.Pot4th = getRoundPoint(nowGameMatch.GameStatus, true)
		nowGameMatch.PotAll = nowGameMatch.Pot1st + nowGameMatch.Pot2nd + nowGameMatch.Pot3rd + nowGameMatch.Pot4th
	}

	logs.Info("end round" + nowGameMatch.GameStatus)
	nowGameMatch.GameStatus = nextRound

	if nowGameMatch.GameStatus == "END" {
		//需比较剩余玩家的卡牌大小
		GameEnd()
	}

	logs.Info(nowGameMatch)
}

func getRoundPoint(round string, clear bool) int {
	var remainRoundPoint int
	for k, v := range roundUserDetail {
		remainRoundPoint += v.RoundPoint
		if clear {
			a := roundUserDetail[k]
			a.RoundPoint = 0
			roundUserDetail[k] = a
		}
	}
	switch nowGameMatch.GameStatus {
	case "ROUND1":
		return nowGameMatch.Pot1st + remainRoundPoint
	case "ROUND2":
		return nowGameMatch.Pot2nd + remainRoundPoint
	case "ROUND3":
		return nowGameMatch.Pot3rd + remainRoundPoint
	case "ROUND4":
		return nowGameMatch.Pot4th + remainRoundPoint
	}
	return 0
}

func startGame() {
	emptySend = 0 //游戏结束时清零 仅用于记录大小盲
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

	//初始化座位
	detailArr = models.GetRoundUserDetail(1)
	roundUserDetail = make(map[int]models.InRoundUserDetail)
	logs.Info("init seat")
	logs.Info(detailArr)
	for _, v := range detailArr {
		roundUserDetail[v.Position] = v
	}
	logs.Info(roundUserDetail)

	//小盲
	positionTurn = nowGameMatch.SmallBindPosition
	uomsg := models.UserOperationMsg{
		Type:     models.EVENT_USER_OPERATION_INFO,
		Position: nowGameMatch.SmallBindPosition,
		Name:     roundUserDetail[nowGameMatch.SmallBindPosition].Name,
		GameMatchLog: models.GameMatchLog{
			GameMatchId: nowGameMatch.Id,
			UserId:      roundUserDetail[nowGameMatch.SmallBindPosition].UserId,
			Operation:   models.GAME_OP_RAISE,
			PointNumber: 5,
		},
	}
	userOperationProcess <- uomsg

	//大盲
	positionTurn = nowGameMatch.BigBindPosition
	uomsg2 := models.UserOperationMsg{
		Type:     models.EVENT_USER_OPERATION_INFO,
		Position: nowGameMatch.BigBindPosition,
		Name:     roundUserDetail[nowGameMatch.BigBindPosition].Name,
		GameMatchLog: models.GameMatchLog{
			GameMatchId: nowGameMatch.Id,
			UserId:      roundUserDetail[nowGameMatch.BigBindPosition].UserId,
			Operation:   models.GAME_OP_RAISE,
			PointNumber: 10,
		},
	}
	userOperationProcess <- uomsg2

	//发送消息 通知小盲，大盲位置，已下注5 10
	LimitPoint = 10
	positionTurn = positionTurn + 1
	if positionTurn > len(roundUserDetail) {
		positionTurn = 1
	}
	if _, ok := roundUserDetail[positionTurn]; !ok {
		positionTurn = nowGameMatch.SmallBindPosition
	}
	nextRoundInfo()
}

/**
一局游戏的结束
重新Init数据
分配上局获胜点数
**/
func GameEnd() {
	winPos := 0
	winName := ""
	winUserId := 0
	if len(roundUserDetail) == 1 {
		for key, v := range roundUserDetail {
			winPos = key
			winName = v.Name
			winUserId = v.UserId
			break
		}
	} else {
		winPos, winName, winUserId = CalWinUser()
	}
	//计算获胜点数
	nowGameMatch.PotAll = nowGameMatch.Pot1st + nowGameMatch.Pot2nd + nowGameMatch.Pot3rd + nowGameMatch.Pot4th
	for _, v := range roundUserDetail {
		nowGameMatch.PotAll += v.RoundPoint
	}

	models.UpdateGameMatchStatus(nowGameMatch, "game_status")
	models.UpdateGameMatchStatus(nowGameMatch, "pot1st")
	models.UpdateGameMatchStatus(nowGameMatch, "pot2nd")
	models.UpdateGameMatchStatus(nowGameMatch, "pot3rd")
	models.UpdateGameMatchStatus(nowGameMatch, "pot4th")
	models.UpdateGameMatchStatus(nowGameMatch, "pot_all")
	models.ChangeUserPoint(winUserId, nowGameMatch.PotAll)
	//提示胜利玩家可以重新开始游戏了
	logs.Info("game end")
	tmp_card := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}
	userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
	sendMsgToSeat(tmp_card)
	sendMsgToSeat(newEvent(models.EVENT_GAME_END, "system", "Game End"))
	logs.Info(winPos)
	logs.Info(winName)
	//将内容初始化
}

func CalWinUser() (int, string, int) {
	return 1, "suchot", 1
}
