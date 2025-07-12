package types

import (
	"time"
)

type VisitEntry struct {
	URL        string    `json:"url"`
	Title      string    `json:"title"`
	VisitCount int       `json:"visit_count"`
	VisitTime  time.Time `json:"visit_time"`
}
