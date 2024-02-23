package cli

import (
	"bufio"
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joshuasprow/log-viewer/k8s"
	"github.com/joshuasprow/log-viewer/pkg"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type logsModel struct {
	clientset *kubernetes.Clientset
	id        k8s.ResourceId
	model     table.Model
	cols      []table.Column
	logsCh    chan pkg.Result[pkg.LogEntry]
}

func newLogsModel(clientset *kubernetes.Clientset, id k8s.ResourceId) logsModel {
	cols := []table.Column{}

	for _, title := range []string{"level", "time", "message"} {
		cols = append(cols, table.Column{
			Title: title,
			Width: len(title),
		})
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithKeyMap(table.DefaultKeyMap()),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return logsModel{
		clientset: clientset,
		id:        id,
		model:     t,
		cols:      cols,
		logsCh:    make(chan pkg.Result[pkg.LogEntry]),
	}
}

func initLogEntries(
	clientset *kubernetes.Clientset,
	id k8s.ResourceId,
	logsCh chan<- pkg.Result[pkg.LogEntry],
) tea.Cmd {
	return func() tea.Msg {
		defer close(logsCh)

		type R = pkg.Result[pkg.LogEntry]

		req := clientset.
			CoreV1().
			Pods(id.Namespace).
			GetLogs(
				id.Pod,
				&v1.PodLogOptions{
					Container: id.Container,
					Follow:    true,
					TailLines: pkg.Ptr(int64(10)),
				},
			)

		stream, err := req.Stream(context.TODO())
		if err != nil {
			logsCh <- R{Err: fmt.Errorf("get stream: %w", err)}
			return nil
		}
		defer func() {
			if err := stream.Close(); err != nil {
				logsCh <- R{Err: fmt.Errorf("close stream: %w", err)}
			}
		}()

		scanner := bufio.NewScanner(stream)

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				logsCh <- R{Err: fmt.Errorf("scan error: %w", err)}
				return nil
			}

			data := scanner.Bytes()

			v, err := pkg.UnmarshalLogEntry(data)
			if err != nil {
				logsCh <- R{Err: fmt.Errorf("unmarshal log entry %q: %w", string(data), err)}
				return nil
			}

			logsCh <- R{V: v}
		}

		return nil
	}
}

func waitForLogEntry(logsCh <-chan pkg.Result[pkg.LogEntry]) tea.Cmd {
	return func() tea.Msg {
		l := <-logsCh

		if l.Err != nil {
			return fmt.Errorf("wait for log entry: %w", l.Err)
		}

		// skip empty log entries
		if l.V == (pkg.LogEntry{}) {
			return nil
		}

		return l.V
	}
}

func (m logsModel) Init() tea.Cmd {
	return tea.Batch(
		initLogEntries(m.clientset, m.id, m.logsCh),
		waitForLogEntry(m.logsCh),
	)
}

func (m logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		panic(fmt.Errorf("logsModel.Update: %w", msg))
	case pkg.LogEntry:
		rows := m.model.Rows()
		row := []string{msg.Level, msg.Time, msg.Msg}

		for i := range m.cols {
			if len(row[i]) > m.cols[i].Width {
				m.cols[i].Width = len(row[i])
			}
		}

		m.model.SetColumns(m.cols)
		m.model.SetRows(append(rows, row))

		m.model.GotoBottom()

		return m, waitForLogEntry(m.logsCh)
	}

	var cmd tea.Cmd
	m.model, cmd = m.model.Update(msg)
	return m, cmd
}

func (m logsModel) View() string {
	return m.model.View()
}
