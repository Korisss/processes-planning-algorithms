package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"plan-algorithms/planner"
	"strconv"
	"strings"
)

func main() {
	openWindow()
}

func main1() {
	planTypes := []string{"fcfs", "rr", "rrsjf", "sjf", "все", "all"}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите тип планирования (FCFS, SJF, RR, RRSJF, все, all): ")

	typ, err := planner.ReadString(scanner)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if !planner.SliceContains(planTypes, strings.ToLower(typ)) {
		fmt.Println("Неправильный тип")
		return
	}

	maxQuant := 0
	processes := make(map[int]int)
	prioritiesMap := map[int]int(nil)
	backedUpFile, err := os.Open("out.csv")
	imported := false
	if err == nil {
		fmt.Print("Был найден сохранённый файл с процессами. Вы хотите его импортировать? y/n (д/н):")
		readFile, _ := planner.ReadString(scanner)
		if readFile == "y" || readFile == "д" {
			prioritiesMap = make(map[int]int)
			imported = true
			reader := csv.NewReader(backedUpFile)
			reader.Read()
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("Ошибка при парсинге файла:", err.Error())
				return
			}

			for _, arr := range records {
				index, err := strconv.Atoi(arr[0])
				if err != nil {
					fmt.Println("Ошибка при парсинге файла:", err.Error())
					return
				}

				processes[index], err = strconv.Atoi(arr[1])
				if err != nil {
					fmt.Println("Ошибка при парсинге файла:", err.Error())
					return
				}

				prioritiesMap[index], err = strconv.Atoi(arr[2])
				if err != nil {
					fmt.Println("Ошибка при парсинге файла:", err.Error())
					return
				}

				if maxQuant < processes[index] {
					maxQuant = processes[index]
				}
			}
		}
	}

	if !imported {
		fmt.Println("Введите количество процессов (3): ")

		processCount, err := planner.ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}

		if processCount < 1 {
			fmt.Println("Количество процессов не может быть меньше 1")
			return
		}

		fmt.Println("Введите максимальное количество квантов (15): ")

		maxQuant, err = planner.ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}

		if maxQuant < 1 {
			fmt.Println("Количество квантов не может быть меньше 1")
			return
		}

		processes = planner.GenerateProcesses(processCount, maxQuant)
		maxQuant = 0
		for _, process := range processes {
			if process > maxQuant {
				maxQuant = process
			}
		}
		fmt.Print("Вы хотите использовать приоритеты? y/n (д/н):")
		priorities, err := planner.ReadString(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}

		if priorities != "n" && priorities != "y" && priorities != "н" && priorities != "д" {
			fmt.Println("Неверный ввод.")
			return
		}

		usePriorities := false
		if priorities == "y" || priorities == "д" {
			usePriorities = true
		}

		if usePriorities {
			prioritiesMap = make(map[int]int)
			for i := 0; i < processCount; i++ {
				prioritiesMap[i] = random.Int()%41 - 20
			}
		}
	}

	separator := strings.Repeat("=", planner.GetSeparatorLength(processes))
	rrQuants := 0

	if strings.ToLower(typ) == "rr" || strings.ToLower(typ) == "rrsjf" || strings.ToLower(typ) == "все" || strings.ToLower(typ) == "all" {
		bestQuant, smallestTime := planner.CalcBestQuantAndTimeForRR(maxQuant, processes)

		fmt.Println("Самое оптимизированное количество квантов на процесс для RR", bestQuant, "при среднем времени ожидания", smallestTime)
		fmt.Println("Введите количество квантов на процесс для планирования RR (Round Robin): ")

		rrQuants, err = planner.ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}
	}

	planners := []planner.Planner{}
	switch strings.ToLower(typ) {
	case "fcfs":
		planners = append(planners, planner.NewFCFSPlanner())
	case "rr":
		planners = append(planners, planner.NewRRPlanner(rrQuants))
	case "rrsjf":
		planners = append(planners, planner.NewRRSJFPlanner(rrQuants))
	case "sjf":
		planners = append(planners, planner.NewSJFPlanner())
	case "все":
	case "all":
		planners = append(planners, planner.NewFCFSPlanner())
		planners = append(planners, planner.NewSJFPlanner())
		planners = append(planners, planner.NewRRPlanner(rrQuants))
		planners = append(planners, planner.NewRRSJFPlanner(rrQuants))
	default:
		break
	}

	for _, p := range planners {
		p.SetProcesses(processes)
		p.GeneratePlans(random, prioritiesMap)
		if prioritiesMap == nil {
			prioritiesMap = make(map[int]int)
		}

		planner.SavePlans(p.GetPlans(), p.GetName())

		fmt.Println(separator)
		fmt.Println(strings.ToUpper(p.GetName()))
		fmt.Println(separator)

		plans := p.GetPlans()
		waitTime := float64(0)
		runTime := float64(0)
		for i := 0; i < len(plans); i++ {
			plans[i].CalcTime()

			waitTime += float64(plans[i].GetWaitTime())
			runTime += float64(plans[i].GetFullRunTime())

			fmt.Println(fmt.Sprint("P", i, " ", plans[i].GetPlanString(), " waitTime: ", plans[i].GetWaitTime(), " fullRunTime: ", plans[i].GetFullRunTime()), "priority", prioritiesMap[i])
		}

		fmt.Println(separator)
		fmt.Println("full waitTime", waitTime)
		fmt.Println("avg waitTime", waitTime/float64(len(processes)))
		fmt.Println("avg runTime", runTime/float64(len(processes)))
		fmt.Println(separator)
	}

	if !imported {
		fmt.Print("Вы хотите записать процессы в файл? y/n (д/н):")
		writeFile, err := planner.ReadString(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}

		if writeFile != "n" && writeFile != "y" && writeFile != "н" && writeFile != "д" {
			fmt.Println("Неверный ввод.")
			return
		}

		if writeFile == "y" || writeFile == "д" {
			file, err := os.Create("out.csv")
			if err != nil {
				fmt.Println("Ошибка при создании файла:", err.Error())
				return
			}

			writer := csv.NewWriter(file)
			writer.Write([]string{"Process", "cpu burst", "priority"})
			for i, process := range processes {
				err = writer.Write([]string{strconv.Itoa(i), strconv.Itoa(process), strconv.Itoa(prioritiesMap[i])})
				if err != nil {
					fmt.Println("Ошибка при записи в файл:", err.Error())
					return
				}
			}

			writer.Flush()
		}
	}
}
