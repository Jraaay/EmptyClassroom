package service

import (
	"EmptyClassroom/cache"
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"testing"
)

func init() {
	logs.Init(false)
	config.InitConfig()
	cache.InitCache()
}

func TestLogin(t *testing.T) {
	err := Login(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestQueryOne(t *testing.T) {
	err := Login(context.Background())
	if err != nil {
		t.Error(err)
	}
	_, err = QueryOne(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryAll(t *testing.T) {
	err := Login(context.Background())
	if err != nil {
		t.Error(err)
	}
	ans, err := QueryAll(context.Background())
	if err != nil {
		t.Error(err)
	}
	marshal, err := json.Marshal(ans)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
