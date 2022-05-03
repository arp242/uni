package unidata

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

type (
	// Codepoint represents a single codepoint.
	Codepoint struct {
		Codepoint rune // The actual codepoint.

		// These are unexported and retrieved with a "getter" so we can change
		// the implementation later. Right now they're all fields on a struct,
		// but might be a good idea to move at least some of them out of there
		// at some point.
		//
		// Note: don't change the order without changing gen_codepoints.go
		width    Width
		category Category
		name     string
	}

	Width        uint8      // Unicode width
	Plane        uint8      // Unicode plane
	Category     uint8      // Unicode category
	Block        uint16     // Unicode block
	Property     uint8      // Unicode property
	PropertyList []Property // Unicode property
)

func (w Width) String() string    { return Widths[w] }
func (c Category) String() string { return Categories[c].Name }
func (p Plane) String() string    { return Planes[p].Name }
func (b Block) String() string    { return Blocks[b].Name }
func (p Property) String() string { return Properties[p].Name }
func (p PropertyList) String() string {
	var b strings.Builder
	for i, pp := range p {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(Properties[pp].Name)
	}
	return b.String()
}

var mName = strings.NewReplacer(
	"&", "",
	" ", "",
	"-", "",
	"_", "",
)

func matchName(s string) string { return strings.ToLower(mName.Replace(s)) }

// FindBlock finds a block by name.
func FindBlock(name string) (Block, bool) {
	var (
		match = matchName(name)
		found []Block
	)
	for k, b := range Blocks {
		if matchName(b.Name) == match {
			return k, true
		}

		if strings.HasPrefix(matchName(b.Name), match) {
			found = append(found, k)
		}
	}

	switch len(found) {
	case 0:
		return 0, false
	case 1:
		return found[0], true
	default:
		return 0, false // TODO: print what we found
	}
}

// FindCategory finds a category by name.
func FindCategory(name string) (Category, bool) {
	var (
		match = matchName(name)
		found []Category
	)
	for k, b := range Categories {
		if matchName(b.Name) == match || matchName(b.ShortName) == match {
			return k, true
		}
		if strings.HasPrefix(matchName(b.Name), match) || matchName(b.ShortName) == match {
			found = append(found, k)
		}
	}

	switch len(found) {
	case 0:
		return 0, false
	case 1:
		return found[0], true
	default:
		return 0, false // TODO: print what we found
	}
}

// FindProperty finds a property by name.
func FindProperty(name string) (Property, bool) {
	var (
		match = matchName(name)
		found []Property
	)
	for k, b := range Properties {
		if matchName(b.Name) == match {
			return k, true
		}
		if strings.HasPrefix(matchName(b.Name), match) {
			found = append(found, k)
		}
	}

	switch len(found) {
	case 0:
		return 0, false
	case 1:
		return found[0], true
	default:
		return 0, false // TODO: print what we found
	}
}

// Find a Codepoint for this rune.
//
// If the second return value is false, the codepoint wasn't found. The
// Codepoint will have only the Codepoint field set.
func Find(cp rune) (Codepoint, bool) {
	info, ok := Codepoints[cp]
	if ok {
		return info, true
	}

	for _, r := range codepointRanges {
		if cp >= r.rng[0] && cp <= r.rng[1] {
			info, ok := Codepoints[r.rng[0]]
			if !ok {
				panic("unidata.Find: '" + string(r.rng[0]) + string(r.rng[1]) +
					"' not found in range; this should never happen")
			}

			info.Codepoint = cp
			info.name = r.name
			return info, true
		}
	}
	return Codepoint{Codepoint: cp, name: "CODEPOINT NOT IN UNICODE"}, false
}

// FromString gets a codepoint from human input.
//
// The input can be as (case-insensitive):
//
//   U+F1, U+00F1, UF1, F1   Unicode codepoint notation (hex)
//   0xF1, xF1               Hex number
//   0d241                   Decimal number
//   0o361, o361             Octal number
//   0b11110001              Binary number
func FromString(s string) (Codepoint, error) {
	os := s
	s = strings.ToUpper(s)
	var base = 16
	switch {
	case strings.HasPrefix(s, "0X") || strings.HasPrefix(s, "U+"):
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

	case strings.HasPrefix(s, "X") || strings.HasPrefix(s, "U"):
		s = s[1:]
	case strings.HasPrefix(s, "O"):
		s = s[1:]
		base = 8
	}
	i, err := strconv.ParseInt(s, base, 32)
	if err != nil {
		if errors.Is(err, strconv.ErrRange) {
			err = fmt.Errorf("out of range: %q", os)
		} else if errors.Is(err, strconv.ErrSyntax) {
			err = fmt.Errorf("not a number or codepoint: %q", os)
		}
		return Codepoint{}, fmt.Errorf("unidata.FromString: %w", err)
	}
	cp, _ := Find(rune(i))
	return cp, nil
}

func (c Codepoint) String() string {
	return c.Display() + ": " + c.FormatCodepoint() + " " + c.name
}

// Display this codepoint. This formats the codepoint as follows:
//
//   - Combining characters are prefixed with ◌ (U+25CC)
//   - C0 control characters use the graphical representation (Control Pictures block, U+2400-2426)
//   - C1 control characters use open box ␣ (U+2423)
//   - Other unprintable characters use the replacement character � (U+FFFD)
//   - Everything else is converted to a string without processing.
func (c Codepoint) Display() string {
	// Display combining characters with ◌.
	if c.Category() == CatNonspacingMark || c.Category() == CatSpacingMark || c.Category() == CatEnclosingMark {
		return "\u25cc" + string(c.Codepoint)
	}

	cp := c.Codepoint
	switch {
	case c.isControl():
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
	case !c.isPrint() && cp != 0x00ad && c.Category() != CatSpaceSeparator:
		cp = 0xfffd
	}

	return string(cp)
}

// Name gets the name for this codepoint.
func (c Codepoint) Name() string { return c.name }

// Width gets this codepoint's width.
func (c Codepoint) Width() Width { return c.width }

// Category gets this codepoint's category.
func (c Codepoint) Category() Category { return c.category }

// Plane gets the Unicode plane.
func (c Codepoint) Plane() Plane {
	for k, v := range Planes {
		if c.Codepoint >= v.Range[0] && c.Codepoint <= v.Range[1] {
			return k
		}
	}
	return PlaneUnknown
}

// Block gets the unicode block.
//
// This only gets the block this codepoint is assigned to; use Codepoint.In() if
// you want to check if a codepoint is within a block (some blocks are a group
// of other blocks; for example Number is DecimalNumber + LetterNumber +
// OtherNumber).
func (c Codepoint) Block() Block {
	for k, v := range Blocks {
		if c.Codepoint >= v.Range[0] && c.Codepoint <= v.Range[1] {
			return k
		}
	}
	return BlockUnknown
}

// Properties gets the unicode properties for this codepoint.
func (c Codepoint) Properties() PropertyList {
	all := make(PropertyList, 0, 1)
	for k, v := range Properties {
		for _, r := range v.Ranges {
			if c.Codepoint >= r[0] && c.Codepoint <= r[1] {
				all = append(all, k)
			}
		}
	}
	return all
}

// FormatCodepoint formats the codepoint in Unicode notation.
func (c Codepoint) FormatCodepoint() string {
	return fmt.Sprintf("U+%04X", c.Codepoint)
}

// Format the codepoint in the given base.
func (c Codepoint) Format(base int) string {
	return strconv.FormatUint(uint64(c.Codepoint), base)
}

// UTF8 gets the UTF-8 representation of this codepoint.
func (c Codepoint) UTF8() []byte {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, c.Codepoint)
	return buf[:n]
}

// UTF16 gets the UTF-16 representation of this codepoint.
//
// The default is to use Little-Endian encoding (which is what Windows uses);
// set bigEndian to use Big-Endian encoding.
func (c Codepoint) UTF16(bigEndian bool) []byte {
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
	return p
}

// JSON gets the JSON representation.
func (c Codepoint) JSON() string {
	u := fmt.Sprintf("%x", c.UTF16(true))
	if len(u) == 4 {
		return `\u` + u
	}
	return `\u` + u[:4] + `\u` + u[4:]
}

// XML formats the codepoint as an XML entity.
func (c Codepoint) XML() string { return "&#x" + strconv.FormatInt(int64(c.Codepoint), 16) + ";" }

// HTML formats the codepoint as an HTML entity, prefering a symbolic name if it
// exists (e.g. &amp; instead of &#x26;)
func (c Codepoint) HTML() string {
	if h := htmlEntities[c.Codepoint]; h != "" {
		return "&" + h + ";"
	}
	return c.XML()
}

// KeySym gets the X11 keysym name.
func (c Codepoint) KeySym() string { return keysyms[c.Codepoint] }

// Digraph gets the digraph sequence.
//
// Digraphs are defined in RFC1345. This adds two digraphs that Vim recognises
// but are not in the RFC:
//
//   =e    €   U+20AC EURO SIGN
//   =R    ₽   U+20BD RUBLE SIGN
func (c Codepoint) Digraph() string { return digraphs[c.Codepoint] }

// in reports if this codepoint is in the given category.
//
// TODO: this is a bit ugly; we should generate this data better.
func (c Codepoint) in(cats ...Category) bool {
	for _, cat := range cats {
		cc := c.Category()
		if cc == cat {
			return true
		}

		t := false
		switch cat {
		case CatCasedLetter: // LC – Lu | Ll | Lt
			t = cc == CatUppercaseLetter || cc == CatLowercaseLetter || cc == CatTitlecaseLetter
		case CatLetter: // L  – Lu | Ll | Lt | Lm | Lo
			t = cc == CatUppercaseLetter || cc == CatLowercaseLetter ||
				cc == CatTitlecaseLetter || cc == CatModifierLetter || cc == CatOtherLetter
		case CatMark: // M  – Mn | Mc | Me
			t = cc == CatNonspacingMark || cc == CatSpacingMark || cc == CatEnclosingMark
		case CatNumber: // N  – Nd | Nl | No
			t = cc == CatDecimalNumber || cc == CatLetterNumber || cc == CatOtherNumber
		case CatPunctuation: // P  – Pc | Pd | Ps | Pe | Pi | Pf | Po
			t = cc == CatConnectorPunctuation || cc == CatDashPunctuation || cc == CatOpenPunctuation ||
				cc == CatClosePunctuation || cc == CatInitialPunctuation || cc == CatFinalPunctuation ||
				cc == CatOtherPunctuation
		case CatSymbol: // S  – Sm | Sc | Sk | So
			t = cc == CatMathSymbol || cc == CatCurrencySymbol || cc == CatModifierSymbol ||
				cc == CatOtherSymbol
		case CatSeparator: // Z  – Zs | Zl | Zp
			t = cc == CatSpaceSeparator || cc == CatLineSeparator || cc == CatParagraphSeparator
		case CatOther: // C  – Cc | Cf | Cs | Co | Cn
			t = cc == CatControl || cc == CatFormat || cc == CatSurrogate ||
				cc == CatPrivateUse || cc == CatUnassigned
		}
		if t {
			return true
		}
	}
	return false
}

func (c Codepoint) isControl() bool {
	return c.Category() == CatControl
}

func (c Codepoint) isPrint() bool {
	if c.isControl() {
		return false
	}
	return c.in(CatLetter, CatMark, CatNumber, CatPunctuation, CatSymbol)
}
