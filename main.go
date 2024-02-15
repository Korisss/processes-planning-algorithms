package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"plan-algorithms/planner"
	"strconv"
	"strings"
	"time"
)

func main() {
	planTypes := []string{"fcfs", "rr", "sjf", "все", "all"}
	scanner := bufio.NewScanner(os.Stdin)
	random := rand.New(rand.NewSource(time.Now().Unix()))

	fmt.Println("Введите тип планирования (FCFS, SJF, RR, все, all): ")

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
	prioritiesMap := make(map[int]int)
	backedUpFile, err := os.Open("out.csv")
	imported := false
	if err == nil {
		fmt.Print("Был найден сохранённый файл с процессами. Вы хотите его импортировать? y/n (д/н):")
		readFile, _ := planner.ReadString(scanner)
		if readFile == "y" || readFile == "д" {
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

		for i := 0; i < processCount; i++ {
			if usePriorities {
				prioritiesMap[i] = random.Int()%41 - 20
			} else {
				prioritiesMap[i] = 0
			}
		}
	}

	separator := strings.Repeat("=", planner.GetSeparatorLength(processes))
	rrQuants := 0

	if strings.ToLower(typ) == "rr" || strings.ToLower(typ) == "все" || strings.ToLower(typ) == "all" {
		bestQuant, smallestTime := planner.CalcBestQuantAndTimeForRR(maxQuant, processes)

		fmt.Println("Самое оптимизированное количество квантов на процесс для RR", bestQuant, "при среднем времени ожидания", smallestTime)
		fmt.Println("Введите количество квантов на процесс для планирования RR (Round Robin): ")

		rrQuants, err = planner.ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}
	}

	planners_ := []planner.Planner{}
	switch strings.ToLower(typ) {
	case "fcfs":
		planners_ = append(planners_, planner.NewFCFSPlanner())
	case "rr":
		planners_ = append(planners_, planner.NewRRPlanner(rrQuants))
		planners_ = append(planners_, planner.NewRRSJFPlanner(rrQuants))
	case "sjf":
		planners_ = append(planners_, planner.NewSJFPlanner())
	case "все":
	case "all":
		planners_ = append(planners_, planner.NewFCFSPlanner())
		planners_ = append(planners_, planner.NewSJFPlanner())
		planners_ = append(planners_, planner.NewRRPlanner(rrQuants))
		planners_ = append(planners_, planner.NewRRSJFPlanner(rrQuants))
	default:
		break
	}

	for _, planner := range planners_ {
		planner.SetProcesses(processes)
		planner.GeneratePlans(random, prioritiesMap)

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
