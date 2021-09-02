package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"poke/models"
	"strconv"
)

type RoomController struct {
	baseController
}

var user models.User

func (r *RoomController) Prepare() {
	s, _ := beego.AppConfig.String("session_name")
	sessionData := r.GetSession(s)
	logs.Info(sessionData)
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
	logs.Info(r.Data["RoomList"])
}

func (r *RoomController) Create() {
	r.TplName = "room/room_add.html"
	r.Data["UserName"] = user.Name
}

func (r *RoomController) EntryRoom() {
	roomID := r.Ctx.Input.Param(":id")

	logs.Info(roomID)
	r.TplName = "websocket.html"
	r.Data["IsWebSocket"] = true
	r.Data["UserName"] = user.Name
	r.Data["Point"] = user.Point
	r.Data["RoomID"] = roomID
}

func (r *RoomController) Post() {
	roomName := r.GetString("room_name")
	roomPassword := r.GetString("room_password")
	room := models.Room{
		CreateUserId: user.Id,
		RoomName:     roomName,
		RoomPassword: roomPassword,
	}
	roomId := models.CreateRoom(&room)
	if roomId > 0 {
		go chatroom_new(int(roomId))
		r.Redirect("/room/entry/"+strconv.FormatInt(roomId, 10), 302)
	} else {
		r.Redirect("/room", 302)
		return
	}
}
