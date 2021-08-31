package controllers

import (
	"encoding/json"
	"net/http"

	"poke/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	baseController
}

func (w *WebSocketController) Get() {
	user := w.GetSession("USER")
	logs.Info(user)
	if user == nil {
		w.Redirect("/", 302)
		return
	}

	w.TplName = "websocket.html"
	w.Data["IsWebSocket"] = true
	w.Data["UserName"] = user.(models.User).Name
	w.Data["Point"] = user.(models.User).Point
}

func (w *WebSocketController) Join() {
	user := w.GetSession("USER")
	logs.Info(user)
	if user == nil {
		w.Redirect("/", 302)
		return
	}

	u := user.(models.User)

	ws, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logs.Error("Cannot setup WebSocket connection:", err)
		return
	}

	_, inGame := models.CheckUserInGame(u.Name)
	if inGame {
		Join(u.Name, ws, models.VIEWER, u) //此名字玩家已在游戏中，观战视角
	} else {
		Join(u.Name, ws, models.POKER_PLAYER, u)
	}

	defer Leave(u)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		// publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
		data := new(models.ClientMessage)
		json.Unmarshal(p, data)
		// logs.Info(data)
		switch data.Type {
		case "game_op":
			gameop <- data.Message
		case "user_op":
			top := models.UserOperationMsg{
				Type:     models.EVENT_USER_OPERATION_INFO,
				Position: data.Position,
				Name:     data.Name,
				GameMatchLog: models.GameMatchLog{
					UserId:      data.UserId,
					Operation:   data.Operation,
					PointNumber: data.Point,
				},
			}
			// logs.Info("user op")
			// logs.Info(top)
			userOperationProcess <- top
		}

	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		logs.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).User
			}
		}
	}
}
