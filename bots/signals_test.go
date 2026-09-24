package bots

import (
	"net/http"
	"strings"
	"testing"
)

const (
	uaChromeDesktop = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	uaChromeAndroid = "Mozilla/5.0 (Linux; Android 13) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"
	uaEdgeDesktop   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"
	uaFirefox       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"
	uaSafariMac     = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15"
	uaChromeIOS     = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.0.0 Mobile/15E148 Safari/604.1"
	uaCurl          = "curl/8.4.0"
	uaGoClient      = "Go-http-client/1.1"
)

func newRealBrowserHeaders() http.Header {
	h := make(http.Header)
	h.Set("Sec-CH-UA", `"Chromium";v="120"`)
	h.Set("Sec-Fetch-Mode", "cors")
	h.Set("Sec-Fetch-Site", "cross-site")
	h.Set("Accept-Language", "en-US,en;q=0.9")
	return h
}

func TestDetectHeaderAnomalies(t *testing.T) {
	t.Run("does not flag legitimate browser traffic", func(t *testing.T) {
		// Clean desktop Chrome -> no reasons
		h := newRealBrowserHeaders()
		reasons := DetectHeaderAnomalies(h, uaChromeDesktop)
		if len(reasons) != 0 {
			t.Errorf("expected 0 reasons for clean Chrome, got %v", reasons)
		}

		// Firefox (no sec-ch-ua, but sends everything else) -> no reasons
		ffHeaders := make(http.Header)
		ffHeaders.Set("Sec-Fetch-Mode", "cors")
		ffHeaders.Set("Sec-Fetch-Site", "cross-site")
		ffHeaders.Set("Accept-Language", "en-US,en;q=0.5")
		reasons = DetectHeaderAnomalies(ffHeaders, uaFirefox)
		if len(reasons) != 0 {
			t.Errorf("expected 0 reasons for Firefox, got %v", reasons)
		}

		// Safari (no sec-ch-ua, WebKit) -> no reasons
		safariHeaders := make(http.Header)
		safariHeaders.Set("Sec-Fetch-Mode", "cors")
		safariHeaders.Set("Sec-Fetch-Site", "cross-site")
		safariHeaders.Set("Accept-Language", "en-US,en;q=0.9")
		reasons = DetectHeaderAnomalies(safariHeaders, uaSafariMac)
		if len(reasons) != 0 {
			t.Errorf("expected 0 reasons for Safari, got %v", reasons)
		}

		// iOS Chrome / CriOS (WebKit, no sec-ch-ua) -> not flagged for sec-ch-ua
		criosHeaders := make(http.Header)
		criosHeaders.Set("Sec-Fetch-Mode", "cors")
		criosHeaders.Set("Sec-Fetch-Site", "cross-site")
		criosHeaders.Set("Accept-Language", "en-US,en;q=0.9")
		reasons = DetectHeaderAnomalies(criosHeaders, uaChromeIOS)
		if len(reasons) != 0 {
			t.Errorf("expected 0 reasons for iOS Chrome, got %v", reasons)
		}
	})

	t.Run("does not apply to non-browser clients", func(t *testing.T) {
		for _, ua := range []string{uaCurl, uaGoClient, ""} {
			reasons := DetectHeaderAnomalies(make(http.Header), ua)
			if len(reasons) != 0 {
				t.Errorf("expected 0 reasons for non-browser UA %q, got %v", ua, reasons)
			}
		}
	})

	t.Run("flags spoofed Chromium user agents", func(t *testing.T) {
		// Chrome UA with no sec-ch-ua
		h := make(http.Header)
		h.Set("Sec-Fetch-Mode", "cors")
		h.Set("Sec-Fetch-Site", "cross-site")
		h.Set("Accept-Language", "en-US")

		reasons := DetectHeaderAnomalies(h, uaChromeDesktop)
		if !contains(reasons, "header:missing_sec_ch_ua") {
			t.Errorf("expected missing_sec_ch_ua, got %v", reasons)
		}

		reasons = DetectHeaderAnomalies(h, uaEdgeDesktop)
		if !contains(reasons, "header:missing_sec_ch_ua") {
			t.Errorf("expected missing_sec_ch_ua for Edge, got %v", reasons)
		}

		reasons = DetectHeaderAnomalies(h, uaChromeAndroid)
		if !contains(reasons, "header:missing_sec_ch_ua") {
			t.Errorf("expected missing_sec_ch_ua for Android Chrome, got %v", reasons)
		}

		// Bare Chrome UA with no browser headers at all -> multiple reasons
		reasons = DetectHeaderAnomalies(make(http.Header), uaChromeDesktop)
		if !contains(reasons, "header:missing_sec_ch_ua") ||
			!contains(reasons, "header:missing_sec_fetch") ||
			!contains(reasons, "header:missing_accept_language") {
			t.Errorf("expected multiple header reasons for bare Chrome UA, got %v", reasons)
		}
	})

	t.Run("individual header rules", func(t *testing.T) {
		// Missing both sec-fetch headers
		h := make(http.Header)
		h.Set("Accept-Language", "en-US")
		reasons := DetectHeaderAnomalies(h, uaFirefox)
		if !contains(reasons, "header:missing_sec_fetch") {
			t.Errorf("expected missing_sec_fetch, got %v", reasons)
		}

		// One sec-fetch header present -> no missing_sec_fetch
		h.Set("Sec-Fetch-Site", "cross-site")
		reasons = DetectHeaderAnomalies(h, uaFirefox)
		if contains(reasons, "header:missing_sec_fetch") {
			t.Errorf("did not expect missing_sec_fetch when Sec-Fetch-Site present, got %v", reasons)
		}
	})
}

func TestSummarizeSignals(t *testing.T) {
	t.Run("no signals", func(t *testing.T) {
		res := SummarizeSignals(nil)
		if res.IsBot || len(res.Reasons) != 0 {
			t.Errorf("expected isBot=false and 0 reasons, got %+v", res)
		}
	})

	t.Run("single category (datacenter only) -> not flagged", func(t *testing.T) {
		res := SummarizeSignals([]string{"datacenter_ip:AS15169"})
		if res.IsBot {
			t.Errorf("expected isBot=false for single category, got true")
		}
		if len(res.Reasons) != 1 {
			t.Errorf("expected 1 reason, got %v", res.Reasons)
		}
	})

	t.Run("multiple reasons in SAME category -> still one category, not flagged", func(t *testing.T) {
		res := SummarizeSignals([]string{
			"header:missing_sec_ch_ua",
			"header:missing_sec_fetch",
			"header:missing_accept_language",
		})
		if res.IsBot {
			t.Errorf("expected isBot=false for same category, got true")
		}
	})

	t.Run("two distinct categories -> flagged", func(t *testing.T) {
		res := SummarizeSignals([]string{
			"datacenter_ip:AS16509",
			"header:missing_sec_ch_ua",
		})
		if !res.IsBot {
			t.Errorf("expected isBot=true for 2 distinct categories, got false")
		}
	})
}

func TestApplyBotSuspicion(t *testing.T) {
	t.Run("strips client-supplied bot properties", func(t *testing.T) {
		props := map[string]string{
			"__bot":         "1",
			"__bot_reasons": "forged",
			"keep":          "me",
		}
		StripBotProperties(props)
		if _, exists := props["__bot"]; exists {
			t.Errorf("expected __bot to be stripped")
		}
		if _, exists := props["__bot_reasons"]; exists {
			t.Errorf("expected __bot_reasons to be stripped")
		}
		if props["keep"] != "me" {
			t.Errorf("expected keep=me preserved")
		}
	})

	t.Run("datacenter IP alone -> recorded but NOT flagged", func(t *testing.T) {
		props := map[string]string{"foo": "bar"}
		res := ApplyBotSuspicion(props, SuspicionOptions{
			IsDatacenter: true,
			ASN:          "16509",
			Headers:      newRealBrowserHeaders(),
			UserAgent:    uaChromeDesktop,
			IsServer:     false,
		})
		if res["__bot_reasons"] != "datacenter_ip:AS16509" {
			t.Errorf("expected datacenter_ip reason, got %q", res["__bot_reasons"])
		}
		if _, exists := res["__bot"]; exists {
			t.Errorf("expected __bot not to be set")
		}
	})

	t.Run("datacenter IP + header anomaly -> flagged __bot=1", func(t *testing.T) {
		props := map[string]string{}
		spoofedHeaders := make(http.Header)
		res := ApplyBotSuspicion(props, SuspicionOptions{
			IsDatacenter: true,
			ASN:          "16509",
			Headers:      spoofedHeaders,
			UserAgent:    uaChromeDesktop,
			IsServer:     false,
		})
		if res["__bot"] != "1" {
			t.Errorf("expected __bot='1', got %q", res["__bot"])
		}
		if !strings.Contains(res["__bot_reasons"], "datacenter_ip:AS16509") ||
			!strings.Contains(res["__bot_reasons"], "header:missing_sec_ch_ua") {
			t.Errorf("expected both datacenter and header reasons, got %q", res["__bot_reasons"])
		}
	})

	t.Run("server-side (ClientSecretAuth or IsServer) never flagged", func(t *testing.T) {
		props := map[string]string{}
		res := ApplyBotSuspicion(props, SuspicionOptions{
			IsDatacenter:     true,
			ASN:              "16509",
			Headers:          make(http.Header),
			UserAgent:        uaChromeDesktop,
			ClientSecretAuth: true,
		})
		if len(res) != 0 {
			t.Errorf("expected clean props for clientSecretAuth, got %v", res)
		}

		res = ApplyBotSuspicion(props, SuspicionOptions{
			IsDatacenter: true,
			ASN:          "16509",
			Headers:      make(http.Header),
			UserAgent:    "Go-http-client/1.1",
			IsServer:     true,
		})
		if len(res) != 0 {
			t.Errorf("expected clean props for isServer, got %v", res)
		}
	})
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
