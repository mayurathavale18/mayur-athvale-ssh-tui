package tabs

import (
	"fmt"
	"strings"

	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderExperience(s theme.Styles, p content.Portfolio, width int, scrollOffset int) string {
	var b strings.Builder
	maxWidth := min(width-6, 80)

	for i, exp := range p.Experience {
		if i > 0 {
			b.WriteString("\n")
		}

		header := fmt.Sprintf("%s @ %s",
			s.Title.Render(exp.Role),
			s.Accent.Render(exp.Company),
		)
		b.WriteString(header)
		b.WriteString("\n")

		meta := fmt.Sprintf("%s  |  %s",
			s.Subtitle.Render(exp.Period),
			s.Muted.Render(exp.Location),
		)
		b.WriteString(meta)
		b.WriteString("\n\n")

		for _, h := range exp.Highlights {
			bullet := fmt.Sprintf("  %s %s",
				s.Accent.Render("->"),
				s.Body.Width(maxWidth-5).Render(h),
			)
			b.WriteString(bullet)
			b.WriteString("\n")
		}
	}

	return ApplyScroll(b.String(), scrollOffset)
}

func ApplyScroll(txt string, offset int) string {
	lines := strings.Split(txt, "\n")
	if offset >= len(lines) {
		offset = max(len(lines)-1, 0)
	}
	if offset < 0 {
		offset = 0
	}
	return strings.Join(lines[offset:], "\n")
}
