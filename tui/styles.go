package main

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// ── Color Palette ────────────────────────────────────────────────────────────
// A modern violet/cyan/emerald palette on a slate dark base.
// Extended with pink/indigo for richer gradients and accents.

var (
	// Primary accents
	Violet   = lipgloss.Color("#a78bfa")
	VioletD  = lipgloss.Color("#7c3aed")
	VioletL  = lipgloss.Color("#c4b5fd")
	Indigo   = lipgloss.Color("#818cf8")
	Cyan     = lipgloss.Color("#22d3ee")
	CyanD    = lipgloss.Color("#0891b2")
	Emerald  = lipgloss.Color("#34d399")
	EmeraldD = lipgloss.Color("#059669")
	Amber    = lipgloss.Color("#fbbf24")
	Rose     = lipgloss.Color("#fb7185")
	Pink     = lipgloss.Color("#f472b6")

	// Slate scale
	Slate100 = lipgloss.Color("#f1f5f9")
	Slate200 = lipgloss.Color("#e2e8f0")
	Slate300 = lipgloss.Color("#cbd5e1")
	Slate400 = lipgloss.Color("#94a3b8")
	Slate500 = lipgloss.Color("#64748b")
	Slate600 = lipgloss.Color("#475569")
	Slate700 = lipgloss.Color("#334155")
	Slate800 = lipgloss.Color("#1e293b")
	Slate900 = lipgloss.Color("#0f172a")
	Slate950 = lipgloss.Color("#020617")
)

// ── Reusable styles ──────────────────────────────────────────────────────────

var (
	// Title bar — a pill-shaped label with violet background.
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Slate100).
		Background(VioletD).
		Padding(0, 3)

	// Section header within a screen.
	SubHeader = lipgloss.NewStyle().
			Foreground(Violet).
			Bold(true)

	// Body text.
	Normal = lipgloss.NewStyle().
		Foreground(Slate300)

	// Dimmed/secondary text.
	Dim = lipgloss.NewStyle().
		Foreground(Slate500)

	// A bordered card container.
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Slate700).
		Padding(0, 2).
		Background(Slate900)

	// A card with a violet accent border.
	CardAccent = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Violet).
			Padding(0, 2).
			Background(Slate900)

	// A card with a thick violet left border (for selected/highlighted items).
	CardHighlight = lipgloss.NewStyle().
			BorderLeft(true).
			BorderForeground(Violet).
			Background(Slate800).
			Padding(0, 1, 0, 2)

	// Score styles.
	MatchHigh   = lipgloss.NewStyle().Bold(true).Foreground(Emerald)
	MatchMedium = lipgloss.NewStyle().Bold(true).Foreground(Violet)
	MatchLow    = lipgloss.NewStyle().Bold(true).Foreground(Amber)

	// Badge — a small pill-shaped tag for attributes like family, desktop, etc.
	Badge = lipgloss.NewStyle().
		Foreground(Slate300).
		Background(Slate800).
		Padding(0, 1)

	// BadgeViolet — a violet-accented badge.
	BadgeViolet = lipgloss.NewStyle().
			Foreground(VioletL).
			Background(Slate800).
			Padding(0, 1)

	// BadgeCyan — a cyan-accented badge.
	BadgeCyan = lipgloss.NewStyle().
			Foreground(Cyan).
			Background(Slate800).
			Padding(0, 1)

	// BadgeEmerald — an emerald-accented badge.
	BadgeEmerald = lipgloss.NewStyle().
			Foreground(Emerald).
			Background(Slate800).
			Padding(0, 1)
)

// matchStyle picks a style based on a score percentage.
func matchStyle(score int) lipgloss.Style {
	if score >= 80 {
		return MatchHigh
	}
	if score >= 60 {
		return MatchMedium
	}
	return MatchLow
}

// matchColor returns the raw color for a score, for use in bars.
func matchColor(score int) color.Color {
	if score >= 80 {
		return Emerald
	}
	if score >= 60 {
		return Violet
	}
	return Amber
}

// ── Gradient text ────────────────────────────────────────────────────────────

// gradientColors is a violet→cyan gradient used for the logo.
var gradientColors = []color.Color{
	color.RGBA{R: 0x7c, G: 0x3a, B: 0xed, A: 0xff},
	color.RGBA{R: 0x8b, G: 0x5c, B: 0xf6, A: 0xff},
	color.RGBA{R: 0xa7, G: 0x8b, B: 0xfa, A: 0xff},
	color.RGBA{R: 0xc4, G: 0xb5, B: 0xfd, A: 0xff},
	color.RGBA{R: 0x22, G: 0xd3, B: 0xee, A: 0xff},
}

// renderGradient applies a per-character color gradient to a string.
func renderGradient(s string, colors []color.Color) string {
	runes := []rune(s)
	var b strings.Builder
	for i, r := range runes {
		if r == ' ' {
			b.WriteRune(r)
			continue
		}
		c := colors[i%len(colors)]
		b.WriteString(lipgloss.NewStyle().Foreground(c).Bold(true).Render(string(r)))
	}
	return b.String()
}

// ── Score bar ────────────────────────────────────────────────────────────────

// renderScoreBar renders a compact horizontal bar representing a 0-100 score.
// The bar is `width` characters wide, filled proportionally.
func renderScoreBar(score, width int) string {
	filled := int(float64(score) / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	c := matchColor(score)

	filledPart := lipgloss.NewStyle().Foreground(c).Render(strings.Repeat("━", filled))
	emptyPart := lipgloss.NewStyle().Foreground(Slate700).Render(strings.Repeat("─", width-filled))
	return filledPart + emptyPart
}

// ── Badge helpers ─────────────────────────────────────────────────────────────

// familyBadge renders a distro family as a colored badge.
func familyBadge(family string) string {
	if family == "" || family == "independent" {
		return ""
	}
	colors := map[string]color.Color{
		"ubuntu":      Amber,
		"debian":      Rose,
		"arch":        Cyan,
		"fedora":      Violet,
		"rhel":        Pink,
		"independent": Emerald,
	}
	c, ok := colors[family]
	if !ok {
		c = Slate500
	}
	return lipgloss.NewStyle().Foreground(c).Render(capitalize(family))
}

// capitalize returns s with the first character uppercased.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
