package planner

import (
	"fmt"
	"math/rand"
	"plan-algorithms/util"
	"sort"
	"strings"
)

type RRPlanner struct {
	name            string
	quantCountForRR int
	plans           map[int]*Plan
	processes       map[int]int
}

func NewRRPlanner(quantCountForRR int) *RRPlanner {
	return &RRPlanner{
		name:            "rr",
		plans:           make(map[int]*Plan),
		processes:       make(map[int]int),
		quantCountForRR: quantCountForRR,
	}
}

func (p *RRPlanner) GetName() string {
	return p.name
}

func (p *RRPlanner) GetPlans() map[int]*Plan {
	return p.plans
}

func (p *RRPlanner) SetProcesses(processes map[int]int) {
	p.processes = processes
}

func (p *RRPlanner) GeneratePlans(random *rand.Rand, prioritiesMap map[int]int) {
	plan := util.CopyMap(p.processes)
	maxLen := 0

	for key := range p.processes {
		p.plans[key] = NewPlan("")
	}

	for len(plan) != 0 {
		keys := util.GetAllMapKeys(p.processes)

		if prioritiesMap != nil {
			sort.SliceStable(keys, func(i, j int) bool {
				return prioritiesMap[keys[i]] > prioritiesMap[keys[j]]
			})
		}

		for _, key := range keys {
			quantCount := 0
			if plan[key] > p.quantCountForRR {
				quantCount = p.quantCountForRR
			} else {
				quantCount = plan[key]
			}

			p.plans[key].PlanString += strings.Repeat("+", quantCount)
			for k := range plan {
				if k != key {
					p.plans[k].PlanString += strings.Repeat("-", quantCount)
				}
			}

			if plan[key] > p.quantCountForRR {
				plan[key] -= p.quantCountForRR
			} else {
				delete(plan, key)
			}

			if maxLen < len(p.plans[key].PlanString) {
				maxLen = len(p.plans[key].PlanString)
			}
		}
	}

	for key := range p.plans {
		fmt.Println(maxLen, len(p.plans[key].PlanString))
		p.plans[key].PlanString += strings.Repeat("-", maxLen-len(p.plans[key].PlanString))
	}
}

func CalcBestQuantAndTimeForRR(maxQuantCountPerProcess int, originalProcesses map[int]int) (int, float64) {
	bestQuantCount := 1
	smallestAvgWaitTime := float64(-1)

	for quantCount := 1; quantCount < maxQuantCountPerProcess; quantCount++ {
		waitTime := 0
		processes := util.CopyMap(originalProcesses)

		planStrings := make(map[int]string, len(processes))
		maxLen := 0

		for len(processes) != 0 {
			keys := make([]int, 0, len(processes))
			for key := range processes {
				keys = append(keys, key)
			}

			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			for _, key := range keys {
				if processes[key] > quantCount {
					planStrings[key] += strings.Repeat("+", quantCount)
					for k := range processes {
						if k != key {
							planStrings[k] += strings.Repeat("-", quantCount)
						}
					}

					processes[key] -= quantCount
				} else {
					planStrings[key] += strings.Repeat("+", processes[key])
					for k := range processes {
						if k != key {
							planStrings[k] += strings.Repeat("-", processes[key])
						}
					}

					delete(processes, key)
				}

				maxLen = len(planStrings[key])
			}
		}

		for key := range planStrings {
			planStrings[key] += strings.Repeat("-", maxLen-len(planStrings[key]))
		}

		for i := 0; i < len(planStrings); i++ {
			waitTime += calcWaitTime(planStrings[i])
		}

		avgWaitTime := float64(waitTime) / float64(len(originalProcesses))
		if avgWaitTime < smallestAvgWaitTime || smallestAvgWaitTime == -1 {
			smallestAvgWaitTime = avgWaitTime
			bestQuantCount = quantCount
		}
	}

	return bestQuantCount, smallestAvgWaitTime
}
