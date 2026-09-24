package referrer

import (
	"net/http"
	"net/url"
	"strings"
)

// Parse resolves a referrer URL into a recognized referral source.
func Parse(rawURL string) *Result {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return &Result{}
	}

	cleanURL := strings.TrimRight(rawURL, "/")
	hostname := getHostname(cleanURL)
	if hostname == "" {
		return &Result{URL: cleanURL}
	}

	hostnameNoWWW := strings.TrimPrefix(hostname, "www.")
	entry, ok := Lookup(hostname)
	if !ok {
		entry, ok = Lookup(hostnameNoWWW)
	}

	domain := hostnameNoWWW
	if !ok {
		return &Result{
			Domain:     domain,
			URL:        cleanURL,
			FaviconURL: FaviconURL(domain),
		}
	}

	return &Result{
		Name:       entry.Name,
		Type:       entry.Type,
		Domain:     domain,
		URL:        cleanURL,
		FaviconURL: FaviconURL(domain),
	}
}

// ParseWithQuery extracts referral information from UTM or query parameters (e.g. utm_source, ref, utm_referrer).
func ParseWithQuery(query map[string]string) *Result {
	if len(query) == 0 {
		return nil
	}

	var source string
	for _, key := range []string{"utm_source", "ref", "utm_referrer"} {
		if val, exists := query[key]; exists && strings.TrimSpace(val) != "" {
			source = strings.ToLower(strings.TrimSpace(val))
			break
		}
	}

	if source == "" {
		return nil
	}

	// 1. Direct domain match
	if entry, ok := Lookup(source); ok {
		return &Result{
			Name:       entry.Name,
			Type:       entry.Type,
			Domain:     source,
			FaviconURL: FaviconURL(source),
		}
	}

	// 2. Try with .com suffix
	withCom := source + ".com"
	if entry, ok := Lookup(withCom); ok {
		return &Result{
			Name:       entry.Name,
			Type:       entry.Type,
			Domain:     withCom,
			FaviconURL: FaviconURL(withCom),
		}
	}

	// 3. Match against known source names (e.g. "google", "twitter", "chatgpt")
	if entry, domain, ok := LookupByName(source); ok {
		return &Result{
			Name:       entry.Name,
			Type:       entry.Type,
			Domain:     domain,
			FaviconURL: FaviconURL(domain),
		}
	}

	// 4. Fallback to custom source
	return &Result{
		Name:       source,
		Type:       TypeUnknown,
		Domain:     source,
		FaviconURL: FaviconURL(source),
	}
}

// ParseRequest inspects an HTTP request for both the Referer header and UTM query parameters.
func ParseRequest(r *http.Request) *Result {
	if r == nil {
		return &Result{}
	}

	// 1. First, check UTM query parameters for campaign attribution
	q := r.URL.Query()
	qMap := make(map[string]string, len(q))
	for k, v := range q {
		if len(v) > 0 {
			qMap[k] = v[0]
		}
	}

	if res := ParseWithQuery(qMap); res != nil {
		if ref := r.Referer(); ref != "" {
			res.URL = ref
		}
		return res
	}

	// 2. Fall back to standard Referer header
	return Parse(r.Referer())
}

func getHostname(rawURL string) string {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(u.Hostname())
}
