package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/analytics"
	"mayur-athavale-tui/internal/ui/components"
	"mayur-athavale-tui/internal/ui/tabs"
	"mayur-athavale-tui/internal/ui/theme"
)

const tabCount = 6 // Added Stats tab

// VisitorTickMsg triggers a refresh of visitor data.
type VisitorTickMsg time.Time

// StatsMsg carries updated analytics data.
type StatsMsg struct {
	Active      int
	Stats       analytics.Stats
	RecentVisits []analytics.RecentVisit
}

type Model struct {
	portfolio    content.Portfolio
	styles       theme.Styles
	keys         KeyMap
	activeTab    int
	width        int
	height       int
	scrollOffset int

	// Analytics
	tracker      *analytics.Tracker
	visitID      int64
	activeCount  int
	stats        analytics.Stats
	recentVisits []analytics.RecentVisit
}

func NewModel(p content.Portfolio, width, height int, tracker *analytics.Tracker, visitID int64, renderer *lipgloss.Renderer) Model {
	return Model{
		portfolio: p,
		styles:    theme.NewStyles(renderer),
		keys:      DefaultKeyMap(),
		activeTab: 0,
		width:     width,
		height:    height,
		tracker:   tracker,
		visitID:   visitID,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.tickVisitors(), m.fetchStats())
}

func (m Model) tickVisitors() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return VisitorTickMsg(t)
	})
}

func (m Model) fetchStats() tea.Cmd {
	return func() tea.Msg {
		stats, _ := m.tracker.GetStats()
		recent, _ := m.tracker.GetRecentVisits(10)
		return StatsMsg{
			Active:       m.tracker.ActiveCount(),
			Stats:        stats,
			RecentVisits: recent,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case VisitorTickMsg:
		return m, tea.Batch(m.tickVisitors(), m.fetchStats())

	case StatsMsg:
		m.activeCount = msg.Active
		m.stats = msg.Stats
		m.recentVisits = msg.RecentVisits

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.tracker.OnDisconnect(m.visitID)
			return m, tea.Quit

		case key.Matches(msg, m.keys.NextTab):
			m.activeTab = (m.activeTab + 1) % tabCount
			m.scrollOffset = 0

		case key.Matches(msg, m.keys.PrevTab):
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
			m.scrollOffset = 0

		case key.Matches(msg, m.keys.Down):
			m.scrollOffset++

		case key.Matches(msg, m.keys.Up):
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}

		case msg.String() == "1":
			m.activeTab = 0
			m.scrollOffset = 0
		case msg.String() == "2":
			m.activeTab = 1
			m.scrollOffset = 0
		case msg.String() == "3":
			m.activeTab = 2
			m.scrollOffset = 0
		case msg.String() == "4":
			m.activeTab = 3
			m.scrollOffset = 0
		case msg.String() == "5":
			m.activeTab = 4
			m.scrollOffset = 0
		case msg.String() == "6":
			m.activeTab = 5
			m.scrollOffset = 0
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	header := components.RenderHeader(m.styles, m.activeTab, m.width, m.portfolio.Name)
	tabContent := m.styles.Content.Render(m.tabContent())
	footer := components.RenderFooter(m.styles, m.width, m.activeCount)

	page := fmt.Sprintf("%s\n%s\n%s", header, tabContent, footer)

	return m.styles.App.
		Width(m.width).
		MaxHeight(m.height).
		Render(page)
}

func (m Model) tabContent() string {
	s := m.styles
	p := m.portfolio
	w := m.width

	switch m.activeTab {
	case 0:
		return tabs.RenderAbout(s, p, w)
	case 1:
		return tabs.RenderExperience(s, p, w, m.scrollOffset)
	case 2:
		return tabs.RenderProjects(s, p, w)
	case 3:
		return tabs.RenderSkills(s, p, w)
	case 4:
		return tabs.RenderContact(s, p)
	case 5:
		return tabs.RenderStats(s, m.activeCount, m.stats, m.recentVisits, w)
	default:
		return lipgloss.NewStyle().Foreground(theme.ColorMuted).Render("Unknown tab")
	}
}
