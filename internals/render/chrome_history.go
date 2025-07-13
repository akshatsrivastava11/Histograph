// render/chrome_history.go
package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/akshatsrivastava11/Histograph/internals/types"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChromeHistoryModel represents the state for Chrome history visualization
type ChromeHistoryModel struct {
	viewport     viewport.Model
	historyData  []types.VisitEntry
	currentView  string // "overview", "timeline", "sites", "details"
	selectedItem int
	ready        bool
	width        int
	height       int
}

// Styles for the UI
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 80).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#F25D94")).
			Padding(0, 1).
			Bold(true)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			MarginBottom(1)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	chartStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#04B575")).
			Padding(1, 2).
			MarginBottom(1)
)

// NewChromeHistoryModel creates a new Chrome history visualization model
func NewChromeHistoryModel(historyData []types.VisitEntry, width, height int) ChromeHistoryModel {
	vp := viewport.New(70-4, 100-6)

	m := ChromeHistoryModel{
		viewport:     vp,
		historyData:  historyData,
		currentView:  "overview",
		selectedItem: 0,
		width:        width,
		height:       height,
	}

	m.updateContent()
	return m
}

func (m ChromeHistoryModel) Init() tea.Cmd {
	return nil
}

func (m ChromeHistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.currentView = "overview"
			m.updateContent()
		case "2":
			m.currentView = "timeline"
			m.updateContent()
		case "3":
			m.currentView = "sites"
			m.updateContent()
		case "4":
			m.currentView = "details"
			m.updateContent()
		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
			}
		case "down", "j":
			if m.selectedItem < len(m.historyData)-1 {
				m.selectedItem++
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 6
		m.updateContent()
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m ChromeHistoryModel) View() string {
	if !m.ready {
		return "Loading Chrome history..."
	}

	header := titleStyle.Render("🌐 Browser History Analyzer") + "\n\n"

	nav := fmt.Sprintf("%s | %s | %s | %s\n\n",
		m.navItem("1", "Overview", m.currentView == "overview"),
		m.navItem("2", "Timeline", m.currentView == "timeline"),
		m.navItem("3", "Top Sites", m.currentView == "sites"),
		m.navItem("4", "Details", m.currentView == "details"))

	footer := dimStyle.Render("Press 1-4 to switch views, ↑/↓ to navigate, q to quit")

	content := header + nav + m.viewport.View() + "\n" + footer
	return content
}

func (m *ChromeHistoryModel) navItem(key, text string, active bool) string {
	if active {
		return highlightStyle.Render(fmt.Sprintf("[%s] %s", key, text))
	}
	return dimStyle.Render(fmt.Sprintf("[%s] %s", key, text))
}

func (m *ChromeHistoryModel) updateContent() {
	var content string

	switch m.currentView {
	case "overview":
		content = m.renderOverview()
	case "timeline":
		content = m.renderTimeline()
	case "sites":
		content = m.renderTopSites()
	case "details":
		content = m.renderDetails()
	}

	m.viewport.SetContent(content)
	m.ready = true
}

func (m ChromeHistoryModel) renderOverview() string {
	if len(m.historyData) == 0 {
		return cardStyle.Render("No Chrome history data found")
	}

	// Calculate statistics
	totalVisits := len(m.historyData)
	totalVisitCount := 0
	domains := make(map[string]int)

	for _, entry := range m.historyData {
		totalVisitCount += entry.VisitCount
		domain := extractDomain(entry.URL)
		domains[domain]++
	}

	// Create overview content
	stats := fmt.Sprintf("📊 Total Entries: %d\n", totalVisits) +
		fmt.Sprintf("🔄 Total Visits: %d\n", totalVisitCount) +
		fmt.Sprintf("🌐 Unique Domains: %d\n", len(domains))

	// Create a simple bar chart for top domains
	chart := m.createDomainChart(domains)

	overview := cardStyle.Render(headerStyle.Render("📈 Statistics") + "\n\n" + stats)
	chartCard := chartStyle.Render(headerStyle.Render("🔝 Top Domains") + "\n\n" + chart)

	return overview + "\n" + chartCard
}

func (m ChromeHistoryModel) renderTimeline() string {
	if len(m.historyData) == 0 {
		return cardStyle.Render("No timeline data available")
	}

	// Group visits by date
	dateGroups := make(map[string][]types.VisitEntry)
	for _, entry := range m.historyData {
		date := entry.VisitTime.Format("2006-01-02")
		dateGroups[date] = append(dateGroups[date], entry)
	}

	// Sort dates
	var dates []string
	for date := range dateGroups {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Create timeline visualization
	var timeline strings.Builder
	timeline.WriteString(headerStyle.Render("📅 Timeline View") + "\n\n")

	// Show last 10 days
	start := len(dates) - 10
	if start < 0 {
		start = 0
	}

	for i := start; i < len(dates); i++ {
		date := dates[i]
		entries := dateGroups[date]

		timeline.WriteString(fmt.Sprintf("📅 %s (%d visits)\n", date, len(entries)))

		// Show activity bar
		activityBar := m.createActivityBar(len(entries), 50)
		timeline.WriteString("   " + activityBar + "\n")

		// Show top sites for this date
		if len(entries) > 0 {
			topSite := entries[0]
			timeline.WriteString(fmt.Sprintf("   🔝 %s\n", truncateString(topSite.Title, 50)))
		}
		timeline.WriteString("\n")
	}

	return cardStyle.Render(timeline.String())
}

func (m ChromeHistoryModel) renderTopSites() string {
	if len(m.historyData) == 0 {
		return cardStyle.Render("No sites data available")
	}

	// Aggregate by domain
	domainData := make(map[string]struct {
		visits int
		count  int
		title  string
	})

	for _, entry := range m.historyData {
		domain := extractDomain(entry.URL)
		data := domainData[domain]
		data.visits += entry.VisitCount
		data.count++
		if data.title == "" {
			data.title = entry.Title
		}
		domainData[domain] = data
	}

	// Sort by visits
	type siteData struct {
		domain string
		visits int
		count  int
		title  string
	}

	var sites []siteData
	for domain, data := range domainData {
		sites = append(sites, siteData{
			domain: domain,
			visits: data.visits,
			count:  data.count,
			title:  data.title,
		})
	}

	sort.Slice(sites, func(i, j int) bool {
		return sites[i].visits > sites[j].visits
	})

	// Create top sites display
	var content strings.Builder
	content.WriteString(headerStyle.Render("🏆 Top Sites") + "\n\n")

	for i, site := range sites {
		if i >= 15 { // Show top 15
			break
		}

		rank := fmt.Sprintf("%2d.", i+1)
		bar := m.createVisitBar(site.visits, sites[0].visits, 20)

		content.WriteString(fmt.Sprintf("%s %s %s\n",
			highlightStyle.Render(rank),
			bar,
			site.domain))
		content.WriteString(fmt.Sprintf("    %s visits • %s entries\n",
			dimStyle.Render(fmt.Sprintf("%d", site.visits)),
			dimStyle.Render(fmt.Sprintf("%d", site.count))))
		content.WriteString("\n")
	}

	return cardStyle.Render(content.String())
}

func (m ChromeHistoryModel) renderDetails() string {
	if len(m.historyData) == 0 {
		return cardStyle.Render("No detailed data available")
	}

	// Show detailed view of recent entries
	var content strings.Builder
	content.WriteString(headerStyle.Render("🔍 Recent History Details") + "\n\n")

	// Sort by visit time (most recent first)
	sortedEntries := make([]types.VisitEntry, len(m.historyData))
	copy(sortedEntries, m.historyData)
	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].VisitTime.After(sortedEntries[j].VisitTime)
	})

	// Show top 20 recent entries
	for i, entry := range sortedEntries {
		if i >= 20 {
			break
		}

		timeStr := entry.VisitTime.Format("Jan 2, 15:04")
		title := truncateString(entry.Title, 60)
		if title == "" {
			title = "Untitled"
		}

		content.WriteString(fmt.Sprintf("🌐 %s\n", highlightStyle.Render(title)))
		content.WriteString(fmt.Sprintf("   %s\n", dimStyle.Render(entry.URL)))
		content.WriteString(fmt.Sprintf("   %s • %s visits\n",
			dimStyle.Render(timeStr),
			dimStyle.Render(fmt.Sprintf("%d", entry.VisitCount))))
		content.WriteString("\n")
	}

	return cardStyle.Render(content.String())
}

// Helper functions

func (m ChromeHistoryModel) createDomainChart(domains map[string]int) string {
	// Sort domains by frequency
	type domainCount struct {
		domain string
		count  int
	}

	var sorted []domainCount
	for domain, count := range domains {
		sorted = append(sorted, domainCount{domain, count})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	var chart strings.Builder
	maxCount := 0
	if len(sorted) > 0 {
		maxCount = sorted[0].count
	}

	// Show top 8 domains
	for i, item := range sorted {
		if i >= 8 {
			break
		}

		bar := m.createVisitBar(item.count, maxCount, 25)
		chart.WriteString(fmt.Sprintf("%-20s %s %d\n",
			truncateString(item.domain, 20),
			bar,
			item.count))
	}

	return chart.String()
}

func (m ChromeHistoryModel) createVisitBar(current, max, width int) string {
	if max == 0 {
		return strings.Repeat("░", width)
	}

	filled := int(float64(current) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}

	return highlightStyle.Render(strings.Repeat("█", filled)) +
		dimStyle.Render(strings.Repeat("░", width-filled))
}

func (m ChromeHistoryModel) createActivityBar(activity, max int) string {
	width := 30
	if max == 0 {
		return strings.Repeat("░", width)
	}

	filled := int(float64(activity) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}

	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func extractDomain(url string) string {
	// Simple domain extraction
	if strings.HasPrefix(url, "http://") {
		url = url[7:]
	} else if strings.HasPrefix(url, "https://") {
		url = url[8:]
	}

	if strings.HasPrefix(url, "www.") {
		url = url[4:]
	}

	parts := strings.Split(url, "/")
	return parts[0]
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// RunChromeHistoryViewer starts the Chrome history visualization
func RunChromeHistoryViewer(historyData []types.VisitEntry) error {
	m := NewChromeHistoryModel(historyData, 120, 40)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
