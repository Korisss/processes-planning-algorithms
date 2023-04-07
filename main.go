package main

import (
	"bufio"
	"fmt"
	"os"
	"plan-algorithms/planners"
	"strings"
)

func main() {
	planTypes := []string{"fcfs", "rr", "sjf", "все", "all"}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите тип планирования (FCFS, SJF, RR, все, all): ")

	typ, err := planners.ReadString(scanner)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if !planners.SliceContains(planTypes, strings.ToLower(typ)) {
		fmt.Println("Неправильный тип")
		return
	}

	fmt.Println("Введите количество процессов (3): ")

	processCount, err := planners.ReadInt(scanner)
	if err != nil {
		fmt.Println("Неверный ввод:", err.Error())
		return
	}

	if processCount < 1 {
		fmt.Println("Количество процессов не может быть меньше 1")
		return
	}

	fmt.Println("Введите максимальное количество квантов (15): ")

	maxQuant, err := planners.ReadInt(scanner)
	if err != nil {
		fmt.Println("Неверный ввод:", err.Error())
		return
	}

	if maxQuant < 1 {
		fmt.Println("Количество квантов не может быть меньше 1")
		return
	}

	processes := planners.GenerateProcesses(processCount, maxQuant)
	separator := strings.Repeat("=", planners.GetSeparatorLength(processes))
	rrQuants := 0

	if strings.ToLower(typ) == "rr" || strings.ToLower(typ) == "все" || strings.ToLower(typ) == "all" {
		bestQuant, smallestTime := planners.CalcBestQuantAndTimeForRR(maxQuant, processes)

		fmt.Println("Самое оптимизированное количество квантов на процесс для RR", bestQuant, "при среднем времени ожидания", smallestTime)
		fmt.Println("Введите количество квантов на процесс для планирования RR (Round Robin): ")

		rrQuants, err = planners.ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}
	}

	planners_ := []planners.Planner{}
	switch strings.ToLower(typ) {
	case "fcfs":
		planners_ = append(planners_, planners.NewFCFSPlanner())
	case "rr":
		planners_ = append(planners_, planners.NewRRPlanner(rrQuants))
		planners_ = append(planners_, planners.NewRRSJFPlanner(rrQuants))
	case "sjf":
		planners_ = append(planners_, planners.NewSJFPlanner())
	case "все":
	case "all":
		planners_ = append(planners_, planners.NewFCFSPlanner())
		planners_ = append(planners_, planners.NewSJFPlanner())
		planners_ = append(planners_, planners.NewRRPlanner(rrQuants))
		planners_ = append(planners_, planners.NewRRSJFPlanner(rrQuants))
	default:
		break
	}

	for _, planner := range planners_ {
		planner.SetProcesses(processes)
		planner.GeneratePlans()

		fmt.Println(separator)
		fmt.Println(strings.ToUpper(planner.GetName()))
		fmt.Println(separator)

		plans := planner.GetPlans()
		waitTime := float64(0)
		runTime := float64(0)
		for i := 0; i < len(plans); i++ {
			plans[i].CalcTime()

			waitTime += float64(plans[i].GetWaitTime())
			runTime += float64(plans[i].GetFullRunTime())

			fmt.Println(fmt.Sprint("P", i, " ", plans[i].GetPlanString(), " waitTime: ", plans[i].GetWaitTime(), " fullRunTime: ", plans[i].GetFullRunTime()))
		}

		fmt.Println(separator)
		fmt.Println("full waitTime", waitTime)
		fmt.Println("avg waitTime", waitTime/float64(len(processes)))
		fmt.Println("avg runTime", runTime/float64(len(processes)))
		fmt.Println(separator)
	}
}
