package main

// QuestionType determines how the user can answer a question.
type QuestionType string

const (
	SingleSelect QuestionType = "single"
	MultiSelect  QuestionType = "multi"
)

// Option represents one selectable answer for a question.
type Option struct {
	Value string
	Label string
	Icon  string
}

// Question represents a single wizard step.
type Question struct {
	ID       string
	Prompt   string
	Type     QuestionType
	Required bool
	Options  []Option
}

// HardwarePrefs groups the hardware-related answers.
type HardwarePrefs struct {
	RAM  string `json:"ram"`
	CPU  string `json:"cpu"`
	Disk string `json:"disk"`
	Type string `json:"type"`
}

// Preferences collects all answers from the wizard.
type Preferences struct {
	ExperienceLevel    string
	UseCase            []string
	Hardware           HardwarePrefs
	DesktopEnvironment []string
	ReleaseModel       string
	PackageManager     string
	SupportLevel       []string
	Philosophy         string
	PrivacyLevel       string
	LearningGoal       string
}

// Questions is the ordered list of wizard steps.
var Questions = []Question{
	{
		ID:       "experienceLevel",
		Prompt:   "What is your experience level with Linux?",
		Type:     SingleSelect,
		Required: true,
		Options: []Option{
			{Value: "beginner", Label: "Beginner - Never used Linux or very limited experience", Icon: "🌱"},
			{Value: "intermediate", Label: "Intermediate - Some Linux experience, comfortable with basic commands", Icon: "🌿"},
			{Value: "advanced", Label: "Advanced - Experienced user, comfortable with system administration", Icon: "🌳"},
			{Value: "expert", Label: "Expert - Deep Linux knowledge, comfortable with complex configurations", Icon: "🌲"},
		},
	},
	{
		ID:       "useCase",
		Prompt:   "What will you primarily use this system for?",
		Type:     MultiSelect,
		Required: true,
		Options: []Option{
			{Value: "general_desktop", Label: "General Desktop - Web browsing, email, documents, media", Icon: "🖥️"},
			{Value: "development", Label: "Development - Programming, coding, web development", Icon: "💻"},
			{Value: "server", Label: "Server - Web servers, databases, containers, cloud", Icon: "🖧"},
			{Value: "security", Label: "Security - Security research, penetration testing", Icon: "🔒"},
			{Value: "gaming", Label: "Gaming - PC gaming with Steam/Proton", Icon: "🎮"},
			{Value: "content_creation", Label: "Content Creation - Video editing, graphic design, audio", Icon: "🎨"},
			{Value: "not_sure", Label: "Not Sure - Multiple use cases", Icon: "🤔"},
		},
	},
	{
		ID:       "ram",
		Prompt:   "How much RAM does your system have?",
		Type:     SingleSelect,
		Required: true,
		Options: []Option{
			{Value: "lt_2gb", Label: "Less than 2GB", Icon: "📉"},
			{Value: "2_4gb", Label: "2GB - 4GB", Icon: "📊"},
			{Value: "4_8gb", Label: "4GB - 8GB", Icon: "📈"},
			{Value: "8_16gb", Label: "8GB - 16GB", Icon: "🚀"},
			{Value: "gt_16gb", Label: "16GB or more", Icon: "⚡"},
			{Value: "not_sure", Label: "Not Sure", Icon: "❓"},
		},
	},
	{
		ID:       "disk",
		Prompt:   "How much disk space is available?",
		Type:     SingleSelect,
		Required: true,
		Options: []Option{
			{Value: "lt_20gb", Label: "Less than 20GB", Icon: "💾"},
			{Value: "20_50gb", Label: "20GB - 50GB", Icon: "💿"},
			{Value: "50_100gb", Label: "50GB - 100GB", Icon: "📀"},
			{Value: "gt_100gb", Label: "100GB or more", Icon: "💽"},
			{Value: "not_sure", Label: "Not Sure", Icon: "❓"},
		},
	},
	{
		ID:       "hardwareType",
		Prompt:   "What type of hardware are you using?",
		Type:     SingleSelect,
		Required: true,
		Options: []Option{
			{Value: "desktop", Label: "Desktop PC - Tower computer", Icon: "🖥️"},
			{Value: "laptop", Label: "Laptop - Portable computer", Icon: "💻"},
			{Value: "server", Label: "Server - Headless machine", Icon: "🖧"},
			{Value: "raspberry_pi", Label: "Raspberry Pi / ARM device", Icon: "🍓"},
			{Value: "virtual_machine", Label: "Virtual Machine / Cloud", Icon: "☁️"},
			{Value: "not_sure", Label: "Not Sure", Icon: "❓"},
		},
	},
	{
		ID:       "desktopEnvironment",
		Prompt:   "Which desktop environment do you prefer?",
		Type:     MultiSelect,
		Required: false,
		Options: []Option{
			{Value: "no_preference", Label: "No Preference - Recommend based on other criteria", Icon: "🎯"},
			{Value: "gnome", Label: "GNOME - Modern, workflow-focused", Icon: "🔷"},
			{Value: "kde", Label: "KDE Plasma - Customizable, feature-rich", Icon: "💠"},
			{Value: "xfce", Label: "XFCE - Lightweight, traditional", Icon: "🔸"},
			{Value: "mate", Label: "MATE - Classic GNOME 2 style", Icon: "🟢"},
			{Value: "cinnamon", Label: "Cinnamon - Traditional desktop", Icon: "🟤"},
			{Value: "pantheon", Label: "Pantheon - macOS-like experience", Icon: "🍎"},
			{Value: "i3", Label: "i3/Sway - Tiling window manager", Icon: "🪟"},
			{Value: "lxqt", Label: "LXQt/LXDE - Very lightweight", Icon: "🟣"},
		},
	},
	{
		ID:       "releaseModel",
		Prompt:   "How do you prefer software updates?",
		Type:     SingleSelect,
		Required: false,
		Options: []Option{
			{Value: "no_preference", Label: "No Preference - Recommend based on experience level", Icon: "🎯"},
			{Value: "stable_lts", Label: "Stable/LTS - Infrequent updates, maximum stability", Icon: "🏛️"},
			{Value: "semi_rolling", Label: "Semi-Rolling - Regular updates with testing period", Icon: "🔄"},
			{Value: "rolling", Label: "Rolling - Always latest software, cutting-edge", Icon: "⚡"},
			{Value: "fixed_release", Label: "Fixed Release - Predictable 6-12 month schedule", Icon: "📅"},
		},
	},
	{
		ID:       "packageManager",
		Prompt:   "Which package manager do you prefer?",
		Type:     SingleSelect,
		Required: false,
		Options: []Option{
			{Value: "no_preference", Label: "No Preference - Recommend based on other criteria", Icon: "🎯"},
			{Value: "apt", Label: "APT - Debian/Ubuntu family", Icon: "📦"},
			{Value: "dnf", Label: "DNF/RPM - Fedora/RHEL family", Icon: "📥"},
			{Value: "pacman", Label: "Pacman - Arch family", Icon: "🎁"},
			{Value: "portage", Label: "Portage - Gentoo", Icon: "⚙️"},
		},
	},
	{
		ID:       "supportLevel",
		Prompt:   "What level of support do you need?",
		Type:     MultiSelect,
		Required: false,
		Options: []Option{
			{Value: "no_preference", Label: "No Preference - Any level acceptable", Icon: "🎯"},
			{Value: "extensive", Label: "Extensive Community - Large forums, beginner-friendly", Icon: "👥"},
			{Value: "professional", Label: "Professional Support - Paid enterprise support available", Icon: "💼"},
			{Value: "documentation", Label: "Documentation - Comprehensive official docs", Icon: "📚"},
			{Value: "minimal", Label: "Minimal - Comfortable figuring things out", Icon: "🚶"},
		},
	},
	{
		ID:       "philosophy",
		Prompt:   "Are there any philosophical requirements?",
		Type:     SingleSelect,
		Required: false,
		Options: []Option{
			{Value: "no_preference", Label: "No Preference - Pragmatic approach is fine", Icon: "🎯"},
			{Value: "free_software", Label: "100% Free Software - FSF-approved, no proprietary code", Icon: "🆓"},
			{Value: "freedom", Label: "User Freedom - Flexible, customizable", Icon: "🦅"},
			{Value: "corporate", Label: "Corporate-Backed - Prefer commercially supported distros", Icon: "🏢"},
		},
	},
	{
		ID:       "privacyLevel",
		Prompt:   "What level of privacy/security do you need?",
		Type:     SingleSelect,
		Required: false,
		Options: []Option{
			{Value: "casual", Label: "Casual - Standard privacy, no special requirements", Icon: "🎯"},
			{Value: "enhanced", Label: "Enhanced - Private browsing, anti-tracking, encryption", Icon: "🔐"},
			{Value: "high", Label: "High - Anonymity-focused, Tor/VPN, anti-forensics", Icon: "🛡️"},
			{Value: "extreme", Label: "Extreme - Maximum isolation, compartmentalization, Qubes-level", Icon: "🔒"},
			{Value: "not_sure", Label: "Not Sure", Icon: "❓"},
		},
	},
	{
		ID:       "learningGoal",
		Prompt:   "What is your primary goal with this Linux installation?",
		Type:     SingleSelect,
		Required: false,
		Options: []Option{
			{Value: "productivity", Label: "Productivity - Just get work done, minimal maintenance", Icon: "⚡"},
			{Value: "balance", Label: "Balance - Learn while being productive", Icon: "⚖️"},
			{Value: "learning", Label: "Learning - Want to understand Linux deeply", Icon: "📚"},
			{Value: "not_sure", Label: "Not Sure", Icon: "❓"},
		},
	},
}

// FormatLabel returns a human-friendly label for a preference value.
func FormatLabel(kind, value string) string {
	switch kind {
	case "experienceLevel":
		switch value {
		case "beginner":
			return "Beginner"
		case "intermediate":
			return "Intermediate"
		case "advanced":
			return "Advanced"
		case "expert":
			return "Expert"
		}
	case "useCase":
		switch value {
		case "general_desktop":
			return "General Desktop"
		case "development":
			return "Development"
		case "server":
			return "Server"
		case "security":
			return "Security"
		case "gaming":
			return "Gaming"
		case "content_creation":
			return "Content Creation"
		case "old_hardware":
			return "Old Hardware"
		case "privacy":
			return "Privacy"
		}
	case "ram":
		switch value {
		case "lt_2gb":
			return "Less than 2GB"
		case "2_4gb":
			return "2GB - 4GB"
		case "4_8gb":
			return "4GB - 8GB"
		case "8_16gb":
			return "8GB - 16GB"
		case "gt_16gb":
			return "16GB or more"
		}
	case "disk":
		switch value {
		case "lt_20gb":
			return "Less than 20GB"
		case "20_50gb":
			return "20GB - 50GB"
		case "50_100gb":
			return "50GB - 100GB"
		case "gt_100gb":
			return "100GB or more"
		}
	case "hardwareType":
		switch value {
		case "desktop":
			return "Desktop PC"
		case "laptop":
			return "Laptop"
		case "server":
			return "Server"
		case "raspberry_pi":
			return "Raspberry Pi / ARM"
		case "virtual_machine":
			return "Virtual Machine / Cloud"
		}
	case "desktopEnvironment":
		switch value {
		case "no_preference":
			return "No Preference"
		case "gnome":
			return "GNOME"
		case "kde":
			return "KDE Plasma"
		case "xfce":
			return "XFCE"
		case "mate":
			return "MATE"
		case "cinnamon":
			return "Cinnamon"
		case "pantheon":
			return "Pantheon"
		case "i3":
			return "i3/Sway"
		case "lxqt":
			return "LXQt/LXDE"
		}
	case "releaseModel":
		switch value {
		case "no_preference":
			return "No Preference"
		case "stable_lts":
			return "Stable/LTS"
		case "semi_rolling":
			return "Semi-Rolling"
		case "rolling":
			return "Rolling Release"
		case "fixed_release":
			return "Fixed Release"
		}
	case "packageManager":
		switch value {
		case "no_preference":
			return "No Preference"
		case "apt":
			return "APT (Debian/Ubuntu)"
		case "dnf":
			return "DNF (Fedora/RHEL)"
		case "pacman":
			return "Pacman (Arch)"
		case "portage":
			return "Portage (Gentoo)"
		}
	case "supportLevel":
		switch value {
		case "no_preference":
			return "No Preference"
		case "extensive":
			return "Extensive Community"
		case "professional":
			return "Professional Support"
		case "documentation":
			return "Documentation"
		case "minimal":
			return "Minimal"
		}
	case "philosophy":
		switch value {
		case "no_preference":
			return "No Preference"
		case "free_software":
			return "100% Free Software"
		case "freedom":
			return "User Freedom"
		case "corporate":
			return "Corporate-Backed"
		}
	case "privacyLevel":
		switch value {
		case "casual":
			return "Casual"
		case "enhanced":
			return "Enhanced"
		case "high":
			return "High"
		case "extreme":
			return "Extreme"
		case "not_sure":
			return "Not Sure"
		}
	case "learningGoal":
		switch value {
		case "productivity":
			return "Productivity"
		case "balance":
			return "Balance"
		case "learning":
			return "Learning"
		case "not_sure":
			return "Not Sure"
		}
	}
	return value
}
