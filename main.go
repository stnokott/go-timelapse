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

	manager := steps.NewManager()

	if _, err := tea.NewProgram(manager).Run(); err != nil {
		panic(err)
	}
}
