package cmd

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)



type model struct {
	table      table.Model
	containers []Container
	width      int
	height     int
	sortCol    string
	sortAsc    bool
}

func buildTable(containers []Container, width, height int) table.Model {
	nameW   := width * 22 / 100
	stackW  := width * 13 / 100
	stateW  := 8
	healthW := 10
	uptimeW := 12
	imageW  := width - nameW - stackW - stateW - healthW - uptimeW - 10

	columns := []table.Column{
		{Title: "Container", Width: nameW},
		{Title: "Stack",     Width: stackW},
		{Title: "State",     Width: stateW},
		{Title: "Health",    Width: healthW},
		{Title: "Uptime",    Width: uptimeW},
		{Title: "Image",     Width: imageW},
	}

	rows := []table.Row{}
	for _, c := range containers {
		rows = append(rows, table.Row{
			c.Name,
			c.Stack,
			c.State,
			c.Health,
			normalizeUptime(c.Uptime),
			truncate(shortImage(c.Image), imageW),
		})
	}

	t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(height-4),
    )

    return styledTable(t)
}

func styledTable(t table.Model) table.Model {
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "240", Dark: "240"}).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("0")).
		Background(lipgloss.Color("12")).
		Bold(true)

	s.Cell = s.Cell.
		Foreground(lipgloss.Color(""))

	t.SetStyles(s)
	return t
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table = buildTable(m.containers, m.width, m.height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "n":
			if m.sortCol == "name" {
				m.sortAsc = !m.sortAsc
			} else {
				m.sortCol = "name"
				m.sortAsc = true
			}
			m.containers = sortContainers(m.containers, m.sortCol, m.sortAsc)
			m.table = buildTable(m.containers, m.width, m.height)

		case "s":
			if m.sortCol == "stack" {
				m.sortAsc = !m.sortAsc
			} else {
				m.sortCol = "stack"
				m.sortAsc = true
			}
			m.containers = sortContainers(m.containers, m.sortCol, m.sortAsc)
			m.table = buildTable(m.containers, m.width, m.height)

		case "h":
			if m.sortCol == "health" {
				m.sortAsc = !m.sortAsc
			} else {
				m.sortCol = "health"
				m.sortAsc = true
			}
			m.containers = sortContainers(m.containers, m.sortCol, m.sortAsc)
			m.table = buildTable(m.containers, m.width, m.height)

		case "u":
			if m.sortCol == "uptime" {
				m.sortAsc = !m.sortAsc
			} else {
				m.sortCol = "uptime"
				m.sortAsc = true
			}
			m.containers = sortContainers(m.containers, m.sortCol, m.sortAsc)
			m.table = buildTable(m.containers, m.width, m.height)
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Render(" labcheck — container status")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render("  ↑/↓ navigate  n name  s stack  h health  u uptime  q quit")

	return title + "\n" + m.table.View() + "\n" + help
}

func runTUI(containers []Container) error {
	containers = sortContainers(containers, "name", true)  // add this

	m := model{
		containers: containers,
		sortCol:    "name",
		sortAsc:    true,
	}

	m.width = 120
	m.height = 40
	m.table = buildTable(containers, m.width, m.height)

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
