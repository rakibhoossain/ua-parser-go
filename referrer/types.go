package referrer

// ReferrerType represents the semantic classification of a referral source.
type ReferrerType string

const (
	TypeDirect   ReferrerType = "direct"
	TypeReferral ReferrerType = "referral"
	TypeSearch   ReferrerType = "search"
	TypeSocial   ReferrerType = "social"
	TypeAI       ReferrerType = "ai"
	TypeTech     ReferrerType = "tech"
	TypeCommerce ReferrerType = "commerce"
	TypeEmail    ReferrerType = "email"
	TypePaid     ReferrerType = "paid"
	TypeContent  ReferrerType = "content"
	TypeUnknown  ReferrerType = "unknown"
)

// ReferrerEntry is the definition stored in the domain dictionary.
type ReferrerEntry struct {
	Type ReferrerType `json:"type"`
	Name string       `json:"name"`
}

// Result is the parsed referral and marketing attribution output.
type Result struct {
	Name       string       `json:"name"`                  // e.g. "Google", "ChatGPT", "Twitter"
	Type       ReferrerType `json:"type"`                  // "search", "social", "ai", "paid", "direct", "referral"
	Domain     string       `json:"domain,omitempty"`      // e.g. "google.com", "chatgpt.com"
	URL        string       `json:"url,omitempty"`         // original referrer URL
	FaviconURL string       `json:"favicon_url,omitempty"` // URL to domain favicon
	Medium     string       `json:"medium,omitempty"`      // utm_medium (e.g. cpc, email, social)
	Campaign   string       `json:"campaign,omitempty"`    // utm_campaign
	Content    string       `json:"content,omitempty"`     // utm_content
	Term       string       `json:"term,omitempty"`        // utm_term
}
