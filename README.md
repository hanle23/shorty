# shorty

A CLI tool to manage and run shortcuts and scripts.

Shorty lets you define aliases for commands (shortcuts) and reusable shell snippets (scripts), then run them by name. When a name matches both a shortcut and a script, the shortcut takes precedence — use `shorty script <name>` to run the script explicitly.

## Supported Platforms

- **macOS** (darwin)
- **Linux**

Windows is not currently supported. Scripts execute via `sh -c`, which is not available on Windows by default.

## Installation

```sh
go install github.com/hanle23/shorty@latest
```

After installing, initialize your configuration:

```sh
shorty init
```

This creates `~/.config/shorty/` with a default `config.yaml` and an empty `runnables.yaml`.

## Usage

```sh
shorty <name> [args...]              # Run a shortcut or script by name
shorty shortcut <name> [args...]     # Run a shortcut explicitly (alias: sc)
shorty script <name> [args...]       # Run a script explicitly (alias: sr)
```

### Commands

| Command | Description |
|---------|-------------|
| `shorty init` | Initialize config or reset to original state |
| `shorty add` | Add a new shortcut or script interactively |
| `shorty list` | List all configured shortcuts and scripts |
| `shorty doctor` | Diagnose configuration health |

### Flags

| Flag | Description |
|------|-------------|
| `--config` | Config file (default: `~/.config/shorty/config.yaml`) |
| `-d, --debug` | Enable debug output |

## Configuration

Shorty stores its configuration in `~/.config/shorty/` (or `$XDG_CONFIG_HOME/.config/shorty/` if set).

- **`config.yaml`** — points to the runnables file location
- **`runnables.yaml`** — defines your shortcuts and scripts

### Example `runnables.yaml`

```yaml
shortcuts:
  k:
    shortcut_name: k
    package_name: kubectl
    args:
      - get
      - pods
    description: List Kubernetes pods
scripts:
  deploy:
    package_name: deploy
    script: git push origin main && echo "Deployed!"
    description: Push to main and confirm
```

## License

GPL-3.0
