# RunHub 🚀

**Mission Control for Your Terminal Workflows**

![RunHub Demo](https://via.placeholder.com/800x400.png?text=RunHub+TUI+Demo+with+Real-Time+Updates)  
_Real-time command output and status monitoring_

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

---

## 🌟 Features

**Core Capabilities**

- ⚡ **Real-Time Updates** - Instant refresh of logs and statuses
- 🚀 **Parallel Execution** - Run commands simultaneously
- 🔄 **One-Click Reruns** - Restart failed/completed commands (press `r`)
- 🚨 **Smart Failure Handling** - Critical command monitoring
- ⌨ **Interactive Input** - Send input to running processes

**Visual Interface**

- 📊 **Live Dashboard** - Unified view of all outputs
- 🎨 **Color-Coded Status** - Instant visual feedback
- 📶 **Output Buffering** - Last 100 lines preserved
- 🔍 **Smart Scrolling** - Auto-track or manual scroll
- 📱 **Responsive Design** - Adapts to terminal size

---

## 📦 Installation

### Homebrew (macOS/Linux)

```bash
brew install comfucios/tap/runhub
```

### Pre-built Binaries

1. Download from [Releases](https://github.com/comfucios/runhub/releases)
2. Unzip and make executable:

```bash
tar -xvzf runhub-*.tar.gz
chmod +x runhub
sudo mv runhub /usr/local/bin/
```

### Linux (Debian/RPM)

```bash
# Debian/Ubuntu
curl -LO https://github.com/comfucios/runhub/releases/latest/download/runhub.deb
sudo dpkg -i runhub.deb

# RHEL/Fedora
curl -LO https://github.com/comfucios/runhub/releases/latest/download/runhub.rpm
sudo rpm -i runhub.rpm
```

---

## 🚀 Quick Start

1. Create configuration file:

```bash
cat > .runhub.yaml <<EOL
commands:
  - name: Web Server
    command: python3 -m http.server 8000
    exit_important: true

  - name: CSS Builder
    command: npm run build:css
    dir: "./frontend"
EOL
```

2. Start RunHub:

```bash
runhub
```

3. Monitor outputs in real-time:

- Use `↑/↓` to navigate commands
- Press `i` to interact with running processes
- Press `r` to restart finished commands

---

## 🔧 Configuration

### Full Schema

```yaml
exit_on_completion: false # Global exit when any command completes

commands:
  - name: "Service Name" # Display name (required)
    command: "cmd args" # Shell command (required)
    dir: "./path" # Working directory (optional)
    exit_important: true # Halt all on failure (optional)
```

### Example Config

```yaml
commands:
  - name: "Database"
    command: "docker-compose up postgres"
    dir: "./infra"

  - name: "API Server"
    command: "go run main.go"
    dir: "./api"
    exit_important: true

  - name: "Frontend"
    command: "npm run dev"
    dir: "./webapp"
```

---

## 🕹 TUI Controls

| Key        | Action                 | Context          |
| ---------- | ---------------------- | ---------------- |
| `↑/↓`      | Select command         | Normal mode      |
| `←/→`      | Scroll horizontally    | Interactive mode |
| `i`        | Enter interactive mode | Normal mode      |
| `r`        | Rerun command          | Finished command |
| `Ctrl+Z`   | Exit interactive mode  | Interactive mode |
| `q/Ctrl+C` | Quit RunHub            | Anywhere         |

**Interactive Mode:**

- Directly send input to selected process
- Arrow keys scroll output history
- Enter sends newline
- Supports most terminal input sequences

---

## 🔨 Building from Source

### Requirements

- Go 1.20+
- GCC compiler
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/comfucios/runhub
cd runhub

# Build optimized binary
go build -ldflags "-s -w" -o runhub ./cmd/runhub

# Install system-wide
sudo install runhub /usr/local/bin
```

### Development Build

```bash
go build -tags dev -o runhub ./cmd/runhub
```

### Cross-Compilation

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o runhub-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o runhub.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o runhub-macos
```

---

## 🚨 Troubleshooting

### Common Issues

**Command Not Found**

```bash
# Verify installation path
which runhub

# Check $PATH environment
echo $PATH
```

**Missing Output**

- Ensure commands produce stdout/stderr
- Test command outside RunHub first
- Check working directory permissions

**Interactive Mode Limitations**

- Some programs require pseudo-TTY (use `script` or `unbuffer`)
- Java applications may need `-Djava.ioprofile=true`

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch:

```bash
git checkout -b feat/amazing-feature
```

3. Commit changes:

```bash
git commit -m 'feat: add amazing feature'
```

4. Push to branch:

```bash
git push origin feat/amazing-feature
```

5. Open Pull Request

---

**RunHub**: _Where distributed development finally feels unified._ 🌐  
_Created with ❤️ by [Ioannis Karasavvaidis](https://github.com/comfucios)_
