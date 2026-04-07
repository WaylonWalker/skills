package ui

import (
	"fmt"
	"strings"

	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputModel struct {
	textInput textinput.Model
	title     string
	value     string
	submitted bool
	quitting  bool
	err       error
}

// Input opens an interactive text input and returns the entered value.
// Returns empty string if the user cancels (esc/ctrl+c).
func Input(title string, placeholder string, validate func(string) error) (string, error) {
	if !IsInteractiveTerminal() {
		return "", fmt.Errorf("interactive input requires a terminal; specify a value as an argument")
	}

	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(theme.ColorAccent)
	ti.TextStyle = lipgloss.NewStyle().Foreground(theme.ColorFg)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.ColorComment)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(theme.ColorAccent)
	if validate != nil {
		ti.Validate = validate
	}

	m := inputModel{
		textInput: ti,
		title:     title,
	}

	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", fmt.Errorf("running input: %w", err)
	}

	im := result.(inputModel)
	if !im.submitted {
		return "", nil
	}
	return im.value, nil
}

func (m inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			val := strings.TrimSpace(m.textInput.Value())
			if val == "" {
				return m, nil
			}
			m.value = val
			m.submitted = true
			m.quitting = true
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	if m.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Foreground(theme.ColorAccent).Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(theme.ColorComment)

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("enter to confirm / esc to cancel"))
	b.WriteString("\n")

	return b.String()
}
