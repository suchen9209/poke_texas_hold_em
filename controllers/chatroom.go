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

type Viewer struct {
	uname string
	Conn  *websocket.Conn
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

	viewerList = list.New()

	// gameMatchPointLog = make(map[int]map[string]int)
	gameMatchAllin = make(map[int]int) //本局中allin的位置和point

	foldPoint = make(map[string]int)

	BigBindPositionTurn = 0
)

// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case op := <-gameop:
			if op == "start" {
				if (nowGameMatch.Id == 0 || nowGameMatch.GameStatus == models.GAME_STATUS_END) && seat.Len() >= 2 {
					startGame()
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

func opChangePoint(point int, position int) {
	a := roundUserDetail[position]
	a.RoundPoint = a.RoundPoint + point
	a.Point = a.Point - point
	roundUserDetail[position] = a
	models.ChangeUserPoint(a.UserId, -point)
	if a.RoundPoint > LimitPoint {
		LimitPoint = a.RoundPoint
	}
	// gameMatchPointLog[position][nowGameMatch.GameStatus] += point
}

func init() {
	models.TruncateGameUser()
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
	for ss := viewerList.Front(); ss != nil; ss = ss.Next() {
		ws := ss.Value.(Viewer).Conn
		if ws != nil {
			msg, _ := json.Marshal(data)
			if ws.WriteMessage(websocket.TextMessage, msg) != nil {
				// User disconnected.
				viewerList.Remove(ss)
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
	case models.GAME_STATUS_ROUND1:
		nextRound = models.GAME_STATUS_ROUND2
		nowGameMatch.Pot1st = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[1])
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[2])
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[3])
	case models.GAME_STATUS_ROUND2:
		nextRound = models.GAME_STATUS_ROUND3
		nowGameMatch.Pot2nd = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[4])
	case models.GAME_STATUS_ROUND3:
		nextRound = models.GAME_STATUS_ROUND4
		nowGameMatch.Pot3rd = getRoundPoint(nowGameMatch.GameStatus, true)
		gameprocess <- sendCard(models.EVENT_PUBLIC_CARD, "", 0, models.PublicCard[5])
	case models.GAME_STATUS_ROUND4:
		nextRound = models.GAME_STATUS_END
		nowGameMatch.Pot4th = getRoundPoint(nowGameMatch.GameStatus, true)
	}

	logs.Info("end round" + nowGameMatch.GameStatus)
	nowGameMatch.GameStatus = nextRound

	if nowGameMatch.GameStatus == models.GAME_STATUS_END {
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
	remainRoundPoint += foldPoint[nowGameMatch.GameStatus]
	if clear {
		foldPoint[nowGameMatch.GameStatus] = 0
	}
	logs.Info(foldPoint)
	logs.Info(remainRoundPoint)

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
	gameMatchAllin = make(map[int]int)
	foldPoint = make(map[string]int)
	BigBindPositionTurn = 0
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
	// logs.Info("init seat")
	// logs.Info(detailArr)
	for _, v := range detailArr {
		roundUserDetail[v.Position] = v
	}
	// logs.Info(roundUserDetail)

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
	// gameMatchPointLog[positionTurn][models.GAME_STATUS_ROUND1] = 5
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
	// gameMatchPointLog[positionTurn][models.GAME_STATUS_ROUND1] = 10
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
	var pointDetail = make(map[int]int)
	var winUserPos []int
	if len(roundUserDetail) == 1 {
		for key := range roundUserDetail {
			winUserPos = append(winUserPos, key)
			break
		}
	} else {
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
			if need_cal_pot > potAll{
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

	} else if len(winUserPos) > 0{
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

func CalWinUser() []int {
	tmpCardC := make(map[int]string)
	bigString := ""
	var winUser []int
	for k := range roundUserDetail {
		tmpArr := models.UsersCard[k]
		for _, v := range models.PublicCard {
			tmpArr = append(tmpArr, v)
		}
		tmpCardC[k] = models.GetString(tmpArr)
		if bigString == "" {
			bigString = models.GetString(tmpArr)
			winUser = append(winUser, k)
		}
	}

	for k, v := range tmpCardC {
		if bigString == v {
			continue
		}
		result := models.Compare(v, bigString)
		if result == 1 {
			winUser = winUser[0:0]
			winUser = append(winUser, k)
			bigString = v
		} else if result == 0 {
			winUser = append(winUser, k)
		}
	}

	var winCard [][]models.Card
	for _, v := range winUser {
		winCard = append(winCard, models.UsersCard[v])
	}

	models.TransMaxHandToCardInfo()
	var tmp []models.Card
	for _, v := range models.PublicCard {
		tmp = append(tmp, v)
	}
	a := models.EndGameCardInfo{
		Type:       models.EVENT_GAME_END_SHOW_CARD,
		WinPos:     winUser,
		WinCard:    winCard,
		PublicCard: tmp,
		// BigCard:    models.StringToCard(bigString),\
		BigCard:   models.ShowMaxCard,
		UserCards: models.UsersCard,
	}
	sendMsgToSeat(a)

	return winUser
}

//传入当前仍在场的用户
func GetBigUser(iru map[int]models.InRoundUserDetail) ([]int, string) {
	tmpCardC := make(map[int]string)
	bigString := ""
	var winUser uint64
	logs.Info(iru)
	for k := range iru {
		tmpArr := models.UsersCard[k]
		for _, v := range models.PublicCard {
			tmpArr = append(tmpArr, v)
		}
		tmpCardC[k] = models.GetString(tmpArr)
		if bigString == "" {
			bigString = models.GetString(tmpArr)
			winUser = 1 << k
		}
		logs.Info(winUser)
	}

	for k, v := range tmpCardC {
		if bigString == v {
			continue
		}
		result := models.Compare(v, bigString)
		if result == 1 {
			winUser = 1 << k
			bigString = v
		} else if result == 0 {
			winUser |= 1 << k
		}
		logs.Info(winUser)
	}

	var winUserArr []int
	i := 0
	for winUser > 0 {
		if winUser&1 > 0 {
			winUserArr = append(winUserArr, i)
		}
		winUser = winUser >> 1
		i++
	}
	logs.Info(winUserArr)

	return winUserArr, bigString
}
