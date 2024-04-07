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
}

func initialTimerModel(title string) timerModel {
    // Initialize progress bar with default work color
    prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
    return timerModel{
        title:        title,
        progress:     prog,
        totalMinutes: *workMinutes,
        isWorkPhase:  true,
        percent:      0,
        timer:        time.Now(),
    }
}

func (m timerModel) Init() tea.Cmd {
    // Start the timer
    return tickCmd()
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tickMsg:
        elapsed := time.Since(m.timer).Minutes()
        if m.isWorkPhase && elapsed >= float64(*workMinutes) {
            m.isWorkPhase = false
            m.totalMinutes = *breakMinutes // Update total minutes for break phase
            m.timer = time.Now() // Reset timer for break
            // Re-initialize the progress bar for the break phase with a new color
            m.progress = progress.New(progress.WithScaledGradient("#76B041", "#A8E05F"))
            m.percent = 0 // Reset percent for the break phase
        } else if !m.isWorkPhase && elapsed >= float64(*breakMinutes) {
            return m, tea.Quit // End after break phase
        }

        // Recalculate percent based on phase
        if m.isWorkPhase {
            m.percent = elapsed / float64(*workMinutes)
        } else {
            m.percent = elapsed / float64(*breakMinutes)
        }
        return m, tickCmd()

    case tea.KeyMsg:
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
