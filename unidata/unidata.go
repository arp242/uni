//go:generate go run ./gen.go

// Package unidata contains information about Unicode characters.
package unidata

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"zgo.at/zstd/zstring"
)

const UnknownCodepoint = "CODEPOINT NOT IN UNICODE"

const (
	GenderNone = 0
	GenderSign = 1
	GenderRole = 2
)

// Codepoint is a single codepoint.
type Codepoint struct {
	Codepoint rune
	Width     uint8
	Cat       uint8
	Name      string
	Digraph   string
	HTML      string
	KeySym    string // TODO: []string?
}

// Emoji is an emoji sequence.
type Emoji struct {
	Codepoints      []rune
	Name            string
	Group, Subgroup int
	CLDR            []string
	SkinTones       bool
	Genders         int
}

func (e Emoji) GroupName() string {
	return EmojiGroups[e.Group]
}

func (e Emoji) SubgroupName() string {
	return EmojiSubgroups[e.GroupName()][e.Subgroup]
}

// Find a codepoint.
func Find(cp rune) (Codepoint, bool) {
	info, ok := Codepoints[cp]
	if ok {
		return info, true
	}

	// The UnicodeData.txt file doesn't list every character; some are included as a
	// range:
	//
	//   3400;<CJK Ideograph Extension A, First>;Lo;0;L;;;;;N;;;;;
	//   4DB5;<CJK Ideograph Extension A, Last>;Lo;0;L;;;;;N;;;;;
	for i, r := range ranges {
		if cp >= r[0] && cp <= r[1] {
			info, ok := Codepoints[r[0]]
			if !ok {
				panic("unidata.Find: '" + string(r) + "' not found; this should never happen")
			}

			info.Codepoint = cp
			info.Name = rangeNames[i]
			return info, true
		}
	}

	return Codepoint{Codepoint: cp, Name: UnknownCodepoint}, false
}

// ToRune converts a human input string to a rune.
//
// The input can be as U+41, U+0041, U41, 0x41, 0o101, 0b1000001
//
// if strings.HasPrefix(strings.ToUpper(s[0]), "0D") {
// 	// Let ToRune deal with any errors.
// 	if n, err := strconv.ParseInt(s[0][2:], 10, 10); err == nil {
// 		s[0] = strconv.FormatInt(n, 16)
// 	}
// }
// if strings.HasPrefix(strings.ToUpper(s[1]), "0D") {
// 	if n, err := strconv.ParseInt(s[1][2:], 10, 10); err == nil {
// 		s[1] = strconv.FormatInt(n, 16)
// 	}
// }
func ToRune(s string) (rune, error) {
	os := s
	s = strings.ToUpper(s)
	var base = 16
	switch {
	case zstring.HasPrefixes(s, "0X", "U+"):
		s = s[2:]
	case strings.HasPrefix(s, "0D"):
		s = s[2:]
		base = 10
	case strings.HasPrefix(s, "0O"):
		s = s[2:]
		base = 8
	case strings.HasPrefix(s, "0B"):
		s = s[2:]
		base = 2

	case zstring.HasPrefixes(s, "X", "U"):
		s = s[1:]
	case strings.HasPrefix(s, "O"):
		s = s[1:]
		base = 8
	}
	i, err := strconv.ParseInt(s, base, 32)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			return 0, fmt.Errorf("out of range: %q", os)
		}
		if errors.Is(err, strconv.ErrSyntax) {
			return 0, fmt.Errorf("not a number or codepoint: %q", os)
		}
		return 0, err
	}
	return rune(i), nil
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

func (c Codepoint) String() string {
	return c.Repr(false) + ": " + c.FormatCodepoint() + " " + c.Name
}

func (c Codepoint) FormatCodepoint() string {
	return fmt.Sprintf("U+%04X", c.Codepoint)
}

func (c Codepoint) Format(base int) string {
	return strconv.FormatUint(uint64(c.Codepoint), base)
}

func (c Codepoint) Plane() string {
	for p, r := range Planes {
		if c.Codepoint >= r[0] && c.Codepoint <= r[1] {
			return p
		}
	}
	return ""
}

func (c Codepoint) WidthName() string {
	return WidthNames[c.Width]
}

func (c Codepoint) Category() string {
	return Catnames[c.Cat]
}

func (c Codepoint) Block() string {
	for b, r := range Blocks {
		if c.Codepoint >= r[0] && c.Codepoint <= r[1] {
			return b
		}
	}
	return ""
}

func (c Codepoint) UTF8() string {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, c.Codepoint)
	return fmt.Sprintf("% x", buf[:n])
}

func (c Codepoint) UTF16(bigEndian bool) string {
	var p []byte
	if c.Codepoint <= 0xffff {
		p = []byte{byte(c.Codepoint % 256), byte(c.Codepoint >> 8)}
		if bigEndian {
			p[1], p[0] = p[0], p[1]
		}
	} else {
		a, b := utf16.EncodeRune(c.Codepoint)
		p = []byte{byte(a % 256), byte(a >> 8), byte(b % 256), byte(b >> 8)}
		if bigEndian {
			p[1], p[0], p[3], p[2] = p[0], p[1], p[2], p[3]
		}
	}
	return fmt.Sprintf(`% x`, p)
}

func (c Codepoint) XMLEntity() string {
	return "&#x" + strconv.FormatInt(int64(c.Codepoint), 16) + ";"
}

func (c Codepoint) JSON() string {
	u := strings.ReplaceAll(c.UTF16(true), " ", "")
	if len(u) == 4 {
		return `\u` + u
	}
	return `\u` + u[:4] + `\u` + u[4:]
}

func (c Codepoint) HTMLEntity() string {
	if c.HTML != "" {
		return "&" + c.HTML + ";"
	}
	return c.XMLEntity()
}

func (c Codepoint) Repr(raw bool) string {
	if raw {
		return string(c.Codepoint)
	}

	cp := c.Codepoint

	// Display combining characters with ◌.
	if unicode.In(cp, unicode.Mn, unicode.Mc, unicode.Me) {
		return "\u25cc" + string(cp)
	}

	switch {
	case unicode.IsControl(cp):
		switch {
		case cp < 0x20: // C0; use "Control Pictures" block
			cp += 0x2400
		case cp == 0x7f: // DEL
			cp = 0x2421
		// No control pictures for C1 or anything else, use "open box".
		default:
			cp = 0x2423
		}
	// "Other, Format" category except the soft hyphen and spaces.
	case !unicode.IsPrint(cp) && cp != 0x00ad && !unicode.In(cp, unicode.Zs):
		cp = 0xfffd
	}

	return string(cp)
}

func (e Emoji) String() string {
	var c string

	// Flags
	// 1F1FF 1F1FC                                 # 🇿🇼 E2.0 flag: Zimbabwe
	// 1F3F4 E0067 E0062 E0065 E006E E0067 E007F   # 🏴󠁧󠁢󠁥󠁮󠁧󠁿 E5.0 flag: England
	if (e.Codepoints[0] >= 0x1f1e6 && e.Codepoints[0] <= 0x1f1ff) ||
		(len(e.Codepoints) > 1 && e.Codepoints[1] == 0xe0067) {
		for _, cp := range e.Codepoints {
			c += string(rune(cp))
		}
		return c
	}

	for i, cp := range e.Codepoints {
		c += string(rune(cp))

		// Don't add ZWJ as last item.
		if i == len(e.Codepoints)-1 {
			continue
		}

		switch e.Codepoints[i+1] {
		// Never add ZWJ before variation selector or skin tone.
		case 0xfe0f, 0x1f3fb, 0x1f3fc, 0x1f3fd, 0x1f3fe, 0x1f3ff:
			continue
		// Keycap: join with 0xfe0f
		case 0x20e3:
			continue
		}

		c += "\u200d"
	}
	return c
}
