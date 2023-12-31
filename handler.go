package main

import (
	"EmptyClassroom/logs"
	"EmptyClassroom/service"
	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	ctx := logs.GetContextFromGinContext(c)
	logs.CtxInfo(ctx, "GetData")
	service.GetData(ctx, c)
}

func Report(c *gin.Context) {
	ctx := logs.GetContextFromGinContext(c)
	service.Report(ctx, c)
}
