package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errorModel struct {
	size  tea.WindowSizeMsg
	err   error
	msgCh chan<- tea.Msg
}

func newErrorModel(
	size tea.WindowSizeMsg,
	err error,
	msgCh chan<- tea.Msg,
) tea.Model {
	return errorModel{
		size:  size,
		err:   err,
		msgCh: msgCh,
	}
}

func (errorModel) Init() tea.Cmd {
	return nil
}

func (e errorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return e, tea.Quit
		case "esc":
			e.msgCh <- namespacesViewMsg{}
		}
	}

	return e, nil
}

var errorStyles = struct {
	title   lipgloss.Style
	message lipgloss.Style
}{
	title: lipgloss.
		NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.Color("#FF0000")),
	message: lipgloss.
		NewStyle().
		PaddingLeft(4),
}

func (e errorModel) View() string {
	title := errorStyles.title.Render("error")
	message := errorStyles.message.Render(e.err.Error())
	help := newErrorHelpView()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		message,
		renderErrorHelpView(e.size.Height, title, message, help),
	)
}

type errorModelKeyMap struct {
	back key.Binding
	quit key.Binding
}

func newErrorModelKeyMap() errorModelKeyMap {
	return errorModelKeyMap{
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (m errorModelKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{m.back, m.quit}
}

func (m errorModelKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{m.back, m.quit}}
}

func newErrorHelpView() string {
	return help.New().View(newErrorModelKeyMap())
}

func renderErrorHelpView(
	height int,
	title string,
	message string,
	help string,
) string {
	return lipgloss.
		NewStyle().
		MarginTop(
			height -
				lipgloss.Height(title) -
				lipgloss.Height(message) -
				lipgloss.Height(help),
		).
		PaddingLeft(4).
		Render(help)
}
