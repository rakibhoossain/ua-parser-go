package referrer

// ReferrerType represents the semantic classification of a referral source.
type ReferrerType string

const (
	TypeSearch   ReferrerType = "search"
	TypeSocial   ReferrerType = "social"
	TypeAI       ReferrerType = "ai"
	TypeTech     ReferrerType = "tech"
	TypeCommerce ReferrerType = "commerce"
	TypeEmail    ReferrerType = "email"
	TypeContent  ReferrerType = "content"
	TypeUnknown  ReferrerType = "unknown"
)

// ReferrerEntry is the definition stored in the domain dictionary.
type ReferrerEntry struct {
	Type ReferrerType `json:"type"`
	Name string       `json:"name"`
}

// Result is the parsed referral metadata output.
type Result struct {
	Name       string       `json:"name"`
	Type       ReferrerType `json:"type"`
	Domain     string       `json:"domain,omitempty"`
	URL        string       `json:"url,omitempty"`
	FaviconURL string       `json:"favicon_url,omitempty"`
}
