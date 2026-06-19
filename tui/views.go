package main

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// ── Layout helpers ───────────────────────────────────────────────────────────

func centerBlock(width int, block string) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, block)
}

func centerLine(width int, s string) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, s)
}

func joinVertical(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, strs...)
}

func wordWrap(text string, maxWidth int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	current := words[0]
	for _, w := range words[1:] {
		if lipgloss.Width(current)+1+lipgloss.Width(w) <= maxWidth {
			current += " " + w
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	lines = append(lines, current)
	return lines
}

// separator renders a thin horizontal divider line.
func separator(width int) string {
	return lipgloss.NewStyle().Foreground(Slate700).Render(strings.Repeat("─", width))
}

// ── Welcome ──────────────────────────────────────────────────────────────────

// asciiLogo is a large stylized penguin in ASCII art.
const asciiLogo = `       .--.
      |o_o |
      |:_/ |
     //   \ \
    (|     | )
   /'\_   _/'\
   \___)=(___/`

func (m appModel) welcomeView() string {
	var b strings.Builder

	b.WriteString("\n")

	// Logo: ASCII art penguin with gradient brand name
	logoLines := strings.Split(asciiLogo, "\n")
	var styledLogo []string
	for _, line := range logoLines {
		styledLogo = append(styledLogo, lipgloss.NewStyle().Foreground(Violet).Render(line))
	}
	penguin := strings.Join(styledLogo, "\n")

	// Large gradient brand name
	brandName := renderGradient("DistroMatch", gradientColors)
	brandStyled := lipgloss.NewStyle().Bold(true).Render(brandName)

	// Tagline under brand
	tagline := lipgloss.NewStyle().Foreground(Slate500).Render("find your perfect linux")

	// Stack brand name + tagline vertically, then place next to penguin
	brandBlock := lipgloss.JoinVertical(lipgloss.Left, brandStyled, "  "+tagline)
	logo := lipgloss.JoinHorizontal(lipgloss.Bottom, penguin, "  ", brandBlock)
	b.WriteString(centerBlock(m.width, logo))
	b.WriteString("\n\n")

	// Decorative gradient line
	accentLine := renderGradient(strings.Repeat("━", 36), gradientColors)
	b.WriteString(centerLine(m.width, accentLine))
	b.WriteString("\n\n")

	b.WriteString(centerLine(m.width, Normal.Render("A guided, privacy-friendly wizard that runs entirely in your terminal.")))
	b.WriteString("\n\n")

	// Feature pills with styled backgrounds
	pills := []string{
		lipgloss.NewStyle().Foreground(Emerald).Background(Slate800).Padding(0, 1).Render("● No tracking"),
		lipgloss.NewStyle().Foreground(Violet).Background(Slate800).Padding(0, 1).Render("● 60+ distros"),
		lipgloss.NewStyle().Foreground(Cyan).Background(Slate800).Padding(0, 1).Render("● No account"),
	}
	b.WriteString(centerLine(m.width, strings.Join(pills, "  ")))
	b.WriteString("\n\n\n")

	// Menu items as a styled block
	menuWidth := 42
	var menuItems []string
	for i, item := range welcomeMenu {
		menuItems = append(menuItems, renderMenuItem(i, item.label, m.welcomeCursor == i, menuWidth))
	}
	menuBlock := joinVertical(menuItems...)
	b.WriteString(centerBlock(m.width, menuBlock))
	b.WriteString("\n\n\n")

	b.WriteString(centerLine(m.width, welcomeHelp()))

	return b.String()
}

func renderMenuItem(index int, label string, selected bool, width int) string {
	num := fmt.Sprintf("%d", index+1)

	var numStr, marker, labelStyled string

	if selected {
		numStr = lipgloss.NewStyle().Foreground(Slate100).Bold(true).Render(num)
		marker = lipgloss.NewStyle().Foreground(Emerald).Bold(true).Render("▸")
		labelStyled = lipgloss.NewStyle().Foreground(Slate100).Bold(true).Render(label)
	} else {
		numStr = lipgloss.NewStyle().Foreground(Violet).Render(num)
		marker = " "
		labelStyled = lipgloss.NewStyle().Foreground(Slate300).Render(label)
	}

	content := fmt.Sprintf(" %s  %s  %s", numStr, marker, labelStyled)

	style := lipgloss.NewStyle().Width(width).Padding(0, 1)
	if selected {
		style = lipgloss.NewStyle().
			Width(width).
			Padding(0, 1).
			Background(Slate800).
			BorderLeft(true).
			BorderForeground(Violet).
			PaddingLeft(1)
	}
	return style.Render(content)
}

// ── Wizard ───────────────────────────────────────────────────────────────────

func (m appModel) wizardView() string {
	q := Questions[m.currentQuestion]
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(centerLine(m.width, Title.Render("DistroMatch Wizard")))
	b.WriteString("\n\n")

	// Step dots — one dot per question, filled = answered, current = highlighted
	var dots []string
	for i := range Questions {
		if i < m.currentQuestion {
			dots = append(dots, lipgloss.NewStyle().Foreground(Violet).Render("●"))
		} else if i == m.currentQuestion {
			dots = append(dots, lipgloss.NewStyle().Foreground(Emerald).Render("◉"))
		} else {
			dots = append(dots, lipgloss.NewStyle().Foreground(Slate700).Render("○"))
		}
	}
	dotLine := strings.Join(dots, " ")
	b.WriteString(centerLine(m.width, dotLine))
	b.WriteString("\n")

	// Progress bar using bubbles/progress with spring animation
	progressWidth := min(m.width-16, 44)
	m.progress.SetWidth(progressWidth)
	pct := float64(m.currentQuestion) / float64(len(Questions))
	bar := m.progress.ViewAs(pct)
	step := fmt.Sprintf("Question %d of %d", m.currentQuestion+1, len(Questions))
	b.WriteString(centerLine(m.width, bar))
	b.WriteString("\n")
	b.WriteString(centerLine(m.width, Dim.Render(step)))
	b.WriteString("\n\n")

	// Question prompt in a styled card
	promptCard := CardAccent.Width(min(m.width-12, 60)).Render(
		SubHeader.Render(q.Prompt))
	b.WriteString(centerBlock(m.width, promptCard))
	b.WriteString("\n\n")

	// Build options as a single block for consistent alignment
	maxOptWidth := 0
	for _, opt := range q.Options {
		w := lipgloss.Width(opt.Icon + "  " + opt.Label)
		if w > maxOptWidth {
			maxOptWidth = w
		}
	}
	blockWidth := min(maxOptWidth+14, m.width-8)

	// Windowing: show only the options that fit in the available height.
	availableHeight := m.height - 16
	if availableHeight < 3 {
		availableHeight = 3
	}

	startIdx := 0
	if m.wizardCursor >= availableHeight {
		startIdx = m.wizardCursor - availableHeight + 1
	}
	endIdx := min(startIdx+availableHeight, len(q.Options))

	var optLines []string
	for i := startIdx; i < endIdx; i++ {
		optLines = append(optLines, renderWizardOption(q, i, q.Options[i], m.wizardCursor == i, m.selectedOptions, blockWidth))
	}
	optBlock := joinVertical(optLines...)
	b.WriteString(centerBlock(m.width, optBlock))

	// Scroll indicator if options are windowed
	if len(q.Options) > availableHeight {
		b.WriteString("\n")
		b.WriteString(centerLine(m.width, Dim.Render(fmt.Sprintf("  showing %d-%d of %d", startIdx+1, endIdx, len(q.Options)))))
	}

	// Selection hint for multi-select
	if q.Type == MultiSelect && len(m.selectedOptions) > 0 {
		b.WriteString("\n")
		labels := formatLabels(q.ID, m.selectedOptions)
		hint := lipgloss.NewStyle().Foreground(Emerald).Render("✓ ") + lipgloss.NewStyle().Foreground(Slate300).Render(strings.Join(labels, ", "))
		b.WriteString(centerLine(m.width, hint))
	}

	// Help line
	b.WriteString("\n\n")
	b.WriteString(centerLine(m.width, wizardHelp(q.Type == MultiSelect, q.Required)))

	return b.String()
}

func renderWizardOption(q Question, index int, opt Option, isCursor bool, selected []string, width int) string {
	num := fmt.Sprintf("%2d", index+1)

	isSelected := false
	for _, s := range selected {
		if s == opt.Value {
			isSelected = true
			break
		}
	}

	var marker, numStr, label string
	if q.Type == MultiSelect {
		if isSelected {
			marker = lipgloss.NewStyle().Foreground(Emerald).Bold(true).Render("◉")
		} else {
			marker = lipgloss.NewStyle().Foreground(Slate600).Render("○")
		}
	} else {
		if isCursor {
			marker = lipgloss.NewStyle().Foreground(Emerald).Bold(true).Render("▸")
		} else {
			marker = lipgloss.NewStyle().Foreground(Slate700).Render("│")
		}
	}

	if isCursor {
		numStr = lipgloss.NewStyle().Foreground(Slate100).Bold(true).Render(num)
	} else {
		numStr = lipgloss.NewStyle().Foreground(Violet).Render(num)
	}

	labelText := opt.Icon + "  " + opt.Label
	if isCursor {
		label = lipgloss.NewStyle().Foreground(Slate100).Bold(true).Render(labelText)
	} else if isSelected {
		label = lipgloss.NewStyle().Foreground(Emerald).Render(labelText)
	} else {
		label = lipgloss.NewStyle().Foreground(Slate300).Render(labelText)
	}

	content := fmt.Sprintf(" %s  %s  %s", numStr, marker, label)

	style := lipgloss.NewStyle().Width(width)
	if isCursor {
		style = lipgloss.NewStyle().
			Width(width).
			Background(Slate800).
			BorderLeft(true).
			BorderForeground(Violet).
			PaddingLeft(1)
	}
	return style.Render(content)
}

// ── Results ──────────────────────────────────────────────────────────────────

func (m appModel) resultsView() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(centerLine(m.width, Title.Render("Your Perfect Match")))
	b.WriteString("\n\n")

	b.WriteString(centerLine(m.width, SubHeader.Render("Top recommendations based on your preferences")))
	b.WriteString("\n\n")

	if len(m.recommendations) == 0 {
		msg := "No distributions matched your preferences.\nTry relaxing some hardware or experience filters."
		cardWidth := min(60, m.width-8)
		card := Card.Width(cardWidth).Render(msg)
		b.WriteString(centerBlock(m.width, card))
		b.WriteString("\n\n")
		b.WriteString(centerLine(m.width, resultsEmptyHelp()))
		return b.String()
	}

	// Preferences summary as a compact one-liner card
	prefLines := m.preferenceLines()
	prefCardHeight := 0
	if len(prefLines) > 0 {
		cardWidth := min(m.width-8, 72)
		innerWidth := cardWidth - 6
		var wrapped []string
		for _, line := range prefLines {
			wrapped = append(wrapped, wordWrap(line, innerWidth)...)
		}
		headerText := lipgloss.NewStyle().Foreground(Violet).Bold(true).Render("Your Preferences")
		detail := strings.Join(wrapped, "\n")
		card := Card.Width(cardWidth).Render(headerText + "\n\n" + detail)
		b.WriteString(centerBlock(m.width, card))
		b.WriteString("\n\n")
		prefCardHeight = lipgloss.Height(card) + 2
	}

	// Recommendation cards — windowed if they don't all fit.
	overhead := 6 + prefCardHeight
	availableHeight := m.height - overhead
	if availableHeight < 3 {
		availableHeight = 3
	}

	// Each result card is 3 lines (name+score, bar+family, blank)
	cardsFit := availableHeight / 3
	numCards := len(m.recommendations)
	startIdx := 0
	if m.resultsCursor >= cardsFit {
		startIdx = m.resultsCursor - cardsFit + 1
	}
	endIdx := min(startIdx+cardsFit, numCards)

	cardWidth := min(m.width-16, 56)
	for i := startIdx; i < endIdx; i++ {
		b.WriteString(centerBlock(m.width, renderResultCard(i, m.recommendations[i], m.resultsCursor == i, cardWidth)))
		b.WriteString("\n")
	}

	// Scroll indicator if cards are windowed
	if numCards > cardsFit {
		b.WriteString("\n")
		b.WriteString(centerLine(m.width, Dim.Render(fmt.Sprintf("  showing %d-%d of %d", startIdx+1, endIdx, numCards))))
	}

	b.WriteString("\n")
	b.WriteString(centerLine(m.width, resultsHelp()))

	return b.String()
}

// renderResultCard renders a single recommendation as a visually rich card.
func renderResultCard(index int, d Distro, isCursor bool, width int) string {
	scoreStr := fmt.Sprintf("%3d%%", d.Score)
	scoreStyled := matchStyle(d.Score).Render(scoreStr)

	nameStyle := lipgloss.NewStyle().Foreground(Slate300).Bold(true)
	if isCursor {
		nameStyle = lipgloss.NewStyle().Foreground(Slate100).Bold(true)
	}

	// Rank badge
	rankBadge := lipgloss.NewStyle().
		Foreground(Slate600).
		Render(fmt.Sprintf("#%d", index+1))

	// Score bar — 12 chars wide
	scoreBar := renderScoreBar(d.Score, 12)

	// Line 1: rank, name, version, score
	line1 := fmt.Sprintf("%s  %s  %s  %s",
		rankBadge,
		nameStyle.Render(d.Name),
		Dim.Render(d.Version),
		scoreStyled)

	// Line 2: score bar + family badge
	family := familyBadge(d.BasedOn)
	line2 := fmt.Sprintf("       %s  %s", scoreBar, Dim.Render("Match"))
	if family != "" {
		line2 += "  " + family
	}

	content := line1 + "\n" + line2

	style := lipgloss.NewStyle().Width(width).Padding(0, 1)
	if isCursor {
		style = lipgloss.NewStyle().
			Width(width).
			Padding(0, 1).
			Background(Slate800).
			BorderLeft(true).
			BorderForeground(Violet).
			PaddingLeft(1)
	}
	return style.Render(content)
}

// preferenceLines returns preference summary as lines of "Label: Value" pairs.
func (m appModel) preferenceLines() []string {
	labelStyle := lipgloss.NewStyle().Foreground(Slate500)
	valueStyle := lipgloss.NewStyle().Foreground(Slate200)

	var lines []string
	add := func(label, value string) {
		if value != "" {
			lines = append(lines, labelStyle.Render(label)+" "+valueStyle.Render(value))
		}
	}

	add("Experience:", FormatLabel("experienceLevel", m.prefs.ExperienceLevel))
	if len(m.prefs.UseCase) > 0 {
		add("Use case:", strings.Join(formatLabels("useCase", m.prefs.UseCase), ", "))
	}
	add("RAM:", FormatLabel("ram", m.prefs.Hardware.RAM))
	add("Disk:", FormatLabel("disk", m.prefs.Hardware.Disk))
	add("Hardware:", FormatLabel("hardwareType", m.prefs.Hardware.Type))
	if len(m.prefs.DesktopEnvironment) > 0 {
		add("Desktop:", strings.Join(formatLabels("desktopEnvironment", m.prefs.DesktopEnvironment), ", "))
	}
	add("Release:", FormatLabel("releaseModel", m.prefs.ReleaseModel))
	add("Package mgr:", FormatLabel("packageManager", m.prefs.PackageManager))
	if len(m.prefs.SupportLevel) > 0 {
		add("Support:", strings.Join(formatLabels("supportLevel", m.prefs.SupportLevel), ", "))
	}
	add("Philosophy:", FormatLabel("philosophy", m.prefs.Philosophy))

	if len(lines) == 0 {
		return nil
	}
	return []string{strings.Join(lines, "  •  ")}
}

// ── Detail ───────────────────────────────────────────────────────────────────

func (m appModel) detailView() string {
	var b strings.Builder

	b.WriteString("\n")
	// Title depends on whether this is a distro detail or the about page
	if m.detailDistroIndex == -1 {
		b.WriteString(centerLine(m.width, Title.Render("About DistroMatch")))
	} else {
		b.WriteString(centerLine(m.width, Title.Render("Distro Details")))
	}
	b.WriteString("\n\n")

	// Cap content width for readability and consistent layout
	contentWidth := min(m.width-6, 72)
	m.viewport.SetWidth(contentWidth)
	m.viewport.SetHeight(m.height - 10)
	// Center each line of the viewport content
	view := m.viewport.View()
	for _, line := range strings.Split(view, "\n") {
		b.WriteString(centerLine(m.width, line))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(centerLine(m.width, detailHelp()))

	return b.String()
}

// ── Explorer ─────────────────────────────────────────────────────────────────

func (m appModel) explorerView() string {
	var b strings.Builder
	l := m.explorerList
	filtering := l.FilterState() == list.Filtering
	filterApplied := l.FilterState() == list.FilterApplied
	visibleCount := len(l.VisibleItems())
	totalCount := len(l.Items())

	// ── Centered title ──────────────────────────────────────────────────────
	b.WriteString("\n")
	b.WriteString(centerLine(m.width, Title.Render("Explore Distros")))
	b.WriteString("\n\n")

	// ── Subtitle / filter bar / status ──────────────────────────────────────
	if filtering {
		// Render the filter input ourselves (the list's title bar is hidden)
		prompt := lipgloss.NewStyle().Foreground(Violet).Bold(true).Render("Filter:")
		input := lipgloss.NewStyle().Foreground(Slate100).Render(l.FilterInput.View())
		filterBar := fmt.Sprintf("%s %s", prompt, input)
		b.WriteString(centerLine(m.width, filterBar))
	} else if filterApplied {
		filterVal := l.FilterInput.Value()
		if len(filterVal) > 20 {
			filterVal = filterVal[:20] + "…"
		}
		info := fmt.Sprintf("“%s”  —  %d of %d distros  (esc to clear)", filterVal, visibleCount, totalCount)
		b.WriteString(centerLine(m.width, Dim.Render(info)))
	} else {
		b.WriteString(centerLine(m.width, SubHeader.Render("Browse all Linux distributions")))
		b.WriteString("\n")
		b.WriteString(centerLine(m.width, Dim.Render(fmt.Sprintf("%d distributions in database", totalCount))))
	}
	b.WriteString("\n\n")

	// ── Distro list ─────────────────────────────────────────────────────────
	// The list's View() now renders only the items (all chrome is hidden).
	// We center the item block in the terminal width.
	listView := l.View()
	b.WriteString(centerBlock(m.width, listView))

	// ── Pagination ──────────────────────────────────────────────────────────
	if l.Paginator.TotalPages > 1 {
		b.WriteString("\n")
		page := l.Paginator.Page + 1
		pagInfo := fmt.Sprintf("page %d of %d  •  %d-%d of %d",
			page, l.Paginator.TotalPages,
			l.Paginator.Page*l.Paginator.PerPage+1,
			min((l.Paginator.Page+1)*l.Paginator.PerPage, visibleCount),
			visibleCount)
		b.WriteString(centerLine(m.width, Dim.Render(pagInfo)))
	}

	// ── Help ────────────────────────────────────────────────────────────────
	b.WriteString("\n\n")
	b.WriteString(centerLine(m.width, explorerHelp(filtering)))

	return b.String()
}

// ── Detail content ───────────────────────────────────────────────────────────

func (m appModel) detailContent(d Distro, showScore bool) string {
	var b strings.Builder
	innerWidth := min(m.width-10, 66)

	// ── Header: name + version in a styled card ─────────────────────────────
	header := lipgloss.NewStyle().Bold(true).Foreground(Slate100).Render(d.Name)
	version := lipgloss.NewStyle().Foreground(Slate500).Render("  " + d.Version)
	b.WriteString(header + version)
	b.WriteString("\n")

	// Score with bar (only shown when arriving from wizard results)
	if showScore {
		scoreStr := fmt.Sprintf("%d%% Match", d.Score)
		scoreBar := renderScoreBar(d.Score, 20)
		b.WriteString(matchStyle(d.Score).Render(scoreStr) + "  " + scoreBar)
		b.WriteString("\n")
	}

	// Family + desktop badges
	var badges []string
	if d.BasedOn != "" && d.BasedOn != "independent" {
		badges = append(badges, familyBadge(d.BasedOn))
	}
	if len(d.DesktopEnvironments) > 0 && d.DesktopEnvironments[0] != "Any" {
		badges = append(badges, lipgloss.NewStyle().Foreground(Cyan).Render(d.DesktopEnvironments[0]))
	}
	if d.ReleaseModel != "" {
		release := strings.ReplaceAll(d.ReleaseModel, "_", " ")
		badges = append(badges, lipgloss.NewStyle().Foreground(Amber).Render(capitalize(release)))
	}
	if len(badges) > 0 {
		b.WriteString(strings.Join(badges, Dim.Render("  •  ")))
	}
	b.WriteString("\n\n")

	// Description
	b.WriteString(Normal.Render(d.Description))
	b.WriteString("\n\n")

	// Separator
	b.WriteString(separator(innerWidth))
	b.WriteString("\n\n")

	// System details in a two-column layout
	b.WriteString(SubHeader.Render("System Details"))
	b.WriteString("\n")

	// System details as a borderless two-column table
	rows := [][]string{
		{"Min RAM:", formatFloat(d.MinRAMGB) + " GB"},
	}
	if d.RecommendedRAMGB > 0 {
		rows = append(rows, []string{"Rec. RAM:", formatFloat(d.RecommendedRAMGB) + " GB"})
	}
	rows = append(rows,
		[]string{"Min Disk:", formatFloat(d.MinDiskGB) + " GB"},
		[]string{"CPU:", strings.Join(d.CPUArchitecture, ", ")},
		[]string{"Pkg Manager:", d.PackageManager},
		[]string{"Release:", strings.ReplaceAll(d.ReleaseModel, "_", " ")},
		[]string{"Desktop:", strings.Join(d.DesktopEnvironments, ", ")},
		[]string{"Support:", d.CommunitySupport},
	)
	if d.ProfessionalSupport {
		rows = append(rows, []string{"Pro Support:", "Available"})
	}

	t := table.New().
		Rows(rows...).
		Border(lipgloss.Border{}).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row < 0 {
				return lipgloss.NewStyle()
			}
			if col == 0 {
				return lipgloss.NewStyle().Foreground(Slate500).Width(14)
			}
			return lipgloss.NewStyle().Foreground(Slate200)
		})
	b.WriteString(t.Render())
	b.WriteString("\n\n")

	// Strengths
	if len(d.Strengths) > 0 {
		b.WriteString(separator(innerWidth))
		b.WriteString("\n\n")
		b.WriteString(SubHeader.Render("Key Strengths"))
		b.WriteString("\n\n")
		for _, s := range d.Strengths {
			check := lipgloss.NewStyle().Foreground(Emerald).Bold(true).Render("✓")
			b.WriteString(fmt.Sprintf("  %s  %s\n", check, Normal.Render(s)))
		}
		b.WriteString("\n")
	}

	// Weaknesses
	if len(d.Weaknesses) > 0 {
		b.WriteString(separator(innerWidth))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(Amber).Bold(true).Render("Things to Consider"))
		b.WriteString("\n\n")
		for _, w := range d.Weaknesses {
			dot := lipgloss.NewStyle().Foreground(Amber).Bold(true).Render("!")
			b.WriteString(fmt.Sprintf("  %s  %s\n", dot, Normal.Render(w)))
		}
		b.WriteString("\n")
	}

	// Links
	if d.OfficialWebsite != "" || d.DownloadPage != "" {
		b.WriteString(separator(innerWidth))
		b.WriteString("\n\n")
		if d.OfficialWebsite != "" {
			linkIcon := lipgloss.NewStyle().Foreground(Cyan).Render("🔗")
			b.WriteString(fmt.Sprintf("  %s  %s\n", linkIcon, lipgloss.NewStyle().Foreground(Cyan).Render(d.OfficialWebsite)))
		}
		if d.DownloadPage != "" {
			dlIcon := lipgloss.NewStyle().Foreground(Emerald).Render("⬇")
			b.WriteString(fmt.Sprintf("  %s  %s\n", dlIcon, lipgloss.NewStyle().Foreground(Emerald).Render(d.DownloadPage)))
		}
	}

	return b.String()
}

// aboutContent renders the about page with rich visual styling.
func aboutContent(width int) string {
	var b strings.Builder
	innerWidth := min(width, 66)

	b.WriteString(Normal.Render("A free, open-source Linux distribution selector."))
	b.WriteString("\n\n")

	// Privacy highlight
	privacyCard := CardAccent.Width(innerWidth).Render(
		lipgloss.NewStyle().Foreground(Emerald).Bold(true).Render("🔒 Privacy First") + "\n\n" +
			Normal.Render("No telemetry, no accounts, no tracking.\nThe distro database is embedded — no external API calls."))
	b.WriteString(privacyCard)
	b.WriteString("\n\n\n")

	// How it works
	b.WriteString(SubHeader.Render("How it works"))
	b.WriteString("\n\n")

	steps := []struct {
		icon      string
		iconColor color.Color
		text      string
	}{
		{"①", Violet, "Answer 12 questions about your experience, hardware, and preferences"},
		{"②", Cyan, "The scoring algorithm matches you to Linux distributions"},
		{"③", Emerald, "Browse detailed profiles of your top recommendations"},
	}

	textWidth := innerWidth - 6
	for _, step := range steps {
		iconStyled := lipgloss.NewStyle().Foreground(step.iconColor).Bold(true).Render(step.icon)
		wrapped := wordWrap(step.text, textWidth)
		for i, line := range wrapped {
			if i == 0 {
				b.WriteString(fmt.Sprintf("  %s  %s\n", iconStyled, Normal.Render(line)))
			} else {
				b.WriteString(fmt.Sprintf("     %s\n", Normal.Render(line)))
			}
		}
	}
	b.WriteString("\n\n")

	// Stats
	b.WriteString(SubHeader.Render("Database"))
	b.WriteString("\n\n")
	stats := []struct {
		label string
		value string
		color color.Color
	}{
		{"Distributions", "60+", Violet},
		{"Questions", "12", Cyan},
		{"Categories", "10", Emerald},
		{"Privacy", "100%", Amber},
	}
	for _, s := range stats {
		labelStyled := lipgloss.NewStyle().Foreground(Slate500).Render(fmt.Sprintf("  %-16s", s.label))
		valueStyled := lipgloss.NewStyle().Foreground(s.color).Bold(true).Render(s.value)
		b.WriteString(labelStyled + valueStyled + "\n")
	}
	b.WriteString("\n\n")

	// Tech
	b.WriteString(SubHeader.Render("Built With"))
	b.WriteString("\n\n")
	tech := []struct {
		name  string
		color color.Color
	}{
		{"Go", Cyan},
		{"Bubble Tea", Violet},
		{"Lip Gloss", Emerald},
		{"Bubbles", Amber},
	}
	var techParts []string
	for _, t := range tech {
		techParts = append(techParts, lipgloss.NewStyle().Foreground(t.color).Render(t.name))
	}
	b.WriteString("  " + strings.Join(techParts, Dim.Render("  •  ")) + "\n")
	b.WriteString("\n\n")

	b.WriteString(renderHelp([]key.Binding{
		key.NewBinding(key.WithKeys("esc", "b", "q", "left"), key.WithHelp("esc/b/q", "back")),
	}))

	return b.String()
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func formatFloat(f float64) string {
	if f == float64(int(f)) {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%g", f)
}
