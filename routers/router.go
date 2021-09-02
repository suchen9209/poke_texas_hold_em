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

	beego.Router("/room", &controllers.RoomController{})
	beego.Router("/room/add", &controllers.RoomController{}, "get:Create")
	beego.Router("/room/entry/:id", &controllers.RoomController{}, "get:EntryRoom")

	beego.Router("/greedisgood", &controllers.AppController{}, "get:GreedIsGood")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")
}
