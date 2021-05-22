package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	ConfigureRouter(r)

	panic(r.Run())
}
