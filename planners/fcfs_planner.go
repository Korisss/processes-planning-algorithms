package planners

import (
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

func (p *FCFSPlanner) GeneratePlans() {
	for key := range p.processes {
		p.plans[key] = NewPlan("")
	}

	for i := 0; i < len(p.processes); i++ {
		for j := 0; j < i; j++ {
			p.plans[i].PlanString += strings.Repeat("-", p.processes[j])
		}

		p.plans[i].PlanString += strings.Repeat("+", p.processes[i])

		for j := i + 1; j < len(p.processes); j++ {
			p.plans[i].PlanString += strings.Repeat("-", p.processes[j])
		}
	}
}
