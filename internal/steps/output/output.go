// Package output lets the user choose the output filename
package output

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
)

type Model struct {
	textInput     textinput.Model
	validationErr error
}

func NewModel() *Model {
	ti := textinput.New()
	ti.Focus()
	ti.Width = 30
	ti.Prompt = "Output file name: "

	return &Model{textInput: ti}
}

func validateOutputFilename(f string) error {
	filename := f + ".mp4"
	if _, err := os.Stat(filepath.Join(config.VideosOutputRootDir, filename)); !os.IsNotExist(err) {
		return errors.New(filename + " already exists")
	}
	return nil
}

func (m *Model) Init() tea.Cmd {
	return m.textInput.Cursor.BlinkCmd()
}

type Msg int

const (
	MsgNext Msg = iota
)

func (*Model) next() tea.Msg {
	return MsgNext
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// omitting "q" because it might be part of desired filename
			return m, tea.Quit
		case "enter":
			// save setting
			val := m.textInput.Value()
			if m.validationErr = validateOutputFilename(val); m.validationErr != nil {
				return m, nil
			}
			config.SetOutputFilename(val + ".mp4")
			return m, m.next
		}
	case error:
		tea.Printf("output error: %v", msg)
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	errString := ""
	if m.validationErr != nil {
		errString = fmt.Sprintf("\n\n%s: %v", style.Err.Render("ERROR"), m.validationErr)
	}
	return style.Base.Render(m.textInput.View()+".mp4", errString) + "\n"
}
