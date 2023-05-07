// Package framerate lets the user choose a framerate using multiple available methods.
package framerate

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
)

// Model lets the user choose a framerate using multiple available methods
type Model struct {
	textInput               textinput.Model
	err                     error
	selectedModelIndex      int
	currentTask             currentTask
	totalDurationPrediction time.Duration
}

type currentTask int

const (
	taskChooseMode currentTask = iota
	taskInput
	taskConfirm
)

type selector struct {
	ToFramerate func(string) (float64, error)
	Name        string
	Description string
	Prompt      string
}

var selectionModels = []selector{
	newNormalSelection(),
	newReverseSelection(),
	newFactorSelection(),
}

// NewModel creates a new instance of Model.
func NewModel() *Model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20
	ti.Prompt = "Framerate (img/s): "

	return &Model{textInput: ti}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	m.textInput.Reset()
	m.currentTask = taskChooseMode
	return m.textInput.Cursor.BlinkCmd()
}

// Msg is used for issuing commands between updates
type Msg int

const (
	// MsgNext signals that this step is finished
	MsgNext Msg = iota
)

func (*Model) next() tea.Msg {
	return MsgNext
}

// Update handles I/O for this step
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			if m.currentTask == taskInput && keyMsg.String() == "q" {
				// pressing q during input will return to mode selection
				return m, m.Init()
			}
			return m, tea.Quit
		case "enter":
			switch m.currentTask {
			case taskChooseMode:
				// selection mode confirmed
				m.currentTask = taskInput
				m.textInput.Prompt = selectionModels[m.selectedModelIndex].Prompt + ": "
				return m, nil
			case taskInput:
				// input for selected mode confirmed
				framerate, err := selectionModels[m.selectedModelIndex].ToFramerate(m.textInput.Value())
				if err != nil {
					m.err = err
					return m, nil
				}
				config.SetImagesPerSecond(framerate)
				m.totalDurationPrediction = config.GetApproxTotalDuration()
				m.currentTask = taskConfirm
				return m, nil
			}
		case "up":
			m.selectedModelIndex--
			if m.selectedModelIndex < 0 {
				m.selectedModelIndex = len(selectionModels) - 1
			}
		case "down":
			m.selectedModelIndex++
			if m.selectedModelIndex >= len(selectionModels) {
				m.selectedModelIndex = 0
			}
		case "y":
			if m.currentTask == taskConfirm {
				// calculated framerate and total duration confirmed
				m.currentTask = taskChooseMode
				return m, m.next
			}
		case "n":
			if m.currentTask == taskConfirm {
				// return to beginning
				return m, m.Init()
			}
		}
	}

	var cmd tea.Cmd
	if m.currentTask == taskInput {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

// View renders the model.
func (m *Model) View() string {
	switch m.currentTask {
	case taskChooseMode:
		var b strings.Builder
		for i, model := range selectionModels {
			// render emphasised if selected
			if i == m.selectedModelIndex {
				b.WriteString(style.Emphasis.Render(model.Name))
			} else {
				b.WriteString(model.Name)
			}
			b.WriteString("\n  " + style.Background.Render(model.Description) + "\n")
		}
		return style.Base.Render(
			"Choose framerate selection mode:\n",
			b.String(),
		)
	case taskInput:
		errString := ""
		if m.err != nil {
			errString = fmt.Sprintf(
				"\n\n%s: %v",
				style.Err.Render("ERROR"),
				m.err,
			)
		}
		return style.Base.Render(
			m.textInput.View(),
			style.Background.Render("\n(<q> to go back to mode selection)"),
			errString,
		)
	default:
		// confirm predicted total duration
		return style.Base.Render(
			fmt.Sprintf(
				"Will continue with this config:\n"+
					"Framerate:                     %s \n"+
					"Approx. output video duration: %s \n"+
					"\nContinue?",
				style.Emphasis.Render(fmt.Sprintf("%.2f", config.Cfg.ImgsPerSecond)),
				style.Emphasis.Render(m.totalDurationPrediction.String()),
			),
			"\n"+style.Background.Render("(press <y> or <n>)"),
		)
	}
}
