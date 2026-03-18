package tabs

import (
	"fmt"
	"strings"

	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderAbout(s theme.Styles, p content.Portfolio, width int) string {
	var b strings.Builder

	nameTitle := fmt.Sprintf("%s  %s",
		s.Title.Render(p.Name),
		s.Muted.Render(p.Title),
	)
	b.WriteString(nameTitle)
	b.WriteString("\n")
	b.WriteString(s.Muted.Render(p.Location))
	b.WriteString("\n\n")

	maxWidth := min(width-6, 80)
	aboutLines := strings.Split(strings.TrimSpace(p.About), "\n")
	for _, line := range aboutLines {
		b.WriteString(s.Body.Width(maxWidth).Render(strings.TrimSpace(line)))
		b.WriteString("\n")
	}

	if p.Resume != "" {
		b.WriteString("\n")
		row := fmt.Sprintf("  %s  %-12s %s",
			s.Accent.Render(">>"),
			s.Subtitle.Render("Resume"),
			s.Link.Render(p.Resume),
		)
		b.WriteString(row)
		b.WriteString("\n")
	}

	return b.String()
}
