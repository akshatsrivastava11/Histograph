package parse

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/types"

	_ "github.com/mattn/go-sqlite3"
)

// Converts Chrome's Webkit timestamp to Unix time
func chromeTimeToUnix(microseconds int64) time.Time {
	const offset = 11644473600
	seconds := microseconds/1000000 - offset
	return time.Unix(seconds, 0)
}

// getChromeHistoryPath returns the path to the Chrome history file for the current OS.
func GetChromeHistoryPath() (string, error) {
	// Allow override via environment variable
	if envPath := os.Getenv("CHROME_HISTORY_PATH"); envPath != "" {
		return envPath, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	switch runtime.GOOS {
	case "linux":
		return filepath.Join(home, ".config", "google-chrome", "Default", "History"), nil
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Google", "Chrome", "Default", "History"), nil
	case "windows":
		return filepath.Join(home, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "History"), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// ParseChromeHistory connects to Chrome's history database and returns a slice of VisitEntry
func ParseChromeHistory() ([]types.VisitEntry, error) {
	fmt.Println("Parsing Chrome's History")

	historyPath, err := GetChromeHistoryPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Chrome history database: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`
	  SELECT urls.url, urls.title, urls.visit_count, visits.visit_time
        FROM urls
        JOIN visits ON urls.id = visits.url
        ORDER BY visits.visit_time DESC
        LIMIT 20;
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query Chrome history: %w", err)
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
			return nil, fmt.Errorf("failed to scan Chrome history row: %w", err)
		}

		convertedTime := chromeTimeToUnix(visitTime)

		history = append(history, types.VisitEntry{
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			VisitTime:  convertedTime,
		})
	}

	return history, nil
}
