# Commons

This folder holds shared resources used by both the web interface and the terminal UI.

## Contents

- `data/distros.json` — the single source of truth for the Linux distribution database. The web app and the TUI both read this file, so edits here apply to both interfaces.

## Notes

- Algorithm and UI question definitions are intentionally kept per-interface (JavaScript for the web, Go for the TUI). This keeps each codebase idiomatic and avoids cross-language build complexity.
- When adding a new distribution, update only `data/distros.json`; both interfaces will see it on the next build.
