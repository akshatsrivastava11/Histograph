// main.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
	"github.com/akshatsrivastava11/Histograph/internals/render"
	"github.com/akshatsrivastava11/Histograph/internals/types"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
	spinner       spinner.Model
	loading       bool
}

func newHistoryModel(browserChoice string) historyModel {
	vp := viewport.New(80, 20)
	sp := spinner.New()
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)

	content := titleStyle.Render("🔍 Fetching " + browserChoice + " history...")
	// timer.New(10000).Init()
	return historyModel{
		browserChoice: browserChoice,
		viewport:      vp,
		content:       content,
		done:          false,
		spinner:       sp,
		loading:       true,
	}
}

func (m historyModel) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.spinner.Tick,
		processHistoryCmd(m.browserChoice),
	)
	// return nil
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
	entries []types.VisitEntry
	count   int
	err     error
}

func processChromeHistory() historyResult {
	historyData, err := parse.ParseChromeHistory()
	if err != nil {
		return historyResult{err: err}
	}
	if len(historyData) == 0 {
		return historyResult{
			err: fmt.Errorf("no Chrome history found or unable to access Chrome history"),
		}
	}

	return historyResult{
		entries: historyData,
		count:   len(historyData),
	}
}

func processFirefoxHistory() historyResult {
	historyData, err := parse.ParseFirefoxHistory()
	if err != nil {
		return historyResult{err: err}
	}
	if len(historyData) == 0 {
		return historyResult{
			err: fmt.Errorf("no Firefox history found or unable to access Firefox history"),
		}
	}

	return historyResult{
		entries: historyData,
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
			if m.err != nil {
				// Retry fetching history
				m.done = false
				m.err = nil
				m.loading = true
				m.content = titleStyle.Render("🔍 Fetching " + m.browserChoice + " history...")
				return m, tea.Batch(m.spinner.Tick, processHistoryCmd(m.browserChoice))
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case historyResult:
		if msg.err != nil {
			m.content = errorStyle.Render("❌ Error: "+msg.err.Error()) + "\n\n" +
				infoStyle.Render("Make sure your browser is closed and try again. If you see a 'database is locked' error, close your browser and retry.") + "\n\n" +
				promptStyle.Render("Press 'q' to quit or 'enter' to retry")
			m.err = msg.err
			m.done = true
			m.loading = false
		} else {
			m.content = successStyle.Render(fmt.Sprintf("✅ Found %d history entries", msg.count)) + "\n\n" +
				infoStyle.Render("🚀 Starting "+m.browserChoice+" History Visualizer...") + "\n\n" +
				promptStyle.Render("Press 'enter' to continue or 'q' to quit")
			// timer.New(20000000).Init()
			// Start the visualizer
			go func() {
				err := render.RunChromeHistoryViewer(msg.entries)
				if err != nil {
					debugLog("Error running history visualizer: %v", err)
				}
			}()

			m.done = true
			m.loading = false
		}
		m.viewport.SetContent(m.content)
		return m, nil

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
		m.viewport.SetContent(m.content)
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
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

	if m.loading {
		spinnerView := m.spinner.View() + " " + titleStyle.Render("Fetching "+m.browserChoice+" history...")
		return cardStyle.Render(spinnerView)
	}
	return cardStyle.Render(m.viewport.View())
}

// debugLog prints debug output only if HISTOGRAPH_DEBUG=1 is set in the environment.
func debugLog(format string, v ...interface{}) {
	if os.Getenv("HISTOGRAPH_DEBUG") == "1" {
		log.Printf(format, v...)
	}
}

func main() {
	// Get user's browser choice
	choice, err := render.GetUserBrowserChoice()
	if err != nil {
		fmt.Println("Error getting browser choice:", err)
		return
	}

	// Create and run the history processing model
	model := newHistoryModel(choice)
	// time.Sleep(3 * time.Second)
	prog := tea.NewProgram(model)

	if _, err := prog.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
