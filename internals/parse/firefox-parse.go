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

// Convert Firefox microsecond timestamp to time.Time
func firefoxTimeToUnix(microseconds int64) time.Time {
	return time.UnixMicro(microseconds)
}

// Get the path to the first available Firefox profile
func GetFirefoxHistoryPath() (string, error) {
	// Allow override via environment variable
	if envPath := os.Getenv("FIREFOX_HISTORY_PATH"); envPath != "" {
		return envPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}

	var profilesDir string
	switch runtime.GOOS {
	case "linux":
		profilesDir = filepath.Join(homeDir, ".mozilla", "firefox")
	case "darwin":
		profilesDir = filepath.Join(homeDir, "Library", "Application Support", "Firefox", "Profiles")
	case "windows":
		profilesDir = filepath.Join(homeDir, "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	dirs, err := os.ReadDir(profilesDir)
	if err != nil {
		return "", fmt.Errorf("failed to read firefox profile dir: %w", err)
	}

	for _, d := range dirs {
		if d.IsDir() && (filepath.Ext(d.Name()) == ".default-release" || filepath.Ext(d.Name()) == ".default") {
			return filepath.Join(profilesDir, d.Name(), "places.sqlite"), nil
		}
	}

	return "", fmt.Errorf("could not find Firefox default(-release) profile")
}

// ParseFirefoxHistory connects to Firefox's history database and returns recent visits
func ParseFirefoxHistory() ([]types.VisitEntry, error) {
	fmt.Println("Parsing Firefox History")

	historyPath, err := GetFirefoxHistoryPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Firefox history database: %w", err)
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
		return nil, fmt.Errorf("failed to query Firefox history: %w", err)
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
			return nil, fmt.Errorf("failed to scan Firefox history row: %w", err)
		}

		convertedTime := firefoxTimeToUnix(visitTime)

		history = append(history, types.VisitEntry{
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			VisitTime:  convertedTime,
		})
	}

	return history, nil
}
