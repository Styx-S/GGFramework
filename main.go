package main

import (
	"GGFramework/Push"
	"github.com/gin-gonic/gin"
)

func main() {
	go Push.CtxManager.Start()

	r := gin.Default()
	ConfigureRouter(r)
	panic(r.Run())
}
