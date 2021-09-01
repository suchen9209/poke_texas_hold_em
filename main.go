package main

import (
	_ "poke/routers"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

const (
	AppVer = "1.0.0"
)

func main() {

	// m123 := make(map[int]int)
	// m123[5] = 100
	// m123[2] = 77
	// m123[666] = 999
	// m123[43] = 722
	// i32 := models.RankByPoint(m123)
	// for _, v := range i32 {
	// 	logs.Info(v)
	// }
	// logs.Info(i32)
	// return

	logs.SetLogger("console")
	logs.Info(beego.BConfig.AppName + AppVer)

	beego.AddFuncMap("i18n", i18n.Tr)

	beego.Run()

}
