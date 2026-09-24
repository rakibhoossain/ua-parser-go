package referrer

import (
	"net/http"
	"net/url"
	"strings"
)

// Parse resolves a referrer URL into a recognized referral source and metadata.
func Parse(rawURL string) *Result {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return &Result{
			Type: TypeDirect,
		}
	}

	cleanURL := strings.TrimRight(rawURL, "/")
	hostname := getHostname(cleanURL)
	if hostname == "" {
		return &Result{
			Name: cleanURL,
			Type: TypeReferral,
			URL:  cleanURL,
		}
	}

	hostnameNoWWW := strings.TrimPrefix(hostname, "www.")

	// 1. Exact domain lookup
	entry, ok := Lookup(hostname)
	if !ok {
		entry, ok = Lookup(hostnameNoWWW)
	}

	// 2. Global search engine & social network prefix fallbacks (catches all international TLDs)
	if !ok {
		switch {
		case strings.Contains(hostnameNoWWW, "google."):
			entry = ReferrerEntry{Type: TypeSearch, Name: "Google"}
			ok = true
		case strings.Contains(hostnameNoWWW, "bing."):
			entry = ReferrerEntry{Type: TypeSearch, Name: "Bing"}
			ok = true
		case strings.Contains(hostnameNoWWW, "duckduckgo."):
			entry = ReferrerEntry{Type: TypeSearch, Name: "DuckDuckGo"}
			ok = true
		case strings.Contains(hostnameNoWWW, "yahoo."):
			entry = ReferrerEntry{Type: TypeSearch, Name: "Yahoo"}
			ok = true
		case strings.Contains(hostnameNoWWW, "youtube.") || hostnameNoWWW == "youtu.be":
			entry = ReferrerEntry{Type: TypeSocial, Name: "YouTube"}
			ok = true
		case strings.Contains(hostnameNoWWW, "facebook.") || hostnameNoWWW == "fb.com":
			entry = ReferrerEntry{Type: TypeSocial, Name: "Facebook"}
			ok = true
		case strings.Contains(hostnameNoWWW, "twitter.") || hostnameNoWWW == "x.com" || hostnameNoWWW == "t.co":
			entry = ReferrerEntry{Type: TypeSocial, Name: "X (Twitter)"}
			ok = true
		case strings.Contains(hostnameNoWWW, "instagram."):
			entry = ReferrerEntry{Type: TypeSocial, Name: "Instagram"}
			ok = true
		case strings.Contains(hostnameNoWWW, "linkedin."):
			entry = ReferrerEntry{Type: TypeSocial, Name: "LinkedIn"}
			ok = true
		case strings.Contains(hostnameNoWWW, "pinterest."):
			entry = ReferrerEntry{Type: TypeSocial, Name: "Pinterest"}
			ok = true
		case strings.Contains(hostnameNoWWW, "tiktok."):
			entry = ReferrerEntry{Type: TypeSocial, Name: "TikTok"}
			ok = true
		}
	}

	domain := hostnameNoWWW
	if !ok {
		// Unknown external referrer
		return &Result{
			Name:       domain,
			Type:       TypeReferral,
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

// ParseWithQuery extracts referral and UTM marketing campaign parameters.
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

	medium := strings.ToLower(strings.TrimSpace(query["utm_medium"]))
	campaign := strings.TrimSpace(query["utm_campaign"])
	content := strings.TrimSpace(query["utm_content"])
	term := strings.TrimSpace(query["utm_term"])

	if source == "" {
		if medium != "" || campaign != "" {
			return &Result{
				Type:     classifyMedium(medium),
				Medium:   medium,
				Campaign: campaign,
				Content:  content,
				Term:     term,
			}
		}
		return nil
	}

	res := &Result{
		Medium:   medium,
		Campaign: campaign,
		Content:  content,
		Term:     term,
	}

	// 1. Direct domain match
	if entry, ok := Lookup(source); ok {
		res.Name = entry.Name
		res.Type = entry.Type
		res.Domain = source
		res.FaviconURL = FaviconURL(source)
	} else if entry, ok := Lookup(source + ".com"); ok { // 2. Try with .com suffix
		res.Name = entry.Name
		res.Type = entry.Type
		res.Domain = source + ".com"
		res.FaviconURL = FaviconURL(source + ".com")
	} else if entry, domain, ok := LookupByName(source); ok { // 3. Match against known source names
		res.Name = entry.Name
		res.Type = entry.Type
		res.Domain = domain
		res.FaviconURL = FaviconURL(domain)
	} else { // 4. Fallback to custom source
		res.Name = source
		res.Type = TypeReferral
		res.Domain = source
		res.FaviconURL = FaviconURL(source)
	}

	// If utm_medium specifies a more specific classification, override Type
	if mediumType := classifyMedium(medium); mediumType != TypeUnknown {
		res.Type = mediumType
	}

	return res
}

func classifyMedium(medium string) ReferrerType {
	if medium == "" {
		return TypeUnknown
	}
	switch {
	case strings.Contains(medium, "cpc") || strings.Contains(medium, "ppc") ||
		strings.Contains(medium, "paid") || strings.Contains(medium, "ad"):
		return TypePaid
	case strings.Contains(medium, "social"):
		return TypeSocial
	case strings.Contains(medium, "email") || strings.Contains(medium, "newsletter"):
		return TypeEmail
	default:
		return TypeUnknown
	}
}

// ParseRequest inspects an HTTP request for both the Referer header and UTM query parameters.
func ParseRequest(r *http.Request) *Result {
	if r == nil {
		return &Result{Type: TypeDirect}
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
