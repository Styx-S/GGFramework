package main

import (
	"GGFramework/Define"
	"GGFramework/Push"
	"GGFramework/Room"
	"github.com/gin-gonic/gin"
)

var svr = Define.Server{
	ModuleMap: make(map[string]Define.Module),
}

func main() {

	svr.AddModules(Define.ModuleNamePush, Push.ModuleInstance)
	svr.AddModules(Define.ModuleNameRoom, Room.ModuleInstance)

	r := gin.Default()
	ConfigureRouter(r)
	panic(r.Run())
}
