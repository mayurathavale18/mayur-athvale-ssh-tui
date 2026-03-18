package server

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"mayur-athavale-tui/content"
	"mayur-athavale-tui/internal/analytics"
	"mayur-athavale-tui/internal/config"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	defer l.Close()
	_, port, _ := net.SplitHostPort(l.Addr().String())
	return port
}

func generateTestKey(t *testing.T, path string) {
	t.Helper()
	out, err := exec.Command("ssh-keygen", "-t", "ed25519", "-f", path, "-N", "", "-q").CombinedOutput()
	if err != nil {
		t.Fatalf("ssh-keygen: %v: %s", err, out)
	}
}

func setupTestServer(t *testing.T) (*Server, *analytics.Tracker, string) {
	t.Helper()
	dir := t.TempDir()

	generateTestKey(t, filepath.Join(dir, "id_ed25519"))

	port := freePort(t)

	cfg := config.Config{
		Host:       "127.0.0.1",
		Port:       port,
		HostKeyDir: dir,
		DBPath:     filepath.Join(dir, "test.db"),
	}

	store, err := analytics.NewStore(cfg.DBPath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	tracker := analytics.NewTracker(store)

	portfolio := content.Portfolio{
		Name:  "Test User",
		Title: "Developer",
		About: "Test about section",
	}

	srv, err := New(cfg, portfolio, tracker)
	if err != nil {
		t.Fatalf("New server: %v", err)
	}

	return srv, tracker, port
}

func TestServerStartsAndAcceptsConnections(t *testing.T) {
	srv, _, port := setupTestServer(t)

	go srv.ListenAndServe()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	})

	// Wait for server to be ready
	time.Sleep(500 * time.Millisecond)

	// Try to connect via SSH
	sshConfig := &ssh.ClientConfig{
		User:            "testuser",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", port), sshConfig)
	if err != nil {
		t.Fatalf("SSH dial: %v", err)
	}
	defer client.Close()

	t.Log("Server accepted SSH connection successfully")
}

func TestServerTracksVisitors(t *testing.T) {
	srv, tracker, port := setupTestServer(t)

	go srv.ListenAndServe()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	})

	time.Sleep(500 * time.Millisecond)

	// Connect and request a PTY session to trigger the tea handler
	sshConfig := &ssh.ClientConfig{
		User:            "visitor1",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", port), sshConfig)
	if err != nil {
		t.Fatalf("SSH dial: %v", err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		t.Fatalf("NewSession: %v", err)
	}

	// Request PTY (needed for Bubble Tea middleware)
	if err := session.RequestPty("xterm-256color", 24, 80, ssh.TerminalModes{}); err != nil {
		session.Close()
		client.Close()
		t.Fatalf("RequestPty: %v", err)
	}

	if err := session.Shell(); err != nil {
		session.Close()
		client.Close()
		t.Fatalf("Shell: %v", err)
	}

	// Give the handler time to process
	time.Sleep(500 * time.Millisecond)

	// Verify active connection is tracked
	if tracker.ActiveCount() < 1 {
		t.Fatalf("expected at least 1 active visitor, got %d", tracker.ActiveCount())
	}

	// Disconnect
	session.Close()
	client.Close()

	time.Sleep(500 * time.Millisecond)

	// Verify stats recorded
	stats, err := tracker.GetStats()
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	if stats.TotalVisits < 1 {
		t.Fatalf("expected at least 1 total visit, got %d", stats.TotalVisits)
	}
}

func TestServerGracefulShutdown(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	go srv.ListenAndServe()

	time.Sleep(300 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown: %v", err)
	}
}
