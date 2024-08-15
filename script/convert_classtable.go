package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Class struct {
	Campus  string    `json:"campus"`
	Seat    string    `json:"seat"`
	Name    string    `json:"name"`
	Classes [][][]int `json:"classes"`
}

type ClassTable struct {
	ClassList []Class           `json:"class"`
	TypeMap   map[string]string `json:"typeMap"`
}

type CampusNameReConfig []struct {
	Src  string
	Repl string
}

type CampusIgnoreCondig []struct {
	Re string
}

type CampusConfig struct {
	Name               string
	CampusNameReConfig CampusNameReConfig
	CampusIgnoreCondig CampusIgnoreCondig
}

func main() {
	xtcConfig := CampusConfig{
		Name: "西土城",
		CampusNameReConfig: CampusNameReConfig{
			{
				Src:  `主-(.+)`,
				Repl: `主楼-$1`,
			},
			{
				Src:  `(\d)-(.+)`,
				Repl: `教$1-$2`,
			},
			{
				Src:  `动画室`,
				Repl: `教2-动画室`,
			},
			{
				Src:  `新科研-(.+)`,
				Repl: `新科研楼-$1`,
			},
			{
				Src:  `明-(.+)`,
				Repl: `明光楼-$1`,
			},
			{
				Src:  `明光楼(\d.+)`,
				Repl: `明光楼-$1`,
			},
			{
				Src:  `经管楼(\d.+)`,
				Repl: `经管楼-$1`,
			},
			{
				Src:  `(东.)-(\d+)`,
				Repl: `本部图书馆-$1$2`,
			},
			{
				Src:  `图书馆一层`,
				Repl: `本部图书馆-一层`,
			},
			{
				Src:  `工程管理仿真中心`,
				Repl: `学10-工程管理仿真中心`,
			},
		},
		CampusIgnoreCondig: CampusIgnoreCondig{
			{
				Re: `网络教室`,
			},
			{
				Re: `自行安排`,
			},
			{
				Re: `虚拟教室`,
			},
			{
				Re: `体育`,
			},
			{
				Re: `工程管理仿真中心`,
			},
		},
	}
	shConfig := CampusConfig{
		Name: "沙河",
		CampusNameReConfig: CampusNameReConfig{
			{
				Src:  `（.*`,
				Repl: ``,
			},
			{
				Src:  `^([NS])(\d+)$`,
				Repl: `$1-$2`,
			},
			{
				Src:  `图-(.+)`,
				Repl: `沙河图书馆-$1`,
			},
			{
				Src:  `沙河图书馆东配楼`,
				Repl: `沙河图书馆-东配楼_`,
			},
			{
				Src:  `D-([NS])(.*)`,
				Repl: `教学实验综合楼-D_$1$2`,
			},
			{
				Src:  `邮政楼D1-(.*)`,
				Repl: `邮政楼-D1_$1`,
			},
			{
				Src:  `^D1-(.*)$`,
				Repl: `邮政楼-D1_$1`,
			},
			{
				Src:  `办-(.*)`,
				Repl: `综合办公楼-$1`,
			},
			{
				Src:  `地下一层自动化`,
				Repl: `教学实验综合楼-地下一层自动化`,
			},
			{
				Src:  `报告厅`,
				Repl: `教学实验综合楼-报告厅`,
			},
			{
				Src:  `学-(.*)`,
				Repl: `学生活动中心-$1`,
			},
			{
				Src:  `(..楼)(\d+)`,
				Repl: `$1-$2`,
			},
			{
				Src:  `电路中心(\d+)`,
				Repl: `电路中心实验楼-$1`,
			},
		},
		CampusIgnoreCondig: CampusIgnoreCondig{
			{
				Re: `网络教室`,
			},
			{
				Re: `自行安排`,
			},
			{
				Re: `虚拟教室`,
			},
			{
				Re: `体育`,
			},
			{
				Re: `旧`,
			},
		},
	}
	hnConfig := CampusConfig{
		Name: "海南",
		CampusNameReConfig: CampusNameReConfig{
			{
				Src:  `(..楼)(\d+)室?`,
				Repl: `$1-$2`,
			},
		},
		CampusIgnoreCondig: CampusIgnoreCondig{
			{
				Re: `网络教室`,
			},
			{
				Re: `体育`,
			},
		},
	}
	convert(xtcConfig)
	convert(shConfig)
	convert(hnConfig)
}

func convert(campusConfig CampusConfig) {
	f, err := excelize.OpenFile("script/raw/" + campusConfig.Name + ".xlsx")
	if err != nil {
		println(err.Error())
		return
	}
	// 获取工作表中指定单元格的值
	cell := f.GetSheetList()
	rows, err := f.GetRows(cell[0])
	classTable := ClassTable{
		ClassList: []Class{},
		TypeMap:   make(map[string]string),
	}
	classInfoRe, _ := regexp.Compile(`.*\s教室.\s?([^\s]+)\s*所属功能区.\s?([^\s(（]+)[(（](.*)[)）].*所属教学楼.\s?([^\s]+).*座位数.\s?([^\s]+)`)
	classWeekRe, _ := regexp.Compile(`([^\s]+)周\[`)
	for i := 0; i < len(rows); i++ {
		class := Class{}
		i++
		classInfo := classInfoRe.FindAllStringSubmatch(rows[i][1], -1)
		if len(classInfo) == 0 {
			fmt.Println(rows[i][1])
		}
		class.Campus = classInfo[0][3]
		class.Name = classInfo[0][1]
		isIgnore := false
		for _, reConfig := range campusConfig.CampusIgnoreCondig {
			ignoreRe, _ := regexp.Compile(reConfig.Re)
			if ignoreRe.MatchString(class.Name) {
				isIgnore = true
				break
			}
		}
		if isIgnore {
			i += 16
			continue
		}
		for _, reConfig := range campusConfig.CampusNameReConfig {
			nameRe, _ := regexp.Compile(reConfig.Src)
			if nameRe.MatchString(class.Name) {
				class.Name = nameRe.ReplaceAllString(class.Name, reConfig.Repl)
			}
		}
		class.Seat = classInfo[0][5]
		classTable.TypeMap[class.Name] = classInfo[0][2]
		i++
		i++
		classes := [][][]int{}
		for j := i; j < i+14; j++ {
			var nowClass [][]int
			for weekday := 1; weekday <= 7; weekday++ {
				if weekday >= len(rows[j]) || rows[j][weekday] == "" {
					nowClass = append(nowClass, []int{})
					continue
				}
				classWeeksRaw := classWeekRe.FindAllStringSubmatch(rows[j][weekday], -1)
				classWeeks := []int{}
				for _, classWeek := range classWeeksRaw {
					fromToList := strings.Split(classWeek[1], ",")
					for _, fromToRaw := range fromToList {
						fromTo := strings.Split(fromToRaw, "-")
						if len(fromTo) == 1 {
							week, err := strconv.ParseInt(fromTo[0], 10, 32)
							if err != nil {
								fmt.Println(err.Error())
							}
							classWeeks = append(classWeeks, int(week))
						} else {
							from, err := strconv.ParseInt(fromTo[0], 10, 32)
							if err != nil {
								fmt.Println(err.Error())
							}
							to, err := strconv.ParseInt(fromTo[1], 10, 32)
							weekType := 0
							if err != nil {
								if strings.HasSuffix(fromTo[1], "单") {
									weekType = 1
									fromTo[1] = strings.Replace(fromTo[1], "单", "", -1)
									to, err = strconv.ParseInt(fromTo[1], 10, 32)
									if err != nil {
										fmt.Println(err.Error())
									}
								} else if strings.HasSuffix(fromTo[1], "双") {
									weekType = 2
									fromTo[1] = strings.Replace(fromTo[1], "双", "", -1)
									if err != nil {
										fmt.Println(err.Error())
									}
								} else {
									fmt.Println(err.Error())
								}
							}
							for week := from; week <= to; week++ {
								if weekType == 1 && week%2 == 0 {
									continue
								}
								if weekType == 2 && week%2 == 1 {
									continue
								}
								classWeeks = append(classWeeks, int(week))
							}
						}
					}
				}
				nowClass = append(nowClass, classWeeks)
			}
			classes = append(classes, nowClass)
		}
		class.Classes = classes
		i += 14
		classTable.ClassList = append(classTable.ClassList, class)
	}
	jsonByte, _ := json.Marshal(classTable)
	os.WriteFile("script/classtable_config/"+campusConfig.Name+".json", jsonByte, os.ModePerm)
}
