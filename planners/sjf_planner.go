package planners

import (
	"math/rand"
	"plan-algorithms/utils"
	"sort"
	"strings"
)

type SJFPlanner struct {
	name      string
	plans     map[int]*Plan
	processes map[int]int
}

func NewSJFPlanner() *SJFPlanner {
	return &SJFPlanner{
		name:      "sjf",
		plans:     make(map[int]*Plan),
		processes: make(map[int]int),
	}
}

func (p *SJFPlanner) GetName() string {
	return p.name
}

func (p *SJFPlanner) GetPlans() map[int]*Plan {
	return p.plans
}

func (p *SJFPlanner) SetProcesses(processes map[int]int) {
	p.processes = processes
}

func (p *SJFPlanner) GeneratePlans(random *rand.Rand, prioritiesMap map[int]int) {
	for key := range p.processes {
		p.plans[key] = NewPlan("")
	}

	keys := utils.GetAllMapKeys(p.processes)

	sort.SliceStable(keys, func(i, j int) bool {
		return p.processes[keys[i]] < p.processes[keys[j]]
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
