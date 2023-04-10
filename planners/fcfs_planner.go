package planners

import (
	"math/rand"
	"plan-algorithms/utils"
	"sort"
	"strings"
)

type FCFSPlanner struct {
	name      string
	plans     map[int]*Plan
	processes map[int]int
}

func NewFCFSPlanner() *FCFSPlanner {
	return &FCFSPlanner{
		name:      "fcfs",
		plans:     make(map[int]*Plan),
		processes: make(map[int]int),
	}
}

func (p *FCFSPlanner) GetName() string {
	return p.name
}

func (p *FCFSPlanner) GetPlans() map[int]*Plan {
	return p.plans
}

func (p *FCFSPlanner) SetProcesses(processes map[int]int) {
	p.processes = processes
}

func (p *FCFSPlanner) GeneratePlans(random *rand.Rand, prioritiesMap map[int]int) {
	for key := range p.processes {
		p.plans[key] = NewPlan("")
	}

	keys := utils.GetAllMapKeys(p.processes)

	sort.SliceStable(keys, func(i, j int) bool {
		return prioritiesMap[keys[i]] > prioritiesMap[keys[j]]
	})

	for i, key := range keys {
		for j := 0; j < i; j++ {
			p.plans[key].PlanString += strings.Repeat("-", p.processes[keys[j]])
		}

		p.plans[key].PlanString += strings.Repeat("+", p.processes[key])

		for j := i + 1; j < len(keys); j++ {
			p.plans[key].PlanString += strings.Repeat("-", p.processes[keys[j]])
		}
	}
}
