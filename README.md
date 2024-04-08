# POM-CLI

POM-CLI is a command line app for pomodoro time management technique. It is a simple and easy to use app that helps you to manage your time effectively.

It is built using the Go programming language and these Charmbrace libraries:

- [bubbletea](github.com/charmbracelet/bubbletea)
- [lipgloss](github.com/charmbracelet/lipgloss)
- [bubbles](github.com/charmbracelet/bubbles/progress)

Major shoutout to the [Charm](https://charm.sh/) team for making these awesome libraries!!

# Current Features

1. Set the duration of work and break sessions through command line flags
2. Asign the session a title

# Current Issues

- When the first session ends, the app crashes. It won't start the next session. I belive the issue is with the Update function in the main.go file. I am currently working on fixing this issue.

# Roadmap(Planned Features)

- [x] Add a session title
- [ ] Add lo-fi music
- [ ] Add notifications for start and end of session
