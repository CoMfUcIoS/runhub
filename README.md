# RunHub 🚀

**Centralize, visualize, and orchestrate your CLI workflows in one sleek dashboard.**

![RunHub Demo](https://via.placeholder.com/800x400.png?text=RunHub+TUI+Demo) _← Add actual demo GIF later_

---

## Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Configuration](#-configuration)
- [TUI Controls](#-tui-controls)
- [Building from Source](#-building-from-source)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🌟 Features

**Terminal Power-Ups:**

- 🛠 **Parallel Process Orchestration** - Run services, scripts, and tools simultaneously
- 📊 **Live Dashboard** - Unified view of all command outputs with scrollback history
- ⚡ **Zero-Config Reruns** - Restart any command with single keypress
- 🚨 **Smart Failure Handling** - Critical command failure detection and reporting
- ⌨ **Interactive Mode** - Direct input to running processes

**Developer Experience:**

- 🔧 **YAML Configuration** - Simple declarative setup
- 📁 **Directory Isolation** - Per-command working directories
- 🎨 **Color-Coded Status** - Instant visual feedback (running/success/failed)
- 📶 **Output Buffering** - Maintains last 100 lines per command
- 🖥 **Responsive Layout** - Adapts to terminal size

---

## 📦 Installation

### Homebrew (macOS/Linux)

```bash
brew install comfucios/tap/runhub
```

### Pre-built Binaries

1. Visit [Releases Page](https://github.com/comfucios/runhub/releases)
2. Download for your OS:
   - **Windows**: `runhub-windows-amd64.exe`
   - **Linux**: `runhub-linux-amd64`
   - **macOS**: `runhub-darwin-amd64`
3. Make executable (Unix):

```bash
chmod +x runhub-*
```

### From Source

Prerequisites: Go 1.20+

```bash
git clone https://github.com/comfucios/runhub
cd runhub
go build -o runhub ./cmd/runhub
sudo mv runhub /usr/local/bin/
```

---

## 🚀 Quick Start

1. Create `runhub.yaml` in your project root:

```yaml
# Minimal configuration
commands:
  - name: "Web Server"
    command: "python3 -m http.server"
```
