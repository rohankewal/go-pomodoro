package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
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

var musicFile embed.FS

func extractMusic() error {
	musicData, err := musicFile.ReadFile("path/to/your/music.mp3")
	if err != nil {
		return err
	}

	// Write to temp file
	tempFile, err := ioutil.TempFile("", "music*.mp3")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(musicData); err != nil {
		return err
	}

	// Save the path for later use
	var musicPath = tempFile.Name()

	return nil
}

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
	isWorkPhase    bool // Indicates if it is the work phase
	timer          time.Time
	currentSession int
	totalSessions  int
}

func initialTimerModel(title string) timerModel {
	// Initialize progress bar with default work color
	prog := progress.New(progress.WithScaledGradient("#FF0000", "#FF4500"))
	return timerModel{
		title:          title,
		progress:       prog,
		totalMinutes:   *workMinutes,
		isWorkPhase:    true,
		percent:        0,
		timer:          time.Now(),
		totalSessions:  *totalSessions,
		currentSession: 5,
	}
}

func (m timerModel) Init() tea.Cmd {
	// Start the timer
	return tickCmd()
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		now := time.Now()
		elapsed := now.Sub(m.timer)

		// Check if the current phase duration is met
		if elapsed.Minutes() >= float64(m.totalMinutes) {
			// Manage the end of a phase
			if m.isWorkPhase {
				// Transition from Work to Break phase within the same session
				m.isWorkPhase = false
				m.totalMinutes = *breakMinutes
				m.timer = now                                                                // Reset timer for the break phase
				m.progress = progress.New(progress.WithScaledGradient("#76B041", "#A8E05F")) // Adjust for break phase
				m.percent = 0                                                                // Reset progress for the break phase
			} else {
				// Handle end of break phase here and check for session continuation
				if m.currentSession < m.totalSessions {
					// Prepare for the next session
					m.currentSession++ // Move to the next session
					m.isWorkPhase = true
					m.totalMinutes = *workMinutes
					m.timer = now                                                                // Reset timer for new work phase
					m.progress = progress.New(progress.WithScaledGradient("#FF0000", "#FF4500")) // Reset for work phase
					m.percent = 0                                                                // Reset progress for the new session
				} else {
					// All sessions completed
					return m, tea.Quit
				}
			}
		} else {
			// Update progress based on elapsed time
			m.percent = elapsed.Minutes() / float64(m.totalMinutes)
		}

		return m, tickCmd() // Ensure continuous ticking

	case tea.KeyMsg:
		// Optionally, handle key messages (like quitting early)
		return m, tea.Quit
	}

	return m, nil
}

func (m timerModel) View() string {
	progressBarView := m.progress.ViewAs(m.percent)
	phase := "Work Time"
	if !m.isWorkPhase {
		phase = "Break Time"
	}
	// Display session title at the top
	return fmt.Sprintf("\nSession: %s\n\n%s\n\n%s\n\n%s\n",
		m.title, phase, progressBarView, helpStyle("Press any key to quit"))
}

func tickCmd() tea.Cmd {
	// Tick every second to update the timer
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time
