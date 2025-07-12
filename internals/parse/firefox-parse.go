package parse

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/types"
	_ "github.com/mattn/go-sqlite3"
)

// Convert Firefox microsecond timestamp to time.Time
func firefoxTimeToUnix(microseconds int64) time.Time {
	return time.UnixMicro(microseconds)
}

// Get the path to the first available Firefox profile
func getFirefoxHistoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home dir:", err)
	}

	// Find the profile dir inside ~/.mozilla/firefox
	profilesDir := filepath.Join(homeDir, ".mozilla", "firefox")
	dirs, err := os.ReadDir(profilesDir)
	if err != nil {
		log.Fatal("Failed to read firefox profile dir:", err)
	}

	for _, d := range dirs {
		if d.IsDir() && filepath.Ext(d.Name()) == ".default-release" {
			return filepath.Join(profilesDir, d.Name(), "places.sqlite")
		}
	}

	log.Fatal("Could not find Firefox default-release profile")
	return ""
}

// ParseFirefoxHistory connects to Firefox's history database and returns recent visits
func ParseFirefoxHistory() []types.VisitEntry {
	fmt.Println("Parsing Firefox History")

	historyPath := getFirefoxHistoryPath()

	db, err := sql.Open("sqlite3", historyPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT p.url, p.title, p.visit_count, v.visit_date
		FROM moz_places p
		JOIN moz_historyvisits v ON p.id = v.place_id
		ORDER BY v.visit_date DESC
		LIMIT 20;
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var history []types.VisitEntry

	for rows.Next() {
		var url string
		var title string
		var visitCount int
		var visitTime int64

		err = rows.Scan(&url, &title, &visitCount, &visitTime)
		if err != nil {
			log.Fatal(err)
		}

		convertedTime := firefoxTimeToUnix(visitTime)

		fmt.Printf("Visited: %s\nTitle: %s\nCount: %d\nTime: %s\n\n",
			url, title, visitCount, convertedTime.Format(time.RFC3339))

		history = append(history, types.VisitEntry{
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			VisitTime:  convertedTime,
		})
	}

	return history
}
