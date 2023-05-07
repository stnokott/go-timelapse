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

func (a *Model) Init() tea.Cmd {
	return tea.Batch(a.spinner.Tick, a.loadInputFiles)
}

type Msg int

const (
	msgInputFilesLoaded Msg = iota
	msgInputFilesSorted
	MsgDone
)

func (a *Model) loadInputFiles() tea.Msg {
	entries, err := os.ReadDir(config.Cfg.AbsInputDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		a.filenames = append(a.filenames, name)
	}
	return msgInputFilesLoaded
}

func (a *Model) orderInputFiles() tea.Msg {
	ctimes := make(map[string]time.Time, len(a.filenames))
	for _, filename := range a.filenames {
		fi, err := os.Stat(filepath.Join(config.Cfg.AbsInputDir, filename))
		if err != nil {
			return err
		}
		ctimes[filename] = fi.ModTime()
	}
	sort.Slice(a.filenames, func(i, j int) bool {
		iName, jName := a.filenames[i], a.filenames[j]
		iCtime, jCtime := ctimes[iName], ctimes[jName]
		if iCtime.Equal(jCtime) {
			return strings.Compare(iName, jName) < 0
		}
		return iCtime.Before(jCtime)
	})
	return msgInputFilesSorted
}

func (a *Model) next() tea.Msg {
	return MsgDone
}

func (a *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		}
	case Msg:
		switch msg {
		case msgInputFilesLoaded:
			a.filesLoaded = true
			return a, a.orderInputFiles
		case msgInputFilesSorted:
			config.Cfg.ImageNamesSorted = a.filenames
			return a, a.next
		}
	case error:
		a.err = msg
		return a, tea.Quit
	}

	var cmd tea.Cmd
	a.spinner, cmd = a.spinner.Update(msg)
	return a, cmd
}

func (a *Model) View() string {
	errorString := ""
	if a.err != nil {
		errorString = fmt.Sprintf(
			"\n\n%s: %v",
			style.Err.Render("ERROR"),
			a.err,
		)
	}
	var task string
	if !a.filesLoaded {
		task = "Assembling image list"
	} else {
		task = "Sorting images"
	}
	return fmt.Sprintf(
		" %s %s...",
		a.spinner.View(),
		task,
	) +
		"\n" +
		errorString
}
