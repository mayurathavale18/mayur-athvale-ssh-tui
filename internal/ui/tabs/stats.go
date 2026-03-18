package tabs

import (
	"fmt"
	"strings"

	"mayur-athavale-tui/internal/analytics"
	"mayur-athavale-tui/internal/ui/theme"
)

func RenderStats(s theme.Styles, activeCount int, stats analytics.Stats, recent []analytics.RecentVisit, width int) string {
	var b strings.Builder

	b.WriteString(s.Title.Render("Visitor Analytics"))
	b.WriteString("\n\n")

	// Live stats
	b.WriteString(fmt.Sprintf("  %s  %s\n",
		s.Success.Render(fmt.Sprintf("%d", activeCount)),
		s.Body.Render("currently connected"),
	))
	b.WriteString(fmt.Sprintf("  %s  %s\n",
		s.Accent.Render(fmt.Sprintf("%d", stats.TotalVisits)),
		s.Body.Render("total visits"),
	))
	b.WriteString(fmt.Sprintf("  %s  %s\n",
		s.Subtitle.Render(fmt.Sprintf("%d", stats.UniqueIPs)),
		s.Body.Render("unique visitors"),
	))

	if stats.AvgDuration > 0 {
		avgStr := formatDuration(int64(stats.AvgDuration))
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			s.Muted.Render(avgStr),
			s.Body.Render("avg session duration"),
		))
	}

	// Recent visitors table
	if len(recent) > 0 {
		b.WriteString("\n")
		b.WriteString(s.Heading.Render("Recent Visitors"))
		b.WriteString("\n\n")

		// Table header
		header := fmt.Sprintf("  %-22s %-16s %-10s",
			s.Subtitle.Render("Time"),
			s.Subtitle.Render("IP"),
			s.Subtitle.Render("Duration"),
		)
		b.WriteString(header)
		b.WriteString("\n")
		b.WriteString(s.Separator.Render("  "+strings.Repeat("─", min(width-8, 50))))
		b.WriteString("\n")

		for _, v := range recent {
			// Mask last octet of IP for privacy
			maskedIP := maskIP(v.IP)
			row := fmt.Sprintf("  %-22s %-16s %-10s",
				s.Muted.Render(v.ConnectedAt),
				s.Body.Render(maskedIP),
				s.Success.Render(v.Duration),
			)
			b.WriteString(row)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func formatDuration(secs int64) string {
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}
	return fmt.Sprintf("%dm %ds", secs/60, secs%60)
}

func maskIP(ip string) string {
	// Mask last segment for privacy: 192.168.1.100:port -> 192.168.1.***
	parts := strings.Split(ip, ".")
	if len(parts) >= 4 {
		// Remove port from last part if present
		lastPart := parts[3]
		if idx := strings.Index(lastPart, ":"); idx != -1 {
			lastPart = lastPart[:idx]
		}
		_ = lastPart
		parts[3] = "***"
		return strings.Join(parts[:4], ".")
	}
	return ip
}
