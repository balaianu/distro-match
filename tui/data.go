package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed distros.json
var distrosJSON []byte

// Distro represents a single Linux distribution.
type Distro struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Version              string   `json:"version"`
	ExperienceLevel      []string `json:"experience_level"`
	UseCases             []string `json:"use_cases"`
	MinRAMGB             float64  `json:"min_ram_gb"`
	RecommendedRAMGB     float64  `json:"recommended_ram_gb,omitempty"`
	MinDiskGB            float64  `json:"min_disk_gb"`
	CPUArchitecture      []string `json:"cpu_architecture"`
	DesktopEnvironments  []string `json:"desktop_environments"`
	ReleaseModel         string   `json:"release_model"`
	PackageManager       string   `json:"package_manager"`
	CommunitySupport     string   `json:"community_support"`
	ProfessionalSupport  bool     `json:"professional_support"`
	DocumentationQuality string   `json:"documentation_quality"`
	Philosophy           string   `json:"philosophy"`
	BasedOn              string   `json:"based_on"`
	Description          string   `json:"description"`
	OfficialWebsite      string   `json:"official_website"`
	DownloadPage         string   `json:"download_page"`
	Strengths            []string `json:"strengths"`
	Weaknesses           []string `json:"weaknesses,omitempty"`
	Score                int      `json:"-"`
}

// Distros loads and returns the embedded distribution database.
func Distros() ([]Distro, error) {
	var distros []Distro
	if err := json.Unmarshal(distrosJSON, &distros); err != nil {
		return nil, fmt.Errorf("failed to parse distros.json: %w", err)
	}
	return distros, nil
}
