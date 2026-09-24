package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Upstream Snowplow referer parser URL
const snowplowURL = "https://s3-eu-west-1.amazonaws.com/snowplow-hosted-assets/third-party/referer-parser/referers-latest.json"

type ReferrerInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Extra curated modern referrers (AI, Tech, Social, Commerce, Email, Content)
var extraReferrers = map[string]ReferrerInfo{
	"zoom.us":                 {Type: "social", Name: "Zoom"},
	"apple.com":               {Type: "tech", Name: "Apple"},
	"adobe.com":               {Type: "tech", Name: "Adobe"},
	"figma.com":               {Type: "tech", Name: "Figma"},
	"wix.com":                 {Type: "commerce", Name: "Wix"},
	"gmail.com":               {Type: "email", Name: "Gmail"},
	"notion.so":               {Type: "tech", Name: "Notion"},
	"ebay.com":                {Type: "commerce", Name: "eBay"},
	"github.com":              {Type: "tech", Name: "GitHub"},
	"gitlab.com":              {Type: "tech", Name: "GitLab"},
	"slack.com":               {Type: "social", Name: "Slack"},
	"etsy.com":                {Type: "commerce", Name: "Etsy"},
	"bsky.app":                {Type: "social", Name: "Bluesky"},
	"twitch.tv":               {Type: "content", Name: "Twitch"},
	"dropbox.com":             {Type: "tech", Name: "Dropbox"},
	"outlook.com":             {Type: "email", Name: "Outlook"},
	"medium.com":              {Type: "content", Name: "Medium"},
	"paypal.com":              {Type: "commerce", Name: "PayPal"},
	"discord.com":             {Type: "social", Name: "Discord"},
	"stripe.com":              {Type: "commerce", Name: "Stripe"},
	"spotify.com":             {Type: "content", Name: "Spotify"},
	"netflix.com":             {Type: "content", Name: "Netflix"},
	"whatsapp.com":            {Type: "social", Name: "WhatsApp"},
	"shopify.com":             {Type: "commerce", Name: "Shopify"},
	"microsoft.com":           {Type: "tech", Name: "Microsoft"},
	"alibaba.com":             {Type: "commerce", Name: "Alibaba"},
	"telegram.org":            {Type: "social", Name: "Telegram"},
	"substack.com":            {Type: "content", Name: "Substack"},
	"salesforce.com":          {Type: "tech", Name: "Salesforce"},
	"instagram.com":           {Type: "social", Name: "Instagram"},
	"wikipedia.org":           {Type: "content", Name: "Wikipedia"},
	"mastodon.social":         {Type: "social", Name: "Mastodon"},
	"office.com":              {Type: "tech", Name: "Microsoft Office"},
	"squarespace.com":         {Type: "commerce", Name: "Squarespace"},
	"stackoverflow.com":       {Type: "tech", Name: "Stack Overflow"},
	"teams.microsoft.com":     {Type: "social", Name: "Microsoft Teams"},
	"chat.com":                {Type: "ai", Name: "Chat.com"},
	"chatgpt.com":             {Type: "ai", Name: "ChatGPT"},
	"openai.com":              {Type: "ai", Name: "OpenAI"},
	"anthropic.com":           {Type: "ai", Name: "Anthropic"},
	"claude.ai":               {Type: "ai", Name: "Claude"},
	"gemini.google.com":       {Type: "ai", Name: "Google Gemini"},
	"bard.google.com":         {Type: "ai", Name: "Google Bard"},
	"copilot.microsoft.com":   {Type: "ai", Name: "Microsoft Copilot"},
	"copilot.cloud.microsoft": {Type: "ai", Name: "Microsoft Copilot"},
	"perplexity.ai":           {Type: "ai", Name: "Perplexity"},
	"you.com":                 {Type: "ai", Name: "You.com"},
	"poe.com":                 {Type: "ai", Name: "Poe"},
	"phind.com":               {Type: "ai", Name: "Phind"},
	"huggingface.co":          {Type: "ai", Name: "Hugging Face"},
	"hf.co":                   {Type: "ai", Name: "Hugging Face"},
	"character.ai":            {Type: "ai", Name: "Character.AI"},
	"meta.ai":                 {Type: "ai", Name: "Meta AI"},
	"mistral.ai":              {Type: "ai", Name: "Mistral"},
	"chat.mistral.ai":         {Type: "ai", Name: "Mistral Le Chat"},
	"deepseek.com":            {Type: "ai", Name: "DeepSeek"},
	"chat.deepseek.com":       {Type: "ai", Name: "DeepSeek Chat"},
	"pi.ai":                   {Type: "ai", Name: "Pi"},
	"inflection.ai":           {Type: "ai", Name: "Inflection"},
	"cohere.com":              {Type: "ai", Name: "Cohere"},
	"coral.cohere.com":        {Type: "ai", Name: "Cohere Coral"},
	"jasper.ai":               {Type: "ai", Name: "Jasper"},
	"writesonic.com":          {Type: "ai", Name: "Writesonic"},
	"copy.ai":                 {Type: "ai", Name: "Copy.ai"},
	"rytr.me":                 {Type: "ai", Name: "Rytr"},
	"notion.ai":               {Type: "ai", Name: "Notion AI"},
	"grammarly.com":           {Type: "ai", Name: "Grammarly"},
	"grok.com":                {Type: "ai", Name: "Grok"},
	"x.ai":                    {Type: "ai", Name: "xAI"},
	"x.com":                   {Type: "social", Name: "X (Twitter)"},
	"fb.com":                  {Type: "social", Name: "Facebook"},
	"youtu.be":                {Type: "social", Name: "YouTube"},
	"youtube.com":             {Type: "social", Name: "YouTube"},
	"aistudio.google.com":     {Type: "ai", Name: "Google AI Studio"},
	"labs.google.com":         {Type: "ai", Name: "Google Labs"},
	"ai.google":               {Type: "ai", Name: "Google AI"},
	"forefront.ai":            {Type: "ai", Name: "Forefront"},
	"together.ai":             {Type: "ai", Name: "Together AI"},
	"groq.com":                {Type: "ai", Name: "Groq"},
	"replicate.com":           {Type: "ai", Name: "Replicate"},
	"vercel.ai":               {Type: "ai", Name: "Vercel AI"},
	"v0.dev":                  {Type: "ai", Name: "v0"},
	"bolt.new":                {Type: "ai", Name: "Bolt"},
	"replit.com":              {Type: "ai", Name: "Replit"},
	"cursor.com":              {Type: "ai", Name: "Cursor"},
	"tabnine.com":             {Type: "ai", Name: "Tabnine"},
	"codeium.com":             {Type: "ai", Name: "Codeium"},
	"sourcegraph.com":         {Type: "ai", Name: "Sourcegraph Cody"},
	"kimi.moonshot.cn":        {Type: "ai", Name: "Kimi"},
	"moonshot.ai":             {Type: "ai", Name: "Moonshot AI"},
	"doubao.com":              {Type: "ai", Name: "Doubao"},
	"tongyi.aliyun.com":       {Type: "ai", Name: "Tongyi Qianwen"},
	"yiyan.baidu.com":         {Type: "ai", Name: "Ernie Bot"},
	"chatglm.cn":              {Type: "ai", Name: "ChatGLM"},
	"zhipu.ai":                {Type: "ai", Name: "Zhipu AI"},
	"minimax.chat":            {Type: "ai", Name: "MiniMax"},
	"lmsys.org":               {Type: "ai", Name: "LMSYS"},
	"chat.lmsys.org":          {Type: "ai", Name: "LMSYS Chat"},
	"llama.meta.com":          {Type: "ai", Name: "Meta Llama"},
}

func main() {
	log.Println("Starting referrer data generator...")

	// Fetch from Snowplow or read from fallback file
	data, err := fetchSnowplowData()
	if err != nil {
		log.Printf("Warning: failed to fetch Snowplow data over network: %v", err)
		log.Println("Using fallback to local OpenPanel index.ts if available...")
		data = readLocalFallback()
	}

	merged := make(map[string]ReferrerInfo)
	for domain, entry := range data {
		merged[domain] = entry
	}
	for domain, entry := range extraReferrers {
		merged[domain] = entry
	}

	log.Printf("Total parsed referrers: %d domains", len(merged))

	// Generate Go file
	outPath := filepath.Join("referrer", "data_generated.go")
	if err := generateGoFile(outPath, merged); err != nil {
		log.Fatalf("Failed to generate %s: %v", outPath, err)
	}

	log.Printf("Successfully generated %s", outPath)
}

func fetchSnowplowData() (map[string]ReferrerInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(snowplowURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw map[string]map[string]struct {
		Domains []string `json:"domains"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	result := make(map[string]ReferrerInfo)
	for refType, names := range raw {
		for name, item := range names {
			for _, domain := range item.Domains {
				result[domain] = ReferrerInfo{
					Type: refType,
					Name: name,
				}
			}
		}
	}

	return result, nil
}

func readLocalFallback() map[string]ReferrerInfo {
	fallbackPath := "/Users/rakib/Projects/analytics/openpanel/packages/common/server/referrers/index.ts"
	content, err := os.ReadFile(fallbackPath)
	if err != nil {
		log.Printf("Fallback read failed: %v", err)
		return nil
	}

	result := make(map[string]ReferrerInfo)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "'") {
			continue
		}
		// 'support.google.com': { type: 'unknown', name: 'Google' },
		parts := strings.SplitN(line, ": {", 2)
		if len(parts) != 2 {
			continue
		}
		domain := strings.Trim(parts[0], "' \t:")

		typePart := ""
		namePart := ""
		if idx := strings.Index(parts[1], "type: '"); idx != -1 {
			rem := parts[1][idx+len("type: '"):]
			if end := strings.Index(rem, "'"); end != -1 {
				typePart = rem[:end]
			}
		}
		if idx := strings.Index(parts[1], "name: '"); idx != -1 {
			rem := parts[1][idx+len("name: '"):]
			if end := strings.Index(rem, "'"); end != -1 {
				namePart = rem[:end]
			}
		}

		if domain != "" && namePart != "" {
			result[domain] = ReferrerInfo{
				Type: typePart,
				Name: namePart,
			}
		}
	}
	return result
}

func generateGoFile(filePath string, referrers map[string]ReferrerInfo) error {
	var buf bytes.Buffer

	buf.WriteString("// Code generated by cmd/generator/main.go; DO NOT EDIT.\n")
	buf.WriteString(fmt.Sprintf("// Generated on: %s\n\n", time.Now().UTC().Format(time.RFC3339)))
	buf.WriteString("package referrer\n\n")
	buf.WriteString("import \"strings\"\n\n")

	// Domain map
	buf.WriteString("var referrers = map[string]ReferrerEntry{\n")

	// Sort keys for deterministic output
	domains := make([]string, 0, len(referrers))
	for d := range referrers {
		domains = append(domains, d)
	}
	sort.Strings(domains)

	for _, d := range domains {
		entry := referrers[d]
		escapedDomain := strings.ReplaceAll(d, `"`, `\"`)
		escapedName := strings.ReplaceAll(entry.Name, `"`, `\"`)
		buf.WriteString(fmt.Sprintf("\t%q: {Type: %q, Name: %q},\n", escapedDomain, entry.Type, escapedName))
	}
	buf.WriteString("}\n\n")

	// Name lookup map
	buf.WriteString("var referrersByName = map[string]struct {\n\tEntry  ReferrerEntry\n\tDomain string\n}{\n")

	// Aggregate unique names
	nameMap := make(map[string]struct {
		entry  ReferrerInfo
		domain string
	})
	for _, d := range domains {
		entry := referrers[d]
		lowerName := strings.ToLower(entry.Name)
		if _, exists := nameMap[lowerName]; !exists {
			nameMap[lowerName] = struct {
				entry  ReferrerInfo
				domain string
			}{entry: entry, domain: d}
		}
	}

	var names []string
	for n := range nameMap {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, n := range names {
		val := nameMap[n]
		escapedName := strings.ReplaceAll(val.entry.Name, `"`, `\"`)
		escapedDomain := strings.ReplaceAll(val.domain, `"`, `\"`)
		buf.WriteString(fmt.Sprintf("\t%q: {Entry: ReferrerEntry{Type: %q, Name: %q}, Domain: %q},\n",
			n, val.entry.Type, escapedName, escapedDomain))
	}
	buf.WriteString("}\n\n")

	// Accessor functions
	buf.WriteString(`// Lookup finds a referrer entry by hostname/domain.
func Lookup(domain string) (ReferrerEntry, bool) {
	entry, ok := referrers[strings.ToLower(domain)]
	return entry, ok
}

// LookupByName finds a referrer entry by service name (e.g. "chatgpt", "google", "github").
func LookupByName(name string) (ReferrerEntry, string, bool) {
	val, ok := referrersByName[strings.ToLower(name)]
	return val.Entry, val.Domain, ok
}
`)

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Write unformatted for debugging if syntax error
		_ = os.WriteFile(filePath, buf.Bytes(), 0644)
		return fmt.Errorf("gofmt error: %w", err)
	}

	return os.WriteFile(filePath, formatted, 0644)
}
