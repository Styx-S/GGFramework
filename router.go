package main

import (
	"GGFramework/Room"
	"github.com/gin-gonic/gin"
)


func ConfigureRouter(engine *gin.Engine) {
	engine.POST("/room/join", Room.JoinRoom)
}
