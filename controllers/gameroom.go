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

	"poke/models"

	"github.com/beego/beego/v2/core/logs"

	"github.com/gorilla/websocket"
)

var (
	UserConnMap = make(map[int] *websocket.Conn)

	publish_map         map[int]chan models.SeatInfo
	gameProcessMap     map[int]chan models.CardInfo
	userInfoChannelMap map[int]chan models.Event
	roundprocess_map   map[int]chan models.RoundInfo
	userOperationProcessMap map[int]chan models.UserOperationMsg
	gameOpMap               map[int]chan string
	seat_map                map[int]*list.List
	positionTurn_map         map[int]int //记录当前轮到的玩家位置
	LimitPoint_map           map[int]int //当前轮最小点数值

	nowGameMatchMap map[int]models.GameMatch

	roundUserDetail_map map[int]map[int]models.InRoundUserDetail

	emptySendMap  map[int]int //游戏结束时清零 仅用于记录大小盲
	detailArr_map map[int][]models.InRoundUserDetail

	roundCheckNumber_map map[int]int

	viewerList_map map[int]*list.List

	gameMatchAllin_map map[int]map[int]int //allin的位置和point

	foldPoint_map map[int]map[string]int

	BigBindPositionTurn_map map[int]int
)

// This function handles all incoming chan messages.
func gameRoom(roomId int) {
	for {
		select {
		case op := <-gameOpMap[roomId]:
			if op == "start" {
				if (nowGameMatchMap[roomId].Id == 0 || nowGameMatchMap[roomId].GameStatus == models.GAME_STATUS_END) && seat_map[roomId].Len() >= 2 {
					startRoomGame(roomId)
				}
			}
			if op == "close"{
				models.CloseRoom(roomId)
				return
			}
		case process := <-gameProcessMap[roomId]:
			switch process.Type {
			case models.EVENT_LICENSING:
				var otherProcess = process
				otherProcess.Card = models.Card{}
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
							msg, _ := json.Marshal(otherProcess)
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
				tmp := roundUserDetail[roundInfo.NowPosition]
				tmp.AllowOp = tmp.AllowOp[0:0]
				if LimitPoint == 0 {
					tmp.AllowOp = append(tmp.AllowOp, "check")
				}
				if tmp.RoundPoint+tmp.Point >= LimitPoint && LimitPoint > 0 {
					tmp.AllowOp = append(tmp.AllowOp, "call")
				}
				if tmp.RoundPoint+tmp.Point > LimitPoint {
					tmp.AllowOp = append(tmp.AllowOp, "raise")
				}
				tmp.AllowOp = append(tmp.AllowOp, "allin")
				tmp.AllowOp = append(tmp.AllowOp, "fold")

				roundUserDetail[roundInfo.NowPosition] = tmp
				roundInfo.Detail = roundUserDetail[roundInfo.NowPosition]
				roundInfo.AllPointInRound = getRoundPoint(nowGameMatch.GameStatus, false)
				sendMsgToSeat(roundInfo)
			}
		case uop := <-userOperationProcess:
			emptySend++
			fromCheck := false
			uop.GameMatchLog.GameMatchId = nowGameMatch.Id
			if uop.Position == nowGameMatch.BigBindPosition {
				BigBindPositionTurn++
			}
			switch uop.GameMatchLog.Operation {
			case models.GAME_OP_RAISE: //raise
				opChangePoint(uop.GameMatchLog.PointNumber, uop.Position)
			case models.GAME_OP_CALL: //call
				uop.GameMatchLog.PointNumber = LimitPoint - roundUserDetail[uop.Position].RoundPoint
				opChangePoint(uop.GameMatchLog.PointNumber, uop.Position)
			case models.GAME_OP_CHECK: //check
				roundCheckNumber++
				fromCheck = true
			case models.GAME_OP_ALLIN: // allin
				userRePoint := models.GetUserPoint(roundUserDetail[uop.Position].UserId)
				uop.GameMatchLog.PointNumber = userRePoint
				opChangePoint(uop.GameMatchLog.PointNumber, uop.Position)
				gameMatchAllin[uop.Position] = roundUserDetail[uop.Position].RoundPoint
			case models.GAME_OP_FOLD:
				foldPoint[nowGameMatch.GameStatus] += roundUserDetail[uop.Position].RoundPoint
				logs.Info(foldPoint)
				delete(roundUserDetail, uop.Position)
				delete(models.UsersCard, uop.Position)
			}
			models.AddGameMatchLog(uop.GameMatchLog)
			sendMsgToSeat(uop)
			userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
			// logs.Info("after user operation")
			// logs.Info(roundUserDetail)
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
						_, ok := gameMatchAllin[v.Position]
						if v.RoundPoint < LimitPoint && !ok {
							have_not_fill_point = true
						}
					}
					if BigBindPositionTurn <= 1 { //只在第一轮中有效
						have_not_fill_point = true
					}

					if !have_not_fill_point && (len(roundUserDetail)-len(gameMatchAllin) <= 1) {
						GameEnd()
					} else if have_not_fill_point {
						//next user
						if emptySend > 2 { //只是为了排除大小盲
							positionTurn = getNextPosition(roundUserDetail, positionTurn)
							var ok bool
							_, ok = gameMatchAllin[positionTurn]
							for ok {
								positionTurn = getNextPosition(roundUserDetail, positionTurn)
								_, ok = gameMatchAllin[positionTurn]
							}
							nextRoundInfo()
						}
					} else {
						//end this round
						endRoundPoint()
						nextRoundInfo()
					}

					// if !have_not_fill_point && (len(roundUserDetail)-len(gameMatchAllin) <= 1) {
					// 	GameEnd()
					// } else if have_not_fill_point && len(roundUserDetail) > (len(gameMatchAllin)+1) {
					// 	//next user
					// 	emptySend++
					// 	if emptySend > 2 {
					// 		positionTurn = getNextPosition(roundUserDetail, positionTurn)
					// 		var ok bool
					// 		_, ok = gameMatchAllin[positionTurn]
					// 		for ok {
					// 			positionTurn = getNextPosition(roundUserDetail, positionTurn)
					// 			_, ok = gameMatchAllin[positionTurn]
					// 		}
					// 		nextRoundInfo()
					// 	}
					// } else {
					// 	//end this round
					// 	endRoundPoint()
					// 	nextRoundInfo()
					// }
				}
			}
		case sub := <-subscribe:
			if nowGameMatch.Id != 0 && nowGameMatch.GameStatus != models.GAME_STATUS_END {
				sub.Conn.WriteMessage(websocket.TextMessage, []byte("wait Game End"))
			} else {
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
				if sub.UserType == models.VIEWER {
					var view = new(Viewer)
					view.uname = sub.Name
					view.Conn = sub.Conn
					viewerList.PushBack(*view)
				}
			}

			// Publish a JOIN event.
			// publish <- newEvent(models.EVENT_JOIN, sub.Name, "")
			// logs.Info("User:", sub.Name, ";WebSocket:", sub.Conn != nil)
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
					delete(roundUserDetail, ss.Value.(Player).Position)
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
					if seat.Len() == 1 {
						GameEnd()
					}
					if seat.Len() == 0 {
						nowGameMatch.GameStatus = "END"
					}
					break
				}
			}
		}
	}
}

func startRoomGame(RoomId int) {
	//游戏结束时清零 仅用于记录大小盲
	emptySendMap[RoomId] = 0
	//清除上局剩余的用户信息，例如卡牌信息等
	userInfoChannelMap[RoomId] <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
	nowGameMatchMap[RoomId] = models.InitGameMatch(RoomId)
	models.InitCardMap()
	gameMatchAllin = make(map[int]int)
	nowGameMatch.GameStatus = models.GAME_STATUS_LICENSING
	models.UpdateGameMatchStatus(nowGameMatch, "game_status")

	tmp_card := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}
	sendMsgToSeat(tmp_card)

	for ss := seat.Front(); ss != nil; ss = ss.Next() {
		var tmp_poker = models.GetOneCard()
		models.UsersCard[ss.Value.(Player).Position] = append(models.UsersCard[ss.Value.(Player).Position], *tmp_poker)
		gameprocess <- sendCard(models.EVENT_LICENSING, ss.Value.(Player).user.Name, ss.Value.(Player).Position, *tmp_poker)
	}
	for ss := seat.Front(); ss != nil; ss = ss.Next() {
		var tmp_poker = models.GetOneCard()
		models.UsersCard[ss.Value.(Player).Position] = append(models.UsersCard[ss.Value.(Player).Position], *tmp_poker)
		gameprocess <- sendCard(models.EVENT_LICENSING, ss.Value.(Player).user.Name, ss.Value.(Player).Position, *tmp_poker)
	}

	for i := 1; i <= 5; i++ {
		var tmp_poker = models.GetOneCard()
		models.PublicCard[i] = *tmp_poker
	}
	nowGameMatch.GameStatus = models.GAME_STATUS_ROUND1
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
