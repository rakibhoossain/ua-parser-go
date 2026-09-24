package referrer

import (
	"fmt"
	"net/url"
	"strings"
)

// FaviconURL generates a favicon icon URL using DuckDuckGo's favicon proxy.
func FaviconURL(domainOrURL string) string {
	domain := extractDomain(domainOrURL)
	if domain == "" {
		return ""
	}
	return fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s.ico", domain)
}

// GoogleFaviconURL generates a favicon icon URL using Google's favicon service with an optional size in pixels.
func GoogleFaviconURL(domainOrURL string, size int) string {
	domain := extractDomain(domainOrURL)
	if domain == "" {
		return ""
	}
	if size <= 0 {
		size = 64
	}
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=%d", domain, size)
}

func extractDomain(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// If it doesn't contain a scheme, prefix with https:// to parse host properly
	if !strings.Contains(input, "://") {
		input = "https://" + input
	}

	u, err := url.Parse(input)
	if err != nil || u.Hostname() == "" {
		return strings.TrimPrefix(input, "https://")
	}

	host := strings.ToLower(u.Hostname())
	return strings.TrimPrefix(host, "www.")
}
