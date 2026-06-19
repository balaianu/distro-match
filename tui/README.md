# DistroMatch TUI

A terminal user interface for [DistroMatch](https://github.com/balaianu/distro-match), built with Go and the [Charm.land](https://charm.land) ecosystem.

## Quick Start

```bash
make run
```

## Build

```bash
make build    # produces ./distro-match
./distro-match
```

## Batch Mode

Non-interactive mode for scripting and automation:

```bash
./distro-match --batch 3,2,3,3,1,2,3,4,1,2,3,1
```

The 12 comma-separated values are 1-based option indices for each wizard question. Results print to stdout.

## Controls

| Screen | Keys |
|--------|------|
| Welcome | `1-4` or `↑/↓`+`Enter` to select, `q` to quit |
| Wizard (single-select) | `1-9` quick pick, `↑/↓` move, `Enter`/`Space` select, `n` next, `b` back, `s` skip |
| Wizard (multi-select) | `1-9` quick toggle, `↑/↓` move, `Space` toggle, `Enter` confirm, `n` next, `b` back, `s` skip |
| Results | `↑/↓` browse, `Enter` view details, `r` restart, `b` back, `q` quit |
| Detail | `↑/↓` scroll, `b`/`esc`/`q` back |
| Explorer | `/` start filter, type to search, `Enter` details, `esc` clear/back, `q` back |

## Architecture

| File | Responsibility |
|------|---------------|
| `main.go` | Entry point, batch mode handler |
| `model.go` | Bubble Tea model, state machine, input handling |
| `views.go` | Screen rendering for all 5 screens |
| `styles.go` | Lip Gloss color palette and reusable styled components |
| `questions.go` | Wizard question definitions and preference types |
| `scoring.go` | Recommendation scoring algorithm |
| `data.go` | Distro data loading from embedded JSON |
| `model_test.go` | Unit tests |

## Data

The distro database lives at `../commons/data/distros.json` and is copied into the build directory by the Makefile, then embedded at compile time via `//go:embed`.

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 — Elm Architecture TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) v2 — Reusable components (list, viewport)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2 — Style/format library
