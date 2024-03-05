package model

import (
	"EmptyClassroom/config"
	"time"
)

type LoginResponse struct {
	Code string `json:"code"`
	Msg  string `json:"Msg"`
	Data struct {
		Birthday     string `json:"birthday"`
		AcademyName  string `json:"academyName"`
		UserNo       string `json:"userNo"`
		EntranceYear string `json:"entranceYear"`
		ClsName      string `json:"clsName"`
		Name         string `json:"name"`
		UserType     string `json:"userType"`
		Token        string `json:"token"`
	} `json:"data"`
}

type JWClassInfo struct {
	Classrooms string `json:"CLASSROOMS"`
	NodeTime   string `json:"NODETIME"`
	NodeName   string `json:"NODENAME"`
}

type QueryResponse struct {
	Code string        `json:"code"`
	Msg  string        `json:"Msg"`
	Data []JWClassInfo `json:"data"`
}

type ClassroomInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	CanTrust   bool   `json:"can_trust"`
	BuildingId int    `json:"building_id"`
	Type       string `json:"type"`
}

type BuildingInfo struct {
	Name             string                 `json:"name"`
	ClassroomInfoMap map[int]*ClassroomInfo `json:"classroom_info_map"`
	ClassroomIdMap   map[string]int         `json:"classroom_id_map"`
	ClassMatrix      [][]int                `json:"class_matrix"`
	MaxClassroomId   int                    `json:"max_classroom_id"`
}

type CampusInfo struct {
	Name            string                `json:"name"`
	BuildingInfoMap map[int]*BuildingInfo `json:"building_info_map"`
	BuildingIdMap   map[string]int        `json:"building_id_map"`
	MaxBuildingId   int                   `json:"max_building_id"`
}

type ClassInfo struct {
	CampusInfoMap map[string]*CampusInfo     `json:"campus_info_map"`
	ClassTable    *config.ClassTableConfig   `json:"class_table"`
	UpdateAt      time.Time                  `json:"update_at"`
	Notification  *config.NotificationConfig `json:"notification"`
	IsFallback    map[string]bool            `json:"is_fallback"`
}
