// main.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
	"github.com/akshatsrivastava11/Histograph/internals/render"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VisitEntry represents a browser history entry
type VisitEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	VisitTime  time.Time `json:"visit_time"`
}

// Styles
var (
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			MarginBottom(1)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("228")).
			Bold(true)
)

// Model for handling browser history processing
type historyModel struct {
	browserChoice string
	viewport      viewport.Model
	content       string
	done          bool
	err           error
}

func newHistoryModel(browserChoice string) historyModel {
	vp := viewport.New(80, 20)

	content := titleStyle.Render("🔍 Fetching " + browserChoice + " history...")

	return historyModel{
		browserChoice: browserChoice,
		viewport:      vp,
		content:       content,
		done:          false,
	}
}

func (m historyModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		processHistoryCmd(m.browserChoice),
	)
}

// Command to process browser history
func processHistoryCmd(browserChoice string) tea.Cmd {
	return func() tea.Msg {
		switch browserChoice {
		case "Firefox":
			return processFirefoxHistory()
		case "Chrome":
			return processChromeHistory()
		default:
			return historyResult{err: fmt.Errorf("invalid browser selection")}
		}
	}
}

type historyResult struct {
	entries []render.VisitEntry
	count   int
	err     error
}

func processChromeHistory() historyResult {
	historyData := parse.ParseChromeHistory()

	if len(historyData) == 0 {
		return historyResult{
			err: fmt.Errorf("no Chrome history found or unable to access Chrome history"),
		}
	}

	var renderEntries []render.VisitEntry
	for _, entry := range historyData {
		renderEntries = append(renderEntries, render.VisitEntry{
			URL:        entry.URL,
			Title:      entry.Title,
			VisitCount: entry.VisitCount,
			VisitTime:  entry.VisitTime,
		})
	}

	return historyResult{
		entries: renderEntries,
		count:   len(historyData),
	}
}

func processFirefoxHistory() historyResult {
	historyData := parse.ParseFirefoxHistory()

	if len(historyData) == 0 {
		return historyResult{
			err: fmt.Errorf("no Firefox history found or unable to access Firefox history"),
		}
	}

	var renderEntries []render.VisitEntry
	for _, entry := range historyData {
		renderEntries = append(renderEntries, render.VisitEntry{
			URL:        entry.URL,
			Title:      entry.Title,
			VisitCount: entry.VisitCount,
			VisitTime:  entry.VisitTime,
		})
	}

	return historyResult{
		entries: renderEntries,
		count:   len(historyData),
	}
}

func (m historyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.done {
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case historyResult:
		if msg.err != nil {
			m.content = errorStyle.Render("❌ Error: "+msg.err.Error()) + "\n\n" +
				infoStyle.Render("Make sure your browser is closed and try again.") + "\n\n" +
				promptStyle.Render("Press 'q' to quit or 'enter' to continue")
			m.err = msg.err
			m.done = true
		} else {
			m.content = successStyle.Render(fmt.Sprintf("✅ Found %d history entries", msg.count)) + "\n\n" +
				infoStyle.Render("🚀 Starting "+m.browserChoice+" History Visualizer...") + "\n\n" +
				promptStyle.Render("Press 'enter' to continue or 'q' to quit")

			// Start the visualizer
			go func() {
				err := render.RunChromeHistoryViewer(msg.entries)
				if err != nil {
					log.Printf("Error running history visualizer: %v", err)
				}
			}()

			m.done = true
		}
		m.viewport.SetContent(m.content)
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.viewport.SetContent(m.content)
		return m, nil
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m historyModel) View() string {
	if m.done {
		footer := "\n" + infoStyle.Render("✨ Thanks for using Histograph!")
		return cardStyle.Render(m.viewport.View() + footer)
	}

	return cardStyle.Render(m.viewport.View())
}

func main() {
	// Get user's browser choice
	choice, err := render.GetUserBrowserChoice()
	if err != nil {
		log.Fatal("Error getting browser choice:", err)
	}

	// Create and run the history processing model
	model := newHistoryModel(choice)
	prog := tea.NewProgram(model)

	if _, err := prog.Run(); err != nil {
		log.Fatal("Error running program:", err)
	}
}
