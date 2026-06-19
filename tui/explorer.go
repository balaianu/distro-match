package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// ── List item ────────────────────────────────────────────────────────────────

// distroItem wraps a Distro to satisfy the list.Item and list.DefaultItem interfaces.
type distroItem struct {
	distro Distro
	// precomputed display fields
	title  string
	desc   string
	filter string
}

func newDistroItem(d Distro) distroItem {
	titleStr := fmt.Sprintf("%s %s", d.Name, d.Version)
	var descParts []string
	if d.BasedOn != "" && d.BasedOn != "independent" {
		descParts = append(descParts, capitalize(d.BasedOn))
	}
	descStr := strings.Join(descParts, "  ")

	// FilterValue is what the user searches against. We include all
	// searchable fields matching the website's filter categories.
	filter := strings.ToLower(strings.Join([]string{
		d.Name,
		d.Version,
		d.BasedOn,
		d.Description,
		d.ReleaseModel,
		d.Philosophy,
		d.PackageManager,
		strings.Join(d.ExperienceLevel, " "),
		strings.Join(d.UseCases, " "),
		strings.Join(d.DesktopEnvironments, " "),
		strings.Join(d.CPUArchitecture, " "),
		strings.Join(d.Strengths, " "),
	}, " "))
	// Normalize underscores so "old_hardware" matches "old hardware".
	filter = strings.ReplaceAll(filter, "_", " ")

	return distroItem{distro: d, title: titleStr, desc: descStr, filter: filter}
}

func (i distroItem) FilterValue() string { return i.filter }
func (i distroItem) Title() string       { return i.title }
func (i distroItem) Description() string { return i.desc }

// ── Custom delegate ──────────────────────────────────────────────────────────

// distroDelegate renders list items matching the app's visual language:
// numbered entries with ▸ cursor markers and family tags.
type distroDelegate struct {
	styles distroDelegateStyles
}

type distroDelegateStyles struct {
	num       lipgloss.Style
	cursor    lipgloss.Style
	name      lipgloss.Style
	nameSel   lipgloss.Style
	version   lipgloss.Style
	dimmed    lipgloss.Style
	filterHit lipgloss.Style
}

func newDistroDelegate() distroDelegate {
	return distroDelegate{
		styles: distroDelegateStyles{
			num:       lipgloss.NewStyle().Foreground(Slate600),
			cursor:    lipgloss.NewStyle().Foreground(Emerald).Bold(true),
			name:      lipgloss.NewStyle().Foreground(Slate300),
			nameSel:   lipgloss.NewStyle().Foreground(Slate100).Bold(true),
			version:   lipgloss.NewStyle().Foreground(Slate500),
			dimmed:    lipgloss.NewStyle().Foreground(Slate700),
			filterHit: lipgloss.NewStyle().Foreground(Violet).Bold(true),
		},
	}
}

func (d distroDelegate) Height() int  { return 1 }
func (d distroDelegate) Spacing() int { return 0 }

func (d distroDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d distroDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(distroItem)
	if !ok {
		return
	}

	s := &d.styles
	isSelected := index == m.Index()
	isFiltering := m.FilterState() == list.Filtering
	emptyFilter := isFiltering && m.FilterValue() == ""

	// Number prefix (e.g. "  1." or " 12.")
	num := fmt.Sprintf("%3d.", index+1)

	// Cursor marker
	marker := " "
	if isSelected {
		marker = s.cursor.Render("▸")
	}

	// Truncate name to fit
	maxName := m.Width() - 16 // leave room for number, marker, version, family
	if maxName < 10 {
		maxName = 10
	}
	name := ansi.Truncate(i.distro.Name, maxName, "…")

	// Version (dimmed)
	version := ansi.Truncate(i.distro.Version, 16, "…")

	// Family tag with color
	family := ""
	if i.distro.BasedOn != "" && i.distro.BasedOn != "independent" {
		family = "  " + familyBadge(i.distro.BasedOn)
	}

	// Dim everything when filter is empty
	if emptyFilter {
		fmt.Fprintf(w, "%s %s %s  %s%s",
			s.dimmed.Render(num),
			s.dimmed.Render(" "),
			s.dimmed.Render(name),
			s.dimmed.Render(version),
			s.dimmed.Render(family))
		return
	}

	// Highlight matched runes in the name when filtering
	matchedRunes := m.MatchesForItem(index)
	if len(matchedRunes) > 0 && isFiltering {
		unmatched := s.name.Inline(true)
		matched := unmatched.Inherit(s.filterHit)
		name = lipgloss.StyleRunes(name, matchedRunes, matched, unmatched)
	}

	// Assemble the line
	nameStyle := s.name
	if isSelected {
		nameStyle = s.nameSel
	}

	fmt.Fprintf(w, "%s %s %s  %s%s",
		s.num.Render(num),
		marker,
		nameStyle.Render(name),
		s.version.Render(version),
		family)
}

// ── Key map overrides ────────────────────────────────────────────────────────

// explorerListKeys returns a custom key map for the explorer list.
// We override quit to go back to welcome instead of quitting the program,
// and disable prev/next page keys that conflict with our back navigation.
func explorerListKeys() list.KeyMap {
	km := list.DefaultKeyMap()
	// Disable the list's quit keys — we handle them ourselves.
	km.Quit.SetEnabled(false)
	// Disable prev/next page (left/right) — we use those for back navigation.
	km.PrevPage = key.NewBinding(key.WithKeys("pgup", "h", "u"))
	km.NextPage = key.NewBinding(key.WithKeys("pgdown", "l", "f"))
	return km
}

// newExplorerList creates a configured list.Model for the distro explorer.
// We hide the list's built-in chrome (title, status, pagination, help) and
// render our own centered versions in explorerView.
func newExplorerList(distros []Distro) list.Model {
	items := make([]list.Item, len(distros))
	for i, d := range distros {
		items[i] = newDistroItem(d)
	}

	delegate := newDistroDelegate()
	l := list.New(items, delegate, 80, 20)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false) // hide title bar (also hides filter input)
	l.SetShowFilter(false)
	l.SetFilteringEnabled(true)
	l.KeyMap = explorerListKeys()
	l.Filter = substringFilter

	// Clear the filter input's built-in prompt — we render our own.
	l.FilterInput.Prompt = ""

	// Style the filter input prompt
	l.Styles.DefaultFilterCharacterMatch = lipgloss.NewStyle().Foreground(Violet).Bold(true)

	return l
}

// substringFilter implements list.FilterFunc using substring matching with
// AND logic for space-separated terms. This is more intuitive than the
// default fuzzy subsequence matcher, which produces too many false positives
// on multi-field filter values.
func substringFilter(term string, targets []string) []list.Rank {
	// Normalize underscores in the search term.
	term = strings.ReplaceAll(strings.ToLower(term), "_", " ")
	terms := strings.Fields(term)

	var ranks []list.Rank
	for i, target := range targets {
		t := strings.ToLower(target)
		matched := true
		for _, term := range terms {
			if !strings.Contains(t, term) {
				matched = false
				break
			}
		}
		if matched {
			ranks = append(ranks, list.Rank{Index: i})
		}
	}
	return ranks
}
