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
	"regexp"
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
	log.Println("=== Starting ua-parser-go data generators ===")

	// 1. Generate Referrers
	log.Println("Generating referrer data...")
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
	refOutPath := filepath.Join("referrer", "data_generated.go")
	if err := generateGoFile(refOutPath, merged); err != nil {
		log.Fatalf("Failed to generate %s: %v", refOutPath, err)
	}
	log.Printf("Successfully generated %s", refOutPath)

	// 2. Generate Bots
	log.Println("Generating bot data...")
	botEntries, err := fetchMatomoBots()
	if err != nil {
		log.Printf("Warning: failed to fetch Matomo bots over network: %v", err)
		log.Println("Using fallback to local OpenPanel bots.ts if available...")
		botEntries, err = readLocalBotsFallback()
		if err != nil {
			log.Fatalf("Failed to read local bots fallback: %v", err)
		}
	}

	log.Printf("Total parsed bot entries: %d", len(botEntries))
	botsOutPath := filepath.Join("bots", "data_generated.go")
	if err := generateBotsFile(botsOutPath, botEntries); err != nil {
		log.Fatalf("Failed to generate %s: %v", botsOutPath, err)
	}
	log.Printf("Successfully generated %s", botsOutPath)
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

const matomoBotsURL = "https://raw.githubusercontent.com/matomo-org/device-detector/master/regexes/bots.yml"

type BotDataEntry struct {
	Includes string
	Regex    string
	Name     string
	Category string
	URL      string
	Producer string
}

var allowlistedBotTokens = map[string]bool{
	"node":      true,
	"Node\\.js": true,
	"Node.js":   true,
}

var (
	anchoredGroupRegex = regexp.MustCompile(`\^\(\?:([^)]+)\)\$(\|?)`)
	regexSpecialChars  = regexp.MustCompile(`[|^$.*+?(){}\[\]\\]`)
)

func stripAllowlistedTokens(r string) string {
	return anchoredGroupRegex.ReplaceAllStringFunc(r, func(match string) string {
		trailingPipe := ""
		if strings.HasSuffix(match, "|") {
			trailingPipe = "|"
		}
		sub := anchoredGroupRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		inner := sub[1]
		alts := strings.Split(inner, "|")
		var kept []string
		for _, alt := range alts {
			if !allowlistedBotTokens[alt] {
				kept = append(kept, alt)
			}
		}
		if len(kept) == 0 {
			return ""
		}
		return fmt.Sprintf("^(?:%s)$%s", strings.Join(kept, "|"), trailingPipe)
	})
}

func fetchMatomoBots() ([]BotDataEntry, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(matomoBotsURL)
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

	return parseMatomoYAML(string(body))
}

func parseMatomoYAML(content string) ([]BotDataEntry, error) {
	lines := strings.Split(content, "\n")
	var entries []BotDataEntry
	var current *BotDataEntry
	inProducer := false

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if strings.HasPrefix(line, "- regex:") {
			if current != nil {
				entries = append(entries, transformBotEntry(*current))
			}
			current = &BotDataEntry{
				Regex: strings.Trim(strings.TrimPrefix(line, "- regex:"), " '\""),
			}
			inProducer = false
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, "producer:") {
			inProducer = true
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), " '\"")

		switch key {
		case "name":
			if inProducer {
				current.Producer = val
			} else {
				current.Name = val
			}
		case "category":
			current.Category = val
		case "url":
			if !inProducer {
				current.URL = val
			}
		}
	}

	if current != nil {
		entries = append(entries, transformBotEntry(*current))
	}

	return entries, nil
}

func readLocalBotsFallback() ([]BotDataEntry, error) {
	filePath := "/Users/rakib/Projects/analytics/openpanel/apps/api/src/bots/bots.ts"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var entries []BotDataEntry
	var current *BotDataEntry
	inProducer := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "{" {
			current = &BotDataEntry{}
			inProducer = false
			continue
		}
		if current == nil {
			continue
		}

		if line == "}," || line == "}" {
			if inProducer {
				inProducer = false
				continue
			}
			entries = append(entries, *current)
			current = nil
			continue
		}

		if strings.HasPrefix(line, "producer:") {
			inProducer = true
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if (val == "" || val == "'" || val == "\"") && i+1 < len(lines) {
			i++
			val = strings.TrimSpace(lines[i])
		}

		val = strings.TrimSuffix(val, ",")
		val = strings.Trim(val, "'\"`")

		switch key {
		case "includes":
			current.Includes = val
		case "regex":
			current.Regex = val
		case "name":
			if inProducer {
				current.Producer = val
			} else {
				current.Name = val
			}
		case "category":
			current.Category = val
		case "url":
			if !inProducer {
				current.URL = val
			}
		}
	}

	return entries, nil
}

func transformBotEntry(entry BotDataEntry) BotDataEntry {
	cleaned := stripAllowlistedTokens(entry.Regex)
	if regexSpecialChars.MatchString(cleaned) {
		entry.Regex = cleaned
		entry.Includes = ""
	} else {
		entry.Includes = cleaned
		entry.Regex = ""
	}
	return entry
}

func generateBotsFile(filePath string, entries []BotDataEntry) error {
	var buf bytes.Buffer

	buf.WriteString("// Code generated by cmd/generator/main.go; DO NOT EDIT.\n")
	buf.WriteString(fmt.Sprintf("// Generated on: %s\n\n", time.Now().UTC().Format(time.RFC3339)))
	buf.WriteString("package bots\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"regexp\"\n")
	buf.WriteString("\t\"strings\"\n")
	buf.WriteString(")\n\n")

	// Split into includes and regexes
	var includesList []BotDataEntry
	var regexList []BotDataEntry

	for _, e := range entries {
		if e.Includes != "" {
			includesList = append(includesList, e)
		} else if e.Regex != "" {
			regexList = append(regexList, e)
		}
	}

	// 1. Write includeBot struct and slice
	buf.WriteString(`type includeBot struct {
	includes string
	name     string
	category string
	url      string
	producer string
}

type regexBot struct {
	pattern  string
	re       *regexp.Regexp
	name     string
	category string
	url      string
	producer string
	special  string
}

`)

	buf.WriteString(fmt.Sprintf("// includesBots contains %d literal substring bot definitions (evaluated first).\n", len(includesList)))
	buf.WriteString("var includesBots = [...]includeBot{\n")
	for _, b := range includesList {
		buf.WriteString(fmt.Sprintf("\t{includes: %q, name: %q, category: %q, url: %q, producer: %q},\n",
			b.Includes, b.Name, b.Category, b.URL, b.Producer))
	}
	buf.WriteString("}\n\n")

	// 2. Write regexBots slice with Go RE2 compatibility
	buf.WriteString(fmt.Sprintf("// regexBots contains %d compiled regex bot patterns.\n", len(regexList)))
	buf.WriteString("var regexBots = [...]regexBot{\n")
	for _, b := range regexList {
		pattern := b.Regex
		special := ""

		// Sanitize lookaround assertions for RE2
		if strings.Contains(pattern, "(?<!HTC)[ _]Butterfly/") {
			pattern = strings.ReplaceAll(pattern, "(?<!HTC)[ _]Butterfly/", "[ _]Butterfly/")
			special = "butterfly"
		} else if strings.Contains(pattern, "Daum(?!(?:Apps|Device))") {
			pattern = strings.ReplaceAll(pattern, "Daum(?!(?:Apps|Device))", "Daum")
			special = "daum"
		} else if strings.Contains(pattern, "zeal(?!ot)") {
			pattern = strings.ReplaceAll(pattern, "zeal(?!ot)", "zeal")
			special = "zealot"
		} else if strings.Contains(pattern, "(?<!cu|Hu|") || strings.Contains(pattern, "(?<!node-") {
			pattern = `[a-z0-9_-]*(?:bot|analyzer|appengine|archiver?|checker|collector|crawl|crawler|fetch(?:er)?|grabber|indexer|inspector|monitor|^parser|project|proxy|research|resolver|robots|scanner|scraper|script|searcher|security|spider|study|transcoder|uptime|user[ _]?agent|validator|-(?:AI|Extended|User)/)(?:[^a-z]|$)`
			special = "heuristic"
		}

		buf.WriteString(fmt.Sprintf("\t{pattern: %q, re: regexp.MustCompile(%q), name: %q, category: %q, url: %q, producer: %q, special: %q},\n",
			pattern, pattern, b.Name, b.Category, b.URL, b.Producer, special))
	}
	buf.WriteString("}\n\n")

	// 3. Helper for negative lookaround assertions
	buf.WriteString(`func isSpecialExcluded(special string, ua string) bool {
	switch special {
	case "butterfly":
		return strings.Contains(ua, "HTC")
	case "daum":
		return strings.Contains(ua, "DaumApps") || strings.Contains(ua, "DaumDevice")
	case "zealot":
		return strings.Contains(ua, "zealot")
	case "heuristic":
		return !isGenericHeuristicValid(ua)
	default:
		return false
	}
}

func isGenericHeuristicValid(ua string) bool {
	lower := strings.ToLower(ua)

	// Exclude server-side fetch libraries and browsers
	if strings.Contains(lower, "node-fetch") || strings.Contains(lower, "uclient-fetch") || strings.Contains(lower, "electron-fetch") {
		return false
	}
	if strings.Contains(lower, "urlgrabber") {
		return false
	}
	if strings.Contains(lower, "projector") || strings.Contains(lower, "microsoft project") || strings.Contains(lower, "banshee-project") {
		return false
	}
	if strings.Contains(lower, "camscanner") {
		return false
	}
	if strings.Contains(lower, "presearch") {
		return false
	}

	// Bot exceptions: cubot, Hubot, power bot, m bot, etc.
	if strings.Contains(lower, "cubot") || strings.Contains(lower, "hubot") ||
		strings.Contains(lower, "power bot") || strings.Contains(lower, "power_bot") ||
		strings.Contains(lower, "m bot") || strings.Contains(lower, "m_bot") {
		return false
	}
	if strings.Contains(lower, "bot tab") || strings.Contains(lower, "bot_tab") ||
		strings.Contains(lower, "bot senior") || strings.Contains(lower, "bot_senior") ||
		strings.Contains(lower, "bot junior") || strings.Contains(lower, "bot_junior") {
		return false
	}

	return true
}

// DetectBot evaluates the User-Agent against all known bot substring and regex signatures.
func DetectBot(ua string) *BotMatch {
	if ua == "" {
		return nil
	}

	// 1. Fast literal substring matching
	for i := range includesBots {
		if strings.Contains(ua, includesBots[i].includes) {
			cat := includesBots[i].category
			if cat == "" {
				cat = "Unknown"
			}
			return &BotMatch{
				Name:     includesBots[i].name,
				Category: cat,
				URL:      includesBots[i].url,
				Producer: includesBots[i].producer,
			}
		}
	}

	// 2. Pre-compiled regex patterns
	for i := range regexBots {
		rb := &regexBots[i]
		if rb.re.MatchString(ua) {
			if rb.special != "" && isSpecialExcluded(rb.special, ua) {
				continue
			}

			name := rb.name
			if strings.Contains(name, "$1") {
				sub := rb.re.FindStringSubmatch(ua)
				if len(sub) > 1 && sub[1] != "" {
					name = sub[1]
				}
			}

			cat := rb.category
			if cat == "" {
				cat = "Unknown"
			}
			return &BotMatch{
				Name:     name,
				Category: cat,
				URL:      rb.url,
				Producer: rb.producer,
			}
		}
	}

	return nil
}

// IsBot reports whether the User-Agent represents a known automated bot or crawler.
func IsBot(ua string) bool {
	return DetectBot(ua) != nil
}

// IsCrawler reports whether the User-Agent represents a known search or SEO spider/crawler.
func IsCrawler(ua string) bool {
	b := DetectBot(ua)
	if b == nil {
		return false
	}
	cat := strings.ToLower(b.Category)
	return cat == "crawler" || cat == "search bot" || cat == "feed fetcher" ||
		strings.Contains(strings.ToLower(b.Name), "crawler") ||
		strings.Contains(strings.ToLower(b.Name), "spider") ||
		strings.Contains(strings.ToLower(b.Name), "bot")
}

// IsSearchBot reports whether the User-Agent is explicitly a search engine crawler.
func IsSearchBot(ua string) bool {
	b := DetectBot(ua)
	if b == nil {
		return false
	}
	return strings.EqualFold(b.Category, "search bot")
}
`)

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		_ = os.WriteFile(filePath, buf.Bytes(), 0644)
		return fmt.Errorf("gofmt error: %w", err)
	}

	return os.WriteFile(filePath, formatted, 0644)
}
