// Package steps contains all steps for the program and handles the flow between them
package steps

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stnokott/go-timelapse/internal/steps/assemble"
	"github.com/stnokott/go-timelapse/internal/steps/framerate"
	"github.com/stnokott/go-timelapse/internal/steps/input"
	"github.com/stnokott/go-timelapse/internal/steps/output"
	"github.com/stnokott/go-timelapse/internal/steps/summary"
	"github.com/stnokott/go-timelapse/internal/steps/timerange"
	"github.com/stnokott/go-timelapse/internal/style"
)

// Manager handles the flow between the steps
type Manager struct {
	steps     []tea.Model
	stepIndex int
}

// NewManager returns a new instance of Manager
func NewManager() *Manager {
	return &Manager{
		steps: []tea.Model{
			input.NewModel(),
			output.NewModel(),
			timerange.NewModel(),
			assemble.NewModel(),
			framerate.NewModel(),
			&summary.Model{},
		},
		stepIndex: 0,
	}
}

// Init initializes the manager.
// The substeps will be initialized when they are called.
func (m *Manager) Init() tea.Cmd {
	return m.steps[0].Init()
}

// Update performs I/O for the manager
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	stepIndexBefore := m.stepIndex
	switch msg := msg.(type) {
	case input.Msg:
		if msg == input.MsgNext {
			m.stepIndex++
		}
	case output.Msg:
		if msg == output.MsgNext {
			m.stepIndex++
		}
	case timerange.Msg:
		if msg == timerange.MsgNext {
			m.stepIndex++
		}
	case assemble.Msg:
		if msg == assemble.MsgNext {
			m.stepIndex++
		}
	case framerate.Msg:
		if msg == framerate.MsgNext {
			m.stepIndex++
		}
	case summary.Msg:
		if msg == summary.MsgNext {
			m.stepIndex++
		}
	}
	if m.stepIndex >= len(m.steps) {
		// no more steps left
		m.stepIndex = len(m.steps) - 1 // needed to render Println
		return m, tea.Sequence(tea.Println("TODO: end of app flow reached"), tea.Quit)
	}
	// changed step, need to initialize new step
	if stepIndexBefore != m.stepIndex {
		return m, m.steps[m.stepIndex].Init()
	}
	var cmd tea.Cmd
	m.steps[m.stepIndex], cmd = m.steps[m.stepIndex].Update(msg)
	return m, cmd
}

// View renders the model
func (m *Manager) View() string {
	return style.Margin.Render(m.steps[m.stepIndex].View(), "\n")
}
