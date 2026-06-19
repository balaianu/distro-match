package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
)

// ── Shared key bindings ──────────────────────────────────────────────────────

var (
	keyUp   = key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up"))
	keyDown = key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down"))
	keyQuit = key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit"))
)

// ── Help model ───────────────────────────────────────────────────────────────

var helpModel = help.New()

// initHelp configures the help model with our color scheme.
func initHelp() {
	helpModel.Styles.ShortKey = helpModel.Styles.ShortKey.Foreground(Violet)
	helpModel.Styles.ShortDesc = helpModel.Styles.ShortDesc.Foreground(Slate600)
	helpModel.Styles.ShortSeparator = helpModel.Styles.ShortSeparator.Foreground(Slate700)
}

// renderHelp renders a help line from a slice of key bindings.
func renderHelp(bindings []key.Binding) string {
	return helpModel.ShortHelpView(bindings)
}

// ── Per-screen help ──────────────────────────────────────────────────────────

// welcomeHelp returns the help string for the welcome screen.
//
// Keys: ↑/k up • ↓/j down • ⏎ select • 1-4 quick pick • q quit
func welcomeHelp() string {
	return renderHelp([]key.Binding{
		keyUp, keyDown,
		key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("⏎", "select")),
		key.NewBinding(key.WithKeys("1", "2", "3", "4"), key.WithHelp("1-4", "quick pick")),
		keyQuit,
	})
}

// wizardHelp returns the help string for the wizard screen.
//
// Single-select keys: ⏎ select • ↑/k up • ↓/j down • 1-9 quick pick • n next • b back • [s skip •] esc/q menu
// Multi-select keys:  space toggle • ⏎ confirm • ↑/k up • ↓/j down • n next • b back • [s skip •] esc/q menu
//
// `s skip` is only shown for non-required questions.
func wizardHelp(multi, required bool) string {
	menu := key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "menu"))
	back := key.NewBinding(key.WithKeys("b", "left", "shift+tab"), key.WithHelp("b", "back"))
	next := key.NewBinding(key.WithKeys("n", "right", "tab"), key.WithHelp("n", "next"))
	skip := key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "skip"))
	quickPick := key.NewBinding(key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"), key.WithHelp("1-9", "quick pick"))

	bindings := make([]key.Binding, 0, 8)
	if multi {
		bindings = append(bindings,
			key.NewBinding(key.WithKeys("space"), key.WithHelp("space", "toggle")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("⏎", "confirm")),
		)
	} else {
		bindings = append(bindings,
			key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("⏎", "select")),
		)
	}
	bindings = append(bindings, keyUp, keyDown)
	if !multi {
		bindings = append(bindings, quickPick)
	}
	bindings = append(bindings, next, back)
	if !required {
		bindings = append(bindings, skip)
	}
	bindings = append(bindings, menu)

	return renderHelp(bindings)
}

// resultsHelp returns the help string for the results screen.
//
// Keys: ↑/k up • ↓/j down • ⏎ details • r restart • esc/b back • q menu
func resultsHelp() string {
	return renderHelp([]key.Binding{
		keyUp, keyDown,
		key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("⏎", "details")),
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restart")),
		key.NewBinding(key.WithKeys("esc", "b", "left"), key.WithHelp("esc/b", "back")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "menu")),
	})
}

// resultsEmptyHelp returns the help string for the results screen with no matches.
//
// Keys: r restart • esc/b back • q menu
func resultsEmptyHelp() string {
	return renderHelp([]key.Binding{
		key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restart")),
		key.NewBinding(key.WithKeys("esc", "b", "left"), key.WithHelp("esc/b", "back")),
		key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "menu")),
	})
}

// detailHelp returns the help string for the detail/about screen.
//
// Keys: ↑/k up • ↓/j down • PgUp page up • PgDn page down • g top • G bottom • esc/b/q back
func detailHelp() string {
	return renderHelp([]key.Binding{
		keyUp, keyDown,
		key.NewBinding(key.WithKeys("pgup", "pageup"), key.WithHelp("PgUp", "page up")),
		key.NewBinding(key.WithKeys("pgdown", "pagedown"), key.WithHelp("PgDn", "page down")),
		key.NewBinding(key.WithKeys("g", "home"), key.WithHelp("g", "top")),
		key.NewBinding(key.WithKeys("G", "end"), key.WithHelp("G", "bottom")),
		key.NewBinding(key.WithKeys("esc", "b", "q", "left"), key.WithHelp("esc/b/q", "back")),
	})
}

// explorerHelp returns the help string for the explorer screen.
//
// Browsing keys: ↑/k up • ↓/j down • / filter • ⏎ details • b/q back
// Filtering keys: ⏎ apply • esc cancel
func explorerHelp(filtering bool) string {
	if filtering {
		return renderHelp([]key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("⏎", "apply")),
			key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
		})
	}
	return renderHelp([]key.Binding{
		keyUp, keyDown,
		key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "filter")),
		key.NewBinding(key.WithKeys("enter", "space"), key.WithHelp("⏎", "details")),
		key.NewBinding(key.WithKeys("b", "q", "left"), key.WithHelp("b/q", "back")),
	})
}
