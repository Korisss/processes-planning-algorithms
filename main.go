package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	planTypes := []string{"fcfs", "rr", "sjf", "все", "all"}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Введите тип планирования (FCFS, SJF, RR, все, all): ")

	typ, err := ReadString(scanner)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if !SliceContains(planTypes, strings.ToLower(typ)) {
		fmt.Println("Неправильный тип")
		return
	}

	fmt.Println("Введите количество процессов (3): ")

	processCount, err := ReadInt(scanner)
	if err != nil {
		fmt.Println("Неверный ввод:", err.Error())
		return
	}

	if processCount < 1 {
		fmt.Println("Количество процессов не может быть меньше 1")
		return
	}

	fmt.Println("Введите максимальное количество квантов (15): ")

	maxQuant, err := ReadInt(scanner)
	if err != nil {
		fmt.Println("Неверный ввод:", err.Error())
		return
	}

	if maxQuant < 1 {
		fmt.Println("Количество квантов не может быть меньше 1")
		return
	}

	planner := NewPlanner(maxQuant, processCount, typ, 0)

	rrQuants := 0
	if strings.ToLower(typ) == "rr" || strings.ToLower(typ) == "все" || strings.ToLower(typ) == "all" {
		bestQuant, smallestTime := planner.CalcBestQuantAndTimeForRR()
		fmt.Println("Самое оптимизированное количество квантов на процесс для RR", bestQuant, "при среднем времени ожидания", smallestTime)
		fmt.Println("Введите количество квантов на процесс для планирования RR (Round Robin): ")
		rrQuants, err = ReadInt(scanner)
		if err != nil {
			fmt.Println("Неверный ввод:", err.Error())
			return
		}

		if maxQuant < 1 {
			fmt.Println("Количество квантов не может быть меньше 1")
			return
		}
	}

	planner.SetQuantForRR(rrQuants)
	planner.GeneratePlan()
	planner.Plan()
}

func GetSeparatorLength(plan map[int]int) int {
	length := 0
	for _, v := range plan {
		length += v
	}

	return length + 3
}

func calcWaitTime(planString string) float64 {
	lastPlus := 0

	for i, v := range planString {
		if v == '+' {
			lastPlus = i
		}
	}

	return float64(strings.Count(planString[0:lastPlus], "-"))
}

func calcFullRunTime(planString string) float64 {
	lastPlus := 0

	for i, v := range planString {
		if v == '+' {
			lastPlus = i
		}
	}

	return float64(lastPlus + 1)
}

func ReadString(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", errors.New("")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return scanner.Text(), nil
}

func ReadInt(scanner *bufio.Scanner) (int, error) {
	if !scanner.Scan() {
		return 0, errors.New("")
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return strconv.Atoi(scanner.Text())
}

func SliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if str == v {
			return true
		}
	}

	return false
}
