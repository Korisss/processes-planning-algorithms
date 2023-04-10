package planners

import (
	"math/rand"
	"plan-algorithms/utils"
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
	plan := utils.CopyMap(p.processes)
	maxLen := 0

	for key := range p.processes {
		p.plans[key] = NewPlan("")
	}

	for len(plan) != 0 {
		keys := utils.GetAllMapKeys(plan)
		sort.SliceStable(keys, func(i, j int) bool {
			return prioritiesMap[keys[i]] > prioritiesMap[keys[j]]
		})

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

			maxLen = len(p.plans[key].PlanString)
		}
	}

	for key := range p.plans {
		p.plans[key].PlanString += strings.Repeat("-", maxLen-len(p.plans[key].PlanString))
	}
}
