// Package assemble collects and prepares the required files for the next steps
package assemble

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
)

type Model struct {
	err         error
	filenames   []string
	spinner     spinner.Model
	filesLoaded bool
}

func NewModel() *Model {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Padding(0, 1)

	return &Model{spinner: spin}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadInputFiles)
}

type Msg int

const (
	msgInputFilesLoaded Msg = iota
	msgInputFilesSorted
	MsgDone
)

func (m *Model) loadInputFiles() tea.Msg {
	entries, err := os.ReadDir(config.Cfg.AbsInputDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		m.filenames = append(m.filenames, name)
	}
	return msgInputFilesLoaded
}

func (m *Model) orderInputFiles() tea.Msg {
	ctimes := make(map[string]time.Time, len(m.filenames))
	for _, filename := range m.filenames {
		fi, err := os.Stat(filepath.Join(config.Cfg.AbsInputDir, filename))
		if err != nil {
			return err
		}
		ctimes[filename] = fi.ModTime()
	}
	sort.Slice(m.filenames, func(i, j int) bool {
		iName, jName := m.filenames[i], m.filenames[j]
		iCtime, jCtime := ctimes[iName], ctimes[jName]
		if iCtime.Equal(jCtime) {
			return strings.Compare(iName, jName) < 0
		}
		return iCtime.Before(jCtime)
	})
	return msgInputFilesSorted
}

func (m *Model) next() tea.Msg {
	return MsgDone
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case Msg:
		switch msg {
		case msgInputFilesLoaded:
			m.filesLoaded = true
			return m, m.orderInputFiles
		case msgInputFilesSorted:
			config.Cfg.ImageNamesSorted = m.filenames
			return m, m.next
		}
	case error:
		m.err = msg
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	errorString := ""
	if m.err != nil {
		errorString = fmt.Sprintf(
			"\n\n%s: %v",
			style.Err.Render("ERROR"),
			m.err,
		)
	}
	var task string
	if !m.filesLoaded {
		task = "Assembling image list"
	} else {
		task = "Sorting images"
	}
	return fmt.Sprintf(
		" %s %s...",
		m.spinner.View(),
		task,
	) +
		"\n" +
		errorString
}
