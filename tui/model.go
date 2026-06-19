package main

import (
	"fmt"
	"slices"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// screenState identifies the current high-level screen.
type screenState int

const (
	stateWelcome screenState = iota
	stateWizard
	stateResults
	stateDetail
	stateExplorer
)

// appModel is the main Bubble Tea model.
type appModel struct {
	state   screenState
	width   int
	height  int
	distros []Distro
	prefs   Preferences

	// wizard
	currentQuestion int
	selectedOptions []string
	wizardCursor    int

	// viewport for scrollable content (detail/about view)
	viewport viewport.Model

	// progress bar for wizard
	progress progress.Model

	// results
	recommendations []Distro
	resultsCursor   int

	// explorer
	explorerList list.Model

	// welcome
	welcomeCursor int

	// detail
	detailDistroIndex int // index into recommendations or distros
	detailFromResults bool

	// status
	err error
}

// initialModel creates the starting model.
func initialModel() (appModel, error) {
	distros, err := Distros()
	if err != nil {
		return appModel{}, err
	}

	m := appModel{
		state:             stateWelcome,
		distros:           distros,
		detailDistroIndex: -1,
	}
	m.viewport = viewport.New(viewport.WithWidth(0), viewport.WithHeight(0))
	m.viewport.SoftWrap = false
	m.progress = progress.New(
		progress.WithColors(gradientColors...),
		progress.WithFillCharacters('━', '─'),
		progress.WithoutPercentage(),
	)
	m.progress.PercentFormat = " %3.0f%%"
	m.explorerList = newExplorerList(m.distros)
	initHelp()
	return m, nil
}

// Minimum terminal dimensions for a usable UI.
const (
	minWidth  = 50
	minHeight = 15
)

// Init is the Bubble Tea Init function.
func (m appModel) Init() tea.Cmd {
	return nil
}

// Update is the Bubble Tea update loop.
func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.SetWidth(min(m.width-6, 72))
		m.viewport.SetHeight(m.height - 8)
		// Explorer list: reserve space for our centered chrome
		// (title, subtitle, spacing, pagination, help ≈ 10 lines)
		m.explorerList.SetSize(msg.Width, msg.Height-10)
		return m, nil

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.state {
		case stateWelcome:
			return m.updateWelcome(msg)
		case stateWizard:
			return m.updateWizard(msg)
		case stateResults:
			return m.updateResults(msg)
		case stateDetail:
			return m.updateDetail(msg)
		case stateExplorer:
			return m.updateExplorer(msg)
		}

	case list.FilterMatchesMsg:
		if m.state == stateExplorer {
			var cmd tea.Cmd
			m.explorerList, cmd = m.explorerList.Update(msg)
			return m, cmd
		}
	}

	// Forward non-key messages to viewport.
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the current screen.
func (m appModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}
	if m.err != nil {
		return tea.NewView(Card.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)))
	}

	// Minimum size guard
	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("Terminal too small\n\nNeed at least %dx%d\nCurrent: %dx%d\n\nPlease resize your terminal.", minWidth, minHeight, m.width, m.height)
		v := tea.NewView(centerLine(m.width, Card.Render(msg)))
		v.AltScreen = true
		v.WindowTitle = "DistroMatch"
		return v
	}

	var content string
	switch m.state {
	case stateWelcome:
		content = m.welcomeView()
	case stateWizard:
		content = m.wizardView()
	case stateResults:
		content = m.resultsView()
	case stateDetail:
		content = m.detailView()
	case stateExplorer:
		content = m.explorerView()
	}
	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "DistroMatch"
	return v
}

// ── Welcome ──────────────────────────────────────────────────────────────────

var welcomeMenu = []struct{ label, action string }{
	{"Start the Wizard", "wizard"},
	{"Explore Distros", "explore"},
	{"About DistroMatch", "about"},
	{"Quit", "quit"},
}

func (m appModel) updateWelcome(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "space":
		return m.welcomeSelect()
	case "up", "k":
		if m.welcomeCursor > 0 {
			m.welcomeCursor--
		}
	case "down", "j":
		if m.welcomeCursor < len(welcomeMenu)-1 {
			m.welcomeCursor++
		}
	case "1":
		m.welcomeCursor = 0
		return m.welcomeSelect()
	case "2":
		m.welcomeCursor = 1
		return m.welcomeSelect()
	case "3":
		m.welcomeCursor = 2
		return m.welcomeSelect()
	case "4", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m appModel) welcomeSelect() (tea.Model, tea.Cmd) {
	switch welcomeMenu[m.welcomeCursor].action {
	case "wizard":
		m.state = stateWizard
		m.currentQuestion = 0
		m.prefs = Preferences{}
		m.selectedOptions = nil
		m.wizardCursor = 0
	case "explore":
		m.state = stateExplorer
		m.explorerList = newExplorerList(m.distros)
	case "about":
		m.showAbout()
	case "quit":
		return m, tea.Quit
	}
	return m, nil
}

// ── Wizard ───────────────────────────────────────────────────────────────────

func (m appModel) updateWizard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	q := Questions[m.currentQuestion]

	// Number keys select an option directly.
	if len(msg.String()) == 1 && msg.String() >= "1" && msg.String() <= "9" {
		idx := int(msg.String()[0] - '1')
		if idx < len(q.Options) {
			m.wizardCursor = idx
			if q.Type == SingleSelect {
				m.selectedOptions = []string{q.Options[idx].Value}
				return m.advanceWizard()
			}
			m.toggleOption(idx)
		}
		return m, nil
	}

	switch msg.String() {
	case "enter":
		if q.Type == SingleSelect {
			m.selectedOptions = []string{q.Options[m.wizardCursor].Value}
			return m.advanceWizard()
		}
		if m.canAdvance() {
			return m.advanceWizard()
		}
	case "space":
		if q.Type == SingleSelect {
			m.selectedOptions = []string{q.Options[m.wizardCursor].Value}
			return m.advanceWizard()
		}
		m.toggleOption(m.wizardCursor)
	case "down", "j":
		if m.wizardCursor < len(q.Options)-1 {
			m.wizardCursor++
		}
	case "up", "k":
		if m.wizardCursor > 0 {
			m.wizardCursor--
		}
	case "n", "right", "tab":
		if m.canAdvance() {
			return m.advanceWizard()
		}
	case "b", "left", "shift+tab":
		if m.currentQuestion > 0 {
			m.currentQuestion--
			m.selectedOptions = m.restoreSelections()
			m.wizardCursor = 0
			pct := float64(m.currentQuestion) / float64(len(Questions))
			return m, m.progress.SetPercent(pct)
		}
	case "s":
		if !q.Required {
			m.selectedOptions = nil
			return m.advanceWizard()
		}
	case "q", "esc":
		m.state = stateWelcome
		m.welcomeCursor = 0
	}
	return m, nil
}

func (m *appModel) toggleOption(idx int) {
	q := Questions[m.currentQuestion]
	if idx < 0 || idx >= len(q.Options) {
		return
	}
	val := q.Options[idx].Value
	for i, v := range m.selectedOptions {
		if v == val {
			m.selectedOptions = append(m.selectedOptions[:i], m.selectedOptions[i+1:]...)
			return
		}
	}
	m.selectedOptions = append(m.selectedOptions, val)
}

func (m appModel) canAdvance() bool {
	q := Questions[m.currentQuestion]
	if q.Required {
		return len(m.selectedOptions) > 0
	}
	return true
}

func (m appModel) advanceWizard() (tea.Model, tea.Cmd) {
	q := Questions[m.currentQuestion]
	m.storeAnswer(q)

	if m.currentQuestion < len(Questions)-1 {
		m.currentQuestion++
		m.selectedOptions = nil
		m.wizardCursor = 0
		pct := float64(m.currentQuestion) / float64(len(Questions))
		return m, m.progress.SetPercent(pct)
	}

	m.recommendations = getRecommendations(m.prefs, m.distros, 5)
	m.state = stateResults
	m.resultsCursor = 0
	return m, nil
}

func (m appModel) restoreSelections() []string {
	q := Questions[m.currentQuestion]
	switch q.ID {
	case "useCase":
		return append([]string{}, m.prefs.UseCase...)
	case "desktopEnvironment":
		return append([]string{}, m.prefs.DesktopEnvironment...)
	case "supportLevel":
		return append([]string{}, m.prefs.SupportLevel...)
	}
	if q.Type == SingleSelect {
		var val string
		switch q.ID {
		case "experienceLevel":
			val = m.prefs.ExperienceLevel
		case "ram":
			val = m.prefs.Hardware.RAM
		case "disk":
			val = m.prefs.Hardware.Disk
		case "hardwareType":
			val = m.prefs.Hardware.Type
		case "releaseModel":
			val = m.prefs.ReleaseModel
		case "packageManager":
			val = m.prefs.PackageManager
		case "philosophy":
			val = m.prefs.Philosophy
		case "privacyLevel":
			val = m.prefs.PrivacyLevel
		case "learningGoal":
			val = m.prefs.LearningGoal
		}
		if val != "" {
			return []string{val}
		}
	}
	return nil
}

func (m *appModel) storeAnswer(q Question) {
	if q.Type == SingleSelect {
		if len(m.selectedOptions) == 0 {
			return
		}
		val := m.selectedOptions[0]
		switch q.ID {
		case "experienceLevel":
			m.prefs.ExperienceLevel = val
		case "ram":
			m.prefs.Hardware.RAM = val
		case "disk":
			m.prefs.Hardware.Disk = val
		case "hardwareType":
			m.prefs.Hardware.Type = val
		case "releaseModel":
			m.prefs.ReleaseModel = val
		case "packageManager":
			m.prefs.PackageManager = val
		case "philosophy":
			m.prefs.Philosophy = val
		case "privacyLevel":
			m.prefs.PrivacyLevel = val
		case "learningGoal":
			m.prefs.LearningGoal = val
		}
	} else {
		switch q.ID {
		case "useCase":
			m.prefs.UseCase = append([]string{}, m.selectedOptions...)
		case "desktopEnvironment":
			m.prefs.DesktopEnvironment = append([]string{}, m.selectedOptions...)
		case "supportLevel":
			m.prefs.SupportLevel = append([]string{}, m.selectedOptions...)
		}
	}
}

// ── Results ──────────────────────────────────────────────────────────────────

func (m appModel) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.recommendations) == 0 {
		switch msg.String() {
		case "r":
			m.state = stateWizard
			m.currentQuestion = 0
			m.prefs = Preferences{}
			m.selectedOptions = nil
			m.wizardCursor = 0
		case "b", "left", "esc":
			m.state = stateWizard
			m.currentQuestion = len(Questions) - 1
			m.selectedOptions = m.restoreSelections()
			m.wizardCursor = 0
		case "q":
			m.state = stateWelcome
			m.welcomeCursor = 0
		}
		return m, nil
	}

	switch msg.String() {
	case "enter", "space":
		if m.resultsCursor >= 0 && m.resultsCursor < len(m.recommendations) {
			m.detailDistroIndex = m.resultsCursor
			m.detailFromResults = true
			m.state = stateDetail
			m.setViewportContent(m.detailContent(m.recommendations[m.resultsCursor], true))
			m.viewport.SetYOffset(0)
		}
	case "up", "k":
		if m.resultsCursor > 0 {
			m.resultsCursor--
		}
	case "down", "j":
		if m.resultsCursor < len(m.recommendations)-1 {
			m.resultsCursor++
		}
	case "b", "left", "esc":
		m.state = stateWizard
		m.currentQuestion = len(Questions) - 1
		m.selectedOptions = m.restoreSelections()
		m.wizardCursor = 0
	case "r":
		m.state = stateWizard
		m.currentQuestion = 0
		m.prefs = Preferences{}
		m.selectedOptions = nil
		m.wizardCursor = 0
	case "q":
		m.state = stateWelcome
		m.welcomeCursor = 0
	}
	return m, nil
}

// ── Detail ───────────────────────────────────────────────────────────────────

func (m appModel) updateDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "b", "left", "q":
		if m.detailFromResults {
			m.state = stateResults
		} else if m.detailDistroIndex == -1 {
			// About page
			m.state = stateWelcome
			m.welcomeCursor = 0
		} else {
			m.state = stateExplorer
		}
	case "up", "k":
		m.viewport.ScrollUp(1)
	case "down", "j":
		m.viewport.ScrollDown(1)
	case "pgup":
		m.viewport.PageUp()
	case "pgdown":
		m.viewport.PageDown()
	case "home", "g":
		m.viewport.SetYOffset(0)
	case "end", "G":
		m.viewport.SetYOffset(m.viewport.TotalLineCount())
	}
	return m, nil
}

// ── Explorer ─────────────────────────────────────────────────────────────────

func (m appModel) updateExplorer(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// While the list is in filtering mode, forward all keys to the list.
	if m.explorerList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.explorerList, cmd = m.explorerList.Update(msg)
		return m, cmd
	}

	switch msg.String() {
	case "enter", "space":
		if item, ok := m.explorerList.SelectedItem().(distroItem); ok {
			d := item.distro
			m.detailDistroIndex = slices.IndexFunc(m.distros, func(d2 Distro) bool { return d2.ID == d.ID })
			m.detailFromResults = false
			m.state = stateDetail
			m.setViewportContent(m.detailContent(d, false))
			m.viewport.SetYOffset(0)
			return m, nil
		}
	case "q", "b", "left":
		m.state = stateWelcome
		m.welcomeCursor = 0
		return m, nil
	case "esc":
		// If a filter is applied, clear it. Otherwise go back to welcome.
		if m.explorerList.FilterState() == list.FilterApplied {
			m.explorerList.ResetFilter()
			return m, nil
		}
		m.state = stateWelcome
		m.welcomeCursor = 0
		return m, nil
	}

	// Forward all other keys to the list (navigation, filter activation, etc.)
	var cmd tea.Cmd
	m.explorerList, cmd = m.explorerList.Update(msg)
	return m, cmd
}

// ── About ────────────────────────────────────────────────────────────────────

func (m *appModel) showAbout() {
	m.state = stateDetail
	m.detailDistroIndex = -1
	m.detailFromResults = false
	m.setViewportContent(aboutContent(m.width - 6))
	m.viewport.SetYOffset(0)
}

// setViewportContent wraps content to the viewport width (word-aware) and sets it.
func (m *appModel) setViewportContent(content string) {
	w := min(m.width-6, 72)
	m.viewport.SetContent(lipgloss.Wrap(content, w, ""))
}
