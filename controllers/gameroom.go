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
	UserConnMap = make(map[int] *websocket.Conn)	//记录用户ws连接

	publish_map             map[int]chan models.SeatInfo
	gameProcessMap          map[int]chan models.CardInfo	//游戏消息对应的channel
	userInfoChannelMap      map[int]chan models.Event	//用户信息
	roundProcessMap         map[int]chan models.RoundInfo
	userOperationProcessMap map[int]chan models.UserOperationMsg	//用户操作
	gameOpMap               map[int]chan string		//游戏操作 控制游戏进程
	seat_map                map[int]*list.List
	PositionTurnMap         map[int]int //记录当前轮到的玩家位置
	LimitPointMap           map[int]int //当前轮最小点数值

	nowGameMatchMap map[int]models.GameMatch	//当前游戏的详情

	roundUserDetailMap map[int]map[int]models.InRoundUserDetail

	emptySendMap  map[int]int //游戏结束时清零 仅用于记录大小盲
	detailArr_map map[int][]models.InRoundUserDetail

	RoundCheckNumberMap map[int]int

	viewerList_map map[int]*list.List

	gameMatchAllinMap map[int]map[int]int //allin的位置和point

	foldPoint_map map[int]map[string]int

	BigBindPositionTurnMap map[int]int

	GameAllCardMap map[int]map[int]models.Card
	PublicCardMap map[int]map[int]models.Card

	UsersCardMap map[int]map[int][]models.Card
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
			case models.EVENT_PUBLIC_CARD:
				sendToRoomUser(roomId,process)
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
		case uop := <-userOperationProcessMap[roomId]:
			emptySendMap[roomId]++
			fromCheck := false
			uop.GameMatchLog.GameMatchId = nowGameMatchMap[roomId].Id
			if uop.Position == nowGameMatchMap[roomId].BigBindPosition {
				BigBindPositionTurnMap[roomId]++
			}
			switch uop.GameMatchLog.Operation {
			case models.GAME_OP_RAISE: //raise
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position,roomId)
			case models.GAME_OP_CALL: //call
				uop.GameMatchLog.PointNumber = LimitPointMap[roomId] - roundUserDetailMap[roomId][uop.Position].RoundPoint
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position,roomId)
			case models.GAME_OP_CHECK: //check
				RoundCheckNumberMap[roomId]++
				fromCheck = true
			case models.GAME_OP_ALLIN: // allin
				userRePoint := models.GetUserPoint(roundUserDetailMap[roomId][uop.Position].UserId)
				uop.GameMatchLog.PointNumber = userRePoint
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position,roomId)
				gameMatchAllinMap[roomId][uop.Position] = roundUserDetailMap[roomId][uop.Position].RoundPoint
			case models.GAME_OP_FOLD:
				foldPoint[nowGameMatchMap[roomId].GameStatus] += roundUserDetailMap[roomId][uop.Position].RoundPoint
				logs.Info(foldPoint)
				delete(roundUserDetailMap[roomId], uop.Position)
				delete(UsersCardMap[roomId], uop.Position)
			}
			models.AddGameMatchLog(uop.GameMatchLog)
			sendToRoomUser(roomId,uop)
			userInfoChannelMap[roomId] <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
			// logs.Info("after user operation")
			// logs.Info(roundUserDetail)
			if len(roundUserDetailMap[roomId]) <= 1 {
				//game end
				RoomGameEnd(roomId)
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
	initCardMap(RoomId)
	gameMatchAllinMap[RoomId] = make(map[int]int)
	changeGameMatch(RoomId,models.GAME_STATUS_LICENSING)

	tmpCard := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}
	sendToRoomUser(RoomId, tmpCard)
	sendCardToUser(RoomId)

	for i := 1; i <= 5; i++ {
		var tmpPoker = models.GetOneCardFromCardMap(GameAllCardMap[RoomId])
		PublicCardMap[RoomId][i] = *tmpPoker
	}
	changeGameMatch(RoomId,models.GAME_STATUS_ROUND1)

	//初始化座位
	detailArr = models.GetRoundUserDetail(RoomId)
	roundUserDetailMap[RoomId] = make(map[int]models.InRoundUserDetail)
	logs.Info("init seat")
	logs.Info(detailArr)
	for _, v := range detailArr {
		roundUserDetailMap[RoomId][v.Position] = v
	}
	logs.Info(roundUserDetailMap[RoomId])

	//小盲
	PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].SmallBindPosition
	userOpMsg := models.UserOperationMsg{
		Type:     models.EVENT_USER_OPERATION_INFO,
		Position: nowGameMatchMap[RoomId].SmallBindPosition,
		Name:     roundUserDetailMap[RoomId][nowGameMatch.SmallBindPosition].Name,
		GameMatchLog: models.GameMatchLog{
			GameMatchId: nowGameMatchMap[RoomId].Id,
			UserId:      roundUserDetailMap[RoomId][nowGameMatch.SmallBindPosition].UserId,
			Operation:   models.GAME_OP_RAISE,
			PointNumber: 5,
		},
	}
	userOperationProcessMap[RoomId] <- userOpMsg

	//大盲
	PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].BigBindPosition
	uomsg2 := models.UserOperationMsg{
		Type:     models.EVENT_USER_OPERATION_INFO,
		Position: nowGameMatchMap[RoomId].BigBindPosition,
		Name:     roundUserDetail[nowGameMatchMap[RoomId].BigBindPosition].Name,
		GameMatchLog: models.GameMatchLog{
			GameMatchId: nowGameMatchMap[RoomId].Id,
			UserId:      roundUserDetail[nowGameMatchMap[RoomId].BigBindPosition].UserId,
			Operation:   models.GAME_OP_RAISE,
			PointNumber: 10,
		},
	}
	userOperationProcessMap[RoomId] <- uomsg2

	//发送消息 通知小盲，大盲位置，已下注5 10
	LimitPointMap[RoomId] = 10
	PositionTurnMap[RoomId] = PositionTurnMap[RoomId] + 1
	if positionTurn > len(roundUserDetailMap[RoomId]) {
		PositionTurnMap[RoomId] = 1
	}
	if _, ok := roundUserDetailMap[RoomId][positionTurn]; !ok {
		PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].SmallBindPosition
	}
	roomNextRoundInfo(RoomId)
}

func roomOpChangePoint(point int, position int,roomId int) {
	a := roundUserDetailMap[roomId][position]
	a.RoundPoint = a.RoundPoint + point
	a.Point = a.Point - point
	roundUserDetailMap[roomId][position] = a
	models.ChangeUserPoint(a.UserId, -point)
	if a.RoundPoint > LimitPointMap[roomId] {
		LimitPointMap[roomId] = a.RoundPoint
	}
	// gameMatchPointLog[position][nowGameMatch.GameStatus] += point
}

func roomNextRoundInfo(RoomId int) {
	roundProcessMap[RoomId] <- models.RoundInfo{
		Type:        models.EVENT_ROUND_INFO,
		GM:          nowGameMatchMap[RoomId],
		NowPosition: PositionTurnMap[RoomId],
		MaxPoint:    LimitPointMap[RoomId],
	}
}


func initCardMap(roomId int){
	GameAllCardMap[roomId] = models.GetNewCardMap()
	PublicCardMap[roomId] = make(map[int]models.Card)
	UsersCardMap[roomId] = make(map[int][]models.Card)
}

func changeGameMatch(roomId int,status string){
	tmp := nowGameMatchMap[roomId]
	tmp.GameStatus = status
	models.UpdateGameMatchStatus(tmp, "game_status")
	nowGameMatchMap[roomId] = tmp
}

func sendToRoomUser(roomId int,data interface{}){
	positionMap := models.GetRoomUserPositionList(roomId)
	for _, uid := range positionMap {
		sendInWs(UserConnMap[uid],data)
	}
}

func sendInWs(ws *websocket.Conn,data interface{}){
	if ws != nil {
		msg, _ := json.Marshal(data)
		if ws.WriteMessage(websocket.TextMessage, msg) != nil {
			// User disconnected.
			//unsubscribe <- ss.Value.(Player).user
		}
	}
}
func sendCardToUser(roomId int){
	positionMap := models.GetRoomUserPositionList(roomId)
	for pos, uid := range positionMap {
		var tmpPoker = models.GetOneCardFromCardMap(GameAllCardMap[roomId])
		UsersCardMap[roomId][pos] = append(UsersCardMap[roomId][pos], *tmpPoker)
		sendInWs(UserConnMap[uid],sendCard(models.EVENT_LICENSING, "", pos, *tmpPoker))

		var tmpPoker2 = models.GetOneCardFromCardMap(GameAllCardMap[roomId])
		UsersCardMap[roomId][pos] = append(UsersCardMap[roomId][pos], *tmpPoker2)
		sendInWs(UserConnMap[uid],sendCard(models.EVENT_LICENSING, "", pos, *tmpPoker2))
	}
}


/**
一局游戏的结束
重新Init数据
分配上局获胜点数
**/
func RoomGameEnd(RoomId int) {
	var pointDetail = make(map[int]int)
	var winUserPos []int
	if len(roundUserDetailMap[RoomId]) == 1 {
		for key := range roundUserDetailMap[RoomId] {
			winUserPos = append(winUserPos, key)
			break
		}
	} else {
		//写到这了
		winUserPos = CalWinUser()
	}
	//计算获胜点数
	nowGameMatch.PotAll = nowGameMatch.Pot1st + nowGameMatch.Pot2nd + nowGameMatch.Pot3rd + nowGameMatch.Pot4th
	for _, v := range roundUserDetail {
		nowGameMatch.PotAll += v.RoundPoint
	}
	logs.Info(foldPoint)
	for _, v := range foldPoint {
		nowGameMatch.PotAll += v
	}

	if len(gameMatchAllin) > 0 {
		potAll := nowGameMatch.PotAll                           //总池
		all_in_point_desc := models.RankByPoint(gameMatchAllin) //根据allin数量进行的排序
		var cal_user_detail = make(map[int]models.InRoundUserDetail)
		for k, v := range roundUserDetail {
			cal_user_detail[k] = v
		}
		cal_all_in_pot := 0 //已结算的allin底池

		logs.Info(potAll)
		logs.Info(all_in_point_desc)
		logs.Info(cal_user_detail)
		logs.Info(cal_all_in_pot)

		for _, v := range all_in_point_desc {

			logs.Info(v)

			need_cal_pot := (v.Value - cal_all_in_pot) * len(cal_user_detail)
			if need_cal_pot > potAll {
				need_cal_pot = potAll
			}
			win_user, _ := GetBigUser(cal_user_detail)
			perPot := need_cal_pot / len(win_user)
			for _, v := range win_user {
				models.ChangeUserPoint(roundUserDetail[v].UserId, perPot)
				pointDetail[v] += perPot
			}
			cal_all_in_pot = v.Value
			potAll -= need_cal_pot
			delete(cal_user_detail, v.Key)

			logs.Info(need_cal_pot)
			logs.Info(perPot)
			logs.Info(cal_all_in_pot)
			logs.Info(potAll)
		}
		logs.Info(pointDetail)
		if potAll > 0 { //cal_user_detail中还剩余多个未allin玩家时未出现
			var win_user2 []int
			logs.Info(cal_user_detail)
			logs.Info(len(cal_user_detail))
			logs.Info(len(cal_user_detail) > 0)
			logs.Info(roundUserDetail)
			// if len(cal_user_detail) > 0 {
			// 	win_user2, _ = GetBigUser(cal_user_detail)
			// } else {
			// 	win_user2, _ = GetBigUser(roundUserDetail)
			// }
			logs.Info(win_user2)
			win_user2, _ = GetBigUser(roundUserDetail)

			perPot := potAll / len(win_user2)
			logs.Info(roundUserDetail)
			for _, v := range win_user2 {
				models.ChangeUserPoint(roundUserDetail[v].UserId, perPot)
				pointDetail[v] += perPot
			}
			logs.Info(pointDetail)
		}

	} else if len(winUserPos) > 0 {
		perPot := nowGameMatch.PotAll / len(winUserPos)
		for _, v := range winUserPos {
			models.ChangeUserPoint(roundUserDetail[v].UserId, perPot)
			pointDetail[v] = perPot
		}
	}

	// if len(winUserPos) == 1 {
	// 	if all_point, pok := gameMatchAllin[winUserPos[0]]; pok {
	// 		models.ChangeUserPoint(roundUserDetail[winUserPos[0]].UserId, all_point*2)

	// 	} else {
	// 		models.ChangeUserPoint(roundUserDetail[winUserPos[0]].UserId, nowGameMatch.PotAll)
	// 	}
	// } else {

	// }

	nowGameMatch.GameStatus = models.GAME_STATUS_END
	models.UpdateGameMatchStatus(nowGameMatch, "game_status")
	models.UpdateGameMatchStatus(nowGameMatch, "pot1st")
	models.UpdateGameMatchStatus(nowGameMatch, "pot2nd")
	models.UpdateGameMatchStatus(nowGameMatch, "pot3rd")
	models.UpdateGameMatchStatus(nowGameMatch, "pot4th")
	models.UpdateGameMatchStatus(nowGameMatch, "pot_all")

	//提示胜利玩家可以重新开始游戏了
	logs.Info("game end")
	tmp_card := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}
	userInfoChannel <- newEvent(models.EVENT_REFRESH_USER_INFO, "", "")
	sendMsgToSeat(tmp_card)
	sendMsgToSeat(newEvent(models.EVENT_GAME_END, "system", "Game End"))

}
