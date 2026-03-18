package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"mayur-athavale-tui/internal/ui/theme"
)

var TabNames = []string{"About", "Experience", "Projects", "Skills", "Contact", "Stats"}

func RenderHeader(styles theme.Styles, activeTab int, width int, name string) string {
	var b strings.Builder

	// Show ASCII art name on the About tab if terminal is wide enough
	asciiWidth := ASCIIWidth(name)
	if activeTab == 0 && asciiWidth <= width-4 {
		ascii := RenderASCII(name)
		styledASCII := styles.Header.Render(ascii)
		b.WriteString(styledASCII)
		b.WriteString("\n\n")
	}

	// Tab bar
	var tabs []string
	for i, tabName := range TabNames {
		if i == activeTab {
			tabs = append(tabs, styles.TabActive.Render(tabName))
		} else {
			tabs = append(tabs, styles.TabInactive.Render(tabName))
		}
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	separator := styles.Separator.Render(strings.Repeat("─", max(width-4, 0)))

	b.WriteString(tabRow)
	b.WriteString("\n")
	b.WriteString(separator)

	return b.String()
}
