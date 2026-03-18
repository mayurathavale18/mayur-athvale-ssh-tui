package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// ANSI 256 colors — works on all terminals including basic SSH sessions
var (
	ColorPrimary   = lipgloss.Color("141") // purple
	ColorSecondary = lipgloss.Color("38")  // cyan
	ColorAccent    = lipgloss.Color("214") // amber
	ColorMuted     = lipgloss.Color("245") // gray
	ColorText      = lipgloss.Color("252") // light gray
	ColorBright    = lipgloss.Color("15")  // white
	ColorSuccess   = lipgloss.Color("42")  // green
	ColorDim       = lipgloss.Color("238") // dark gray
	ColorBorder    = lipgloss.Color("240") // border gray
)

type Styles struct {
	App        lipgloss.Style
	Header     lipgloss.Style
	TabActive  lipgloss.Style
	TabInactive lipgloss.Style
	TabBar     lipgloss.Style
	Content    lipgloss.Style
	Footer     lipgloss.Style
	Title      lipgloss.Style
	Subtitle   lipgloss.Style
	Heading    lipgloss.Style
	Body       lipgloss.Style
	Muted      lipgloss.Style
	Accent     lipgloss.Style
	Success    lipgloss.Style
	Link       lipgloss.Style
	Tag        lipgloss.Style
	Separator  lipgloss.Style
	KeyHint    lipgloss.Style
	KeyBinding lipgloss.Style
}

// NewStyles creates styles using the given renderer (tied to an SSH session).
func NewStyles(r *lipgloss.Renderer) Styles {
	return Styles{
		App: r.NewStyle().
			Padding(1, 2),

		Header: r.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1),

		TabActive: r.NewStyle().
			Bold(true).
			Foreground(ColorBright).
			Background(ColorPrimary).
			Padding(0, 2),

		TabInactive: r.NewStyle().
			Foreground(ColorMuted).
			Padding(0, 2),

		TabBar: r.NewStyle().
			MarginBottom(1),

		Content: r.NewStyle().
			Padding(1, 0),

		Footer: r.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorDim).
			PaddingTop(1),

		Title: r.NewStyle().
			Bold(true).
			Foreground(ColorBright),

		Subtitle: r.NewStyle().
			Foreground(ColorSecondary).
			Bold(true),

		Heading: r.NewStyle().
			Bold(true).
			Foreground(ColorAccent).
			MarginTop(1),

		Body: r.NewStyle().
			Foreground(ColorText),

		Muted: r.NewStyle().
			Foreground(ColorMuted),

		Accent: r.NewStyle().
			Foreground(ColorAccent).
			Bold(true),

		Success: r.NewStyle().
			Foreground(ColorSuccess),

		Link: r.NewStyle().
			Foreground(ColorSecondary).
			Underline(true),

		Tag: r.NewStyle().
			Foreground(ColorPrimary).
			Background(lipgloss.Color("236")).
			Padding(0, 1),

		Separator: r.NewStyle().
			Foreground(ColorDim),

		KeyHint: r.NewStyle().
			Foreground(ColorMuted),

		KeyBinding: r.NewStyle().
			Foreground(ColorSecondary).
			Bold(true),
	}
}
