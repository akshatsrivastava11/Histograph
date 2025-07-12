package parse

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/types" // Replace `project` with your actual module name

	_ "github.com/mattn/go-sqlite3" // Required for sqlite
)

// Converts Chrome's Webkit timestamp to Unix time
func chromeTimeToUnix(microseconds int64) time.Time {
	const offset = 11644473600
	seconds := microseconds/1000000 - offset
	return time.Unix(seconds, 0)
}

// ParseChromeHistory connects to Chrome's history database and returns a slice of VisitEntry
func ParseChromeHistory() []types.VisitEntry {
	fmt.Println("Parsing Chrome's History")

	db, err := sql.Open("sqlite3", "/home/zeek1108/.config/google-chrome/Default/History")
	if err != nil {
		log.Fatal(err)
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

		convertedTime := chromeTimeToUnix(visitTime)

		// Print for debug/logging

		// Add to result slice
		history = append(history, types.VisitEntry{
			URL:        url,
			Title:      title,
			VisitCount: visitCount,
			VisitTime:  convertedTime,
		})
	}

	return history
}
