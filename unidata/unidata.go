package unidata

import (
	"fmt"
	"strconv"
	"strings"
)

const UnknownCodepoint = "CODEPOINT NOT IN UNICODE"

// FindCodepoint finds a codepoint
func FindCodepoint(c rune) (Codepoint, bool) {
	info, ok := Codepoints[fmt.Sprintf("%.4X", c)]
	if ok {
		return info, true
	}

	// The UnicodeData.txt file doesn't list every character; some are included as a
	// range:
	//
	//   3400;<CJK Ideograph Extension A, First>;Lo;0;L;;;;;N;;;;;
	//   4DB5;<CJK Ideograph Extension A, Last>;Lo;0;L;;;;;N;;;;;
	for i, r := range ranges {
		if c >= r[0] && c <= r[1] {
			info, ok := Codepoints[fmt.Sprintf("%.4X", r[0])]
			if !ok {
				panic(fmt.Sprintf("FindCodepoint: %#v not found; this should never happen", r[0]))
			}

			info.Codepoint = uint32(c)
			info.Name = rangeNames[i]
			return info, true
		}
	}

	return Codepoint{Codepoint: uint32(c), Name: UnknownCodepoint}, false
}

// ToCodepoint converts a human input string to a codepoint.
//
// The input can be as U+41, U+0041, U41, 0x41, 0o101, 0b1000001
func ToCodepoint(s string) (int64, error) {
	s = strings.ToUpper(s)
	var base = 16
	switch {
	case strings.HasPrefix(s, "0X"), strings.HasPrefix(s, "U+"):
		s = s[2:]
	case strings.HasPrefix(s, "U"):
		s = s[1:]
	case strings.HasPrefix(s, "0O"):
		s = s[2:]
		base = 8
	case strings.HasPrefix(s, "0B"):
		s = s[2:]
		base = 2
	}
	return strconv.ParseInt(s, base, 64)
}

// CanonicalCategory transforms a category name to the canonical representation.
func CanonicalCategory(cat string) string {
	// TODO: improve.
	cat = strings.Replace(cat, " ", "", -1)
	cat = strings.Replace(cat, ",", "", -1)
	cat = strings.Replace(cat, "_", "", -1)
	cat = strings.ToLower(cat)
	return cat
}
