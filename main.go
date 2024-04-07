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

const (
    padding  = 2
    maxWidth = 80
)

var (
    helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
    workMinutes = flag.Int("work", 25, "Number of minutes to work for")
    breakMinutes = flag.Int("break", 5, "Number of minutes for break")
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
    title        string
    percent      float64
    progress     progress.Model
    totalMinutes int
    isWorkPhase  bool // Indicates if it is the work phase
    timer        time.Time
    currentSession int
    totalSessions int
}

func initialTimerModel(title string) timerModel {
    // Initialize progress bar with default work color
    prog := progress.New(progress.WithScaledGradient("#FF0000", "#FF4500")) // TODO: Change colors to red gradient
    return timerModel{
        title:        title,
        progress:     prog,
        totalMinutes: *workMinutes,
        isWorkPhase:  true,
        percent:      0,
        timer:        time.Now(),
        totalSessions: *totalSessions,
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

        // Assuming each work/break phase duration is stored in m.totalMinutes
        // and m.timer is reset at the start of each phase.
        if elapsed.Minutes() >= float64(m.totalMinutes) {
            if m.isWorkPhase {
                // Work phase just finished
                if m.currentSession < m.totalSessions {
                    // If not on the last session, start break phase
                    m.isWorkPhase = false
                    m.totalMinutes = *breakMinutes
                    m.timer = now // Reset timer for break phase
                    m.progress = progress.New(progress.WithScaledGradient("#76B041", "#A8E05F")) // Optional: Adjust for break phase
                    m.percent = 0 // Reset progress for the break phase
                } else {
                    // Last work phase of the last session finished
                    return m, tea.Quit // All sessions completed
                }
            } else {
                // Break phase just finished
                m.currentSession++ // Move to the next session
                if m.currentSession <= m.totalSessions {
                    // Start the next session's work phase
                    m.isWorkPhase = true
                    m.totalMinutes = *workMinutes
                    m.timer = now // Reset timer for new work phase
                    m.progress = progress.New(progress.WithScaledGradient("#FF0000", "#FF4500")) // Reset progress for work phase
                    m.percent = 0 // Reset progress for the new session
                } else {
                    // All sessions completed
                    return m, tea.Quit
                }
            }
        } else {
            // Update progress based on elapsed time
            percentComplete := elapsed.Minutes() / float64(m.totalMinutes)
            m.percent = percentComplete
        }

        return m, tickCmd() // Continue updating every tick

    case tea.KeyMsg:
        // Quit on any key press
        return m, tea.Quit
    }

    // Handle other messages like window resize
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
