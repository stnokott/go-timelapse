// Package input lets the user define the input directory
package input

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stnokott/go-timelapse/internal/config"
	"github.com/stnokott/go-timelapse/internal/style"
)

// Model is used to select the directory where the source images for the timelapse are found.
type Model struct {
	table       table.Model
	initSpinner spinner.Model
	dirsLoaded  bool
}

// NewModel creates a new model instance
func NewModel() *Model {
	spin := spinner.New()
	spin.Spinner = spinner.Points

	cols := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Files", Width: 8},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
	)
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(style)

	return &Model{
		table:       t,
		initSpinner: spin,
	}
}

type Msg int

const (
	msgInputDirsLoaded Msg = iota
	MsgNext
)

func (s *Model) loadInputDir() tea.Msg {
	entries, err := os.ReadDir(config.ImagesInputRootDir)
	if err != nil {
		return fmt.Errorf("cannot read input dir %s: %w", config.ImagesInputRootDir, err)
	}
	rows := make([]table.Row, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			count := 0
			subentries, err := os.ReadDir(filepath.Join(config.ImagesInputRootDir, name))
			if err != nil {
				return fmt.Errorf("cannot read dir %s: %w", name, err)
			}
			for _, subentry := range subentries {
				if !subentry.IsDir() {
					count++
				}
			}
			rows = append(rows, table.Row{name, strconv.Itoa(count)})
		}
	}

	s.table.SetRows(rows)

	return msgInputDirsLoaded
}

// Init initializes the step.
func (s *Model) Init() tea.Cmd {
	return tea.Batch(s.initSpinner.Tick, s.loadInputDir)
}

func (s *Model) next() tea.Msg {
	return MsgNext
}

// Update handles I/O for this step.
func (s *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return s, tea.Quit
		case "r":
			// reload directories
			s.dirsLoaded = false
			return s, s.Init()
		case "enter":
			if s.dirsLoaded {
				// save setting
				config.SetInputDir(s.table.SelectedRow()[0])
				return s, s.next
			}
		}
	case Msg:
		if msg == msgInputDirsLoaded {
			s.dirsLoaded = true
		}
	case error:
		return s, tea.Printf("input error: %v", msg)
	}
	if s.dirsLoaded {
		s.table, cmd = s.table.Update(msg)
	} else {
		s.initSpinner, cmd = s.initSpinner.Update(msg)
	}
	return s, cmd
}

// View renders the step.
func (s *Model) View() string {
	var component string
	if s.dirsLoaded {
		component = s.table.View() + "\n(r to reload)"
	} else {
		component = " " + s.initSpinner.View() + " Analyzing input folders... "
	}
	return style.Border.Render("Folders in "+config.ImagesInputRootDir) +
		"\n" +
		style.Border.Render(component)
}
