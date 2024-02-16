package planner

import "strings"

type Plan struct {
	PlanString  string
	fullRunTime int
	waitTime    int
}

func NewPlan(planString string) *Plan {
	return &Plan{
		PlanString:  planString,
		fullRunTime: calcFullRunTime(planString),
		waitTime:    calcWaitTime(planString),
	}
}

func (p *Plan) CalcTime() {
	p.fullRunTime = calcFullRunTime(p.PlanString)
	p.waitTime = calcWaitTime(p.PlanString)
}

func (p Plan) GetPlanString() string {
	return p.PlanString
}

func (p Plan) GetFullRunTime() int {
	return p.fullRunTime
}

func (p Plan) GetWaitTime() int {
	return p.waitTime
}

func calcFullRunTime(planString string) int {
	lastPlus := 0

	for i, v := range planString {
		if v == '+' {
			lastPlus = i
		}
	}

	return lastPlus + 1
}

func calcWaitTime(planString string) int {
	lastPlus := 0

	for i, v := range planString {
		if v == '+' {
			lastPlus = i
		}
	}

	return strings.Count(planString[0:lastPlus], "-")
}
