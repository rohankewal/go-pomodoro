package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	workMinutes = flag.Int("work", 25, "Number of minutes to work for")
	breakMinutes = flag.Int("break", 5, "Number of minutes for break")
)

func main() {
	flag.Parse() // Parse the command-line flags

	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	percent      float64
	progress     progress.Model
	totalMinutes int
	isWorkPhase  bool // Indicates if it is the work phase
}

func initialModel() model {
	// Initialize progress bar with default work color
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	return model{progress: prog, totalMinutes: *workMinutes, isWorkPhase: true, percent: 0}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		increment := 1.0 / float64(m.totalMinutes)
		m.percent += increment / 60 // Adjusted for a tick every second
		if m.percent >= 1 {
			if m.isWorkPhase {
				// Transition to break phase
				m.isWorkPhase = false
				m.totalMinutes = *breakMinutes
				m.percent = 0 // Reset percent for the break phase
				// Initialize a new progress bar for the break phase with green color
				m.progress = progress.New(progress.WithScaledGradient("#76B041", "#A8E05F"))
				m.progress.Width = maxWidth - padding*2 - 4 // Ensure the width is correctly set
			} else {
				return m, tea.Quit // End after break phase
			}
		}
		return m, tickCmd()

	default:
		return m, nil
	}
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	phase := "Work Time"
	if !m.isWorkPhase {
		phase = "Break Time"
	}
	return "\n" + pad + phase + "\n\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	// Tick every second to simulate minute progress
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}