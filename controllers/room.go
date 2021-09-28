package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/gorilla/websocket"
	"net/http"
	"poke/models"
	"strconv"
)

type RoomController struct {
	baseController
}

type JsonResponse struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var user models.User

func (r *RoomController) Prepare() {
	s, _ := beego.AppConfig.String("session_name")
	sessionData := r.GetSession(s)
	//logs.Info(sessionData)
	if sessionData == nil {
		r.Redirect("/", 302)
		return
	} else {
		user = sessionData.(models.User)
	}

	// Reset language option.
	r.Lang = "" // This field is from i18n.Locale.

	// 1. Get language information from 'Accept-Language'.
	al := r.Ctx.Request.Header.Get("Accept-Language")
	if len(al) > 4 {
		al = al[:5] // Only compare first 5 letters.
		if i18n.IsExist(al) {
			r.Lang = al
		}
	}

	// 2. Default language is English.
	if len(r.Lang) == 0 {
		r.Lang = "en-US"
	}

	// Set template level language option.
	r.Data["Lang"] = r.Lang
}

func (r *RoomController) Get() {
	r.TplName = "room/room_list.html"
	r.Data["UserName"] = user.Name
	r.Data["RoomList"] = models.GetOnlineRoom()
	//logs.Info(r.Data["RoomList"])
}

func (r *RoomController) Create() {
	r.TplName = "room/room_add.html"
	r.Data["UserName"] = user.Name
}

func (r *RoomController) Close() {
	roomID, err := strconv.Atoi(r.Ctx.Input.Param(":id"))

	var jsonData JsonResponse
	if err != nil {
		jsonData.Code = 100
		jsonData.Msg = "error input"
		r.Data["json"] = &jsonData
		r.ServeJSON()
	}

	roomManageCloseList <- roomID
	jsonData.Code = 0
	jsonData.Msg = "success"
	r.Data["json"] = &jsonData
	r.ServeJSON()
}

func (r *RoomController) EntryRoom() {
	roomID := r.Ctx.Input.Param(":id")

	//logs.Info(roomID)
	r.TplName = "room/game_room.html"
	r.Data["IsGameRoom"] = true
	r.Data["User"] = user
	r.Data["RoomID"] = roomID
}

func (r *RoomController) RoomSocket() {
	user := r.GetSession("USER")
	roomID, err := strconv.Atoi(r.Ctx.Input.Param(":id"))
	//logs.Info(user)
	if user == nil || err != nil {
		r.Redirect("/", 302)
		return
	}

	u := user.(models.User)

	//ws, err := websocket.Upgrade(r.Ctx.ResponseWriter, r.Ctx.Request, nil, 1024, 1024)
	upgrade := websocket.Upgrader{
		HandshakeTimeout: 10,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}
	ws, err := upgrade.Upgrade(r.Ctx.ResponseWriter, r.Ctx.Request, nil)
	if err != nil {
		return
	}

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(r.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logs.Error("Cannot setup WebSocket connection:", err)
		return
	}

	UserConnMap[u.Id] = ws

	gu := models.SetUserIntoRoom(u, roomID)
	msg, _ := json.Marshal(models.SeatInfo{
		Type:     models.EVENT_JOIN,
		GameUser: gu,
		User:     "",
	})

	ws.WriteMessage(websocket.TextMessage, msg)

	defer Leave(u)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		// publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
		data := new(models.ClientMessage)
		err2 := json.Unmarshal(p, data)
		if err2 != nil {
			return
		}
		logs.Info(data)
		switch data.Type {
		case "game_op":
			gameOpMap[roomID] <- data.Message
		case "user_op":
			top := models.UserOperationMsg{
				Type:     models.EVENT_USER_OPERATION_INFO,
				Position: data.Position,
				Name:     data.Name,
				GameMatchLog: models.GameMatchLog{
					UserId:      data.UserId,
					Operation:   data.Operation,
					PointNumber: data.Point,
				},
			}
			// logs.Info("user op")
			// logs.Info(top)
			userOperationProcessMap[roomID] <- top
		}

	}
}

func (r *RoomController) Post() {
	roomName := r.GetString("room_name")
	roomPassword := r.GetString("room_password")
	roomCardType := r.GetString("room_card_type")
	room := models.Room{
		CreateUserId: user.Id,
		RoomName:     roomName,
		RoomPassword: roomPassword,
	}
	roomId := models.CreateRoom(&room)
	if roomId > 0 {
		roomManageOpenList <- int(roomId)
		r.Redirect("/room/entry/"+strconv.FormatInt(roomId, 10), 302)
	} else {
		r.Redirect("/room", 302)
		return
	}
}

func sendChan(roomId int) {
	//logs.Info("164")
	roomManageOpenList <- roomId
}
