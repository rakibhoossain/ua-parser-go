package uaparser

// European regulatory and geographic timezone sets.
// Based on IANA timezones corresponding to EU, EEA, EFTA, and Schengen territories.

var euTimezones = map[string]struct{}{
	"Europe/Vienna":         {},
	"Europe/Brussels":       {},
	"Europe/Sofia":          {},
	"Europe/Zagreb":         {},
	"Europe/Nicosia":        {},
	"Asia/Nicosia":          {},
	"Asia/Famagusta":        {},
	"Europe/Prague":         {},
	"Europe/Copenhagen":     {},
	"Atlantic/Faroe":        {},
	"Europe/Tallinn":        {},
	"Europe/Helsinki":       {},
	"Europe/Mariehamn":      {},
	"Europe/Paris":          {},
	"Europe/Berlin":         {},
	"Europe/Busingen":       {},
	"Europe/Athens":         {},
	"Europe/Budapest":       {},
	"Europe/Dublin":         {},
	"Europe/Rome":           {},
	"Europe/Riga":           {},
	"Europe/Vilnius":        {},
	"Europe/Luxembourg":     {},
	"Europe/Malta":          {},
	"Europe/Amsterdam":      {},
	"Europe/Warsaw":         {},
	"Europe/Lisbon":         {},
	"Europe/Bucharest":      {},
	"Europe/Bratislava":     {},
	"Europe/Ljubljana":      {},
	"Europe/Madrid":         {},
	"Europe/Stockholm":      {},
	"America/Cayenne":       {},
	"America/Guadeloupe":    {},
	"America/Marigot":       {},
	"America/Martinique":    {},
	"Indian/Mayotte":        {},
	"Indian/Reunion":        {},
	"Atlantic/Azores":       {},
	"Atlantic/Madeira":      {},
	"Atlantic/Canary":       {},
	"Africa/Ceuta":          {},
}

var eeaEftaTimezones = map[string]struct{}{
	"Atlantic/Reykjavik":  {},
	"Europe/Vaduz":        {},
	"Europe/Oslo":         {},
	"Atlantic/Jan_Mayen":  {},
	"Arctic/Longyearbyen": {},
}

var eftaSpecificTimezones = map[string]struct{}{
	"Europe/Zurich": {},
}

var schengenMicrostates = map[string]struct{}{
	"Europe/Andorra":    {},
	"Europe/Monaco":     {},
	"Europe/San_Marino": {},
	"Europe/Vatican":    {},
}

// IsFromEU reports whether a given IANA timezone belongs to an EU member state or its outermost regions.
func IsFromEU(timezone string) bool {
	_, ok := euTimezones[timezone]
	return ok
}

// IsFromEEA reports whether a given IANA timezone belongs to the European Economic Area (EU + Iceland, Liechtenstein, Norway).
func IsFromEEA(timezone string) bool {
	if IsFromEU(timezone) {
		return true
	}
	_, ok := eeaEftaTimezones[timezone]
	return ok
}

// IsFromEFTA reports whether a given IANA timezone belongs to the European Free Trade Association (Switzerland + Iceland, Liechtenstein, Norway).
func IsFromEFTA(timezone string) bool {
	if _, ok := eftaSpecificTimezones[timezone]; ok {
		return true
	}
	_, ok := eeaEftaTimezones[timezone]
	return ok
}

// IsFromSchengen reports whether a given IANA timezone belongs to the Schengen border-free area.
func IsFromSchengen(timezone string) bool {
	if IsFromEEA(timezone) || IsFromEFTA(timezone) {
		// Note: Ireland is EU but not Schengen, Cyprus is pending full accession
		if timezone == "Europe/Dublin" || timezone == "Europe/Nicosia" || timezone == "Asia/Nicosia" || timezone == "Asia/Famagusta" {
			return false
		}
		return true
	}
	_, ok := schengenMicrostates[timezone]
	return ok
}
