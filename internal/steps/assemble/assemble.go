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

// Model is the assembly model
type Model struct {
	err             error
	fileModTimes    map[string]time.Time
	sortedFilenames []string
	spinner         spinner.Model
	filesLoaded     bool
}

// NewModel returns a new Model instance
func NewModel() *Model {
	spin := spinner.New()
	spin.Spinner = spinner.Points
	spin.Style = lipgloss.NewStyle().Padding(0, 1)

	return &Model{spinner: spin}
}

// Init initializes this step
func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.loadInputFiles)
}

// Msg is used for issuing commands between updates
type Msg int

const (
	msgInputFilesLoaded Msg = iota
	msgInputFilesSorted
	// MsgDone signals that this step is finished
	MsgDone
)

func (m *Model) loadInputFiles() tea.Msg {
	entries, err := os.ReadDir(config.Cfg.AbsInputDir)
	if err != nil {
		return err
	}
	m.fileModTimes = make(map[string]time.Time, 1000)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		fi, err := os.Stat(filepath.Join(config.Cfg.AbsInputDir, name))
		if err != nil {
			return err
		}
		modTime := fi.ModTime()
		// filter out files not in configured time range
		if (modTime.After(config.Cfg.TimeFrom) && modTime.Before(config.Cfg.TimeTo)) || modTime.Equal(config.Cfg.TimeFrom) || modTime.Equal(config.Cfg.TimeTo) {
			m.fileModTimes[name] = modTime
		}
	}
	return msgInputFilesLoaded
}

func (m *Model) sortInputFiles() tea.Msg {
	m.sortedFilenames = make([]string, len(m.fileModTimes))
	i := 0
	for k := range m.fileModTimes {
		m.sortedFilenames[i] = k
		i++
	}
	sort.Slice(m.sortedFilenames, func(i, j int) bool {
		iName, jName := m.sortedFilenames[i], m.sortedFilenames[j]
		iCtime, jCtime := m.fileModTimes[iName], m.fileModTimes[jName]
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

// Update handles I/O for the model
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
			return m, m.sortInputFiles
		case msgInputFilesSorted:
			config.Cfg.ImageNamesSorted = m.sortedFilenames
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

// View renders the model
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
