// Package timerange lets the user choose the from- and to-times for the input files
package timerange

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
	"github.com/stnokott/go-timelapse/internal/util"
)

// Model is the bubbletea model
type Model struct {
	textInput textinput.Model
	err       error
	from      time.Time
	to        time.Time
	// confirming is true while the timeframe has been input and the program is awaiting user confirmation
	confirming bool
}

// NewModel returns a new instance of Model
func NewModel() *Model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 20

	return &Model{textInput: ti}
}

func (m *Model) help() string {
	if m.confirming {
		return "(press <y> or <n>)"
	}
	return `(timeframe in last 24hrs, e.g. "10:00 - 06:00")`
}

var regexTimeframe = regexp.MustCompile(`(\d{2}:\d{2})\s?-\s?(\d{2}:\d{2})`)

func (m *Model) parseTimeframeInput() (from time.Time, to time.Time, err error) {
	matches := regexTimeframe.FindStringSubmatch(m.textInput.Value())
	if matches == nil {
		err = errors.New("input does not match required format")
		return
	}
	from, to, err = util.GetClosestPastTimeRange(matches[1], matches[2], time.Now())
	if err != nil {
		err = fmt.Errorf("invalid input: %w", err)
		return
	}
	return
}

// Init initializes this step
func (m *Model) Init() tea.Cmd {
	m.textInput.Prompt = "Time range: "
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

// Update handles step I/O
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if !m.confirming {
				m.from, m.to, m.err = m.parseTimeframeInput()
				if m.err == nil {
					m.confirming = true
				}
			}
			m.textInput.Reset()
			return m, nil
		case "y":
			if m.confirming {
				config.SetTimerange(m.from, m.to)
				return m, m.next
			}
		case "n":
			if m.confirming {
				m.confirming = false
				return m, nil
			}
		}
	}
	var cmd tea.Cmd
	if !m.confirming {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

// View renders the model
func (m *Model) View() string {
	errString := ""
	if m.err != nil {
		errString = fmt.Sprintf("\n\n%s: %v", style.Err.Render("ERROR"), m.err)
	}
	component := ""
	if !m.confirming {
		component = m.textInput.View()
	} else {
		component = fmt.Sprintf(
			"Will continue with this timeframe:\n"+
				"From: %s \n"+
				"To:   %s \n"+
				"\nContinue?",
			style.Emphasis.Render(m.from.Format(time.DateTime)),
			style.Emphasis.Render(m.to.Format(time.DateTime)),
		)
	}
	return style.Base.Render(
		component,
		"\n"+style.Background.Render(m.help()),
		errString,
	)
}
