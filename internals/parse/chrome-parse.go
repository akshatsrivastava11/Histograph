package parse

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // make sure this is imported in the same file
	"log"
	"time"
)

func chromeTimeToUnix(microseconds int64) time.Time {
	// WebKit epoch offset to Unix (in seconds)
	const offset = 11644473600
	seconds := microseconds/1000000 - offset
	return time.Unix(seconds, 0)
}

func ParseChromeHistory() {
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
	for rows.Next() {
		var url string
		var title string
		var visitCount int
		var visitTime int64

		err = rows.Scan(&url, &title, &visitCount, &visitTime)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Visited: %s\nTitle: %s\nCount: %d\nTime: %s\n\n",
			url, title, visitCount, chromeTimeToUnix(visitTime).Format(time.RFC3339))
	}

}
