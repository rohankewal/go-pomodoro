# POM-CLI

POM-CLI is a command line app for pomodoro time management technique. It is a simple and easy to use app that helps you to manage your time effectively.

It is built using the Go programming language and these Charmbrace libraries:

- [bubbletea](github.com/charmbracelet/bubbletea)
- [lipgloss](github.com/charmbracelet/lipgloss)
- [bubbles](github.com/charmbracelet/bubbles/progress)

Major shoutout to the [Charm](https://charm.sh/) team for making these awesome libraries!!

## How To Use

```bash
./main -work 25 -break 5 -sessions 5
```

## Current Features

1. Set the duration of work and break sessions through command line flags
2. Asign the session a title
3. Pause and resume the timer at any time should you need to

## Current Issues

- When the first session ends, the app crashes. It won't start the next session. I belive the issue is with the Update function in the main.go file. I am currently working on fixing this issue.

<hr>

## License

The MIT License (MIT)

Copyright (c) 2024 Rohan Kewalramani

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
