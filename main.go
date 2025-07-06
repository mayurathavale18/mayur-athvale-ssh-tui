package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "224"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	mainStyle := renderer.NewStyle().MarginLeft(2)
	checkboxStyle := renderer.NewStyle().Bold(false).Foreground(lipgloss.Color("213"))
	aboutStyle := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("246"))
	aboutNameStyle := renderer.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	subtleStyle := renderer.NewStyle().Foreground(lipgloss.Color("241"))
	dotStyle := renderer.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)

	m := model{
		Width:          pty.Window.Width,
		Height:         pty.Window.Height,
		Choice:         0,
		Chosen:         false,
		mainStyle:      mainStyle,
		aboutStyle:     aboutStyle,
		aboutNameStyle: aboutNameStyle,
		checkboxStyle:  checkboxStyle,
		subtleStyle:    subtleStyle,
		dotStyle:       dotStyle,
		sess:           s,
		runtime:        "",
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

const (
	dotChar = " • "
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	Width          int
	Height         int
	Choice         int
	Chosen         bool
	mainStyle      lipgloss.Style
	aboutStyle     lipgloss.Style
	aboutNameStyle lipgloss.Style
	checkboxStyle  lipgloss.Style
	subtleStyle    lipgloss.Style
	dotStyle       string
	sess           ssh.Session
	runtime        string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {

	about := m.aboutStyle.Render(fmt.Sprintf(strings.TrimSpace(`
Hi I'm %s,
A self taught developer specialized in many software domains
including Mobile Apps, Web, Backend, Gen AI.
I'm currently working at an AI startup as a FullStack 
Engineer.
I'm fluent in Python, Go, Typescript, Javascript.
`), m.aboutNameStyle.Render("Mayur Athavale")))

	tpl := m.subtleStyle.Render("Hint: q, ctrl+c: quit")

	choices := fmt.Sprintf(
		"%s\n%s",
		m.subtleStyle.Copy().Foreground(lipgloss.Color("13")).Render("GitHub         https://github.com/mayurthavale18"),
		m.subtleStyle.Copy().Foreground(lipgloss.Color("33")).Render("Linkedin       https://linkedin.com/in/mayurathavale1729"),
	)

	s := fmt.Sprintf("%s\n\n%s\n\n%s", about, choices, tpl)
	return m.mainStyle.Render("\n" + s + "\n\n")
}
