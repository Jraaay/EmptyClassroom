package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	Init()
	r := gin.Default()
	SetRouter(r)
	r.Run()
}
