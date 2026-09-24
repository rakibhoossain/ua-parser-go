package bots

import (
	"net/http"
	"regexp"
	"strings"
)

var (
	chromiumUARegex    = regexp.MustCompile(`(?:Chrome|Chromium|Edg|OPR)/\d`)
	appleWebKitUARegex = regexp.MustCompile(`(CriOS|EdgiOS|OPiOS|FxiOS|iPhone|iPad|iPod)`)
)

// looksLikeBrowser checks if the User-Agent claims to be a standard web browser.
func looksLikeBrowser(ua string) bool {
	return strings.Contains(ua, "Mozilla/")
}

// isBlinkChromium checks if the User-Agent is a Desktop or Android Chromium browser
// that is expected to send Client Hints and Sec-Fetch headers. iOS Chromium (WebKit) is excluded.
func isBlinkChromium(ua string) bool {
	return chromiumUARegex.MatchString(ua) && !appleWebKitUARegex.MatchString(ua)
}

// DetectHeaderAnomalies checks an incoming HTTP request's headers against the claimed User-Agent.
// Returns namespaced reasons (e.g. "header:missing_sec_ch_ua", "header:missing_sec_fetch").
func DetectHeaderAnomalies(h http.Header, ua string) []string {
	if ua == "" || !looksLikeBrowser(ua) || h == nil {
		return nil
	}

	var reasons []string

	// 1. Blink/Chromium claiming to be Chrome/Edge on Desktop/Android but omitting Sec-CH-UA
	if isBlinkChromium(ua) && strings.TrimSpace(h.Get("Sec-CH-UA")) == "" {
		reasons = append(reasons, "header:missing_sec_ch_ua")
	}

	// 2. Fetch-metadata headers sent by all modern browsers on fetch requests
	if strings.TrimSpace(h.Get("Sec-Fetch-Mode")) == "" && strings.TrimSpace(h.Get("Sec-Fetch-Site")) == "" {
		reasons = append(reasons, "header:missing_sec_fetch")
	}

	// 3. Browsers send Accept-Language by default; total absence is suspect
	if strings.TrimSpace(h.Get("Accept-Language")) == "" {
		reasons = append(reasons, "header:missing_accept_language")
	}

	return reasons
}

// SummarizeSignals combines multi-factor reasons and applies the distinct-category threshold.
// Flags as bot (isBot=true) ONLY when >= BotCategoryThreshold (2) distinct categories fire.
func SummarizeSignals(reasons []string) BotSuspicion {
	if len(reasons) == 0 {
		return BotSuspicion{IsBot: false}
	}

	categories := make(map[string]struct{})
	for _, r := range reasons {
		parts := strings.SplitN(r, ":", 2)
		if len(parts) > 0 && parts[0] != "" {
			categories[parts[0]] = struct{}{}
		}
	}

	return BotSuspicion{
		IsBot:   len(categories) >= BotCategoryThreshold,
		Reasons: reasons,
	}
}

// StripBotProperties deletes any client-supplied spoofed __bot and __bot_reasons properties.
func StripBotProperties(props map[string]string) {
	if props != nil {
		delete(props, "__bot")
		delete(props, "__bot_reasons")
	}
}

// ApplyBotSuspicion annotates event properties with authoritative server-side bot flags without blocking.
func ApplyBotSuspicion(props map[string]string, opts SuspicionOptions) map[string]string {
	if props == nil {
		props = make(map[string]string)
	}

	StripBotProperties(props)

	if opts.ClientSecretAuth || opts.IsServer {
		return props
	}

	var reasons []string
	if opts.IsDatacenter {
		asnLabel := opts.ASN
		if asnLabel == "" {
			asnLabel = "unknown"
		}
		reasons = append(reasons, "datacenter_ip:AS"+asnLabel)
	}

	reasons = append(reasons, DetectHeaderAnomalies(opts.Headers, opts.UserAgent)...)
	if len(reasons) == 0 {
		return props
	}

	verdict := SummarizeSignals(reasons)
	props["__bot_reasons"] = strings.Join(verdict.Reasons, ",")
	if verdict.IsBot {
		props["__bot"] = "1"
	}

	return props
}
