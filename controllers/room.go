package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"poke/models"
)

type RoomController struct {
	baseController
}

func (r *RoomController) Get() {
	s, _ := beego.AppConfig.String("session_name")
	user := r.GetSession(s)
	logs.Info(user)
	if user == nil {
		r.Redirect("/", 302)
		return
	}

	r.TplName = "room_list.html"
	r.Data["UserName"] = user.(models.User).Name
}
