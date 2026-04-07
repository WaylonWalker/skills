package ui

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/skills/internal/skills"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// skillItem implements list.Item for the bubbles list component.
type skillItem struct {
	skill skills.Skill
}

func (i skillItem) Title() string       { return i.skill.Name }
func (i skillItem) Description() string { return i.skill.Description }
func (i skillItem) FilterValue() string { return i.skill.Name + " " + i.skill.Description }

type pickerModel struct {
	list     list.Model
	selected *skills.Skill
	quitting bool
}

// Pick opens an interactive fuzzy picker and returns the selected skill.
// Returns nil if the user cancels.
func Pick(available []skills.Skill, title string) (*skills.Skill, error) {
	if fi, _ := os.Stdin.Stat(); fi != nil && fi.Mode()&os.ModeCharDevice == 0 {
		return nil, fmt.Errorf("interactive picker requires a terminal; specify a skill name as an argument")
	}

	items := make([]list.Item, len(available))
	for i, s := range available {
		items[i] = skillItem{skill: s}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(theme.ColorAccent).
		BorderLeftForeground(theme.ColorAccent)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(theme.ColorComment).
		BorderLeftForeground(theme.ColorAccent)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(theme.ColorFg)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(theme.ColorComment)

	l := list.New(items, delegate, 50, 16)
	l.Title = title
	l.Styles.Title = theme.Title
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(theme.ColorBlue)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(theme.ColorAccent)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(true)

	m := pickerModel{list: l}

	result, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return nil, fmt.Errorf("running picker: %w", err)
	}

	pm := result.(pickerModel)
	return pm.selected, nil
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(skillItem); ok {
				s := item.skill
				m.selected = &s
			}
			m.quitting = true
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pickerModel) View() string {
	if m.quitting {
		return ""
	}
	return m.list.View()
}
