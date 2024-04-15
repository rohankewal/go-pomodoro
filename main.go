package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	workMinutes   = flag.Int("work", 25, "Number of minutes to work for")
	breakMinutes  = flag.Int("break", 5, "Number of minutes for break")
	totalSessions = flag.Int("sessions", 5, "Total number of sessions")
)

func main() {
	flag.Parse() // Parse the command-line flags

	p := tea.NewProgram(initialTitleModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// Initial title model to capture session title
type titleModel struct {
	textInput textinput.Model
}

func initialTitleModel() titleModel {
	ti := textinput.New()
	ti.Placeholder = "Enter session title"
	ti.Focus()
	return titleModel{textInput: ti}
}

func (m titleModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m titleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return initialTimerModel(m.textInput.Value()), tickCmd()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m titleModel) View() string {
	return "\n" + m.textInput.View() + "\n\nPress Enter to start session..."
}

// Timer model including session title
type timerModel struct {
	title          string
	percent        float64
	progress       progress.Model
	totalMinutes   int
	isWorkPhase    bool
	timer          time.Time
	pausedTime     time.Time
	isPaused       bool
	currentSession int
	totalSessions  int
}

func initialTimerModel(title string) timerModel {
	prog := progress.New(progress.WithScaledGradient("#FF0000", "#FF4500"))
	return timerModel{
		title:          title,
		progress:       prog,
		totalMinutes:   *workMinutes,
		isWorkPhase:    true,
		percent:        0,
		timer:          time.Now(),
		totalSessions:  *totalSessions,
		currentSession: 1,
	}
}

func (m timerModel) Init() tea.Cmd {
	return tickCmd()
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.isPaused {
			return m, nil
		}

		now := time.Now()
		elapsed := now.Sub(m.timer)
		if elapsed.Minutes() >= float64(m.totalMinutes) {
			if m.isWorkPhase {
				m.isWorkPhase = false
				m.totalMinutes = *breakMinutes
				m.timer = time.Now()
			} else {
				if m.currentSession < m.totalSessions {
					m.currentSession++
					m.isWorkPhase = true
					m.totalMinutes = *workMinutes
					m.timer = time.Now()
				} else {
					return m, tea.Quit
				}
			}
			m.percent = 0
			return m, tickCmd()
		} else {
			m.percent = elapsed.Minutes() / float64(m.totalMinutes)
		}

		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "p":
			if !m.isPaused {
				m.isPaused = true
				m.pausedTime = time.Now()
			}
			return m, nil
		case "r":
			if m.isPaused {
				m.isPaused = false
				pauseDuration := time.Since(m.pausedTime)
				m.timer = m.timer.Add(pauseDuration)
				return m, tickCmd()
			}
		case "q", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m timerModel) View() string {
	progressBarView := m.progress.ViewAs(m.percent)
	phase := "Work Time"
	if !m.isWorkPhase {
		phase = "Break Time"
	}
	status := ""
	if m.isPaused {
		status = " (Paused)"
	}
	return fmt.Sprintf("\nSession: %s\n\n%s%s\n\n%s\n\n%s\n",
		m.title, phase, status, progressBarView, helpStyle("Press 'p' to pause, 'r' to resume, 'q' to quit"))
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time
