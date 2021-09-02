package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"poke/models"
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
}

func (r *RoomController) Get() {
	r.TplName = "room/room_list.html"
	r.Data["UserName"] = user.Name
}

func (r *RoomController) Create() {
	r.TplName = "room/room_add.html"
	r.Data["UserName"] = user.Name
}

func (r *RoomController) Post() {
	roomName := r.GetString("room_name")
	roomPassword := r.GetString("room_password")
}
