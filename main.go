// Package main runs the program, starting with the first configured step in the settings flow
package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/steps"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	// TODO: remove
	/*
		config.Cfg.AbsInputDir = "F:/timelapse/input/MeeriTimelapse"
		config.Cfg.AbsOutpuFilepath = "F:/timelapse/output/test.mp4"
		config.Cfg.TimeFrom = time.Date(2023, time.April, 26, 20, 0, 0, 0, time.Local)
		config.Cfg.TimeTo = time.Date(2023, time.April, 27, 11, 0, 0, 0, time.Local)
	*/

	manager := steps.NewManager()

	if _, err := tea.NewProgram(manager).Run(); err != nil {
		panic(err)
	}
}
