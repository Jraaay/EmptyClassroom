package service

import (
	"EmptyClassroom/logs"
	"EmptyClassroom/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	LarkWebhookKey = "LARK_WEBHOOK"
)

type ElementType struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type TitleType struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type HeaderType struct {
	Template string    `json:"template"`
	Title    TitleType `json:"title"`
}

type CardType struct {
	Elements []ElementType `json:"elements"`
	Header   HeaderType    `json:"header"`
}

type LarkReq struct {
	MsgType string   `json:"msg_type"`
	Card    CardType `json:"card"`
}

type LarkResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Report(ctx context.Context, c *gin.Context) {
	err := ReportToLark(ctx, c)
	if err != nil {
		logs.CtxError(ctx, "ReportToLark error: %v", err)
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "ReportToLark error",
		})
		return
	}
}

func ReportToLark(ctx context.Context, c *gin.Context) error {
	reqJson := make(map[string]interface{})
	err := c.BindJSON(&reqJson)
	if err != nil {
		logs.CtxError(ctx, "c.BindJSON error: %v", err)
		return err
	}
	// 取出text
	text, ok := reqJson["text"]
	if !ok {
		logs.CtxError(ctx, "bodyJson[\"text\"] not found")
		return err
	}
	textString, ok := text.(string)
	if !ok {
		logs.CtxError(ctx, "text.(string) error")
		return err
	}
	logs.CtxInfo(ctx, "Report text: %s", textString)

	// 从环境变量获取飞书机器人webhook
	webhook := os.Getenv(LarkWebhookKey)
	if webhook == "" {
		logs.CtxError(ctx, "webhook == \"\"")
		return errors.New("webhook == \"\"")
	}
	req := LarkReq{
		MsgType: "interactive",
		Card: CardType{
			Elements: []ElementType{
				{
					Tag:     "markdown",
					Content: textString,
				},
			},
			Header: HeaderType{
				Template: "blue",
				Title: TitleType{
					Content: "用户反馈",
					Tag:     "plain_text",
				},
			},
		},
	}
	jsonStr, _ := json.Marshal(req)
	_, _, body, err := utils.HttpPostJson(ctx, webhook, jsonStr)
	if err != nil {
		logs.CtxError(ctx, "utils.HttpPostJson error: %v", err)
		return err
	}
	resp := LarkResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.CtxError(ctx, "json.Unmarshal error: %v", err)
		return err
	}
	if resp.Code != 0 {
		logs.CtxError(ctx, "resp.Code != 0, resp: %v", resp)
		return errors.New(resp.Msg)
	}
	return nil
}
