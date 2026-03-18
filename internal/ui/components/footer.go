package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderFooter(styles theme.Styles, width int, visitors int) string {
	keys := fmt.Sprintf(
		"%s %s  %s %s  %s %s  %s %s",
		styles.KeyBinding.Render("tab/arrows"),
		styles.KeyHint.Render("navigate"),
		styles.KeyBinding.Render("j/k"),
		styles.KeyHint.Render("scroll"),
		styles.KeyBinding.Render("1-6"),
		styles.KeyHint.Render("jump to tab"),
		styles.KeyBinding.Render("q"),
		styles.KeyHint.Render("quit"),
	)

	var visitorInfo string
	if visitors > 0 {
		visitorInfo = styles.Success.Render(fmt.Sprintf("%d visitor(s) connected", visitors))
	}

	footer := lipgloss.JoinHorizontal(
		lipgloss.Top,
		keys,
		lipgloss.NewStyle().Width(max(width-4-lipgloss.Width(keys)-lipgloss.Width(visitorInfo), 1)).Render(""),
		visitorInfo,
	)

	return styles.Footer.Render(footer)
}
