package tabs

import (
	"strings"

	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderSkills(s theme.Styles, p content.Portfolio, width int) string {
	var b strings.Builder
	maxWidth := min(width-6, 80)

	for i, cat := range p.Skills.Categories {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(s.Heading.Render(cat.Name))
		b.WriteString("\n")

		var tags []string
		for _, item := range cat.Items {
			tags = append(tags, s.Tag.Render(item))
		}

		b.WriteString(wrapTags(tags, maxWidth))
		b.WriteString("\n")
	}

	return b.String()
}

func wrapTags(tags []string, maxWidth int) string {
	var lines []string
	var currentLine string
	currentWidth := 0

	for _, tag := range tags {
		tagWidth := len(tag) + 1
		if currentWidth+tagWidth > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = tag
			currentWidth = tagWidth
		} else {
			if currentLine != "" {
				currentLine += " " + tag
			} else {
				currentLine = tag
			}
			currentWidth += tagWidth
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
