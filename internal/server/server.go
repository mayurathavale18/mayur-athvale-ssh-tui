package server

import (
	"context"
	"fmt"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/analytics"
	"mayur-athavale-tui/internal/config"
	"mayur-athavale-tui/internal/ui"
)

type Server struct {
	wishServer *ssh.Server
	portfolio  content.Portfolio
	cfg        config.Config
	tracker    *analytics.Tracker
}

func New(cfg config.Config, portfolio content.Portfolio, tracker *analytics.Tracker) (*Server, error) {
	s := &Server{
		portfolio: portfolio,
		cfg:       cfg,
		tracker:   tracker,
	}

	keyPath := filepath.Join(cfg.HostKeyDir, "id_ed25519")

	srv, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(s.teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating SSH server: %w", err)
	}

	s.wishServer = srv
	return s, nil
}

func (s *Server) ListenAndServe() error {
	log.Info("Starting SSH server", "host", s.cfg.Host, "port", s.cfg.Port)
	return s.wishServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.wishServer.Shutdown(ctx)
}

func (s *Server) teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := sess.Pty()

	// Create a renderer tied to this SSH session's terminal
	renderer := bubbletea.MakeRenderer(sess)

	// Extract visitor info
	ip := sess.RemoteAddr().String()
	clientID := sess.User()

	// Track this connection
	visitID := s.tracker.OnConnect(ip, clientID)

	m := ui.NewModel(
		s.portfolio,
		pty.Window.Width,
		pty.Window.Height,
		s.tracker,
		visitID,
		renderer,
	)

	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
