package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	batchFlag := flag.String("batch", "", "Run in batch mode with comma-separated answer indices (1-based)")
	serveFlag := flag.String("serve", "", "Start SSH server (e.g. :2323 or 0.0.0.0:2323)")
	hostKeyFlag := flag.String("host-key", ".ssh/host_ed25519", "Path to SSH host key (generated if missing)")
	flag.Parse()

	if *serveFlag != "" {
		if err := runSSHServer(*serveFlag, *hostKeyFlag); err != nil {
			fmt.Fprintf(os.Stderr, "SSH server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	m, err := initialModel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading distros: %v\n", err)
		os.Exit(1)
	}

	if *batchFlag != "" {
		if err := runBatch(m, *batchFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Batch error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// runBatch executes the wizard with the provided answers and prints the top recommendations.
func runBatch(m appModel, answers string) error {
	parts := strings.Split(answers, ",")
	if len(parts) != len(Questions) {
		return fmt.Errorf("expected %d answers, got %d", len(Questions), len(parts))
	}

	m.state = stateWizard
	m.currentQuestion = 0
	m.prefs = Preferences{}
	m.selectedOptions = nil

	for i, q := range Questions {
		idx, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil || idx < 1 || idx > len(q.Options) {
			return fmt.Errorf("invalid answer for question %d: %s", i+1, parts[i])
		}
		idx-- // convert to 0-based

		if q.Type == SingleSelect {
			m.selectedOptions = []string{q.Options[idx].Value}
		} else {
			m.selectedOptions = []string{q.Options[idx].Value}
		}
		m.storeAnswer(q)
		m.currentQuestion++
	}

	m.currentQuestion = len(Questions) - 1
	m.recommendations = getRecommendations(m.prefs, m.distros, 5)
	m.state = stateResults

	fmt.Println(renderBatchResults(m))
	return nil
}

// renderBatchResults prints a simple text summary of the recommendations.
func renderBatchResults(m appModel) string {
	var b strings.Builder
	b.WriteString("DistroMatch Batch Results\n")
	b.WriteString("=========================\n\n")
	if len(m.recommendations) == 0 {
		b.WriteString("No distributions matched your preferences.\n")
		return b.String()
	}
	for i, d := range m.recommendations {
		b.WriteString(fmt.Sprintf("%d. %s %s - %d%% match\n", i+1, d.Name, d.Version, d.Score))
		b.WriteString(fmt.Sprintf("   %s\n", d.Description))
		b.WriteString(fmt.Sprintf("   Website: %s\n", d.OfficialWebsite))
		b.WriteString("\n")
	}
	return b.String()
}
