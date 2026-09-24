package uaparser

import (
	"net/http"
	"strings"
)

// Header constants for User-Agent Client Hints.
const (
	HeaderSecCHUA             = "Sec-CH-UA"
	HeaderSecCHUAFullVersion  = "Sec-CH-UA-Full-Version-List"
	HeaderSecCHUAMobile       = "Sec-CH-UA-Mobile"
	HeaderSecCHUAModel        = "Sec-CH-UA-Model"
	HeaderSecCHUAPlatform     = "Sec-CH-UA-Platform"
	HeaderSecCHUAPlatformVer  = "Sec-CH-UA-Platform-Version"
	HeaderSecCHUAArch         = "Sec-CH-UA-Arch"
	HeaderSecCHUABitness      = "Sec-CH-UA-Bitness"
	HeaderSecCHUAFormFactors  = "Sec-CH-UA-Form-Factors"
)

// ParseClientHints extracts ClientHints from standard http.Header.
func ParseClientHints(h http.Header) ClientHints {
	if h == nil {
		return ClientHints{}
	}

	var ch ClientHints
	if val := h.Get(HeaderSecCHUA); val != "" {
		ch.Brands = parseBrandsHeader(val)
	}
	if val := h.Get(HeaderSecCHUAFullVersion); val != "" {
		ch.FullVersionList = parseBrandsHeader(val)
	}
	if val := h.Get(HeaderSecCHUAMobile); val != "" {
		isMobile := strings.Contains(val, "?1")
		ch.Mobile = &isMobile
	}
	if val := h.Get(HeaderSecCHUAModel); val != "" {
		ch.Model = cleanHeaderValue(val)
	}
	if val := h.Get(HeaderSecCHUAPlatform); val != "" {
		ch.Platform = cleanHeaderValue(val)
	}
	if val := h.Get(HeaderSecCHUAPlatformVer); val != "" {
		ch.PlatformVersion = cleanHeaderValue(val)
	}
	if val := h.Get(HeaderSecCHUAArch); val != "" {
		ch.Architecture = cleanHeaderValue(val)
	}
	if val := h.Get(HeaderSecCHUABitness); val != "" {
		ch.Bitness = cleanHeaderValue(val)
	}
	if val := h.Get(HeaderSecCHUAFormFactors); val != "" {
		ch.FormFactors = parseListHeader(val)
	}

	return ch
}

// ParseClientHintsMap extracts ClientHints from a key-value header map (case-insensitive).
func ParseClientHintsMap(m map[string]string) ClientHints {
	if m == nil {
		return ClientHints{}
	}
	h := make(http.Header, len(m))
	for k, v := range m {
		h.Set(k, v)
	}
	return ParseClientHints(h)
}

func cleanHeaderValue(val string) string {
	val = strings.TrimSpace(val)
	val = strings.Trim(val, `"`)
	return strings.TrimSpace(val)
}

func parseBrandsHeader(headerVal string) []Brand {
	if headerVal == "" {
		return nil
	}

	var brands []Brand
	tokens := strings.Split(headerVal, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		parts := strings.Split(token, ";")
		brandName := cleanHeaderValue(parts[0])

		version := ""
		for _, part := range parts[1:] {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "v=") {
				version = cleanHeaderValue(strings.TrimPrefix(part, "v="))
			}
		}

		brands = append(brands, Brand{
			Brand:   brandName,
			Version: version,
		})
	}
	return brands
}

func parseListHeader(headerVal string) []string {
	if headerVal == "" {
		return nil
	}
	tokens := strings.Split(headerVal, ",")
	res := make([]string, 0, len(tokens))
	for _, t := range tokens {
		t = cleanHeaderValue(t)
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}
