# labcheck

A CLI tool for checking the health of your self-hosted homelab over SSH.

## Features

- Docker container status
- Disk usage (mergerFS)
- SnapRAID sync status

## Installation

```bash
go install github.com/rjdailey/labcheck@latest
```

## Usage

```bash
labcheck init        # configure your homelab connection
labcheck run         # run all health checks
```

## Configuration

Config is stored at `~/.labcheck/config.yaml`.

## Planned Improvements

- [ ] Tab completion for file paths in `labcheck init` (huh)
- [ ] Interactive TUI dashboard (Bubbletea)
- [ ] Multi-host support
