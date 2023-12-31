package config

import (
	"EmptyClassroom/logs"
	"context"
	"encoding/json"
	"io"
	"os"
)

const (
	ConfigPathKey = "CONFIG_PATH"
)

type CampusConfig struct {
	Name         string `json:"name"`
	Id           int    `json:"id,omitempty"`
	HasRealtime  bool   `json:"has_realtime"`
	ReplaceRegex []struct {
		Regex   string `json:"regex"`
		Replace string `json:"replace"`
	} `json:"replace_regex"`
}

type NotificationConfig struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	Duration         int    `json:"duration"`
	Type             string `json:"type"`
	ShowNotification bool   `json:"showNotification"`
	Start            string `json:"start"`
	End              string `json:"end"`
}

type Config struct {
	ClassTable   ClassTableConfig   `json:"class_table"`
	Campus       []CampusConfig     `json:"campus"`
	Notification NotificationConfig `json:"notification"`
}

type ClassTableConfig struct {
	StartWeek     string                `json:"start_week"`
	EndWeek       string                `json:"end_week"`
	UnableReason  string                `json:"unable_reason"`
	IsAvailable   bool                  `json:"is_available"`
	ClassTableMap map[string]ClassTable `json:"class_table_map"`
}

type ClassTable struct {
	Class   []ClassTableClassroomInfo `json:"class"`
	TypeMap map[string]string         `json:"typeMap"`
}

type ClassTableClassroomInfo struct {
	Campus  string    `json:"campus"`
	Seat    string    `json:"seat"`
	Name    string    `json:"name"`
	Classes [][][]int `json:"classes"`
}

var GlobalConfig *Config

func InitConfig() {
	configPath := os.Getenv(ConfigPathKey)
	if configPath == "" {
		configPath = "config"
	}
	_, err := os.Stat(configPath + "/config.json")
	if err != nil {
		logs.CtxError(context.Background(), "stat config file failed: %v", err)
		panic(err)
	}
	configFile, err := os.Open(configPath + "/config.json")
	if err != nil {
		logs.CtxError(context.Background(), "open config file failed: %v", err)
		panic(err)
	}
	configContent, err := io.ReadAll(configFile)
	if err != nil {
		logs.CtxError(context.Background(), "read config file failed: %v", err)
		panic(err)
	}
	GlobalConfig = new(Config)
	err = json.Unmarshal(configContent, &GlobalConfig)
	if err != nil {
		logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
		panic(err)
	}
	for _, building := range GlobalConfig.Campus {
		_, err := os.Stat(configPath + "/" + building.Name + ".json")
		if err != nil {
			logs.CtxError(context.Background(), "stat config file failed: %v", err)
			panic(err)
		}
		configFile, err := os.Open(configPath + "/" + building.Name + ".json")
		if err != nil {
			logs.CtxError(context.Background(), "open configPathconfig file failed: %v", err)
			panic(err)
		}
		configContent, err := io.ReadAll(configFile)
		if err != nil {
			logs.CtxError(context.Background(), "read config file failed: %v", err)
			panic(err)
		}
		buildingConfig := new(ClassTable)
		err = json.Unmarshal(configContent, &buildingConfig)
		if err != nil {
			logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
			panic(err)
		}
		if GlobalConfig.ClassTable.ClassTableMap == nil {
			GlobalConfig.ClassTable.ClassTableMap = make(map[string]ClassTable)
		}
		GlobalConfig.ClassTable.ClassTableMap[building.Name] = *buildingConfig
	}
	_, err = os.Stat(configPath + "/notification.json")
	if err != nil {
		logs.CtxError(context.Background(), "stat config file failed: %v", err)
		panic(err)
	}
	configFile, err = os.Open(configPath + "/notification.json")
	if err != nil {
		logs.CtxError(context.Background(), "open config file failed: %v", err)
		panic(err)
	}
	configContent, err = io.ReadAll(configFile)
	if err != nil {
		logs.CtxError(context.Background(), "read config file failed: %v", err)
		panic(err)
	}
	notificationConfig := new(NotificationConfig)
	err = json.Unmarshal(configContent, &notificationConfig)
	if err != nil {
		logs.CtxError(context.Background(), "unmarshal config file failed: %v", err)
		panic(err)
	}
	GlobalConfig.Notification = *notificationConfig
}

func GetConfig() Config {
	if GlobalConfig == nil {
		InitConfig()
	}
	return *GlobalConfig
}
