package referrer

import (
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		url          string
		expectedName string
		expectedType ReferrerType
	}{
		{"", "", TypeDirect},
		{"https://www.google.com/search?q=openpanel", "Google", TypeSearch},
		{"https://google.de/search?q=test", "Google", TypeSearch},
		{"https://google.com/", "Google", TypeSearch},
		{"https://chatgpt.com/c/123", "ChatGPT", TypeAI},
		{"https://claude.ai/chat/456", "Claude", TypeAI},
		{"https://github.com/rakibhoossain/ua-parser-go", "GitHub", TypeTech},
		{"https://twitter.com/rakib", "Twitter", TypeSocial},
		{"https://x.com/rakib", "X (Twitter)", TypeSocial},
		{"https://t.co/xyz123", "Twitter", TypeSocial},
		{"https://fb.com/page", "Facebook", TypeSocial},
		{"https://youtu.be/video123", "YouTube", TypeSocial},
		{"https://unknown-domain-12345.xyz/test", "unknown-domain-12345.xyz", TypeReferral},
	}

	for _, tt := range tests {
		res := Parse(tt.url)
		if tt.expectedName != "" && res.Name != tt.expectedName {
			t.Errorf("URL %q: expected Name %q, got %q", tt.url, tt.expectedName, res.Name)
		}
		if tt.expectedType != "" && res.Type != tt.expectedType {
			t.Errorf("URL %q: expected Type %q, got %q", tt.url, tt.expectedType, res.Type)
		}
		if tt.expectedName != "" && res.FaviconURL == "" {
			t.Errorf("URL %q: expected non-empty FaviconURL", tt.url)
		}
	}
}

func TestParseWithQuery(t *testing.T) {
	tests := []struct {
		query        map[string]string
		expectedName string
		expectedType ReferrerType
		expectMedium string
	}{
		{map[string]string{"utm_source": "google"}, "Google", TypeSearch, ""},
		{map[string]string{"ref": "twitter"}, "Twitter", TypeSocial, ""},
		{map[string]string{"utm_referrer": "chatgpt"}, "ChatGPT", TypeAI, ""},
		{map[string]string{"utm_source": "newsletter", "utm_medium": "email"}, "newsletter", TypeEmail, "email"},
		{map[string]string{"utm_source": "google", "utm_medium": "cpc", "utm_campaign": "brand"}, "Google", TypePaid, "cpc"},
		{map[string]string{"utm_source": "facebook", "utm_medium": "social"}, "Facebook", TypeSocial, "social"},
	}

	for _, tt := range tests {
		res := ParseWithQuery(tt.query)
		if res == nil {
			t.Fatalf("query %v: unexpected nil result", tt.query)
		}
		if res.Name != tt.expectedName {
			t.Errorf("query %v: expected Name %q, got %q", tt.query, tt.expectedName, res.Name)
		}
		if tt.expectedType != "" && res.Type != tt.expectedType {
			t.Errorf("query %v: expected Type %q, got %q", tt.query, tt.expectedType, res.Type)
		}
		if tt.expectMedium != "" && res.Medium != tt.expectMedium {
			t.Errorf("query %v: expected Medium %q, got %q", tt.query, tt.expectMedium, res.Medium)
		}
	}
}

func TestParseRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://example.com/pricing?utm_source=chatgpt&utm_medium=ai_referral", nil)
	req.Header.Set("Referer", "https://chatgpt.com/")

	res := ParseRequest(req)
	if res.Name != "ChatGPT" {
		t.Fatalf("expected ChatGPT, got %s", res.Name)
	}
	if res.Type != TypeAI {
		t.Fatalf("expected AI, got %s", res.Type)
	}
	if res.Medium != "ai_referral" {
		t.Fatalf("expected Medium ai_referral, got %s", res.Medium)
	}
}

func TestFaviconURL(t *testing.T) {
	fav := FaviconURL("https://github.com/something")
	expected := "https://icons.duckduckgo.com/ip3/github.com.ico"
	if fav != expected {
		t.Fatalf("expected %s, got %s", expected, fav)
	}

	gFav := GoogleFaviconURL("github.com", 32)
	expectedG := "https://www.google.com/s2/favicons?domain=github.com&sz=32"
	if gFav != expectedG {
		t.Fatalf("expected %s, got %s", expectedG, gFav)
	}
}
