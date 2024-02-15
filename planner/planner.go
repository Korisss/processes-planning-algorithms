package planner

import (
	"bufio"
	"errors"
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Planner interface {
	GeneratePlans(random *rand.Rand, prioritiesMap map[int]int)
	GetName() string
	GetPlans() map[int]*Plan
	SetProcesses(processes map[int]int)
}

func SavePlans(plans map[int]*Plan, label string) {
	p := plot.New()
	p.Title.Text = strings.ToUpper(label)
	p.X.Label.TextStyle.Font.Size = 500

	for planNum, plan := range plans {
		fmt.Println(123)
		str := plan.PlanString
		start, end := -1, -1
		for i := 0; i < len(str); i++ {
			if start != -1 && end != -1 && str[i] == '-' {
				xyer := plotter.XYs{
					{
						X: float64(start),
						Y: float64(planNum),
					},
					{
						X: float64(end + 1),
						Y: float64(planNum),
					},
				}

				line, _ := plotter.NewLine(xyer)
				line.Color = color.RGBA{uint8(255 * planNum / len(plans)), uint8(255 * (len(plans) - planNum) / len(plans)), 255, 255}
				line.Width = 15
				p.Add(line)

				start = -1
				end = -1
			}

			if str[i] == '+' {
				if start == -1 {
					start = i
				}
				end = i
			}
		}
	}

	grid := plotter.NewGrid()
	p.Add(grid)

	tf := plot.TickerFunc(func(min, max float64) []plot.Tick {
		ticks := make([]plot.Tick, int(max-min))
		for i := min; i < max; i++ {
			ticks[int(i-min)] = plot.Tick{Value: i, Label: strconv.Itoa(int(i))}
		}
		return ticks
	})

	// if int(ticks[int(i-min)].Value)%2 != 0 {
	// 	ticks[int(i-min)].Label = ""
	// }

	p.X.Tick.Marker = tf
	p.Y.Tick.Marker = tf

	p.Save(vg.Inch*10, vg.Inch*5, label+".png")
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
