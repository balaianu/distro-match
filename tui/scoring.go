package main

import (
	"math"
	"slices"
	"sort"
)

// weights mirror the website's scoring in src/pages/index.astro.
var weights = map[string]float64{
	"experienceLevel":    0.30,
	"useCase":            0.25,
	"hardware":           0.20,
	"desktopEnvironment": 0.10,
	"releaseModel":       0.05,
	"packageManager":     0.05,
	"supportLevel":       0.03,
	"philosophy":         0.02,
}

// scoreExperienceLevel mirrors the website's logic.
func scoreExperienceLevel(d Distro, pref string) float64 {
	if pref == "" || pref == "not_sure" {
		return 0.5
	}
	experienceMap := map[string][]string{
		"beginner":     {"beginner"},
		"intermediate": {"beginner", "intermediate"},
		"advanced":     {"intermediate", "advanced"},
		"expert":       {"advanced", "expert"},
	}
	acceptableLevels := experienceMap[pref]
	exactMatch := slices.Contains(d.ExperienceLevel, pref)
	acceptableMatch := false
	for _, level := range d.ExperienceLevel {
		if slices.Contains(acceptableLevels, level) {
			acceptableMatch = true
			break
		}
	}
	if exactMatch {
		return 1
	}
	if acceptableMatch {
		return 0.8
	}
	return 0.2
}

// scoreUseCase mirrors the website's logic.
func scoreUseCase(d Distro, pref []string) float64 {
	if len(pref) == 0 || (len(pref) == 1 && pref[0] == "not_sure") {
		return 0.5
	}
	useCaseMap := map[string][]string{
		"general_desktop":  {"general_desktop"},
		"development":      {"general_desktop", "development"},
		"server":           {"server"},
		"security":         {"security"},
		"gaming":           {"general_desktop", "gaming"},
		"content_creation": {"general_desktop", "content_creation"},
		"old_hardware":     {"general_desktop", "old_hardware"},
		"privacy":          {"general_desktop", "privacy"},
	}
	totalScore := 0.0
	exactMatches := 0
	for _, p := range pref {
		if p == "not_sure" {
			continue
		}
		acceptableUseCases := useCaseMap[p]
		if acceptableUseCases == nil {
			acceptableUseCases = []string{p}
		}
		exactMatch := slices.Contains(d.UseCases, p)
		acceptableMatch := false
		for _, uc := range d.UseCases {
			if slices.Contains(acceptableUseCases, uc) {
				acceptableMatch = true
				break
			}
		}
		if exactMatch {
			totalScore += 1
			exactMatches++
		} else if acceptableMatch {
			totalScore += 0.8
		} else {
			totalScore += 0.2
		}
	}
	avgScore := math.Min(1, totalScore/float64(len(pref)))
	if exactMatches == len(pref) && len(pref) > 1 {
		return math.Min(1, avgScore+0.1)
	}
	return avgScore
}

// scoreHardware mirrors the website's logic. Only RAM and disk are used in the website.
func scoreHardware(d Distro, h HardwarePrefs) float64 {
	ramScore := 1.0
	cpuScore := 1.0
	diskScore := 1.0

	if h.RAM != "" && h.RAM != "not_sure" {
		ramMap := map[string]int{
			"lt_2gb":  0,
			"2_4gb":   2,
			"4_8gb":   4,
			"8_16gb":  8,
			"gt_16gb": 16,
		}
		requiredRAM := float64(ramMap[h.RAM])
		if requiredRAM == 0 {
			requiredRAM = 4
		}
		if d.MinRAMGB <= requiredRAM {
			if d.MinRAMGB <= requiredRAM/2 {
				ramScore = 1
			} else if d.MinRAMGB <= requiredRAM*0.75 {
				ramScore = 0.9
			} else {
				ramScore = 0.7
			}
			if d.RecommendedRAMGB > 0 && d.RecommendedRAMGB <= requiredRAM {
				ramScore = math.Min(1, ramScore+0.1)
			}
		} else {
			ramScore = 0.2
		}
	}

	if h.CPU != "" && h.CPU != "not_sure" {
		cpuMap := map[string][]string{
			"x86_64": {"x86_64"},
			"i386":   {"x86_64", "i386"},
			"arm64":  {"x86_64", "arm64"},
		}
		acceptableArch := cpuMap[h.CPU]
		if acceptableArch == nil {
			acceptableArch = []string{"x86_64"}
		}
		match := false
		for _, arch := range d.CPUArchitecture {
			if slices.Contains(acceptableArch, arch) {
				match = true
				break
			}
		}
		if match {
			cpuScore = 1
		} else {
			cpuScore = 0.1
		}
	}

	if h.Disk != "" && h.Disk != "not_sure" {
		diskMap := map[string]int{
			"lt_20gb":  0,
			"20_50gb":  20,
			"50_100gb": 50,
			"gt_100gb": 100,
		}
		requiredDisk := float64(diskMap[h.Disk])
		if requiredDisk == 0 {
			requiredDisk = 20
		}
		if d.MinDiskGB <= requiredDisk {
			if d.MinDiskGB <= requiredDisk*0.5 {
				diskScore = 1
			} else if d.MinDiskGB <= requiredDisk*0.75 {
				diskScore = 0.9
			} else {
				diskScore = 0.7
			}
		} else {
			diskScore = 0.2
		}
	}

	score := ramScore*0.4 + cpuScore*0.3 + diskScore*0.3
	return math.Max(0.1, math.Min(1, score))
}

// scoreDesktopEnvironment mirrors the website's logic.
func scoreDesktopEnvironment(d Distro, pref []string) float64 {
	if len(pref) == 0 || (len(pref) == 1 && pref[0] == "no_preference") {
		return 0.5
	}
	deMap := map[string][]string{
		"gnome":    {"GNOME"},
		"kde":      {"KDE Plasma"},
		"xfce":     {"XFCE"},
		"mate":     {"MATE"},
		"cinnamon": {"Cinnamon"},
		"pantheon": {"Pantheon"},
		"i3":       {"i3", "Sway"},
		"lxqt":     {"LXQt", "LXDE"},
	}
	totalScore := 0.0
	count := 0
	for _, p := range pref {
		if p == "no_preference" {
			continue
		}
		count++
		acceptableDEs := deMap[p]
		if acceptableDEs == nil || len(acceptableDEs) == 0 {
			totalScore += 0.5
			continue
		}
		if slices.Contains(d.DesktopEnvironments, "Any") {
			totalScore += 0.9
			continue
		}
		exactMatch := false
		for _, de := range d.DesktopEnvironments {
			if de == p || slices.Contains(acceptableDEs, de) {
				exactMatch = true
				break
			}
		}
		if exactMatch {
			totalScore += 1
		} else {
			totalScore += 0.3
		}
	}
	if count == 0 {
		return 0.5
	}
	return math.Min(1, totalScore/float64(count))
}

// scoreReleaseModel mirrors the website's logic.
func scoreReleaseModel(d Distro, pref string) float64 {
	if pref == "" || pref == "no_preference" {
		return 0.5
	}
	modelMap := map[string][]string{
		"stable_lts":    {"stable_lts"},
		"semi_rolling":  {"semi_rolling"},
		"rolling":       {"rolling", "semi_rolling"},
		"fixed_release": {"fixed_release", "stable_lts"},
	}
	acceptableModels := modelMap[pref]
	exactMatch := d.ReleaseModel == pref
	acceptableMatch := slices.Contains(acceptableModels, d.ReleaseModel)
	if exactMatch {
		return 1
	}
	if acceptableMatch {
		return 0.7
	}
	return 0.3
}

// scorePackageManager mirrors the website's logic.
func scorePackageManager(d Distro, pref string) float64 {
	if pref == "" || pref == "no_preference" {
		return 0.5
	}
	pmMap := map[string][]string{
		"apt":     {"apt"},
		"dnf":     {"dnf"},
		"pacman":  {"pacman"},
		"portage": {"portage"},
		"flatpak": {"apt", "dnf", "pacman"},
		"snap":    {"apt"},
		"zypper":  {"zypper"},
	}
	acceptablePMs := pmMap[pref]
	exactMatch := d.PackageManager == pref
	acceptableMatch := slices.Contains(acceptablePMs, d.PackageManager)
	if exactMatch {
		return 1
	}
	if acceptableMatch {
		return 0.7
	}
	return 0.3
}

// scoreSupportLevel mirrors the website's logic.
func scoreSupportLevel(d Distro, pref []string) float64 {
	if len(pref) == 0 || (len(pref) == 1 && pref[0] == "no_preference") {
		return 0.5
	}
	supportMap := map[string][]string{
		"extensive":     {"extensive", "good"},
		"professional":  {"extensive", "good"},
		"documentation": {"high", "excellent"},
		"minimal":       {"good", "high"},
	}
	totalScore := 0.0
	count := 0
	for _, p := range pref {
		if p == "no_preference" {
			continue
		}
		count++
		acceptableSupport := supportMap[p]
		score := 0.0
		if slices.Contains(acceptableSupport, d.CommunitySupport) {
			score += 0.5
		}
		if d.ProfessionalSupport && p == "professional" {
			score += 0.3
		}
		if slices.Contains(acceptableSupport, d.DocumentationQuality) {
			score += 0.3
		}
		totalScore += math.Min(score, 1)
	}
	if count == 0 {
		return 0.5
	}
	return math.Min(1, totalScore/float64(count))
}

// scorePhilosophy mirrors the website's logic.
func scorePhilosophy(d Distro, pref string) float64 {
	if pref == "" || pref == "no_preference" {
		return 0.5
	}
	philosophyMap := map[string][]string{
		"free_software": {"user_freedom", "pragmatic"},
		"privacy":       {"user_freedom", "pragmatic"},
		"freedom":       {"user_freedom"},
		"corporate":     {"pragmatic"},
	}
	acceptablePhilosophies := philosophyMap[pref]
	exactMatch := d.Philosophy == pref
	acceptableMatch := slices.Contains(acceptablePhilosophies, d.Philosophy)
	if exactMatch {
		return 1
	}
	if acceptableMatch {
		return 0.7
	}
	return 0.3
}

// calculateScores mirrors the website's calculateScores.
func calculateScores(prefs Preferences, distros []Distro) []Distro {
	scored := make([]Distro, len(distros))
	copy(scored, distros)
	for i := range scored {
		total := 0.0
		total += scoreExperienceLevel(scored[i], prefs.ExperienceLevel) * weights["experienceLevel"]
		total += scoreUseCase(scored[i], prefs.UseCase) * weights["useCase"]
		total += scoreHardware(scored[i], prefs.Hardware) * weights["hardware"]
		total += scoreDesktopEnvironment(scored[i], prefs.DesktopEnvironment) * weights["desktopEnvironment"]
		total += scoreReleaseModel(scored[i], prefs.ReleaseModel) * weights["releaseModel"]
		total += scorePackageManager(scored[i], prefs.PackageManager) * weights["packageManager"]
		total += scoreSupportLevel(scored[i], prefs.SupportLevel) * weights["supportLevel"]
		total += scorePhilosophy(scored[i], prefs.Philosophy) * weights["philosophy"]
		scored[i].Score = int(math.Round(total * 100))
	}
	return scored
}

// getRecommendations returns the top matching distros, mirroring the website.
func getRecommendations(prefs Preferences, distros []Distro, limit int) []Distro {
	scored := calculateScores(prefs, distros)
	filtered := make([]Distro, 0, len(scored))
	for _, d := range scored {
		if d.Score >= 30 {
			filtered = append(filtered, d)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Score > filtered[j].Score
	})
	if limit > len(filtered) {
		limit = len(filtered)
	}
	return filtered[:limit]
}

func formatLabels(kind string, values []string) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = FormatLabel(kind, v)
	}
	return out
}
