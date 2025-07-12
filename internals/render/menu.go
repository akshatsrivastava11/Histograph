// render/menu.go
package render

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type browserItem struct {
	name string
	desc string
}

func (i browserItem) Title() string       { return i.name }
func (i browserItem) Description() string { return i.desc }
func (i browserItem) FilterValue() string { return i.name }

type model struct {
	list     list.Model
	selected string
	done     bool
	viewport viewport.Model
}

func NewModel() model {
	items := []list.Item{
		browserItem{name: "Chrome", desc: "Google Chrome history file"},
		browserItem{name: "Firefox", desc: "Mozilla Firefox history file"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 120, 40)
	l.Title = "Choose your browser"
	vp := viewport.New(40, 20)
	vp.SetContent(l.View())
	return model{list: l, viewport: vp}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item := m.list.SelectedItem().(browserItem)
			m.selected = item.name
			m.done = true
			return m, tea.Quit
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.viewport.SetContent(m.list.View())
	return m, cmd
}

func (m model) View() string {
	if m.done {
		return cardStyle.Render(fmt.Sprintf("You chose: %s", m.selected))
	}

	// wrap the list view in the card style
	return cardStyle.Render(m.list.View())
}

func GetUserBrowserChoice() (string, error) {
	prog := tea.NewProgram(NewModel())
	finalModel, err := prog.Run()
	if err != nil {
		return "", err
	}
	return finalModel.(model).selected, nil
}
