package planner

import (
	"bufio"
	"errors"
	"math/rand"
	"strconv"
	"time"
)

type Planner interface {
	GeneratePlans(random *rand.Rand, prioritiesMap map[int]int)
	GetName() string
	GetPlans() map[int]*Plan
	SetProcesses(processes map[int]int)
}

func GenerateProcesses(processCount, maxQuantCountPerProcess int) map[int]int {
	random := rand.New(rand.NewSource(time.Now().Unix()))

	processes := make(map[int]int)
	for i := 0; i < processCount; i++ {
		processes[i] = random.Int()%maxQuantCountPerProcess + 1
	}

	return processes
}

func GetSeparatorLength(plan map[int]int) int {
	length := 0
	for _, v := range plan {
		length += v
	}

	return length + 3
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
