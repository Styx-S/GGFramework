package main

import (
	"GGFramework/Push"
	"GGFramework/Room"
	"github.com/gin-gonic/gin"
)

func ConfigureRouter(engine *gin.Engine) {

	engine.GET("/websocket", Push.ConnectWebsocket)

	engine.POST("/room/join", Room.JoinRoom)
	engine.POST("/room/leave", Room.LeaveRoom)
	engine.POST("/room/ready", Room.ReadyRoom)
}
