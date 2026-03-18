package tabs

import (
	"strings"

	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderProjects(s theme.Styles, p content.Portfolio, width int) string {
	var b strings.Builder
	maxWidth := min(width-6, 80)

	for i, proj := range p.Projects {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(s.Title.Render(proj.Name))
		b.WriteString("\n")
		b.WriteString(s.Tag.Render(proj.Tech))
		b.WriteString("\n\n")
		b.WriteString(s.Body.Width(maxWidth).Render(proj.Description))
		b.WriteString("\n")
	}

	return b.String()
}
