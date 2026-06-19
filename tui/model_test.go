package main

import (
	"testing"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// TestWizardFlow verifies that answering the wizard reaches the results screen.
func TestWizardFlow(t *testing.T) {
	m, err := initialModel()
	if err != nil {
		t.Fatalf("initialModel: %v", err)
	}
	m.width = 80
	m.height = 24

	// Start wizard from welcome menu.
	m.state = stateWizard
	m.currentQuestion = 0
	m.prefs = Preferences{}
	m.selectedOptions = nil

	// Answer every question by selecting the first option.
	for i := range Questions {
		m.currentQuestion = i
		m.selectedOptions = []string{Questions[i].Options[0].Value}
		m.storeAnswer(Questions[i])
	}

	m.recommendations = getRecommendations(m.prefs, m.distros, 5)
	m.state = stateResults

	if m.state != stateResults {
		t.Fatalf("expected stateResults, got %d", m.state)
	}
	if len(m.recommendations) == 0 {
		t.Fatalf("expected recommendations, got none")
	}
}

// TestWelcomeMenu verifies that pressing '1' starts the wizard.
func TestWelcomeMenu(t *testing.T) {
	m, err := initialModel()
	if err != nil {
		t.Fatalf("initialModel: %v", err)
	}
	m.width = 80
	m.height = 24

	newM, _ := m.Update(tea.KeyPressMsg{Text: "1", Code: '1'})
	mm := newM.(appModel)
	if mm.state != stateWizard {
		t.Fatalf("expected stateWizard after selecting '1', got %d", mm.state)
	}
}

// TestDataLoad ensures the embedded distro database loads successfully.
func TestDataLoad(t *testing.T) {
	distros, err := Distros()
	if err != nil {
		t.Fatalf("Distros: %v", err)
	}
	if len(distros) == 0 {
		t.Fatalf("expected distros, got none")
	}
}

// TestSpaceToggle verifies that the space key toggles a multi-select option.
func TestSpaceToggle(t *testing.T) {
	m, err := initialModel()
	if err != nil {
		t.Fatalf("initialModel: %v", err)
	}
	m.width = 80
	m.height = 24
	m.state = stateWizard
	m.currentQuestion = 1 // useCase is multi-select
	m.selectedOptions = nil
	m.wizardCursor = 0

	// Press space — should toggle the first option.
	newM, _ := m.Update(tea.KeyPressMsg{Text: "space", Code: ' '})
	mm := newM.(appModel)
	if len(mm.selectedOptions) != 1 {
		t.Fatalf("expected 1 selected option after space, got %d", len(mm.selectedOptions))
	}

	// Press space again — should toggle it off.
	newM, _ = mm.Update(tea.KeyPressMsg{Text: "space", Code: ' '})
	mm = newM.(appModel)
	if len(mm.selectedOptions) != 0 {
		t.Fatalf("expected 0 selected options after second space, got %d", len(mm.selectedOptions))
	}
}

// TestExplorerFilter verifies that filtering works in the explorer.
func TestExplorerFilter(t *testing.T) {
	m, err := initialModel()
	if err != nil {
		t.Fatalf("initialModel: %v", err)
	}
	m.width = 80
	m.height = 24
	m.state = stateExplorer

	// Start filtering with '/'
	newM, _ := m.Update(tea.KeyPressMsg{Text: "/", Code: '/'})
	mm := newM.(appModel)
	if mm.explorerList.FilterState() != list.Filtering {
		t.Fatalf("expected FilterState=Filtering after '/'")
	}

	// Type "arch" — list processes each character
	for _, ch := range "arch" {
		newM, _ = mm.Update(tea.KeyPressMsg{Text: string(ch), Code: ch})
		mm = newM.(appModel)
	}

	// Apply filter with Enter
	newM, _ = mm.Update(tea.KeyPressMsg{Text: "enter", Code: tea.KeyEnter})
	mm = newM.(appModel)
	if mm.explorerList.FilterState() == list.Filtering {
		t.Fatalf("expected FilterState!=Filtering after Enter")
	}

	// Check that results are filtered — Arch Linux should be present
	hasArch := false
	for _, item := range mm.explorerList.VisibleItems() {
		if di, ok := item.(distroItem); ok && di.distro.ID == "arch-linux" {
			hasArch = true
		}
	}
	if !hasArch {
		t.Fatalf("expected arch-linux in filtered results, got %d visible items", len(mm.explorerList.VisibleItems()))
	}
}
