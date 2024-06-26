/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/11 下午6:12.
 *  * Author: guojia(https://github.com/guojia99)
 */

package model

type Project string
type ProjectClass string
type RouteType string

const (
	ProjectClassWCA        = "WCA项目"
	ProjectClassWCABF      = "WCA盲拧"
	ProjectClassWCAAlien   = "WCA异型"
	ProjectClassWCACube    = "WCA正阶"
	ProjectClassXCube      = "趣味项目"
	ProjectClassXCubeBF    = "趣味盲拧"
	ProjectClassXCubeOH    = "单手"
	ProjectClassXCubeFm    = "最少步"
	ProjectClassXCubeRelay = "趣味连拧"
	ProjectClassNotCube    = "趣味非魔方"
	ProjectClassDigit      = "记字"
	ProjectClassSuperHigh  = "超高阶"
	ProjectClassBFGroup    = "盲拧群系列"
)

var ProjectClassList = []ProjectClass{
	ProjectClassWCA,
	ProjectClassWCABF,
	ProjectClassWCAAlien,
	ProjectClassWCACube,
	ProjectClassXCube,
	ProjectClassXCubeBF,
	ProjectClassXCubeOH,
	ProjectClassXCubeFm,
	ProjectClassXCubeRelay,
	ProjectClassNotCube,
	ProjectClassDigit,
	ProjectClassSuperHigh,
	ProjectClassBFGroup,
}

const (
	RouteType1rounds      RouteType = "1_r"      // 单轮项目
	RouteType3roundsBest  RouteType = "3_r_b"    // 三轮取最佳
	RouteType3roundsAvg   RouteType = "3_r_a"    // 三轮取平均
	RouteType5roundsBest  RouteType = "5_r_b"    // 五轮取最佳
	RouteType5roundsAvg   RouteType = "5_r_a"    // 五轮取平均
	RouteType5RoundsAvgHT RouteType = "5_r_a_ht" // 五轮去头尾取平均
	RouteTypeRepeatedly   RouteType = "ry"       // 单轮多次还原项目, 成绩1:还原数; 成绩2:尝试数; 成绩3:时间;
)

func (r RouteType) String() string {
	switch r {
	case RouteType1rounds:
		return "单轮项目"
	case RouteType3roundsBest:
		return "三轮取最佳"
	case RouteType3roundsAvg:
		return "三轮取平均"
	case RouteType5roundsBest:
		return "五轮取最佳"
	case RouteType5roundsAvg:
		return "五轮取平均"
	case RouteType5RoundsAvgHT:
		return "五轮去头尾取平均"
	case RouteTypeRepeatedly:
		return "单轮多次还原项目"
	}
	return string(r)
}

func AllProjectRoute() []Project       { return allProjectRoute }
func WCAProjectRoute() []Project       { return wcaProjectRoute }
func XCubeProjectRoute() []Project     { return xCubeProjectRoute }
func (p Project) String() string       { return string(p) }
func (p Project) Cn() string           { return projectItemsMap[p].Cn }
func (p Project) RouteType() RouteType { return projectItemsMap[p].RouteType }
func (p Project) IsWca() bool          { return projectItemsMap[p].IsWca }
func (p Project) Class() []string      { return projectItemsMap[p].Class }

func AllProjectItem() []projectItem { return projectsItems }

type SorStatisticsKey = string

const (
	SorWCA             = "wca"
	SorXCube           = "xcube"
	SorWCACubeLowLevel = "wca2345"
	SorWCACubeAllLevel = "wca234567"
	SorWCAAlien        = "wca_alien"
	SorWCA333          = "wca333"
	SorWCABf           = "wca_bf"
)

func SorKeyMap() map[SorStatisticsKey][]Project {
	return map[SorStatisticsKey][]Project{
		SorWCA:             wcaProjectRoute,
		SorXCube:           xCubeProjectRoute,
		SorWCACubeLowLevel: {Cube333, Cube222, Cube444, Cube555},
		SorWCACubeAllLevel: {Cube333, Cube222, Cube444, Cube555, Cube666, Cube777},
		SorWCAAlien:        {CubeSk, CubePy, CubeSq1, CubeMinx, CubeClock},
		SorWCA333:          {Cube333, Cube333OH, Cube333Ft, Cube333BF, Cube333MBF, Cube333FM},
		SorWCABf:           {Cube333BF, Cube444BF, Cube555BF, Cube333MBF},
	}
}
