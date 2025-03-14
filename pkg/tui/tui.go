package tui

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"syscall"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/comfucios/runhub/pkg/config"
)

var Program *tea.Program

type Command struct {
	Name          string
	Cmd           *exec.Cmd
	CommandString string
	Dir           string
	Output        []string
	Stdin         io.WriteCloser
	Mu            sync.Mutex
	Running       bool
	Finished      bool
	ExitCode      int
	ExitImportant bool
	ScrollOffset  int
}

type Model struct {
	Config        *config.Config
	Commands      []*Command
	sortedIndices []int
	selected      int
	width         int
	height        int
	contentHeight int
	ready         bool
	interactive   bool
	shouldQuit    bool
	finalOutput   string
}

type commandFinishedMsg struct {
	name   string
	output []string
	cmd    *Command
}

var (
	sidebarStyle = lipgloss.NewStyle().
			Width(20).
			MaxWidth(20).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#4A90E2")).
				MarginLeft(1).
				Padding(0, 1).
				MaxWidth(18).
				Inline(true)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D3D3D3")).
				MarginLeft(2).
				MaxWidth(18).
				Inline(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4A90E2"))

	finalOutputStyle = lipgloss.NewStyle().
				Padding(1, 2).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF5555")).
				Width(120)
)

func (c *Command) Run() {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	// Create new command instance for reruns
	c.Cmd = exec.Command("sh", "-c", c.CommandString)
	if c.Dir != "" {
		c.Cmd.Dir = c.Dir
	}

	stdout, _ := c.Cmd.StdoutPipe()
	stdin, err := c.Cmd.StdinPipe()
	if err != nil {
		c.Output = append(c.Output, "ERROR: "+err.Error())
		return
	}
	c.Stdin = stdin

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			c.Mu.Lock()
			c.Output = append(c.Output, scanner.Text())
			if len(c.Output) > 100 {
				c.Output = c.Output[1:]
			}
			c.Mu.Unlock()
		}
	}()

	c.Running = true
	c.Finished = false
	if err := c.Cmd.Start(); err != nil {
		c.Output = append(c.Output, "ERROR: "+err.Error())
	}

	go func() {
		err := c.Cmd.Wait()
		c.Mu.Lock()
		defer c.Mu.Unlock()
		c.Running = false
		c.Finished = true

		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				c.ExitCode = status.ExitStatus()
			} else {
				c.ExitCode = 1
			}
		} else if err != nil {
			c.ExitCode = 1
		}

		if Program != nil {
			Program.Send(commandFinishedMsg{
				name:   c.Name,
				output: c.Output,
				cmd:    c,
			})
		}
	}()
}

func (m *Model) Init() tea.Cmd {
	m.updateSortedIndices()
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.updateSortedIndices()

	switch msg := msg.(type) {
	case commandFinishedMsg:
		shouldExit := false
		if msg.cmd.ExitCode != 0 {
			shouldExit = m.Config.ExitOnCompletion
		} else {
			for _, cmd := range m.Commands {
				if cmd.Finished && cmd.ExitCode == 0 && cmd.ExitImportant {
					shouldExit = true
					break
				}
			}
		}

		if shouldExit {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("🚨 %s exited with code %d 🚨\n\n",
				msg.name, msg.cmd.ExitCode))
			sb.WriteString(strings.Join(msg.output, "\n"))

			for _, cmd := range m.Commands {
				if cmd.Name != msg.name {
					cmd.Mu.Lock()
					sb.WriteString(fmt.Sprintf("\n\n=== %s ===\n", cmd.Name))
					sb.WriteString(strings.Join(cmd.Output, "\n"))
					cmd.Mu.Unlock()
				}
			}

			m.finalOutput = sb.String()
			return m, tea.Quit
		}

	case tea.KeyMsg:
		if m.interactive {
			switch msg.Type {
			case tea.KeyUp:
				if m.selected >= 0 && m.selected < len(m.sortedIndices) {
					cmdIdx := m.sortedIndices[m.selected]
					cmd := m.Commands[cmdIdx]
					cmd.Mu.Lock()
					maxScroll := len(cmd.Output) - m.contentHeight
					if maxScroll < 0 {
						maxScroll = 0
					}
					if cmd.ScrollOffset < maxScroll {
						cmd.ScrollOffset++
					}
					cmd.Mu.Unlock()
				}
			case tea.KeyDown:
				if m.selected >= 0 && m.selected < len(m.sortedIndices) {
					cmdIdx := m.sortedIndices[m.selected]
					cmd := m.Commands[cmdIdx]
					cmd.Mu.Lock()
					if cmd.ScrollOffset > 0 {
						cmd.ScrollOffset--
					}
					cmd.Mu.Unlock()
				}
			case tea.KeyCtrlZ:
				m.interactive = false
			default:
				m.handleInteractiveInput(msg)
			}
		} else {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up":
				if m.selected > 0 {
					m.selected--
				}
			case "down":
				if m.selected < len(m.sortedIndices)-1 {
					m.selected++
				}
			case "i":
				if len(m.sortedIndices) > 0 && m.selected < len(m.sortedIndices) {
					cmdIdx := m.sortedIndices[m.selected]
					if m.Commands[cmdIdx].Running {
						m.interactive = true
					}
				}
			case "r": // New rerun functionality
				if len(m.sortedIndices) > 0 && m.selected < len(m.sortedIndices) {
					cmdIdx := m.sortedIndices[m.selected]
					cmd := m.Commands[cmdIdx]

					cmd.Mu.Lock()
					if cmd.Finished {
						// Reset command state
						cmd.Output = nil
						cmd.ScrollOffset = 0
						cmd.Finished = false
						cmd.ExitCode = 0
						cmd.Mu.Unlock()

						// Start new instance
						go cmd.Run()
					} else {
						cmd.Mu.Unlock()
					}
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.contentHeight = m.height - 3
		m.ready = true
	}

	return m, nil
}

func (m *Model) updateSortedIndices() {
	indices := make([]int, len(m.Commands))
	for i := range m.Commands {
		indices[i] = i
	}

	sort.SliceStable(indices, func(i, j int) bool {
		a, b := m.Commands[indices[i]], m.Commands[indices[j]]

		// Running commands first
		if a.Running != b.Running {
			return a.Running
		}

		// Then failed commands
		if a.Finished && b.Finished {
			if a.ExitCode != 0 && b.ExitCode == 0 {
				return true
			}
		}

		// Finally alphabetical order
		return a.Name < b.Name
	})

	// Preserve selection
	if len(m.sortedIndices) > 0 && m.selected < len(m.sortedIndices) {
		currentCmd := m.sortedIndices[m.selected]
		for newIdx, cmdIdx := range indices {
			if cmdIdx == currentCmd {
				m.selected = newIdx
				break
			}
		}
	}

	m.sortedIndices = indices
}

func (m *Model) handleInteractiveInput(msg tea.KeyMsg) {
	switch msg.Type {
	case tea.KeyEnter:
		m.writeToStdin("\n")
	case tea.KeySpace:
		m.writeToStdin(" ")
	case tea.KeyEscape:
		m.writeToStdin("\x1b")
	case tea.KeyRunes:
		m.writeToStdin(string(msg.Runes))
	case tea.KeyBackspace:
		m.writeToStdin("\b")
	case tea.KeyTab:
		m.writeToStdin("\t")
	}
}

func (m *Model) writeToStdin(text string) {
	if len(m.sortedIndices) == 0 || m.selected >= len(m.sortedIndices) {
		return
	}

	cmdIdx := m.sortedIndices[m.selected]
	cmd := m.Commands[cmdIdx]
	cmd.Mu.Lock()
	defer cmd.Mu.Unlock()

	if cmd.Stdin != nil && cmd.Running {
		cmd.Stdin.Write([]byte(text))
	}
}

func (m *Model) View() string {
	if m.shouldQuit {
		return finalOutputStyle.Render(m.finalOutput)
	}

	if !m.ready {
		return "Initializing..."
	}

	sidebarWidth := 20
	contentWidth := m.width - sidebarWidth - 4

	// Build sidebar items
	var sidebarItems []string
	for idx, cmdIdx := range m.sortedIndices {
		cmd := m.Commands[cmdIdx]
		var status string
		switch {
		case cmd.Running:
			status = runningStyle.Render("▶")
		case cmd.Finished && cmd.ExitCode == 0:
			status = successStyle.Render("✓")
		case cmd.Finished && cmd.ExitCode != 0:
			status = errorStyle.Render("✗")
		default:
			status = "?"
		}

		// Truncate name to 16 characters with ellipsis
		name := cmd.Name
		if len(name) > 16 {
			name = name[:16] + "..."
		}
		line := fmt.Sprintf("%s %s", status, name)

		// Apply selection style
		if idx == m.selected {
			line = selectedItemStyle.Render(line)
		} else {
			line = unselectedItemStyle.Render(line)
		}

		sidebarItems = append(sidebarItems, line)
	}

	// Build content
	var content strings.Builder
	if len(m.sortedIndices) > 0 && m.selected < len(m.sortedIndices) {
		cmdIdx := m.sortedIndices[m.selected]
		cmd := m.Commands[cmdIdx]
		cmd.Mu.Lock()
		defer cmd.Mu.Unlock()

		lines := cmd.Output
		var displayLines []string
		start := len(lines) - m.contentHeight - cmd.ScrollOffset
		if start < 0 {
			start = 0
		}
		end := len(lines) - cmd.ScrollOffset
		if end > len(lines) {
			end = len(lines)
		}
		if start < end {
			displayLines = lines[start:end]
		}

		content.WriteString(strings.Join(displayLines, "\n"))
	}

	// Compose layout with sidebar on the right
	sidebar := sidebarStyle.Render(strings.Join(sidebarItems, "\n"))

	contentPanel := lipgloss.NewStyle().
		Width(contentWidth).
		Height(m.contentHeight).
		Padding(1, 2).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		Render(content.String())

	statusBar := lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2E2E2E")).
		Padding(0, 1).
		Render(m.getStatusBar())

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, contentPanel, sidebar),
		statusBar,
	)
}

func (m *Model) getStatusBar() string {
	if m.interactive {
		return "↑/↓: Scroll logs | Ctrl+Z: Exit interactive mode"
	}
	return "i: Interact | ↑/↓: Select command | r: Rerun | q: Quit"
}
