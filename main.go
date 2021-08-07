package main

import (
	_ "poke/routers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

const (
	APP_VER = "1.0.0"
)

func main() {
	logs.SetLogger("console")
	logs.Info(beego.BConfig.AppName + APP_VER)

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()
}
