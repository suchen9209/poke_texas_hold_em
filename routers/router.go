package routers

import (
	"github.com/beego/beego/v2/server/web/filter/cors"
	"poke/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		// 允许访问所有源
		//AllowAllOrigins: true,
		// 可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 指的是允许的Header的种类
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		// 如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
	}))
	// Register routers.
	beego.Router("/", &controllers.AppController{}, "get:Index")
	// beego.Router("/", &controllers.AppController{})
	// Indicate AppController.Join method to handle POST requests.
	beego.Router("/join", &controllers.AppController{}, "post:Join")
	beego.Router("/login", &controllers.AppController{}, "post:Login")

	beego.Router("/v2/user/login", &controllers.AppController{}, "post:JsonLogin")
	beego.Router("/v2/user/register", &controllers.AppController{}, "post:JsonRegister")

	beego.Router("/v2/room/list", &controllers.RoomController{}, "post:RoomList")

	beego.Router("/room", &controllers.RoomController{})
	beego.Router("/room/add", &controllers.RoomController{}, "get:Create")
	beego.Router("/room/entry/:id", &controllers.RoomController{}, "get:EntryRoom")
	beego.Router("/room/join/:id", &controllers.RoomController{}, "get:RoomSocket")
	beego.Router("/room/close/:id", &controllers.RoomController{}, "post:Close")

	beego.Router("/greedisgood", &controllers.AppController{}, "get:GreedIsGood")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")
}
