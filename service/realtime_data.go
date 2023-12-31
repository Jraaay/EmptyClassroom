package service

import (
	"EmptyClassroom/cache"
	"EmptyClassroom/config"
	"EmptyClassroom/logs"
	"EmptyClassroom/service/model"
	"EmptyClassroom/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LoginUrl = "http://jwglweixin.bupt.edu.cn/bjyddx/login"
	QueryUrl = "http://jwglweixin.bupt.edu.cn/bjyddx/todayClassrooms?campusId=0"

	LoginUsernameKey = "JW_USERNAME"
	LoginPasswordKey = "JW_PASSWORD"

	TodayCacheKey = "TODAY_CACHE"
)

var (
	Token string
)

func Login(ctx context.Context) error {
	userNo := os.Getenv(LoginUsernameKey)
	pwd := os.Getenv(LoginPasswordKey)
	req := map[string]string{
		"userNo":      userNo,
		"pwd":         pwd,
		"encode":      "1",
		"captchaData": "",
		"codeVal":     "",
	}
	code, _, body, err := utils.HttpPostForm(ctx, LoginUrl, req)
	if err != nil {
		logs.CtxError(ctx, "login failed: %v", err)
		return err
	}
	if code != 200 {
		logs.CtxError(ctx, "login failed - code not 200: %v", err)
		return errors.New("login failed")
	}
	var resp model.LoginResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.CtxError(ctx, "login failed - resp unmarshal failed: %v", err)
		return err
	}
	if resp.Code != "1" {
		logs.CtxError(ctx, "login failed - code not 1: %v", err)
		return errors.New("login failed")
	}
	Token = resp.Data.Token
	return nil
}

func QueryOne(ctx context.Context, id int) ([]model.JWClassInfo, error) {
	errorTime := 0
	err := Login(ctx)
	// 重试3次
	for err != nil && errorTime < 3 {
		time.Sleep(10 * time.Second)
		err = Login(ctx)
		errorTime++
	}
	if err != nil {
		logs.CtxError(ctx, "login failed: %v", err)
	}
	if err == nil && errorTime > 0 {
		logs.CtxWarn(ctx, "login retry success, error time: %v", errorTime)
	}
	header := map[string]string{
		"token": Token,
	}
	code, _, body, err := utils.HttpGetWithHeader(ctx, QueryUrl+strconv.FormatInt(int64(id), 10), header)
	if err != nil {
		logs.CtxError(ctx, "query failed: %v", err)
		return nil, err
	}
	if code != 200 {
		logs.CtxError(ctx, "query failed - code not 200: %v", err)
		return nil, errors.New("query failed")
	}
	var resp model.QueryResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		logs.CtxError(ctx, "query failed - resp unmarshal failed: %v", err)
		return nil, err
	}
	if resp.Code != "1" {
		logs.CtxError(ctx, "query failed - code not 1: %v", err)
		return nil, errors.New("query failed")
	}
	return resp.Data, nil
}

func QueryAll(ctx context.Context) (classInfo *model.ClassInfo, err error) {
	classInfo = &model.ClassInfo{
		UpdateAt: time.Now(),
	}
	sysConfig := config.GetConfig()
	for _, campus := range sysConfig.Campus {
		err = ProcessClassTableInfo(ctx, classInfo, campus.Name)
		if err != nil {
			logs.CtxError(ctx, "process failed: %v", err)
			return nil, err
		}
		if campus.HasRealtime {
			errorTime := 0
			jwClassInfo, err := QueryOne(ctx, campus.Id)
			// 重试3次
			for err != nil && errorTime < 3 {
				time.Sleep(10 * time.Second)
				jwClassInfo, err = QueryOne(ctx, campus.Id)
				errorTime++
			}
			if err != nil {
				logs.CtxError(ctx, "query failed: %v", err)
			}
			if err == nil && errorTime > 0 {
				logs.CtxWarn(ctx, "query retry success, error time: %v", errorTime)
			}
			// 即使查询报错也不返回，用课表数据进行兜底
			err = ProcessJWClassInfo(ctx, jwClassInfo, classInfo, campus)
			if err != nil {
				logs.CtxError(ctx, "process failed: %v", err)
				return nil, err
			}
		}
	}
	startTime, _ := time.Parse("2006-01-02 15:04:05", sysConfig.Notification.Start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", sysConfig.Notification.End)
	if time.Now().After(startTime) && time.Now().Before(endTime) {
		classInfo.Notification = &sysConfig.Notification
	} else {
		classInfo.Notification = nil
	}
	classTableStartWeek, _ := time.Parse("2006-01-02", sysConfig.ClassTable.StartWeek)
	classTableEndWeek, _ := time.Parse("2006-01-02", sysConfig.ClassTable.EndWeek)
	if time.Now().Before(classTableStartWeek) || time.Now().After(classTableEndWeek.AddDate(0, 0, 1)) {
		classInfo.ClassTable = nil
	} else {
		classInfo.ClassTable = &sysConfig.ClassTable
	}
	oldClassInfoRaw, cacheTime, ok := cache.GetCacheWithExpiration(TodayCacheKey)
	if ok {
		if len(oldClassInfoRaw.(*model.ClassInfo).CampusInfoMap) > len(classInfo.CampusInfoMap) && cacheTime.After(time.Now().Add(10*time.Minute)) {
			return oldClassInfoRaw.(*model.ClassInfo), nil
		}
	}
	cache.SetCache(TodayCacheKey, classInfo, 60*time.Minute)
	return classInfo, nil
}

func ProcessJWClassInfo(ctx context.Context, jwClassInfo []model.JWClassInfo, classInfo *model.ClassInfo, campusConfig config.CampusConfig) error {
	sysConfig := config.GetConfig()
	if jwClassInfo == nil {
		return nil
	}
	campusInfo := model.CampusInfo{
		Name:            campusConfig.Name,
		BuildingInfoMap: map[int]*model.BuildingInfo{},
		BuildingIdMap:   map[string]int{},
		MaxBuildingId:   0,
	}
	if classInfo.CampusInfoMap != nil && classInfo.CampusInfoMap[campusConfig.Name] != nil {
		campusInfo = *classInfo.CampusInfoMap[campusConfig.Name]
	}
	campusClassTableConfig := sysConfig.ClassTable.ClassTableMap[campusConfig.Name]
	for _, info := range jwClassInfo {
		classroomList := strings.Split(info.Classrooms, ",")
		for _, classroom := range classroomList {
			for _, replaceConfig := range campusConfig.ReplaceRegex {
				re, err := regexp.Compile(replaceConfig.Regex)
				if err != nil {
					logs.CtxError(ctx, "regex compile failed: %v", err)
					return err
				}
				classroom = re.ReplaceAllString(classroom, replaceConfig.Replace)
			}
			classroomInfo := model.ClassroomInfo{}
			if len(strings.Split(strings.Split(classroom, "(")[0], "-")) != 2 {
				logs.CtxWarn(ctx, "classroom format error: %v", classroom)
				continue
			}
			classroomInfo.Name = strings.Split(strings.Split(classroom, "(")[0], "-")[1]
			classroomInfo.Size, _ = strconv.ParseInt(strings.Split(strings.Split(classroom, "(")[1], ")")[0], 10, 32)
			classroomInfo.CanTrust = true
			buildingName := strings.Split(classroom, "-")[0]
			if _, ok := campusInfo.BuildingIdMap[buildingName]; !ok {
				campusInfo.BuildingIdMap[buildingName] = campusInfo.MaxBuildingId
				campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId] = &model.BuildingInfo{
					Name:             buildingName,
					ClassroomInfoMap: map[int]*model.ClassroomInfo{},
					ClassroomIdMap:   map[string]int{},
					ClassMatrix:      [][]int{},
					MaxClassroomId:   0,
				}
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix = append(campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix, []int{})
				}
				campusInfo.MaxBuildingId++
			}
			buildingId := campusInfo.BuildingIdMap[buildingName]
			if _, ok := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]; !ok {
				campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name] = campusInfo.BuildingInfoMap[buildingId].MaxClassroomId
				classroomInfo.BuildingId = buildingId
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId] = &classroomInfo
				campusInfo.BuildingInfoMap[buildingId].MaxClassroomId++
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i] = append(campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i], 1)
				}
			} else if !campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]].CanTrust {
				// 覆盖
				classroomInfo.BuildingId = buildingId
				classroomType, typeOk := campusClassTableConfig.TypeMap[classroomInfo.Name]
				if typeOk {
					classroomInfo.Type = classroomType
				} else {
					classroomInfo.Type = campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]].Type
				}
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]] = &classroomInfo
				for i := 0; i < 14; i++ {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i][campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]] = 1
				}
			}
			classroomId := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomInfo.Name]
			nodeName, err := strconv.ParseInt(info.NodeName, 10, 32)
			if err != nil {
				logs.CtxWarn(ctx, "node name parse failed: %v", err)
				continue
			}
			campusInfo.BuildingInfoMap[buildingId].ClassMatrix[nodeName-1][classroomId] = 0
		}
	}
	if classInfo.CampusInfoMap == nil {
		classInfo.CampusInfoMap = map[string]*model.CampusInfo{}
	}
	classInfo.CampusInfoMap[campusConfig.Name] = &campusInfo
	return nil
}

func ProcessClassTableInfo(ctx context.Context, classInfo *model.ClassInfo, campusName string) error {
	sysConfig := config.GetConfig()
	classTableStartWeek, err := time.Parse("2006-01-02", sysConfig.ClassTable.StartWeek)
	if err != nil {
		logs.CtxError(ctx, "start week parse failed: %v", err)
		return err
	}
	classTableEndWeek, err := time.Parse("2006-01-02", sysConfig.ClassTable.EndWeek)
	if err != nil {
		logs.CtxError(ctx, "end week parse failed: %v", err)
		return err
	}
	if time.Now().Before(classTableStartWeek) || time.Now().After(classTableEndWeek.AddDate(0, 0, 1)) {
		return nil
	}
	nowWeek := int((time.Now().Unix() - classTableStartWeek.Unix()) / 604800)
	today := int(time.Now().Weekday())

	campusClassTableConfig := sysConfig.ClassTable.ClassTableMap[campusName]
	campusInfo := model.CampusInfo{
		Name:            campusName,
		BuildingInfoMap: map[int]*model.BuildingInfo{},
		BuildingIdMap:   map[string]int{},
		MaxBuildingId:   0,
	}
	if classInfo.CampusInfoMap != nil && classInfo.CampusInfoMap[campusName] != nil {
		campusInfo = *classInfo.CampusInfoMap[campusName]
	}
	for _, classItemInfo := range campusClassTableConfig.Class {
		buildingName := strings.Split(classItemInfo.Name, "-")[0]
		if _, ok := campusInfo.BuildingIdMap[buildingName]; !ok {
			campusInfo.BuildingIdMap[buildingName] = campusInfo.MaxBuildingId
			campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId] = &model.BuildingInfo{
				Name:             buildingName,
				ClassroomInfoMap: map[int]*model.ClassroomInfo{},
				ClassroomIdMap:   map[string]int{},
				ClassMatrix:      [][]int{},
				MaxClassroomId:   0,
			}
			for i := 0; i < 14; i++ {
				campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix = append(campusInfo.BuildingInfoMap[campusInfo.MaxBuildingId].ClassMatrix, []int{})
			}
			campusInfo.MaxBuildingId++
		}
		buildingId := campusInfo.BuildingIdMap[buildingName]
		classroomName := strings.Split(classItemInfo.Name, "-")[1]
		if _, ok := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName]; !ok {
			campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName] = campusInfo.BuildingInfoMap[buildingId].MaxClassroomId
			campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId] = &model.ClassroomInfo{
				Name:       classroomName,
				Size:       0,
				CanTrust:   false,
				BuildingId: buildingId,
			}
			classroomSize, err := strconv.ParseInt(classItemInfo.Seat, 10, 32)
			if err == nil {
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId].Size = classroomSize
			}
			classroomType, typeOk := campusClassTableConfig.TypeMap[classItemInfo.Name]
			if typeOk {
				campusInfo.BuildingInfoMap[buildingId].ClassroomInfoMap[campusInfo.BuildingInfoMap[buildingId].MaxClassroomId].Type = classroomType
			}
			campusInfo.BuildingInfoMap[buildingId].MaxClassroomId++
			for i := 0; i < 14; i++ {
				campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i] = append(campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i], 0)
			}
		} else {
			// 跳过
			continue
		}
		classroomId := campusInfo.BuildingInfoMap[buildingId].ClassroomIdMap[classroomName]
		for i := 0; i < 14; i++ {
			for _, week := range classItemInfo.Classes[i][today] {
				if week == nowWeek {
					campusInfo.BuildingInfoMap[buildingId].ClassMatrix[i][classroomId] = 1
				}
			}
		}
	}
	if classInfo.CampusInfoMap == nil {
		classInfo.CampusInfoMap = map[string]*model.CampusInfo{}
	}
	classInfo.CampusInfoMap[campusName] = &campusInfo
	return nil
}

func GetData(ctx context.Context, c *gin.Context) {
	classInfoRaw, ok := cache.GetCache(TodayCacheKey)
	if ok {
		classInfo := classInfoRaw.(*model.ClassInfo)
		c.JSON(200, gin.H{
			"code": 0,
			"data": classInfo,
		})
		return
	} else {
		classInfo, err := QueryAll(ctx)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "query failed",
				"data": nil,
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"data": classInfo,
		})
		return
	}
}
