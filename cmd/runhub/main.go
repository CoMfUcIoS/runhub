package main

import (
	"log"
	"os/exec"

	"github.com/charmbracelet/bubbletea"
	"github.com/comfucios/runhub/pkg/config"
	"github.com/comfucios/runhub/pkg/tui"
)

var program *tea.Program

func main() {
	cfg, err := config.Load(".runhub.yaml")
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	var commands []*tui.Command
	for _, c := range cfg.Commands {
		cmd := exec.Command("sh", "-c", c.Command)
		if c.Dir != "" {
			cmd.Dir = c.Dir
		}

		co := &tui.Command{
			Name:          c.Name,
			CommandString: c.Command,
			Dir:           c.Dir,
			ExitImportant: c.ExitImportant,
		}
		go co.Run()
		commands = append(commands, co)
	}

	p := tea.NewProgram(
		&tui.Model{
			Config:   cfg,
			Commands: commands,
		},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	program = p
	if _, err := p.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
