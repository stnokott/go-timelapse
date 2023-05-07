// Package style contains styling for the TUI widgets
package style

import "github.com/charmbracelet/lipgloss"

// Border adds a border to the rendered string
var Border = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Margin adds a 1 margin to all edges
var Margin = lipgloss.NewStyle().Margin(1)

// Err is used for displaying error messages
var Err = lipgloss.NewStyle().
	Background(lipgloss.Color("196")).
	Foreground(lipgloss.Color("231"))

// Background moves the user's attention away from what is rendered
var Background = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240"))

// Emphasis emphasises the rendered text
var Emphasis = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Underline(true)
