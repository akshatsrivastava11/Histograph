package main

import (
	"fmt"
	"log"
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
	"github.com/akshatsrivastava11/Histograph/internals/render"
)

type VisitEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	VisitTime  time.Time `json:"visit_time"`
}

func main() {
	choice, err := render.GetUserBrowserChoice()
	if err != nil {
		log.Fatal(err)
	}
	if choice == "Firefox" {
		fmt.Println("User seleceted ", choice)

	}
	if choice == "Chrome" {
		historyData := parse.ParseChromeHistory()

		for _, entry := range historyData {
			fmt.Println("From return:", entry.Title, entry.URL, entry.VisitCount, entry.VisitTime)
		}
	}

}
