package models

import (
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

var o orm.Ormer

func init() {
	dbuser, _ := beego.AppConfig.String("dbuser")
	pwd, _ := beego.AppConfig.String("dbpassword")
	if pwd != "" {
		dbuser = dbuser + ":" + pwd
	}
	dbhost, _ := beego.AppConfig.String("dbhost")
	dbname, _ := beego.AppConfig.String("dbname")
	// set default database
	err := orm.RegisterDataBase("default", "mysql", dbuser+"@tcp("+dbhost+":3306)/"+dbname+"?charset=utf8&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	// register model
	// orm.RegisterModel(new(Game))

	// create table
	// orm.RunSyncdb("default", false, true)
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Game))
	orm.RegisterModel(new(GameMatch))
	orm.RegisterModel(new(GameMatchLog))
	orm.RegisterModel(new(GameUser))
	orm.RegisterModel(new(Room))
	o = orm.NewOrm()
}
