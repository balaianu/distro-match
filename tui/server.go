package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
)

// runSSHServer starts a wish SSH server that serves the TUI to any caller.
// No authentication is required — any username is accepted and the TUI
// is launched directly. This is the same pattern used by terminal.shop.
func runSSHServer(addr, hostKeyPath string) error {
	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		return err
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			done <- nil
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.Shutdown(ctx)
}

// teaHandler is called for each SSH session. It creates a new appModel
// with the terminal dimensions from the SSH PTY.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m, err := initialModel()
	if err != nil {
		// If data fails to load, return a minimal error model.
		m = appModel{err: err}
	}

	pty, _, _ := s.Pty()
	m.width = pty.Window.Width
	m.height = pty.Window.Height
	if m.width < 1 {
		m.width = 80
	}
	if m.height < 1 {
		m.height = 24
	}
	m.viewport.SetWidth(min(m.width-6, 72))
	m.viewport.SetHeight(m.height - 8)

	return m, nil
}
