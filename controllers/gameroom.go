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
	"encoding/json"

	"poke/models"

	"github.com/beego/beego/v2/core/logs"

	"github.com/gorilla/websocket"
)

var (
	roomManageOpenList  = make(chan int)
	roomManageCloseList = make(chan int)

	UserConnMap = make(map[int]*websocket.Conn) //记录用户ws连接

	userInfoChannelMap      = make(map[int]chan models.Event) //用户信息
	roundProcessMap         = make(map[int]chan models.RoundInfo)
	userOperationProcessMap = make(map[int]chan models.UserOperationMsg) //用户操作
	gameOpMap               = make(map[int]chan string)                  //游戏操作 控制游戏进程
	PositionTurnMap         = make(map[int]int)                          //记录当前轮到的玩家位置
	LimitPointMap           = make(map[int]int)                          //当前轮最小点数值

	nowGameMatchMap = make(map[int]models.GameMatch) //当前游戏的详情

	roundUserDetailMap = make(map[int]map[int]models.InRoundUserDetail)

	emptySendMap = make(map[int]int) //游戏结束时清零 仅用于记录大小盲
	DetailArrMap = make(map[int][]models.InRoundUserDetail)

	RoundCheckNumberMap = make(map[int]int)

	gameMatchAllinMap = make(map[int]map[int]int) //allin的位置和point

	FoldPointMap = make(map[int]map[string]int)

	BigBindPositionTurnMap = make(map[int]int)

	GameAllCardMap = make(map[int]map[int]models.Card)
	PublicCardMap  = make(map[int]map[int]models.Card)

	UsersCardMap = make(map[int]map[int][]models.Card)
)

func roomManage() {
	for {
		select {
		case openRoomId := <-roomManageOpenList:
			//userInfoChannelMap[openRoomId] = make(chan models.Event)
			gameOpMap[openRoomId] = make(chan string, 1000)
			userOperationProcessMap[openRoomId] = make(chan models.UserOperationMsg, 1000)
			userInfoChannelMap[openRoomId] = make(chan models.Event, 1000)
			roundProcessMap[openRoomId] = make(chan models.RoundInfo, 1000)
			roundUserDetailMap[openRoomId] = make(map[int]models.InRoundUserDetail)
			PositionTurnMap[openRoomId] = 0
			LimitPointMap[openRoomId] = 0
			nowGameMatchMap[openRoomId] = models.GameMatch{}
			emptySendMap[openRoomId] = 0
			DetailArrMap[openRoomId] = []models.InRoundUserDetail{}
			RoundCheckNumberMap[openRoomId] = 0
			gameMatchAllinMap[openRoomId] = make(map[int]int)
			FoldPointMap[openRoomId] = make(map[string]int)
			GameAllCardMap[openRoomId] = make(map[int]models.Card)
			PublicCardMap[openRoomId] = make(map[int]models.Card)
			UsersCardMap[openRoomId] = make(map[int][]models.Card)
			go gameRoom(openRoomId)
		case closeRoomId := <-roomManageCloseList:
			str := "close"
			gameOpMap[closeRoomId] <- str
		}
	}
}

func init() {
	logs.Info("77")
	go roomManage()
}

// This function handles all incoming chan messages.
func gameRoom(roomId int) {
	logs.Info("abd")
	for {
		select {
		case op := <-gameOpMap[roomId]:
			logs.Info("in here")
			if op == "start" {
				if nowGameMatchMap[roomId].Id == 0 || nowGameMatchMap[roomId].GameStatus == models.GAME_STATUS_END {
					startRoomGame(roomId)
				}
			}
			if op == "close" {
				models.CloseRoom(roomId)
				return
			}
		case roundInfo := <-roundProcessMap[roomId]:
			switch roundInfo.Type {
			case models.EVENT_ROUND_INFO:
				tmp := roundUserDetailMap[roomId][roundInfo.NowPosition]
				tmp.AllowOp = tmp.AllowOp[0:0]
				if LimitPointMap[roomId] == 0 {
					tmp.AllowOp = append(tmp.AllowOp, "check")
				}
				if tmp.RoundPoint+tmp.Point >= LimitPointMap[roomId] && LimitPointMap[roomId] > 0 {
					tmp.AllowOp = append(tmp.AllowOp, "call")
				}
				if tmp.RoundPoint+tmp.Point > LimitPointMap[roomId] {
					tmp.AllowOp = append(tmp.AllowOp, "raise")
				}
				tmp.AllowOp = append(tmp.AllowOp, "allin")
				tmp.AllowOp = append(tmp.AllowOp, "fold")

				roundUserDetailMap[roomId][roundInfo.NowPosition] = tmp
				roundInfo.Detail = roundUserDetailMap[roomId][roundInfo.NowPosition]
				roundInfo.AllPointInRound = getRoomRoundPoint(roomId, false)
				sendToRoomUser(roomId, roundInfo)
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
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position, roomId)
			case models.GAME_OP_CALL: //call
				uop.GameMatchLog.PointNumber = LimitPointMap[roomId] - roundUserDetailMap[roomId][uop.Position].RoundPoint
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position, roomId)
			case models.GAME_OP_CHECK: //check
				RoundCheckNumberMap[roomId]++
				fromCheck = true
			case models.GAME_OP_ALLIN: // allin
				userRePoint := models.GetUserPoint(roundUserDetailMap[roomId][uop.Position].UserId)
				uop.GameMatchLog.PointNumber = userRePoint
				roomOpChangePoint(uop.GameMatchLog.PointNumber, uop.Position, roomId)
				gameMatchAllinMap[roomId][uop.Position] = roundUserDetailMap[roomId][uop.Position].RoundPoint
			case models.GAME_OP_FOLD:
				FoldPointMap[roomId][nowGameMatchMap[roomId].GameStatus] += roundUserDetailMap[roomId][uop.Position].RoundPoint
				logs.Info(foldPoint)
				delete(roundUserDetailMap[roomId], uop.Position)
				delete(UsersCardMap[roomId], uop.Position)
			}
			models.AddGameMatchLog(uop.GameMatchLog)
			sendToRoomUser(roomId, uop)
			SendUserPointInfoToRoom(roomId)
			// logs.Info("after user operation")
			// logs.Info(roundUserDetail)
			if len(roundUserDetailMap[roomId]) <= 1 {
				//game end
				RoomGameEnd(roomId)
			} else {
				if fromCheck {
					if RoundCheckNumberMap[roomId] < len(roundUserDetailMap[roomId]) {
						PositionTurnMap[roomId] = getNextPositionInRoom(roomId, PositionTurnMap[roomId])
					} else {
						endRoomRoundPoint(roomId)
					}
					nextRoomRoundInfo(roomId)
				} else {
					var haveNotFillPoint bool
					haveNotFillPoint = false
					for _, v := range roundUserDetailMap[roomId] {
						_, ok := gameMatchAllinMap[roomId][v.Position]
						if v.RoundPoint < LimitPointMap[roomId] && !ok {
							haveNotFillPoint = true
						}
					}
					if BigBindPositionTurnMap[roomId] <= 1 { //只在第一轮中有效
						haveNotFillPoint = true
					}

					if !haveNotFillPoint && (len(roundUserDetailMap[roomId])-len(gameMatchAllinMap[roomId]) <= 1) {
						RoomGameEnd(roomId)
					} else if haveNotFillPoint {
						//next user
						logs.Info(emptySendMap[roomId])
						logs.Info(PositionTurnMap[roomId])
						if emptySendMap[roomId] > 2 { //只是为了排除大小盲
							PositionTurnMap[roomId] = getNextPositionInRoom(roomId, PositionTurnMap[roomId])
							var ok bool
							_, ok = gameMatchAllinMap[roomId][positionTurn]
							for ok {
								PositionTurnMap[roomId] = getNextPositionInRoom(roomId, PositionTurnMap[roomId])
								_, ok = gameMatchAllinMap[roomId][positionTurn]
							}
							nextRoomRoundInfo(roomId)
						}
					} else {
						//end this round
						endRoomRoundPoint(roomId)
						nextRoomRoundInfo(roomId)
					}
				}
			}
		}
	}
}

func startRoomGame(RoomId int) {
	logs.Info(RoomId)
	//游戏结束时清零 仅用于记录大小盲
	emptySendMap[RoomId] = 0
	BigBindPositionTurnMap[RoomId] = 0
	//清除上局剩余的用户信息，例如卡牌信息等
	SendUserPointInfoToRoom(RoomId)
	nowGameMatchMap[RoomId] = models.InitGameMatch(RoomId)
	initCardMap(RoomId)
	gameMatchAllinMap[RoomId] = make(map[int]int)

	changeGameMatch(RoomId, models.GAME_STATUS_LICENSING)

	tmpCard := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}
	sendToRoomUser(RoomId, tmpCard)
	sendCardToUser(RoomId)

	for i := 1; i <= 5; i++ {
		var tmpPoker = models.GetOneCardFromCardMap(GameAllCardMap[RoomId])
		PublicCardMap[RoomId][i] = *tmpPoker
	}
	changeGameMatch(RoomId, models.GAME_STATUS_ROUND1)

	//初始化座位
	DetailArrMap[RoomId] = models.GetRoundUserDetail(RoomId)
	roundUserDetailMap[RoomId] = make(map[int]models.InRoundUserDetail)
	logs.Info("init seat")
	for _, v := range DetailArrMap[RoomId] {
		roundUserDetailMap[RoomId][v.Position] = v
	}
	logs.Info(roundUserDetailMap[RoomId])

	//小盲
	PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].SmallBindPosition
	userOpMsg := models.UserOperationMsg{
		Type:     models.EVENT_USER_OPERATION_INFO,
		Position: nowGameMatchMap[RoomId].SmallBindPosition,
		Name:     roundUserDetailMap[RoomId][nowGameMatchMap[RoomId].SmallBindPosition].Name,
		GameMatchLog: models.GameMatchLog{
			GameMatchId: nowGameMatchMap[RoomId].Id,
			UserId:      roundUserDetailMap[RoomId][nowGameMatchMap[RoomId].SmallBindPosition].UserId,
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
	PositionTurnMap[RoomId]++
	logs.Info(PositionTurnMap)
	logs.Info(roundUserDetailMap)
	if PositionTurnMap[RoomId] > len(roundUserDetailMap[RoomId]) {
		PositionTurnMap[RoomId] = 1
	}
	if _, ok := roundUserDetailMap[RoomId][PositionTurnMap[RoomId]]; !ok {
		PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].SmallBindPosition
	}
	logs.Info(PositionTurnMap)
	nextRoomRoundInfo(RoomId)
}

func roomOpChangePoint(point int, position int, roomId int) {
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

func initCardMap(roomId int) {
	GameAllCardMap[roomId] = models.GetNewCardMap("short")
	PublicCardMap[roomId] = make(map[int]models.Card)
	UsersCardMap[roomId] = make(map[int][]models.Card)
}

func changeGameMatch(roomId int, status string) {
	tmp := nowGameMatchMap[roomId]
	tmp.GameStatus = status
	models.UpdateGameMatchStatus(tmp, "game_status")
	nowGameMatchMap[roomId] = tmp
}

func sendToRoomUser(roomId int, data interface{}) {
	positionMap := models.GetRoomUserPositionList(roomId)
	for _, uid := range positionMap {
		sendInWs(UserConnMap[uid], data)
	}
}

func sendInWs(ws *websocket.Conn, data interface{}) {
	if ws != nil {
		msg, _ := json.Marshal(data)
		if ws.WriteMessage(websocket.TextMessage, msg) != nil {
			// User disconnected.
			//unsubscribe <- ss.Value.(Player).user
		}
	}
}
func sendCardToUser(roomId int) {
	positionMap := models.GetRoomUserPositionList(roomId)
	for pos, uid := range positionMap {
		var tmpPoker = models.GetOneCardFromCardMap(GameAllCardMap[roomId])
		UsersCardMap[roomId][pos] = append(UsersCardMap[roomId][pos], *tmpPoker)
		sendInWs(UserConnMap[uid], sendCard(models.EVENT_LICENSING, "", pos, *tmpPoker))

		var tmpPoker2 = models.GetOneCardFromCardMap(GameAllCardMap[roomId])
		UsersCardMap[roomId][pos] = append(UsersCardMap[roomId][pos], *tmpPoker2)
		sendInWs(UserConnMap[uid], sendCard(models.EVENT_LICENSING, "", pos, *tmpPoker2))
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
		winUserPos = CalWinUserInRoom(RoomId)
	}
	//计算获胜点数
	nngmm := nowGameMatchMap[RoomId]
	nngmm.PotAll = nngmm.Pot1st + nngmm.Pot2nd + nngmm.Pot3rd + nngmm.Pot4th
	for _, v := range roundUserDetailMap[RoomId] {
		nngmm.PotAll += v.RoundPoint
	}
	logs.Info(foldPoint)
	for _, v := range foldPoint {
		nngmm.PotAll += v
	}

	if len(gameMatchAllinMap[RoomId]) > 0 {
		potAll := nngmm.PotAll                                          //总池
		allInPointDesc := models.RankByPoint(gameMatchAllinMap[RoomId]) //根据allin数量进行的排序
		var calUserDetail = make(map[int]models.InRoundUserDetail)
		for k, v := range roundUserDetailMap[RoomId] {
			calUserDetail[k] = v
		}
		calAllInPot := 0 //已结算的allin底池
		for _, v := range allInPointDesc {
			needCalPot := (v.Value - calAllInPot) * len(calUserDetail)
			if needCalPot > potAll {
				needCalPot = potAll
			}
			winUser, _ := GetBigUserInRoom(calUserDetail, RoomId)
			perPot := needCalPot / len(winUser)
			for _, v := range winUser {
				models.ChangeUserPoint(roundUserDetailMap[RoomId][v].UserId, perPot)
				pointDetail[v] += perPot
			}
			calAllInPot = v.Value
			potAll -= needCalPot
			delete(calUserDetail, v.Key)

		}
		if potAll > 0 { //cal_user_detail中还剩余多个未allin玩家时未出现
			var winUser2 []int
			winUser2, _ = GetBigUserInRoom(roundUserDetailMap[RoomId], RoomId)

			perPot := potAll / len(winUser2)
			for _, v := range winUser2 {
				models.ChangeUserPoint(roundUserDetailMap[RoomId][v].UserId, perPot)
				pointDetail[v] += perPot
			}
		}

	} else if len(winUserPos) > 0 {
		perPot := nngmm.PotAll / len(winUserPos)
		for _, v := range winUserPos {
			models.ChangeUserPoint(roundUserDetailMap[RoomId][v].UserId, perPot)
			pointDetail[v] = perPot
		}
	}

	nngmm.GameStatus = models.GAME_STATUS_END
	models.UpdateGameMatchStatus(nngmm, "game_status")
	models.UpdateGameMatchStatus(nngmm, "pot1st")
	models.UpdateGameMatchStatus(nngmm, "pot2nd")
	models.UpdateGameMatchStatus(nngmm, "pot3rd")
	models.UpdateGameMatchStatus(nngmm, "pot4th")
	models.UpdateGameMatchStatus(nngmm, "pot_all")
	nowGameMatchMap[RoomId] = nngmm
	//提示胜利玩家可以重新开始游戏了
	logs.Info("game end")
	tmpCard := models.CardInfo{
		Type: models.EVENT_CLEAR_CARD,
	}

	SendUserPointInfoToRoom(RoomId)
	sendToRoomUser(RoomId, tmpCard)
	sendToRoomUser(RoomId, newEvent(models.EVENT_GAME_END, "system", "Game End"))

}

func CalWinUserInRoom(RoomId int) []int {
	tmpCardC := make(map[int]string)
	bigString := ""
	var winUser []int
	for k := range roundUserDetailMap[RoomId] {
		tmpArr := UsersCardMap[RoomId][k]
		for _, v := range PublicCardMap[RoomId] {
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
		result := models.Compare(v, bigString, "short")
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
		winCard = append(winCard, UsersCardMap[RoomId][v])
	}

	showMaxCard := models.TransMaxHandToCardInfo(bigString, "short")
	var tmp []models.Card
	for _, v := range PublicCardMap[RoomId] {
		tmp = append(tmp, v)
	}
	a := models.EndGameCardInfo{
		Type:       models.EVENT_GAME_END_SHOW_CARD,
		WinPos:     winUser,
		WinCard:    winCard,
		PublicCard: tmp,
		// BigCard:    models.StringToCard(bigString),\
		BigCard:   showMaxCard,
		UserCards: UsersCardMap[RoomId],
	}
	sendToRoomUser(RoomId, a)

	return winUser
}

//传入当前仍在场的用户
func GetBigUserInRoom(iru map[int]models.InRoundUserDetail, RoomId int) ([]int, string) {
	tmpCardC := make(map[int]string)
	bigString := ""
	var winUser uint64
	for k := range iru {
		tmpArr := UsersCardMap[RoomId][k]
		for _, v := range PublicCardMap[RoomId] {
			tmpArr = append(tmpArr, v)
		}
		tmpCardC[k] = models.GetString(tmpArr)
		if bigString == "" {
			bigString = models.GetString(tmpArr)
			winUser = 1 << k
		}
	}

	for k, v := range tmpCardC {
		if bigString == v {
			continue
		}
		result := models.Compare(v, bigString, "short")
		if result == 1 {
			winUser = 1 << k
			bigString = v
		} else if result == 0 {
			winUser |= 1 << k
		}
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
	return winUserArr, bigString
}

/**
标志着一轮的结束
**/
func endRoomRoundPoint(RoomId int) {
	RoundCheckNumberMap[RoomId] = 0
	if _, ok := roundUserDetailMap[RoomId][nowGameMatchMap[RoomId].SmallBindPosition]; ok {
		PositionTurnMap[RoomId] = nowGameMatchMap[RoomId].SmallBindPosition
	} else {
		PositionTurnMap[RoomId] = getNextPositionInRoom(RoomId, nowGameMatchMap[RoomId].SmallBindPosition)
	}
	LimitPointMap[RoomId] = 0
	var nextRound string
	tmpGameMatch := nowGameMatchMap[RoomId]
	switch tmpGameMatch.GameStatus {
	case models.GAME_STATUS_ROUND1:
		nextRound = models.GAME_STATUS_ROUND2
		tmpGameMatch.Pot1st = getRoomRoundPoint(RoomId, true)
		sendToRoomUser(RoomId, sendCard(models.EVENT_PUBLIC_CARD, "", 0, PublicCardMap[RoomId][1]))
		sendToRoomUser(RoomId, sendCard(models.EVENT_PUBLIC_CARD, "", 0, PublicCardMap[RoomId][2]))
		sendToRoomUser(RoomId, sendCard(models.EVENT_PUBLIC_CARD, "", 0, PublicCardMap[RoomId][3]))
	case models.GAME_STATUS_ROUND2:
		nextRound = models.GAME_STATUS_ROUND3
		tmpGameMatch.Pot2nd = getRoomRoundPoint(RoomId, true)
		sendToRoomUser(RoomId, sendCard(models.EVENT_PUBLIC_CARD, "", 0, PublicCardMap[RoomId][4]))
	case models.GAME_STATUS_ROUND3:
		nextRound = models.GAME_STATUS_ROUND4
		tmpGameMatch.Pot3rd = getRoomRoundPoint(RoomId, true)
		sendToRoomUser(RoomId, sendCard(models.EVENT_PUBLIC_CARD, "", 0, PublicCardMap[RoomId][5]))
	case models.GAME_STATUS_ROUND4:
		nextRound = models.GAME_STATUS_END
		tmpGameMatch.Pot4th = getRoomRoundPoint(RoomId, true)
	}

	logs.Info("end round" + nowGameMatch.GameStatus)
	tmpGameMatch.GameStatus = nextRound

	if tmpGameMatch.GameStatus == models.GAME_STATUS_END {
		//需比较剩余玩家的卡牌大小
		RoomGameEnd(RoomId)
	}
	nowGameMatchMap[RoomId] = tmpGameMatch
}

func getRoomRoundPoint(RoomId int, clear bool) int {
	var remainRoundPoint int
	for k, v := range roundUserDetailMap[RoomId] {
		remainRoundPoint += v.RoundPoint
		if clear {
			a := roundUserDetailMap[RoomId][k]
			a.RoundPoint = 0
			roundUserDetailMap[RoomId][k] = a
		}
	}
	remainRoundPoint += FoldPointMap[RoomId][nowGameMatchMap[RoomId].GameStatus]
	if clear {
		FoldPointMap[RoomId][nowGameMatch.GameStatus] = 0
	}

	switch nowGameMatchMap[RoomId].GameStatus {
	case "ROUND1":
		return nowGameMatchMap[RoomId].Pot1st + remainRoundPoint
	case "ROUND2":
		return nowGameMatchMap[RoomId].Pot2nd + remainRoundPoint
	case "ROUND3":
		return nowGameMatchMap[RoomId].Pot3rd + remainRoundPoint
	case "ROUND4":
		return nowGameMatchMap[RoomId].Pot4th + remainRoundPoint
	}
	return 0
}

func nextRoomRoundInfo(RoomId int) {
	tmp := roundUserDetailMap[RoomId][PositionTurnMap[RoomId]]
	tmp.AllowOp = tmp.AllowOp[0:0]
	if LimitPointMap[RoomId] == 0 {
		tmp.AllowOp = append(tmp.AllowOp, "check")
	}
	if tmp.RoundPoint+tmp.Point >= LimitPointMap[RoomId] && LimitPointMap[RoomId] > 0 {
		tmp.AllowOp = append(tmp.AllowOp, "call")
	}
	if tmp.RoundPoint+tmp.Point > LimitPointMap[RoomId] {
		tmp.AllowOp = append(tmp.AllowOp, "raise")
	}
	tmp.AllowOp = append(tmp.AllowOp, "allin")
	tmp.AllowOp = append(tmp.AllowOp, "fold")

	roundUserDetailMap[RoomId][PositionTurnMap[RoomId]] = tmp
	logs.Info(roundUserDetailMap)
	data := models.RoundInfo{
		Type:        models.EVENT_ROUND_INFO,
		GM:          nowGameMatchMap[RoomId],
		NowPosition: PositionTurnMap[RoomId],
		MaxPoint:    LimitPointMap[RoomId],
	}

	data.Detail = tmp
	data.AllPointInRound = getRoomRoundPoint(RoomId, false)
	sendToRoomUser(RoomId, data)
}

//用于查找下一个位置的人，有则返回座位号，无则返回false
func getNextPositionInRoom(RoomId int, InitPosition int) int {
	userNumber := len(DetailArrMap[RoomId])
	for i := InitPosition + 1; i <= InitPosition+userNumber; i++ {
		key := i % userNumber
		if key == 0 {
			key = userNumber
		}
		if _, ok := roundUserDetailMap[RoomId][key]; ok {
			return key
		}
	}
	return 0
}

func SendUserPointInfoToRoom(RoomId int) {
	var data UserInfoMsg
	data.Type = models.EVENT_REFRESH_USER_INFO
	data.Info = models.GetUserPointWithSeat(RoomId)
	sendToRoomUser(RoomId, data)
}
