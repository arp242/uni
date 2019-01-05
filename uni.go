// Command uni prints Unicode information about characters.
package main // import "arp242.net/uni"

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

var ranges = [][]rune{
	{0x3400, 0x4DB5},
	{0x4E00, 0x9FEF},
	{0xAC00, 0xD7A3},
	{0xD800, 0xDB7F},
	{0xDB80, 0xDBFF},
	{0xDC00, 0xDFFF},
	{0xE000, 0xF8FF},
	{0x17000, 0x187F1},
	{0x20000, 0x2A6D6},
	{0x2A700, 0x2B734},
	{0x2B740, 0x2B81D},
	{0x2B820, 0x2CEA1},
	{0x2CEB0, 0x2EBE0},
	{0xF0000, 0xFFFFD},
	{0x100000, 0x10FFFD},
}

var rangeNames = []string{
	"<CJK Ideograph Extension A>",
	"<CJK Ideograph>",
	"<Hangul Syllable>",
	"<Non Private Use High Surrogate>",
	"<Private Use High Surrogate>",
	"<Low Surrogate>",
	"<Private Use>",
	"<Tangut Ideograph>",
	"<CJK Ideograph Extension B>",
	"<CJK Ideograph Extension C>",
	"<CJK Ideograph Extension D>",
	"<CJK Ideograph Extension E>",
	"<CJK Ideograph Extension F>",
	"<Plane 15 Private Use>",
	"<Plane 16 Private Use>",
}

func main() {
	// TODO: Better argument parsing.
	// TODO: Add option to switch off header.
	// TODO: Add option for TSV and/or JSON output.
	if len(os.Args) < 2 {
		fatal(errors.New("wrong arguments"))
	}

	// TODO: idea: add command to print/search emojis?
	switch strings.ToLower(os.Args[1]) {
	default:
		fatal(errors.New("wrong arguments"))

	case "identify", "i":
		identify(strings.Join(os.Args[2:], ""))

	// TODO: stable output, ordered by code point (due to map it's not stable
	// now).
	case "search", "s":
		if len(os.Args) < 3 {
			fatal(errors.New("need search term"))
		}

		// TODO: don't print with 0 matches.
		header()
		words := make([]string, len(os.Args)-2)
		for i := range os.Args[2:] {
			words[i] = strings.ToUpper(os.Args[i+2])
		}
		for cp, name := range uniData {
			m := 0
			for _, w := range words {
				if strings.Contains(name, w) {
					m++
				}
			}
			if m == len(words) {
				printEntry(cp, name)
			}
		}
	}
}

func fromFile(in string) string {
	x := strings.Split(in, ":#")
	if len(x) != 2 {
		return in
	}

	var seek, read int64
	var err error

	// Can be as :#42 or #42-50
	if strings.Contains(x[1], "-") {
		rng := strings.Split(x[1], "-")
		seek, err = strconv.ParseInt(rng[0], 10, 32)
		if err != nil {
			return in
		}

		to, err := strconv.ParseInt(rng[1], 10, 32)
		if err != nil {
			return in
		}

		if to < seek {
			return in
		}
		read = to - seek + 1
	} else {
		seek, err = strconv.ParseInt(x[1], 10, 32)
		if err != nil {
			return in
		}
		read = 1
	}

	fp, err := os.Open(x[0])
	if err != nil {
		return in
	}

	_, _ = fp.Seek(seek, io.SeekStart)
	b := make([]byte, read)
	_, _ = fp.Read(b)

	return string(b)
}

func identify(in string) {
	if strings.Contains(in, ":#") {
		in = fromFile(in)
	}

	if !utf8.ValidString(in) {
		_, _ = fmt.Fprintf(os.Stderr, "uni: WARNING: input string is not valid UTF-8\n")
	}

	// TODO: don't print with 0 matches.
	header()
	for _, c := range in {
		find := fmt.Sprintf("%04X", c)
		m := false

		if name := inRange(c); name != "" {
			printEntry(find, name)
			m = true
		}

		if !m {
			for cp, name := range uniData {
				if cp == find {
					printEntry(cp, name)
					m = true
					break
				}
			}
			if !m {
				fmt.Printf("NO MATCH: %s\n", find)
			}
		}
	}
}

func inRange(c rune) string {
	for i, r := range ranges {
		if c >= r[0] && c <= r[1] {
			return rangeNames[i]
		}
	}
	return ""
}

// TODO: check terminal size; don't print full descriptin if it doesn't fit.
func printEntry(cp, name string) {
	r, _ := strconv.ParseInt(cp, 16, 64)
	rs := strconv.FormatInt(r, 10)
	char := string(r)
	html := entities[rune(r)]
	if html == "" {
		html = fmt.Sprintf("#x%x", r)
	}
	html = "&" + html + ";"

	var line strings.Builder

	// Write literal character with extra space for alignment if it's a
	// half-width character.
	// TODO: this is an inperfect check; e.g. many emojis (U+1F973) are
	// full-width.
	line.WriteString(fmt.Sprintf("'%s' ", char))
	if inRange(rune(r)) == "" && r != 129395 {
		line.WriteString(" ")
	}

	line.WriteString(fmt.Sprintf("U+%s %s %s %s %s",
		fill(cp, 5), fill(rs, 6), fill(asutf8(rune(r)), 11), fill(html, 10), name))

	fmt.Println(line.String())
}

func asutf8(r rune) string {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	return fmt.Sprintf("% x", buf[:n])
}

func header() {
	fmt.Printf("     cpoint  dec    utf-8       html       name\n")
}

func fill(s string, n int) string {
	if len(s) >= n {
		return s
	}

	return s + strings.Repeat(" ", n-len(s))
}

func fatal(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
