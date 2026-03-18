package tabs

import (
	"fmt"
	"strings"

	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderContact(s theme.Styles, p content.Portfolio) string {
	var b strings.Builder

	b.WriteString(s.Title.Render("Let's Connect"))
	b.WriteString("\n\n")

	links := []struct {
		label string
		url   string
		icon  string
	}{
		{"GitHub", p.Contact.GitHub, ">>"},
		{"LinkedIn", p.Contact.LinkedIn, ">>"},
		{"Email", p.Contact.Email, ">>"},
		{"Blog", p.Contact.Blog, ">>"},
		{"Website", p.Contact.Website, ">>"},
	}

	for _, link := range links {
		if link.url == "" {
			continue
		}
		row := fmt.Sprintf("  %s  %-12s %s",
			s.Accent.Render(link.icon),
			s.Subtitle.Render(link.label),
			s.Link.Render(link.url),
		)
		b.WriteString(row)
		b.WriteString("\n\n")
	}

	b.WriteString("\n")
	b.WriteString(s.Muted.Render("Thanks for visiting! Feel free to reach out."))

	return b.String()
}
