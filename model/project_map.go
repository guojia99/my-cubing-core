package model

type projectItem struct {
	Project   Project
	Cn        string
	RouteType RouteType
	IsWca     bool
	Class     []string
}

const (
	// 特殊
	JuBaoHaoHao Project = "jhh"
	OtherCola   Project = "o_cola"

	// wca
	Cube222    Project = "222"
	Cube333    Project = "333"
	Cube444    Project = "444"
	Cube555    Project = "555"
	Cube666    Project = "666"
	Cube777    Project = "777"
	CubeSk     Project = "skewb"
	CubePy     Project = "pyram"
	CubeSq1    Project = "sq1"
	CubeMinx   Project = "minx"
	CubeClock  Project = "clock"
	Cube333OH  Project = "333oh"
	Cube333FM  Project = "333fm"
	Cube333BF  Project = "333bf"
	Cube444BF  Project = "444bf"
	Cube555BF  Project = "555bf"
	Cube333MBF Project = "333mbf"
	Cube333Ft  Project = "333ft"

	// 趣味最少步
	XCubePyFm    Project = "pyram_fm"
	XCubeSkFm    Project = "skewb_fm"
	XCubeClockFm Project = "clock_fm"
	XCube222Fm   Project = "222fm"

	// 盲
	XCube222BF   Project = "222bf"
	XCube666BF   Project = "666bf"
	XCube777BF   Project = "777bf"
	XCubePyBF    Project = "pyram_bf"
	XCubeSkBF    Project = "skewb_bf"
	XCubeMinxBf  Project = "minx_bf"
	XCubeClockBf Project = "clock_bf"
	XCubeSQ1Bf   Project = "sq1_bf"

	// 单手
	XCube333BfOH   Project = "333bf_oh" // 三单盲
	XCube444BfOH   Project = "444bf_oh"
	XCube555BfOH   Project = "555bf_oh"
	XCube222OH     Project = "222oh"
	XCube333MiniOH Project = "333mini_oh"
	XCube444OH     Project = "444oh"
	XCube555OH     Project = "555oh"
	XCube666OH     Project = "666oh"
	XCube777OH     Project = "777oh"
	XCubeSkOH      Project = "skewb_oh"
	XCubePyOH      Project = "pyram_oh"
	XCubeSq1OH     Project = "sq1_oh"
	XCubeMinxOH    Project = "minx_oh"
	XCubeClockOH   Project = "clock_oh"

	// 连
	XCube333Multiple5      Project = "333multiple5"
	XCube333OHMultiple5    Project = "333oh_multiple5"
	XCube333Multiple10     Project = "333multiple10"
	XCube333OHMultiple10   Project = "333oh_multiple10"
	XCube333Multiple15     Project = "333multiple15"
	XCube333OHMultiple15   Project = "333oh_multiple15"
	XCube333Multiple20     Project = "333multiple20"
	XCube333OHMultiple20   Project = "333oh_multiple20"
	XCube2345Relay         Project = "2345relay"
	XCube2345OHRelay       Project = "2345oh_relay"
	XCube27Relay           Project = "2_7relay"
	XCube27OHRelay         Project = "2_7oh_relay"
	XCube345RelayBF        Project = "345relay_bf"
	XCube345OHRelayBF      Project = "345oh_relay_bf"
	XCubeAlienRelay        Project = "alien_relay" // 塔斜表五Q
	XCubeAlienOHRelay      Project = "alien_oh_relay"
	XCube27AlienRelayAll   Project = "27alien_relay"
	XCube27AlienOHRelayAll Project = "27alien_oh_relay"

	// 特殊魔方
	XCube333Mirror    Project = "333mirror"
	XCube333Mirroring Project = "333mirroring"
	XCube333Ghost     Project = "333ghost"
	XCube333ZongZi    Project = "333Zongzi"
	XCubeHotWheels    Project = "hot_wheels"
	XCubeFisher       Project = "fisher"
	XCubeGear         Project = "gear"
	Xcube333Clone     Project = "333clone"
	XCubeMapleLeaf    Project = "maple_leaf"
	XCube222Minx      Project = "222minx"
	XCube444Minx      Project = "444minx"
	XCube555Minx      Project = "555minx"
	XCube333Mini      Project = "333mini"
	XCube444Py        Project = "444pyram"
	XCube888          Project = "888"
	XCube999          Project = "999"
	XCube10L          Project = "10level"
	XCube11L          Project = "11level"
	XCube12L          Project = "12level"
	XCube13L          Project = "13level"
	XCube14L          Project = "14level"
	XCube15L          Project = "15level"
	XCube16L          Project = "16level"
	XCube17L          Project = "17level"
	XCube21L          Project = "21level"
	XCube133          Project = "133"
	XCube223          Project = "223"
	XCube233          Project = "233"
	XCube334          Project = "334"
	XCube335          Project = "335"
	XCube336          Project = "336"
	XCube337          Project = "337"
	XCubeHelicopter   Project = "helicopter"
	XCubeRedi         Project = "redi"

	// 数独系列
	NotCubeSuDoKuVeryEasy Project = "sudoku_very_easy"
	NotCubeSuDoKuEasy     Project = "sudoku_easy"
	NotCubeSuDoKuModerate Project = "sudoku_moderate"
	NotCubeSuDoKuAdvanced Project = "sudoku_advanced"
	NotCubeSuDoKuHard     Project = "sudoku_hard"
	NotCubeSuDoKuMaster   Project = "sudoku_master"

	// 数字华容道系列
	NotCube8Puzzle  Project = "8puzzle"
	NotCube15Puzzle Project = "15puzzle"
	NotCube24Puzzle Project = "24puzzle"
	NotCube35Puzzle Project = "35puzzle"
	NotCube48Puzzle Project = "48puzzle"
	NotCube63Puzzle Project = "63puzzle"
	NotCube80Puzzle Project = "80puzzle"

	// 记字
	NotCubeDigitVeryEasy           Project = "digit_very_easy"     // 20个
	NotCubeDigitEasy               Project = "digit_easy"          // 40个
	NotCubeDigitModerate           Project = "digit_moderate"      // 80个
	NotCubeDigitAdvanced           Project = "digit_advanced"      // 120个
	NotCubeDigitHard               Project = "digit_hard"          // 200个
	NotCubeDigitMaster             Project = "digit_master"        // 300个
	NotCubeDigitOnlyNumberVeryEasy Project = "digit_num_very_easy" // 纯数字20个
	NotCubeDigitOnlyNumberEasy     Project = "digit_num_easy"      // 纯数字40个
	NotCubeDigitOnlyNumberModerate Project = "digit_num_moderate"  // 纯数字80个
	NotCubeDigitOnlyNumberAdvanced Project = "digit_num_advanced"  // 纯数字120个
	NotCubeDigitOnlyNumberHard     Project = "digit_num_hard"      // 纯数字200个
	NotCubeDigitOnlyNumberMaster   Project = "digit_num_master"    // 纯数字300个
)

var projectsItems = []projectItem{
	{Project: JuBaoHaoHao, Cn: "菊爆浩浩", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: OtherCola, Cn: "速可乐", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCube}},

	{Project: Cube222, Cn: "二阶", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333, Cn: "三阶", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube444, Cn: "四阶", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube555, Cn: "五阶", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube666, Cn: "六阶", RouteType: RouteType3roundsAvg, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube777, Cn: "七阶", RouteType: RouteType3roundsAvg, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: CubeSk, Cn: "斜转", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: CubePy, Cn: "金字塔", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: CubeSq1, Cn: "SQ1", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: CubeMinx, Cn: "五魔方", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: CubeClock, Cn: "魔表", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333OH, Cn: "单手", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333FM, Cn: "最少步", RouteType: RouteType3roundsAvg, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333BF, Cn: "三盲", RouteType: RouteType3roundsBest, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube444BF, Cn: "四盲", RouteType: RouteType3roundsBest, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube555BF, Cn: "五盲", RouteType: RouteType3roundsBest, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333MBF, Cn: "多盲", RouteType: RouteTypeRepeatedly, IsWca: true, Class: []string{ProjectClassWCA}},
	{Project: Cube333Ft, Cn: "脚拧", RouteType: RouteType5RoundsAvgHT, IsWca: true, Class: []string{ProjectClassWCA}},

	{Project: XCubePyFm, Cn: "塔少步", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeFm}},
	{Project: XCubeSkFm, Cn: "斜少步", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeFm}},
	{Project: XCube222Fm, Cn: "二少步", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeFm}},
	{Project: XCubeClockFm, Cn: "表少步", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeFm}},

	{Project: XCube222BF, Cn: "二盲", RouteType: RouteType5roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCube666BF, Cn: "六盲", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCube777BF, Cn: "七盲", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCubePyBF, Cn: "塔盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCubeSkBF, Cn: "斜盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCubeMinxBf, Cn: "五魔盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCubeClockBf, Cn: "表盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},
	{Project: XCubeSQ1Bf, Cn: "SQ1盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF}},

	{Project: XCube333BfOH, Cn: "三单盲", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassXCubeBF, ProjectClassXCubeOH}},
	{Project: XCube444BfOH, Cn: "四单盲", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeBF, ProjectClassXCubeOH}},
	{Project: XCube555BfOH, Cn: "五单盲", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeBF, ProjectClassXCubeOH}},

	{Project: XCube333MiniOH, Cn: "三阶迷你单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCube222OH, Cn: "二单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCube444OH, Cn: "四单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCube555OH, Cn: "五单", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCube666OH, Cn: "六单", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCube777OH, Cn: "七单", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCubeSkOH, Cn: "斜转单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCubePyOH, Cn: "金字塔单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCubeSq1OH, Cn: "SQ1单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCubeMinxOH, Cn: "五魔单", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCubeOH}},
	{Project: XCubeClockOH, Cn: "表单", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCubeOH}},

	{Project: XCube333Multiple5, Cn: "三阶五连", RouteType: RouteType5roundsBest, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333Multiple10, Cn: "三阶十连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333Multiple15, Cn: "三阶十五连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333Multiple20, Cn: "三阶二十连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333OHMultiple5, Cn: "三单五连", RouteType: RouteType5roundsBest, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333OHMultiple10, Cn: "三单十连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333OHMultiple15, Cn: "三单十五连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube333OHMultiple20, Cn: "三单二十连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube2345Relay, Cn: "二三四五连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube2345OHRelay, Cn: "二三四五单连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube27Relay, Cn: "正阶连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube27OHRelay, Cn: "正阶单手连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube345RelayBF, Cn: "盲连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube345OHRelayBF, Cn: "盲单手连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCubeAlienRelay, Cn: "异形连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCubeAlienOHRelay, Cn: "异形单手连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube27AlienRelayAll, Cn: "全项目连拧", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},
	{Project: XCube27AlienOHRelayAll, Cn: "全项目单手连", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassXCubeRelay}},

	{Project: XCube333Mirror, Cn: "镜面魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube333Mirroring, Cn: "镜向三阶", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube333Mini, Cn: "三阶迷你", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube333Ghost, Cn: "鬼魔", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube333ZongZi, Cn: "粽子魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: Xcube333Clone, Cn: "三阶克隆", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeMapleLeaf, Cn: "枫叶魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeHotWheels, Cn: "风火轮", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeFisher, Cn: "移棱魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeGear, Cn: "齿轮魔方", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube222Minx, Cn: "二阶五魔", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube444Minx, Cn: "四阶五魔", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube555Minx, Cn: "五阶五魔", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube444Py, Cn: "四阶塔", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeHelicopter, Cn: "直升机", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCubeRedi, Cn: "redi", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},

	{Project: XCube888, Cn: "八阶", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube999, Cn: "九阶", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube10L, Cn: "十阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube11L, Cn: "十一阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube12L, Cn: "十二阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube13L, Cn: "十三阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube14L, Cn: "十四阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube15L, Cn: "十五阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube16L, Cn: "十六阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube17L, Cn: "十七阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},
	{Project: XCube21L, Cn: "二十一阶", RouteType: RouteType1rounds, IsWca: false, Class: []string{ProjectClassSuperHigh}},

	{Project: XCube133, Cn: "一三三", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube223, Cn: "二二三", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube233, Cn: "二三三", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube334, Cn: "三三四", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube335, Cn: "三三五", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube336, Cn: "三三六", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},
	{Project: XCube337, Cn: "三三七", RouteType: RouteType3roundsAvg, IsWca: false, Class: []string{ProjectClassXCube}},

	{Project: NotCubeSuDoKuVeryEasy, Cn: "数独入门", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCubeSuDoKuEasy, Cn: "数独初级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCubeSuDoKuModerate, Cn: "数独中级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCubeSuDoKuAdvanced, Cn: "数独高级", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCubeSuDoKuHard, Cn: "数独困难", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCubeSuDoKuMaster, Cn: "数独大师", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},

	{Project: NotCube8Puzzle, Cn: "3阶数字华容道", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube15Puzzle, Cn: "4阶数字华容道", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube24Puzzle, Cn: "5阶数字华容道", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube35Puzzle, Cn: "6阶数字华容道", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube48Puzzle, Cn: "7阶数字华容道", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube63Puzzle, Cn: "8阶数字华容道", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},
	{Project: NotCube80Puzzle, Cn: "9阶数字华容道", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassNotCube}},

	{Project: NotCubeDigitVeryEasy, Cn: "记字入门", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitEasy, Cn: "记字初级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitModerate, Cn: "记字中级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitAdvanced, Cn: "记字高级", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitHard, Cn: "记字困难", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitMaster, Cn: "记字大师", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberVeryEasy, Cn: "记数字入门", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberEasy, Cn: "记数字初级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberModerate, Cn: "记数字中级", RouteType: RouteType5RoundsAvgHT, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberAdvanced, Cn: "记数字高级", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberHard, Cn: "记数字困难", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
	{Project: NotCubeDigitOnlyNumberMaster, Cn: "记数字大师", RouteType: RouteType3roundsBest, IsWca: false, Class: []string{ProjectClassDigit}},
}

var projectItemsMap = make(map[Project]projectItem, len(projectsItems))
var allProjectRoute []Project
var xCubeProjectRoute []Project
var wcaProjectRoute []Project

func init() {
	for _, val := range projectsItems {
		allProjectRoute = append(allProjectRoute, val.Project)
		if val.IsWca {
			wcaProjectRoute = append(wcaProjectRoute, val.Project)
		} else {
			xCubeProjectRoute = append(xCubeProjectRoute, val.Project)
		}
		projectItemsMap[val.Project] = val
	}
}
