package main

import (
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
)

type VisitEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	VisitTime  time.Time `json:"visit_time"`
}

func main() {
	parse.ParseChromeHistory()
}

