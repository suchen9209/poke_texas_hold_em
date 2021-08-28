// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"encoding/json"
	"net/http"

	"poke/models"

	"github.com/beego/beego/v2/core/logs"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	baseController
}

// // Get method handles GET requests for WebSocketController.
// func (w *WebSocketController) Get() {
// 	// Safe check.
// 	// uname := this.GetString("uname")
// 	user := w.GetSession("USER")
// 	logs.Info(user)
// 	if user == nil {
// 		w.Redirect("/", 302)
// 		return
// 	}

// 	w.TplName = "websocket.html"
// 	w.Data["IsWebSocket"] = true
// 	w.Data["UserName"] = user.(models.User).Name
// 	w.Data["Point"] = user.(models.User).Point
// }

// for test
func (w *WebSocketController) Get() {
	// Safe check.
	uname := w.GetString("uname")
	if len(uname) == 0 {
		w.Redirect("/", 302)
		return
	}

	w.TplName = "websocket.html"
	w.Data["IsWebSocket"] = true
	w.Data["UserName"] = uname
}

//for test
func (w *WebSocketController) Join() {
	uname := w.GetString("uname")
	if len(uname) == 0 {
		w.Redirect("/", 302)
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		logs.Error("Cannot setup WebSocket connection:", err)
		return
	}

	user, reg := models.CheckUser(uname)
	if reg {
		_, inGame := models.CheckUserInGame(uname)
		if inGame {
			Join(uname, ws, models.VIEWER, *user) //观战视角
		} else {
			Join(uname, ws, models.POKER_PLAYER, *user)
		}

	} else {
		//游客退出暂无
		Join(uname, ws, models.VIEWER, *user) //观战视角
	}

	defer Leave(*user)

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

// Join method handles WebSocket requests for WebSocketController.
// func (w *WebSocketController) Join() {
// 	user := w.GetSession("USER")
// 	logs.Info(user)
// 	if user == nil {
// 		w.Redirect("/", 302)
// 		return
// 	}

// 	u := user.(models.User)

// 	// Upgrade from http request to WebSocket.
// 	ws, err := websocket.Upgrade(w.Ctx.ResponseWriter, w.Ctx.Request, nil, 1024, 1024)
// 	if _, ok := err.(websocket.HandshakeError); ok {
// 		http.Error(w.Ctx.ResponseWriter, "Not a websocket handshake", 400)
// 		return
// 	} else if err != nil {
// 		logs.Error("Cannot setup WebSocket connection:", err)
// 		return
// 	}

// 	_, inGame := models.CheckUserInGame(u.Name)
// 	if inGame {
// 		Join(u.Name, ws, models.VIEWER, u) //此名字玩家已在游戏中，观战视角
// 	} else {
// 		Join(u.Name, ws, models.POKER_PLAYER, u)
// 	}

// 	defer Leave(u)

// 	for {
// 		_, p, err := ws.ReadMessage()
// 		if err != nil {
// 			return
// 		}
// 		// publish <- newEvent(models.EVENT_MESSAGE, uname, string(p))
// 		data := new(models.ClientMessage)
// 		json.Unmarshal(p, data)
// 		// logs.Info(data)
// 		switch data.Type {
// 		case "game_op":
// 			gameop <- data.Message
// 		case "user_op":
// 			top := models.UserOperationMsg{
// 				Type:     models.EVENT_USER_OPERATION_INFO,
// 				Position: data.Position,
// 				Name:     data.Name,
// 				GameMatchLog: models.GameMatchLog{
// 					UserId:      data.UserId,
// 					Operation:   data.Operation,
// 					PointNumber: data.Point,
// 				},
// 			}
// 			// logs.Info("user op")
// 			// logs.Info(top)
// 			userOperationProcess <- top
// 		}

// 	}
// }

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
