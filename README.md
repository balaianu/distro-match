# DistroMatch

<div align="center">

**Find your perfect Linux distribution**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Astro](https://img.shields.io/badge/Astro-FF5D01?style=flat&logo=astro&logoColor=white)](https://astro.build)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind-06B6D4?style=flat&logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![Bubble Tea](https://img.shields.io/badge/Bubble_Tea-FF69B4?style=flat)](https://github.com/charmbracelet/bubbletea)

[Live Demo](https://distro-match.com) · [Report Bug](https://github.com/balaianu/distro-match/issues) · [Request Feature](https://github.com/balaianu/distro-match/issues)

</div>

A guided wizard that helps users find the perfect Linux distribution based on their experience level, use case, hardware specifications, and personal preferences. Available as both a web application and a terminal UI.

## Repository Structure

```
distro-match/
├── commons/                # Shared data and assets
│   ├── data/
│   │   └── distros.json    # Linux distribution database (shared by both apps)
│   └── README.md
├── web/                    # Web application (Astro + Tailwind CSS)
│   ├── src/
│   │   ├── data/           # (symlink or copy of commons data at build time)
│   │   ├── layouts/        # Page layouts
│   │   ├── pages/          # Routes (index, distros, about)
│   │   ├── scripts/        # Wizard state + scoring algorithm
│   │   └── styles/         # Global CSS
│   ├── public/             # Static assets
│   ├── astro.config.mjs
│   ├── tailwind.config.mjs
│   └── package.json
├── tui/                    # Terminal UI (Go + Charm.land ecosystem)
│   ├── main.go             # Entry point (local, batch, SSH server modes)
│   ├── server.go           # Wish SSH server for `ssh distro-match.com`
│   ├── model.go            # Bubble Tea model and state machine
│   ├── views.go            # Screen rendering (welcome, wizard, results, detail, explorer)
│   ├── keybindings.go      # Structured key bindings (bubbles/key + bubbles/help)
│   ├── styles.go           # Lip Gloss color palette and styled components
│   ├── questions.go        # Wizard questions and preference types
│   ├── scoring.go          # Recommendation scoring algorithm
│   ├── data.go             # Distro data loading (embedded JSON)
│   ├── model_test.go       # Unit tests
│   ├── Makefile            # Build/run/serve/test/cross-compile targets
│   ├── distromatch.service # systemd service file for SSH deployment
│   ├── go.mod
│   └── go.sum
├── .github/
│   └── workflows/
│       └── release.yml     # Auto-builds binaries on git tag push
├── install.sh              # One-liner installer (Linux/macOS)
├── install.ps1             # One-liner installer (Windows PowerShell)
├── setup-server.sh         # SSH server setup script (for VM deployment)
└── README.md
```

## Features

- **Smart Recommendations** - Weighted scoring algorithm matches users to suitable Linux distributions
- **Guided Wizard** - Step-by-step questionnaire with 12 decision points
- **Hardware-Aware** - Filters distros based on RAM, CPU architecture, and disk space
- **Distro Database** - Comprehensive database of 60+ Linux distributions
- **Two Interfaces**:
  - **Web** - Modern, responsive UI built with Astro and Tailwind CSS
  - **TUI** - Full terminal UI built with Go and the Charm.land ecosystem (Bubble Tea, Bubbles, Lip Gloss)
- **Open Source** - Fully MIT licensed

## Web Application

### Prerequisites

- Node.js 18.x or higher
- npm or yarn

### Quick Start

```bash
cd web
npm install
npm run dev
```

Open [http://localhost:4321](http://localhost:4321) in your browser.

### Building for Production

```bash
cd web
npm run build
```

The built files will be in the `web/dist/` directory.

## Terminal UI (TUI)

### Install (one-liner)

**Linux / macOS:**

```bash
curl -fsSL https://github.com/balaianu/distro-match/raw/master/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://github.com/balaianu/distro-match/raw/master/install.ps1 | iex
```

The install script auto-detects your OS and architecture, downloads the correct pre-built binary, and installs it to `~/.local/bin` (or `$LOCALAPPDATA\Programs\distromatch` on Windows). Flags:

```bash
# Install a specific version
curl -fsSL https://github.com/balaianu/distro-match/raw/master/install.sh | sh -s -- --version v0.1.0

# Install to a custom directory
curl -fsSL https://github.com/balaianu/distro-match/raw/master/install.sh | sh -s -- --dir /usr/local/bin

# Preview without downloading
curl -fsSL https://github.com/balaianu/distro-match/raw/master/install.sh | sh -s -- --dry-run
```

### Build from Source

#### Prerequisites

- Go 1.25 or higher

#### Quick Start

```bash
cd tui
make run
```

#### Building

```bash
cd tui
make build    # produces ./distro-match binary
./distro-match
```

#### Cross-compiling

```bash
cd tui
make dist     # builds binaries for linux/darwin/windows × amd64/arm64 into dist/
```

### Batch Mode

The TUI also supports a non-interactive batch mode for scripting:

```bash
./distro-match --batch 3,2,3,3,1,2,3,4,1,2,3,1
```

The comma-separated values are 1-based option indices for each of the 12 wizard questions. Results are printed to stdout as plain text.

### TUI Controls

| Screen | Keys |
|--------|------|
| Welcome | `1-4` select, `Enter` confirm, `q` quit |
| Wizard | `1-9` quick pick, `↑/↓` move, `Space` toggle (multi), `Enter` select/confirm, `n` next, `b` back, `s` skip, `q` quit |
| Results | `Enter` view details, `r` restart, `b` back, `q` quit |
| Detail | `↑/↓` scroll, `b`/`esc`/`q` back |
| Explorer | `/` filter, `Enter` details, `esc` clear filter/back, `q` back |

### TUI Tech Stack

- **Language**: [Go](https://go.dev) 1.25
- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 (Elm Architecture)
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) v2 (viewport, textinput, progress, key, help)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2 (including table and gradient)
- **SSH Server**: [Wish](https://github.com/charmbracelet/wish) v2 (custom SSH server for `ssh distro-match.com`)

### SSH Server

The TUI can be served over SSH, so users can run `ssh distro-match.com` and get the wizard directly in their terminal — no install required. This uses [Wish](https://github.com/charmbracelet/wish) (the same library terminal.shop uses), not OpenSSH. Any username is accepted, no password needed, and the TUI launches immediately.

#### Local testing

```bash
cd tui
make serve                              # starts on :2323
ssh -p 2323 anyone@localhost            # connect from another terminal
```

#### Production deployment (Oracle Cloud Free Tier)

The ARM Ampere A1 always-free tier (2 OCPU, 12 GB RAM) is ideal. Set up a fresh Debian/Ubuntu VM and run:

```bash
curl -fsSL https://github.com/balaianu/distro-match/raw/master/setup-server.sh | sh
```

This will:
1. Create a dedicated `distromatch` system user
2. Download the `linux-arm64` binary
3. Generate an SSH host key
4. Install a hardened systemd service on port 22
5. Open the firewall

Then point your DNS A record for `distro-match.com` to the VM's IP, and users can connect with:

```bash
ssh distro-match.com
```

The systemd service runs with security hardening (no new privileges, protected system, private tmp). See `tui/distromatch.service` for the service definition.

## Adding Linux Distributions

To add a new Linux distribution to the database:

1. Edit `commons/data/distros.json`
2. Add a new object following the existing schema
3. Rebuild the web app (`cd web && npm run build`) or the TUI (`cd tui && make build`)

## Deployment

### Web Application

#### Cloudflare Pages (Workers)

The site uses Cloudflare Workers via the `@astrojs/cloudflare` adapter. The `web/wrangler.jsonc` config is already set up for this.

**Cloudflare Pages dashboard settings:**

| Setting | Value |
|---------|-------|
| Framework preset | None |
| Build command | `cd web && npm install && npm run build` |
| Build output directory | `web/dist` |
| Root directory | `/` (repo root) |

The `prebuild` npm script automatically copies `commons/data/distros.json` into `web/src/data/` before each build, so no symlink is needed in the CI environment.

**Local deploy with wrangler:**

```bash
cd web
npm install
npm run deploy    # builds and deploys via wrangler
```

#### Other platforms

- **Netlify** - Build command: `cd web && npm install && npm run build`, Publish directory: `web/dist`
- **Vercel** - Build command: `cd web && npm install && npm run build`, Output directory: `web/dist`

### TUI

The TUI is a single binary with no runtime dependencies. Pre-built binaries are available for Linux, macOS, and Windows (amd64 and arm64) via the [install script](#install-one-liner) or [GitHub Releases](https://github.com/balaianu/distro-match/releases).

To create a new release, push a tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The GitHub Actions workflow will automatically build all platform binaries and create a GitHub Release with the install scripts attached.

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](docs/CONTRIBUTING.md) for details on how to contribute.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Linux distributions data sourced from official project websites
- Web app built with [Astro](https://astro.build) and [Tailwind CSS](https://tailwindcss.com)
- TUI built with the [Charm.land](https://charm.land) ecosystem
- Design inspiration from DistroWatch and various Linux comparison sites

## Support

- **Issues**: [Report bugs or request features](https://github.com/balaianu/distro-match/issues)
- **Discussions**: [Ask questions or share ideas](https://github.com/balaianu/distro-match/discussions)

---

<div align="center">

Made with ❤️ by [balaianu](https://github.com/balaianu)

</div>
