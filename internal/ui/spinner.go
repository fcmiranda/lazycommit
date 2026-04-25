package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type spinnerDoneMsg struct{}

type spinnerModel struct {
	spinner spinner.Model
	message string
	doneCh  <-chan struct{}
}

func newSpinnerModel(message string, doneCh <-chan struct{}) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return spinnerModel{
		spinner: s,
		message: message,
		doneCh:  doneCh,
	}
}

func waitForSpinnerDone(doneCh <-chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-doneCh
		return spinnerDoneMsg{}
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, waitForSpinnerDone(m.doneCh))
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case spinnerDoneMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("\r%s %s", m.spinner.View(), m.message)
}

func RunWithSpinner[T any](message string, fn func() (T, error)) (T, error) {
	var zero T
	doneCh := make(chan struct{})
	resultCh := make(chan struct {
		value T
		err   error
	}, 1)

	go func() {
		defer close(doneCh)
		value, err := fn()
		resultCh <- struct {
			value T
			err   error
		}{value: value, err: err}
	}()

	p := tea.NewProgram(newSpinnerModel(message, doneCh), tea.WithOutput(os.Stderr))
	if _, err := p.Run(); err != nil {
		return zero, err
	}

	fmt.Fprintln(os.Stderr)

	result := <-resultCh
	return result.value, result.err
}
