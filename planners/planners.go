package planners

import (
	"bufio"
	"errors"
	"math/rand"
	"plan-algorithms/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Planner interface {
	GeneratePlans()
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

func CalcBestQuantAndTimeForRR(maxQuantCountPerProcess int, originalProcesses map[int]int) (int, float64) {
	bestQuantCount := 1
	smallestAvgWaitTime := float64(-1)

	for quantCount := 1; quantCount <= maxQuantCountPerProcess; quantCount++ {
		waitTime := 0
		processes := utils.CopyMap(originalProcesses)

		planStrings := make(map[int]string, len(processes))
		maxLen := 0

		for len(processes) != 0 {
			keys := make([]int, 0, len(processes))
			for key := range processes {
				keys = append(keys, key)
			}

			sort.SliceStable(keys, func(i, j int) bool {
				return keys[i] < keys[j]
			})

			for _, key := range keys {
				if processes[key] > quantCount {
					planStrings[key] += strings.Repeat("+", quantCount)
					for k := range processes {
						if k != key {
							planStrings[k] += strings.Repeat("-", quantCount)
						}
					}

					processes[key] -= quantCount
				} else {
					planStrings[key] += strings.Repeat("+", processes[key])
					for k := range processes {
						if k != key {
							planStrings[k] += strings.Repeat("-", processes[key])
						}
					}

					delete(processes, key)
				}

				maxLen = len(planStrings[key])
			}
		}

		for key := range planStrings {
			planStrings[key] += strings.Repeat("-", maxLen-len(planStrings[key]))
		}

		for i := 0; i < len(planStrings); i++ {
			waitTime += calcWaitTime(planStrings[i])
		}

		avgWaitTime := float64(waitTime) / float64(len(originalProcesses))
		if avgWaitTime < smallestAvgWaitTime || smallestAvgWaitTime == -1 {
			smallestAvgWaitTime = avgWaitTime
			bestQuantCount = quantCount
		}
	}

	return bestQuantCount, smallestAvgWaitTime
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
