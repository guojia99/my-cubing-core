package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/guojia99/my-cubing-core/model"
)

// ScoresParser 解析成绩详情
func ScoresParser(score model.Score) string {

	switch score.Project.RouteType() {
	case model.RouteType1rounds:
		return fmt.Sprintf("%s", scoreParser(score.Result1))
	case model.RouteTypeRepeatedly:
		return BestOrAvgParser(score, false)
	case model.RouteType3roundsAvg, model.RouteType3roundsBest:
		return fmt.Sprintf(
			"%s\t%s\t%s",
			paddingScoreParser(score.Result1), paddingScoreParser(score.Result2), paddingScoreParser(score.Result3),
		)
	case model.RouteType5roundsBest, model.RouteType5roundsAvg, model.RouteType5RoundsAvgHT:
		return fmt.Sprintf(
			"%s\t%s\t%s\t%s\t%s",
			paddingScoreParser(score.Result1), paddingScoreParser(score.Result2), paddingScoreParser(score.Result3), paddingScoreParser(score.Result4), paddingScoreParser(score.Result5),
		)
	}
	return ""
}

func paddingScoreParser(in float64) string {
	s := scoreParser(in)

	if len(s) >= 10 {
		return s
	}

	return s + strings.Repeat(" ", 10-len(s))
}

func scoreParser(in float64) string {
	if in == model.DNS {
		return "DNS"
	}
	if in <= model.DNF {
		return "DNF"
	}

	if in < 60 {
		return fmt.Sprintf("%0.2f", in)
	}
	m := int(math.Floor(in) / 60)
	s := in - float64(m*60)

	ss := fmt.Sprintf("%0.2f", s)
	if s < 10 {
		ss = fmt.Sprintf("0%0.2f", s)
	}

	return fmt.Sprintf("%d:%s", m, ss)
}

// BestOrAvgParser 解析最佳成绩或者平均成绩
func BestOrAvgParser(score model.Score, isAvg bool) string {
	switch score.Project.RouteType() {
	case model.RouteTypeRepeatedly:
		if isAvg {
			return ""
		}
		return fmt.Sprintf(
			"%2.0f / %2.0f %s", score.Result1, score.Result2, BestOrAvgParser(
				model.Score{Best: score.Result3, Project: model.Cube333, Avg: score.Result3}, true,
			),
		)
	case model.RouteType1rounds:
		if isAvg {
			return ""
		}
	}

	in := score.Best
	if isAvg {
		in = score.Avg
	}
	return scoreParser(in)
}

func ParserTimeToSeconds(t string) float64 {
	if t == "DNF" || strings.ContainsAny(t, "dD") {
		return model.DNF
	}
	if t == "DNS" || strings.Contains(t, "s") {
		return model.DNS
	}
	// 解析纯秒数格式
	if regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(t) {
		seconds, _ := strconv.ParseFloat(t, 64)
		return seconds
	}

	// 解析分+秒格式
	if regexp.MustCompile(`^\d{1,3}:\d{1,3}(\.\d+)?$`).MatchString(t) {
		parts := strings.Split(t, ":")
		minutes, _ := strconv.ParseFloat(parts[0], 64)
		seconds, _ := strconv.ParseFloat(parts[1], 64)
		return minutes*60 + seconds
	}

	// 解析时+分+秒格式
	if regexp.MustCompile(`^\d{1,3}:\d{1,3}:\d{1,3}(\.\d+)?$`).MatchString(t) {
		parts := strings.Split(t, ":")
		hours, _ := strconv.ParseFloat(parts[0], 64)
		minutes, _ := strconv.ParseFloat(parts[1], 64)
		seconds, _ := strconv.ParseFloat(parts[2], 64)
		return hours*3600 + minutes*60 + seconds
	}

	return model.DNF
}
