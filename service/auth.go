package service

import (
	"crypto/md5"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"poke/models"
)

var o orm.Ormer

func init() {

}

type Auth struct {
	Name   string
	Strict uint64
}

func (a *Auth) AddUser(c *beego.Controller) models.JsonData {
	uname := c.GetString("uname")
	pwd := c.GetString("password")
	pwd2 := c.GetString("password2")

	data := models.JsonData{
		Code: 0,
		Msg:  "",
	}

	if pwd != pwd2 {
		data.Code = 0
		data.Msg = "Two Password Not Equal"
		return data
	}

	salt, _ := beego.AppConfig.String("password_salt")

	u := models.User{
		Name:     uname,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(pwd+salt))),
	}

	err := models.AddUser(&u)
	if err != nil {
		return models.JsonData{
			Code: 10050,
			Msg:  "Add Failed",
		}
	}
	err2 := c.SetSession("USER", u)
	if err2 != nil {
		return models.JsonData{
			Code: 10060,
			Msg:  "Session Failed",
		}
	}
	return models.JsonData{
		Code: 0,
		Msg:  "Success",
		Data: u,
	}

}

func (a *Auth) CheckUser(c *beego.Controller) models.JsonData {
	uname := c.GetString("uname")
	pwd := c.GetString("password")

	salt, _ := beego.AppConfig.String("password_salt")

	u := models.User{
		Name:     uname,
		Password: fmt.Sprintf("%x", md5.Sum([]byte(pwd+salt))),
	}
	//logs.Info(u)
	data := models.JsonData{
		Code: 10040,
		Msg:  "Not Found",
	}
	err := models.CheckUserLogin(&u)
	if err != nil {
		data.Code = 10060
		data.Msg = err.Error()
		return data
	}
	//logs.Info(u)
	if u.Id > 0 {
		err := c.SetSession("USER", u)
		if err != nil {
			data.Code = 10070
			data.Msg = err.Error()
			return data
		}
		c.Data["json"] = models.JsonData{
			Code: 0,
			Msg:  "Success",
		}
		data.Code = 0
		data.Msg = "Success"
		data.Data = u
		return data
	}
	return data
}
