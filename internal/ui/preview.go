package ui

import (
	"fmt"
	"strings"

	"github.com/WaylonWalker/skills/internal/skills"
	"github.com/WaylonWalker/skills/internal/theme"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type previewModel struct {
	list     list.Model
	viewport viewport.Model
	skills   []skills.Skill
	ready    bool
	width    int
	height   int
	lastIdx  int
	quitting bool
}

// Preview opens an interactive split-pane browser with a skill list on the
// left and a rendered markdown preview on the right.
func Preview(available []skills.Skill, global bool, installed map[string]bool) error {
	if !IsInteractiveTerminal() {
		// Non-interactive: print a simple list instead.
		for _, s := range available {
			fmt.Printf("%-30s %s\n", s.Name, skillItem{skill: s, installed: installed[s.Name]}.Description())
		}
		return nil
	}

	items := make([]list.Item, len(available))
	for i, s := range available {
		items[i] = skillItem{skill: s, installed: installed[s.Name]}
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

	l := list.New(items, delegate, 30, 20)
	scope := "project"
	if global {
		scope = "global"
	}
	l.Title = fmt.Sprintf("Skills (%s)", scope)
	l.Styles.Title = theme.Title
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(theme.ColorBlue)
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(theme.ColorAccent)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	m := previewModel{
		list:    l,
		skills:  available,
		lastIdx: -1,
	}

	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func (m previewModel) Init() tea.Cmd {
	return nil
}

func (m previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		listWidth := m.width / 3
		if listWidth < 25 {
			listWidth = 25
		}
		if listWidth > 45 {
			listWidth = 45
		}
		previewWidth := m.width - listWidth - 4

		m.list.SetSize(listWidth, m.height)

		if !m.ready {
			m.viewport = viewport.New(previewWidth, m.height-4)
			m.ready = true
		} else {
			m.viewport.Width = previewWidth
			m.viewport.Height = m.height - 4
		}

		m.lastIdx = -1 // force preview update
	}

	prevIdx := m.list.Index()
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	// Update preview when selection changes.
	if m.list.Index() != m.lastIdx {
		m.lastIdx = m.list.Index()
		if item, ok := m.list.SelectedItem().(skillItem); ok {
			rendered := renderMarkdown(item.skill.Content, m.viewport.Width-2)
			m.viewport.SetContent(rendered)
			m.viewport.GotoTop()
		}
	}
	_ = prevIdx

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m previewModel) View() string {
	if m.quitting {
		return ""
	}
	if !m.ready {
		return theme.Subtle.Render("Loading...")
	}

	listView := m.list.View()

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.ColorComment).
		PaddingLeft(1)

	header := theme.Title.Render("Preview")
	separator := theme.Subtle.Render(strings.Repeat("─", m.viewport.Width))
	previewView := borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			header,
			separator,
			m.viewport.View(),
		),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, listView, previewView)
}

func renderMarkdown(content string, width int) string {
	if width < 20 {
		width = 20
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithWordWrap(width),
		glamour.WithStandardStyle("dark"),
	)
	if err != nil {
		return content
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return content
	}

	return rendered
}
