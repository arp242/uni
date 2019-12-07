package unidata

import (
	"fmt"
	"strings"
)

// FindCodepoint finds a codepoint
func FindCodepoint(c rune) (Codepoint, bool) {
	info, ok := Codepoints[fmt.Sprintf("%.4X", c)]
	if ok {
		return info, ok
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
				// Should never happen.
				panic(fmt.Sprintf("FindCodepoint: %#v not found", r[0]))
			}

			info.Codepoint = uint32(c)
			info.Name = rangeNames[i]
			return info, true
		}
	}

	return Codepoint{}, false
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
