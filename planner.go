package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Planner struct {
	maxQuantCountPerProcess int
	quantCountForRR         int
	processCount            int
	plan                    map[int]int
	planType                string
	rand                    *rand.Rand
}

func NewPlanner(maxQuantCountPerProcess, processCount int, planType string, quantCountForRR int) *Planner {
	return &Planner{
		maxQuantCountPerProcess: maxQuantCountPerProcess,
		quantCountForRR:         quantCountForRR,
		processCount:            processCount,
		plan:                    make(map[int]int),
		planType:                planType,
		rand:                    rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *Planner) GeneratePlan() {
	for i := 0; i < p.processCount; i++ {
		p.plan[i] = p.rand.Int()%p.maxQuantCountPerProcess + 1
	}
}

func (p *Planner) SetQuantForRR(quantCountForRR int) {
	p.quantCountForRR = quantCountForRR
}

func (p *Planner) Plan() {
	switch strings.ToLower(p.planType) {
	case "fcfs":
		p.FCFS()
	case "rr":
		p.RR()
		p.RRSJF()
	case "sjf":
		p.SJF()
	case "все":
	case "all":
		p.FCFS()
		p.SJF()
		p.RR()
		p.RRSJF()
	default:
		break
	}
}

func (p *Planner) FCFS() {
	var fullExecTime float64 = 0
	var waitTime float64 = 0
	separatorlength := GetSeparatorLength(p.plan)
	execTime := 0

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("FCFS")
	fmt.Println(strings.Repeat("=", separatorlength))

	for i := 0; i < len(p.plan); i++ {
		fmt.Print("P", i, " ")

		planString := ""

		execTime += p.plan[i]
		for j := 0; j < i; j++ {
			waitTime += float64(p.plan[j])
			planString += strings.Repeat("-", p.plan[j])
		}

		planString += strings.Repeat("+", p.plan[i])

		for j := i + 1; j < len(p.plan); j++ {
			planString += strings.Repeat("-", p.plan[j])
		}

		fullExecTime += calcFullRunTime(planString)
		fmt.Println(planString)
	}

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("Время ожидания:", waitTime)
	fmt.Println("Время исполнения:", execTime)
	fmt.Println("Среднее время ожидания:", waitTime/float64(len(p.plan)))
	fmt.Println("Среднее время исполнения:", fullExecTime/float64(len(p.plan)))
}

func (p *Planner) RR() {
	plan := make(map[int]int, len(p.plan))
	for i, v := range p.plan {
		plan[i] = v
	}

	length := float64(len(plan))
	separatorlength := GetSeparatorLength(plan)

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("RR")
	fmt.Println(strings.Repeat("=", separatorlength))

	var fullExecTime float64 = 0
	var waitTime float64 = 0
	planStrings := make(map[int]string, len(plan))
	maxLen := 0

	for len(plan) != 0 {
		keys := make([]int, 0, len(plan))
		for key := range plan {
			keys = append(keys, key)
		}

		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, key := range keys {
			if plan[key] > p.quantCountForRR {
				planStrings[key] += strings.Repeat("+", p.quantCountForRR)
				for k := range plan {
					if k != key {
						planStrings[k] += strings.Repeat("-", p.quantCountForRR)
					}
				}

				plan[key] -= p.quantCountForRR
			} else {
				planStrings[key] += strings.Repeat("+", plan[key])
				for k := range plan {
					if k != key {
						planStrings[k] += strings.Repeat("-", plan[key])
					}
				}

				delete(plan, key)
			}

			maxLen = len(planStrings[key])
		}
	}

	for key := range planStrings {
		planStrings[key] += strings.Repeat("-", maxLen-len(planStrings[key]))
	}

	for i := 0; i < len(planStrings); i++ {
		fmt.Print("P", i, " ")
		fmt.Println(planStrings[i])
		waitTime += calcWaitTime(planStrings[i])
		fullExecTime += calcFullRunTime(planStrings[i])
	}

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("Время ожидания:", waitTime)
	fmt.Println("Время исполнения:", maxLen)
	fmt.Println("Среднее время ожидания:", waitTime/length)
	fmt.Println("Среднее время исполнения:", fullExecTime/length)
	fmt.Println(strings.Repeat("=", separatorlength))
}

func (p *Planner) RRSJF() {
	plan := make(map[int]int, len(p.plan))
	for i, v := range p.plan {
		plan[i] = v
	}

	length := float64(len(plan))
	separatorlength := GetSeparatorLength(plan)

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("RR + SJF")
	fmt.Println(strings.Repeat("=", separatorlength))

	var fullExecTime float64 = 0
	var waitTime float64 = 0
	planStrings := make(map[int]string, len(plan))
	maxLen := 0

	keys := make([]int, 0, len(plan))

	for key := range plan {
		keys = append(keys, key)
	}

	for len(plan) != 0 {
		sort.SliceStable(keys, func(i, j int) bool {
			return plan[keys[i]] < plan[keys[j]]
		})

		for _, key := range keys {
			if plan[key] > p.quantCountForRR {
				planStrings[key] += strings.Repeat("+", p.quantCountForRR)
				for k := range plan {
					if k != key {
						planStrings[k] += strings.Repeat("-", p.quantCountForRR)
					}
				}

				plan[key] -= p.quantCountForRR
			} else {
				planStrings[key] += strings.Repeat("+", plan[key])
				for k := range plan {
					if k != key {
						planStrings[k] += strings.Repeat("-", plan[key])
					}
				}

				delete(plan, key)
			}

			maxLen = len(planStrings[key])
		}
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for key := range planStrings {
		planStrings[key] += strings.Repeat("-", maxLen-len(planStrings[key]))
	}

	for i := 0; i < len(planStrings); i++ {
		fmt.Print("P", i, " ")
		fmt.Println(planStrings[i])
		waitTime += calcWaitTime(planStrings[i])
		fullExecTime += calcFullRunTime(planStrings[i])
	}

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("Время ожидания:", waitTime)
	fmt.Println("Время исполнения:", maxLen)
	fmt.Println("Среднее время ожидания:", waitTime/length)
	fmt.Println("Среднее время исполнения:", fullExecTime/length)
	fmt.Println(strings.Repeat("=", separatorlength))
}

func (p *Planner) SJF() {
	var fullExecTime float64 = 0
	var waitTime float64 = 0
	separatorlength := GetSeparatorLength(p.plan)
	keys := make([]int, 0, len(p.plan))
	execTime := 0

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("SJF")
	fmt.Println(strings.Repeat("=", separatorlength))

	for key := range p.plan {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return p.plan[keys[i]] < p.plan[keys[j]]
	})

	for i, v := range keys {
		fmt.Print("P", keys[i], " ")

		planString := ""

		execTime += p.plan[v]
		for j := 0; j < i; j++ {
			waitTime += float64(p.plan[keys[j]])
			planString += strings.Repeat("-", p.plan[keys[j]])
		}

		planString += strings.Repeat("+", p.plan[v])

		for j := i + 1; j < len(p.plan); j++ {
			planString += strings.Repeat("-", p.plan[keys[j]])
		}

		fullExecTime += calcFullRunTime(planString)
		fmt.Println(planString)
	}

	fmt.Println(strings.Repeat("=", separatorlength))
	fmt.Println("Время ожидания:", waitTime)
	fmt.Println("Время исполнения:", execTime)
	fmt.Println("Среднее время ожидания:", waitTime/float64(len(p.plan)))
	fmt.Println("Среднее время исполнения:", fullExecTime/float64(len(p.plan)))
}

func (p *Planner) CalcBestQuantAndTimeForRR() (int, float64) {
	bestQuantCount := 1
	var smallestAvgWaitTime float64 = -1

	for quantCount := 1; quantCount <= p.maxQuantCountPerProcess; quantCount++ {
		var waitTime float64 = 0
		plan := make(map[int]int, len(p.plan))
		for i, v := range p.plan {
			plan[i] = v
		}

		planStrings := make(map[int]string, len(plan))
		maxLen := 0

		for len(plan) != 0 {
			keys := make([]int, 0, len(plan))
			for key := range plan {
				keys = append(keys, key)
			}

			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			for _, key := range keys {
				if plan[key] > quantCount {
					planStrings[key] += strings.Repeat("+", quantCount)
					for k := range plan {
						if k != key {
							planStrings[k] += strings.Repeat("-", quantCount)
						}
					}

					plan[key] -= quantCount
				} else {
					planStrings[key] += strings.Repeat("+", plan[key])
					for k := range plan {
						if k != key {
							planStrings[k] += strings.Repeat("-", plan[key])
						}
					}

					delete(plan, key)
				}

				maxLen = len(planStrings[key])
			}
		}

		for key := range planStrings {
			planStrings[key] += strings.Repeat("-", maxLen-len(planStrings[key]))
		}

		for i := 0; i < len(planStrings); i++ {
			waitTime += float64(calcWaitTime(planStrings[i]))
		}

		if waitTime/float64(len(p.plan)) < smallestAvgWaitTime || smallestAvgWaitTime == -1 {
			smallestAvgWaitTime = waitTime / float64(len(p.plan))
			bestQuantCount = quantCount
		}
	}

	return bestQuantCount, smallestAvgWaitTime
}
