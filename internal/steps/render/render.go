// Package render renders the video and views the progress
package render

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Model is the rendering model
type Model struct {
	currentTask string
	spinner     spinner.Model
	tmpDir      string
}

// NewModel returns a new Model instance
func NewModel() *Model {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Padding(0, 1)

	return &Model{
		spinner: spin,
	}
}

// Msg is used for issuing commands between updates
type Msg int

const (
	msgPrepareDone Msg = iota
	msgRenderDone
	// MsgDone signals that this step is finished
	MsgDone
)

// creates symlinks input images to get filenames that start with 0
// (better accepted by ffmpeg)
func (m *Model) prepare() tea.Msg {
	var err error
	m.tmpDir, err = os.MkdirTemp("", "*")
	if err != nil {
		return fmt.Errorf("could not create temp dir: %w", err)
	}
	for i, filename := range config.Cfg.ImageNamesSorted {
		// create symlink
		symlinkFilename := fmt.Sprintf("%010d.jpg", i)
		err = os.Symlink(
			filepath.Join(config.Cfg.AbsInputDir, filename),
			filepath.Join(m.tmpDir, symlinkFilename),
		)
		if err != nil {
			return err
		}
	}
	return msgPrepareDone
}

func (m *Model) render() tea.Msg {
	input := filepath.Join(m.tmpDir, "%010d.jpg")
	output := config.Cfg.AbsOutpuFilepath
	outputArgs := ffmpeg.KwArgs{
		"c:v":       "libx264",
		"framerate": config.Cfg.ImgsPerSecond,
		"crf":       23,
	}
	if err := ffmpeg.
		Input(input).
		Output(output, outputArgs).
		Silent(true).
		Run(); err != nil {
		return err
	}
	return msgRenderDone
}

func (m *Model) cleanup() tea.Msg {
	if err := os.RemoveAll(m.tmpDir); err != nil {
		return err
	}
	return MsgDone
}

// Init initializes this step
func (m *Model) Init() tea.Cmd {
	m.currentTask = "Preparing"
	return tea.Batch(m.spinner.Tick, m.prepare)
}

// Update handles I/O for this step
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case Msg:
		switch msg {
		case msgPrepareDone:
			m.currentTask = "Rendering"
			return m, m.render
		case msgRenderDone:
			m.currentTask = "Cleaning up"
			return m, m.cleanup
		}
	case error:
		tea.Printf("render error: %v", msg)
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the model.
func (m *Model) View() string {
	return style.Border.Render(m.currentTask, m.spinner.View())
}
