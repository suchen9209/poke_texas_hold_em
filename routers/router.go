package routers

import (
	"poke/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Register routers.
	beego.Router("/", &controllers.AppController{})
	// Indicate AppController.Join method to handle POST requests.
	beego.Router("/join", &controllers.AppController{}, "post:Join")

	beego.Router("/greedisgood", &controllers.AppController{}, "get:GreedIsGood")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
	beego.Router("/ws/join", &controllers.WebSocketController{}, "get:Join")
}
