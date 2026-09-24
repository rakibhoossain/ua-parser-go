package bots

import "net/http"

// BotMatch represents a successful bot identification with its classification and producer details.
type BotMatch struct {
	Name     string `json:"name"`
	Category string `json:"category"` // e.g. "Search bot", "AI Crawler", "Feed Fetcher", "Security checker", "Site Monitor", "Service Agent", "Crawler"
	URL      string `json:"url,omitempty"`
	Producer string `json:"producer,omitempty"`
}

// BotSuspicion summarizes multi-factor heuristic fraud and crawler signals for an event.
type BotSuspicion struct {
	IsBot   bool     `json:"is_bot"`
	Reasons []string `json:"reasons,omitempty"`
}

// BotCategoryThreshold defines the minimum number of DISTINCT signal categories
// required to flag an event as a bot (default 2, matching OpenPanel's suspicion threshold).
const BotCategoryThreshold = 2

// SuspicionOptions specifies inputs for multi-factor bot evaluation.
type SuspicionOptions struct {
	IsDatacenter     bool
	ASN              string
	Headers          http.Header
	UserAgent        string
	IsServer         bool
	ClientSecretAuth bool
}

