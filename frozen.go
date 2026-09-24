package uaparser

import "regexp"

// frozenUARegex matches User-Agent strings produced by Chromium's User-Agent Reduction initiative.
// Under UA reduction, browsers report a frozen OS version, fixed device tokens, and a major version followed by ".0.0.0".
var frozenUARegex = regexp.MustCompile(`^Mozilla/5\.0 \((Windows NT 10\.0; Win64; x64|Macintosh; Intel Mac OS X 10_15_7|X11; Linux x86_64|X11; CrOS x86_64 14541\.0\.0|Fuchsia|Linux; Android 10; K)\) AppleWebKit/537\.36 \(KHTML, like Gecko\) Chrome/\d+\.0\.0\.0 (Mobile )?Safari/537\.36`)

// IsFrozenUA checks if a User-Agent string is a frozen Chromium User-Agent.
// When frozen, fine-grained client information (OS version, device model, architecture)
// should be retrieved from Client Hints (Sec-CH-UA-*).
func IsFrozenUA(ua string) bool {
	return frozenUARegex.MatchString(ua)
}
