package models

import (
	"time"

	"github.com/akshatsrivastava11/Histograph/internals/types"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	viewport        viewport.Model
	state           bool
	tabs            []string
	activeTab       string
	animationTicker *time.Ticker
	HistoryData     types.VisitEntry
}
