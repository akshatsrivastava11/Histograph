// main.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
	"github.com/akshatsrivastava11/Histograph/internals/render"
)

// VisitEntry represents a browser history entry
type VisitEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	VisitTime  time.Time `json:"visit_time"`
}

func main() {
	// Get user's browser choice

	choice, err := render.GetUserBrowserChoice()
	if err != nil {
		log.Fatal("Error getting browser choice:", err)
	}

	fmt.Printf("You selected: %s\n", choice)

	switch choice {
	case "Firefox":
		handleFirefoxHistory()
	case "Chrome":
		handleChromeHistory()
	default:
		fmt.Println("Invalid browser selection")
	}
}

func handleChromeHistory() {
	fmt.Println("🔍 Fetching Chrome history...")

	// Parse Chrome history using your existing function
	historyData := parse.ParseChromeHistory()

	if len(historyData) == 0 {
		fmt.Println("❌ No Chrome history found or unable to access Chrome history.")
		fmt.Println("Make sure Chrome is closed and try again.")
		return
	}

	fmt.Printf("✅ Found %d history entries\n", len(historyData))
	fmt.Println("🚀 Starting Chrome History Visualizer...")
	fmt.Println("Press any key to continue...")

	// Wait for user input
	var input string
	fmt.Scanln(&input)

	// Convert your VisitEntry to the render package's VisitEntry
	var renderEntries []render.VisitEntry
	for _, entry := range historyData {
		renderEntries = append(renderEntries, render.VisitEntry{
			URL:        entry.URL,
			Title:      entry.Title,
			VisitCount: entry.VisitCount,
			VisitTime:  entry.VisitTime,
		})
	}

	// Run the Chrome history visualizer
	err := render.RunChromeHistoryViewer(renderEntries)
	if err != nil {
		log.Fatal("Error running Chrome history visualizer:", err)
	}

	fmt.Println("\n✨ Thanks for using Histograph!")
}

func handleFirefoxHistory() {
	fmt.Println("🔍 Fetching Chrome history...")

	// Parse Chrome history using your existing function
	historyData := parse.ParseFirefoxHistory()

	if len(historyData) == 0 {
		fmt.Println("❌ No Chrome history found or unable to access Chrome history.")
		fmt.Println("Make sure Chrome is closed and try again.")
		return
	}

	fmt.Printf("✅ Found %d history entries\n", len(historyData))
	fmt.Println("🚀 Starting Chrome History Visualizer...")
	fmt.Println("Press any key to continue...")

	// Wait for user input
	var input string
	fmt.Scanln(&input)

	// Convert your VisitEntry to the render package's VisitEntry
	var renderEntries []render.VisitEntry
	for _, entry := range historyData {
		renderEntries = append(renderEntries, render.VisitEntry{
			URL:        entry.URL,
			Title:      entry.Title,
			VisitCount: entry.VisitCount,
			VisitTime:  entry.VisitTime,
		})
	}

	// Run the Chrome history visualizer
	err := render.RunChromeHistoryViewer(renderEntries)
	if err != nil {
		log.Fatal("Error running Chrome history visualizer:", err)
	}

	fmt.Println("\n✨ Thanks for using Histograph!")
}
