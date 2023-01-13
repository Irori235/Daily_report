package main

import (
	"fmt"
	"math"
	"time"
)

type TimeLog struct {
	ID          int       `json:"id"`
	Title       string    `json:"name"`
	ProjectName string    `json:"project_name"`
	Start       time.Time `json:"start"`
	Stop        time.Time `json:"stop"`
}

type ProjectTime map[string](float64)
type TaskTime map[string](float64)
type ScoreTable map[string](int)

type Analyzer struct {
	sts ScoreTable

	output output
}

type output struct {
	msg string
}

func NewAnalyzer(sts ScoreTable) Analyzer {

	out := output{
		msg: "",
	}

	return Analyzer{
		sts:    sts,
		output: out,
	}
}

func (a *Analyzer) Total(tls []TimeLog) (ProjectTime, TaskTime) {
	var pt ProjectTime
	for _, tl := range tls {
		if _, ok := pt[tl.ProjectName]; !ok {
			pt[tl.ProjectName] = 0
		}
		pt[tl.ProjectName] += tl.Stop.Sub(tl.Start).Minutes()
	}

	var tt TaskTime
	for _, tl := range tls {
		if _, ok := tt[tl.Title]; !ok {
			tt[tl.Title] = 0
		}
		tt[tl.Title] += tl.Stop.Sub(tl.Start).Minutes()
	}

	return pt, tt
}

func (a *Analyzer) TotalScore(pt ProjectTime) {
	var totalScore int

	for proj, time := range pt {
		if _, ok := a.sts[proj]; !ok {
			a.sts[proj] = 0
		}
		totalScore += a.sts[proj] * int(time)
	}

	a.CreateMessage("Total score : %v", totalScore)

	return
}

func (a *Analyzer) Coverage(tls []TimeLog, pt ProjectTime) {
	var init time.Time
	var final time.Time

	for _, tl := range tls {
		if init.IsZero() || tl.Start.Before(init) {
			init = tl.Start
		}
		if final.IsZero() || tl.Stop.After(final) {
			final = tl.Stop
		}
	}

	var totalMinutes float64
	for _, t := range pt {
		totalMinutes += t
	}

	coverage := totalMinutes / final.Sub(init).Minutes() * 100

	a.CreateMessage("Coverage : %v", coverage)

	return
}

type TargetType int

const (
	TargetTypeProject TargetType = iota
	TargetTypeTask
)

type Terget struct {
	TargetType TargetType
	Name       string
}

func (a *Analyzer) EmojiSum(name string, target []Terget) {
	var sum float64
	for _, t := range target {
		switch t.TargetType {
		case TargetTypeProject:
			sum += a.pt[t.Name]
		case TargetTypeTask:
			sum += a.tt[t.Name]
		}
	}

	/*
		sum > 10h -> 🚀🚀🚀🚀🚀 (>10h)
		sum > 6h -> 🔥🔥🔥 (>6h)
		sum > 4h -> ⛰️⛰️ (>4h)
		sum > 2h -> 🌳 (>2h)
	*/

	var emoji string
	switch {
	case sum > 10*60:
		emoji = "🚀🚀🚀🚀🚀"
	case sum > 6*60:
		emoji = "🔥🔥🔥"
	case sum > 4*60:
		emoji = "⛰️⛰️"
	case sum > 2*60:
		emoji = "🌳"
	default:
		emoji = ""
	}

	// 四捨五入して、 ?h にする
	hour := math.Floor((sum / 60) + .5)

	a.CreateMessage("%v : %v (%dh)", name, emoji, int(hour))

	return
}

func (a *Analyzer) CreateMessage(fstr string, args ...any) {

	a.output.msg += fmt.Sprintf(fstr+"\n", args...)

	return
}
