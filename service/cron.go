package service

import (
	"EmptyClassroom/logs"
	"github.com/robfig/cron/v3"
)

var (
	GlobalCron *cron.Cron
)

func Cronjob() {
	ctx := logs.GenNewContext()
	_, err := QueryAll(ctx)
	if err != nil {
		logs.CtxError(ctx, "QueryAll error: %v", err)
	} else {
		logs.CtxInfo(ctx, "QueryAll success")
	}
}

func StartCron() {
	GlobalCron = cron.New()
	// 5分钟执行一次
	_, err := GlobalCron.AddFunc("*/5 * * * *", Cronjob)
	if err != nil {
		logs.CtxError(nil, "GlobalCron.AddFunc error: %v", err)
		panic(err)
	}
	GlobalCron.Start()
	// 立即异步执行一次
	go Cronjob()
}
