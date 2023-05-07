// Package summary shows all configured details and lets the user confirm
package summary

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
)

// Model is the summary model
type Model struct {
	content string
}

// Msg is used for issuing commands between updates
type Msg int

const (
	// MsgDone signals that this step is finished
	MsgDone Msg = iota
)

func (*Model) next() tea.Msg {
	return MsgDone
}

func (m *Model) body() string {
	cfg := config.Cfg
	header := style.Border.Width(width).AlignHorizontal(lipgloss.Center).Render("Summary")
	body := style.Border.Width(width).AlignHorizontal(lipgloss.Left).Render(
		fmt.Sprintf(
			"Input folder: %s\n"+
				"Output file:  %s\n"+
				"From:         %s\n"+
				"To:           %s\n"+
				"              (= %s)\n"+
				"Framerate:    %.2f\n"+
				"Files:        %d\n\n"+
				"Total approx. video length: %s",
			cfg.AbsInputDir,
			cfg.AbsOutpuFilepath,
			cfg.TimeFrom.Format(time.DateTime),
			cfg.TimeTo.Format(time.DateTime),
			cfg.TimeTo.Sub(cfg.TimeFrom),
			cfg.ImgsPerSecond,
			len(cfg.ImageNamesSorted),
			config.GetApproxTotalDuration(),
		),
	)
	prompt := style.Border.Inherit(style.Emphasis).AlignHorizontal(lipgloss.Center).Render(
		"Press <ENTER> to start rendering",
	)
	return header + "\n" + body + "\n" + prompt
}

// Init initializes this step
func (m *Model) Init() tea.Cmd {
	m.content = m.body()
	return nil
}

// Update handles I/O for this step
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			return m, m.next
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

const width int = 60

// View renders the model.
func (m *Model) View() string {
	return m.content
}
