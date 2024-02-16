package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"plan-algorithms/planner"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
)

//go:generate go-winres make
const title = "Планировщик процессов"

var (
	mainWindow              *goey.Window
	img                     image.Image
	random                  = rand.New(rand.NewSource(time.Now().Unix()))
	usePriorities           = false
	processCount            = 5
	maxQuantCountPreProcess = 15
	rrQuantCount            = 5
	planType                = 0
	planTypes               = []string{
		"FCFS",
		"RR",
		"RRSJF",
		"SJF",
	}
)

func openWindow() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

func createWindow() error {
	w, err := goey.NewWindow(title, renderWindow())

	mainWindow = w
	return err
}

func updateWindow() {
	err := mainWindow.SetChild(renderWindow())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}

// тип
// кол-во процессов
// кол-во квантов макс
// кол-во квантов на процесс rr sjf

//самое оптимизированное кол-во квантов

func renderWindow() base.Widget {
	f, err := os.Open("./out/" + strings.ToLower(planTypes[planType]) + ".png")
	if err != nil {
		fmt.Println(err)
	}

	img, _, err = image.Decode(f)
	if err != nil {
		fmt.Println(err)
	}

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child: &goey.VBox{
			Children: []base.Widget{
				&goey.SelectInput{
					Value: planType,
					Items: planTypes,
					OnChange: func(value int) {
						planType = value
					},
				},
				&goey.HBox{
					Children: []base.Widget{
						&goey.Label{Text: "Кол-во процессов"},
						&goey.TextInput{
							Value: strconv.Itoa(processCount),
							OnChange: func(value string) {
								v, err := strconv.Atoi(value)
								if err != nil {
									updateWindow()
									return
								}

								processCount = v
							},
						},
					},
				},
				&goey.HBox{
					Children: []base.Widget{
						&goey.Label{Text: "Максимальное количество квантов на процесс"},
						&goey.TextInput{
							Value: strconv.Itoa(maxQuantCountPreProcess),
							OnChange: func(value string) {
								v, err := strconv.Atoi(value)
								if err != nil {
									updateWindow()
									return
								}

								maxQuantCountPreProcess = v
							},
						},
					},
				},
				&goey.HBox{
					Children: []base.Widget{
						&goey.Label{Text: "Количество квантов до переключения (для RR и RRSJF)"},
						&goey.TextInput{
							Value: strconv.Itoa(rrQuantCount),
							OnChange: func(value string) {
								v, err := strconv.Atoi(value)
								if err != nil {
									updateWindow()
									return
								}

								rrQuantCount = v
							},
						},
					},
				},
				&goey.Checkbox{
					Text:  "Использовать приоритеты (для RR)",
					Value: usePriorities,
					OnChange: func(value bool) {
						usePriorities = value
					},
				},
				&goey.Button{
					Text: "Сгенерировать план",
					OnClick: func() {
						generatePlan()
						updateWindow()
					},
				},
				&goey.Img{
					Image: img,
				},
			},
		},
	}
}

// processCount            = 5
// maxQuantCountPreProcess = 15
// rrQuantCount            = 5
// planType                = 0
// planTypes               = []string{
// 	"FCFS",
// 	"RR",
// 	"RRSJF",
// 	"SJF",
// }

func generatePlan() {
	processes := planner.GenerateProcesses(processCount, maxQuantCountPreProcess)
	prioritiesMap := make(map[int]int)

	for i := 0; i < processCount; i++ {
		if usePriorities {
			prioritiesMap[i] = random.Int()%41 - 20
		} else {
			prioritiesMap[i] = 0
		}
	}

	var p planner.Planner
	switch planType {
	case 0:
		p = planner.NewFCFSPlanner()
	case 1:
		p = planner.NewRRPlanner(rrQuantCount)
	case 2:
		p = planner.NewRRSJFPlanner(rrQuantCount)
	case 3:
		p = planner.NewSJFPlanner()
	default:
		break
	}

	p.SetProcesses(processes)
	p.GeneratePlans(random, prioritiesMap)

	planner.SavePlans(p.GetPlans(), p.GetName())
}
