package routers

import (
	"poke/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Register routers.
	beego.Router("/", &controllers.AppController{}, "get:Index")
	// beego.Router("/", &controllers.AppController{})
	// Indicate AppController.Join method to handle POST requests.
	beego.Router("/join", &controllers.AppController{}, "post:Join")
	beego.Router("/login", &controllers.AppController{}, "post:Login")

	beego.Router("/v2/user/login", &controllers.AppController{}, "post:JsonLogin")
	beego.Router("/v2/user/register", &controllers.AppController{}, "post:JsonRegister")

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
