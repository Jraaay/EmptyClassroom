package main

import (
	"EmptyClassroom/cache"
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"EmptyClassroom/service"
	"EmptyClassroom/service/model"
	"encoding/gob"
)

func Init() {
	gob.Register(&model.ClassInfo{})
	logs.Init(true)
	config.InitConfig()
	cache.InitCache()
	service.StartCron()
}
