package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"arp242.net/uni/unidata"
	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

var (
	isTerm    = zli.IsTerminal(os.Stdout.Fd())
	termWidth = func() int {
		if !isTerm {
			return 0
		}

		w, _, err := zli.TerminalSize(os.Stdout.Fd())
		if err != nil || w < 50 {
			return 0
		}
		return w
	}()
)

type Format struct {
	out         io.Writer
	format      string
	printHeader bool
	re          *regexp.Regexp
	autoalign   map[string]int
	buf         *strings.Builder
}

// TODO: show better error on using %(unknown...)
func NewFormat(out io.Writer, format string, printHeader bool) Format {
	f := Format{
		out:         out,
		format:      format,
		buf:         new(strings.Builder),
		printHeader: printHeader,
		autoalign:   make(map[string]int),
	}

	if printHeader {
		f.Header(map[string]string{
			"emoji":        "",
			"char":         "",
			"tab":          tabOrSpace(),
			"wide_padding": " ",

			// TODO: we can do this automatically
			"name":     "name",
			"group":    "group",
			"subgroup": "subgroup",
			"cpoint":   "cpoint",
			"dec":      "dec",
			"cat":      "cat",
			"xml":      "xml",
			"keysym":   "keysym",
			"utf8":     "utf8",
			"html":     "html",
		})
	}

	return f
}

// TODO: much of this should be refactored; it's ugly and all kinds of
// inefficient. should store this properly rather than using the $$...$$
// markers, for example:
//
//    []struct {
//        text      string
//        autoAlign bool
//        trim      bool
//     }
//
// Then we can loop over that, and not have to resort to regexing the $$..$$
// markers out of there.
func (f *Format) Flush() {

	text := f.buf.String()
	for k, v := range f.autoalign {
		text = regexp.MustCompile(fmt.Sprintf(`\$\$auto-%s-\d+\$\$`, k)).
			ReplaceAllStringFunc(text, func(m string) string {
				l, _ := strconv.Atoi(strings.Split(strings.TrimRight(m, "$"), "-")[2])
				if v-l > 0 {
					return strings.Repeat(" ", v-l)
				}
				return ""
			})
	}

	re := regexp.MustCompile(`\$\$trim-start\$\$(.*?)\$\$trim-end\$\$`)
	var nt strings.Builder
	for _, line := range strings.Split(text, "\n") {
		if line == "" {
			continue
		}
		if !isTerm || termWidth == 0 || utf8.RuneCountInString(line) <= termWidth {
			nt.WriteString(re.ReplaceAllString(line, `$1`))
			nt.WriteRune('\n')
			continue
		}

		ntrim := len(re.FindAllString(line, -1))
		tooMany := utf8.RuneCountInString(line) - termWidth - 24*ntrim /* $$trim markers */ + ntrim /* … chars */

		// Divide amount to trim between columns
		tooMany = int(math.Ceil(float64(tooMany) / float64(ntrim)))

		line = re.ReplaceAllStringFunc(line, func(m string) string {
			t := m[14 : utf8.RuneCountInString(m)-12]
			if tooMany <= 0 {
				return t
			}

			x := utf8.RuneCountInString(t) - tooMany
			// TODO: sometimes the "name" col is so long that it's longer than
			// all of the category.
			if x <= 0 {
				return "…"
			}
			return t[:x] + "…"
		})

		nt.WriteString(line)
		nt.WriteRune('\n')
	}

	text = nt.String()
	f.out.Write([]byte(text))
}

func (f *Format) Line(columns map[string]func() string) {
	if f.re == nil {
		keys := make([]string, 0, len(columns))
		for k := range columns {
			keys = append(keys, k)
		}
		f.re = regexp.MustCompile(`%\((` + strings.Join(keys, "|") + `)(?: .+?)?\)`)
	}

	f.buf.WriteString(f.re.ReplaceAllStringFunc(f.format, func(m string) string {
		s := zstring.Fields(m[2:len(m)-1], " ") // name, flags
		return f.col(s[0], columns[s[0]](), s[1:], false)
	}))
	f.buf.WriteRune('\n')
}

func (f *Format) Header(columns map[string]string) {
	if f.re == nil {
		keys := make([]string, 0, len(columns))
		for k := range columns {
			keys = append(keys, k)
		}
		f.re = regexp.MustCompile(`%\((` + strings.Join(keys, "|") + `)(?: .+?)?\)`)
	}

	f.buf.WriteString(f.re.ReplaceAllStringFunc(f.format, func(m string) string {
		s := zstring.Fields(m[2:len(m)-1], " ") // name, flags
		return f.col(s[0], columns[s[0]], s[1:], true)
	}))
	f.buf.WriteRune('\n')
}

func (f Format) col(name, in string, flags []string, header bool) string {
	if len(flags) == 0 {
		return in
	}

	var (
		quote      bool
		trim       bool
		alignAuto  bool
		alignLeft  int
		alignRight int
	)
	for _, f := range flags {
		switch {
		default:
			return fmt.Sprintf("%%(ERROR format %q: unknown flag %q", in, f)
		case f == "":
			continue
		case f == "q":
			quote = true
		case f == "t":
			trim = true
		case f[0] == 'l' || f[0] == 'r':
			n := strings.Split(f, ":")
			if len(n) != 2 {
				return fmt.Sprintf("%%(ERROR format %q: need width after :)", in)
			}

			if n[1] == "auto" {
				alignAuto = true
				continue
			}

			num, err := strconv.ParseInt(n[1], 10, 8)
			if err != nil {
				return fmt.Sprintf("%%(ERROR format %q: %s)", in, err)
			}

			if f[0] == 'l' {
				alignLeft = int(num)
			} else {
				alignRight = int(num)
			}
		}
	}

	if quote {
		if header {
			in = " " + in + " "
		} else {
			in = "'" + in + "'"
		}
	}
	if alignLeft > 0 {
		in = zstring.AlignLeft(in, alignLeft)
	}
	if alignRight > 0 {
		in = zstring.AlignRight(in, alignRight)
	}

	if trim {
		in = "$$trim-start$$" + in + "$$trim-end$$"
	}

	// TODO: this is kind of an ugly way to do it...
	if alignAuto {
		if utf8.RuneCountInString(in) > f.autoalign[name] {
			f.autoalign[name] = utf8.RuneCountInString(in)
		}
		in = in + "$$auto-" + name + "-" + strconv.Itoa(utf8.RuneCountInString(in)) + "$$"
	}

	return in
}

func toLine(info unidata.Codepoint, raw bool) map[string]func() string {
	// func (p printer) entity(c rune, cp uint32) string {
	c := rune(info.Codepoint)

	return map[string]func() string{
		"char":         func() string { return fmtChar(c, raw) },
		"wide_padding": func() string { return widePadding(info) },
		"cpoint":       func() string { return fmt.Sprintf("U+%04X", info.Codepoint) },
		"dec":          func() string { return strconv.FormatUint(uint64(info.Codepoint), 10) },
		"utf8":         func() string { return utf8Bytes(c) },
		"html":         func() string { return htmlEntity(c, info.Codepoint) },
		"xml":          func() string { return fmt.Sprintf("#x%x", info.Codepoint) },
		"keysym":       func() string { return keysym(c) },
		"name":         func() string { return info.Name },
		"cat":          func() string { return unidata.Catnames[info.Cat] },
	}
}

// Alignment with spaces is tricky, as some emojis are double-width and some are
// not. As far as I can tell, there is no good way to predict this as it will
// depend on the font. Unicode recommends "emoji presentation sequences behave
// as though they were East Asian Wide", but that's too simplistic too. So use a
// tab character for this, which aligns correctly, even though it adds some
// unnecessary whitespace.
//
// Don't do this when piping, since dmenu doesn't display tabs well :-/ This
// seems like a problem in Xft as near as I can determine.
func tabOrSpace() string {
	if zli.IsTerminal(os.Stdout.Fd()) {
		return "\t"
	}
	return "    "
}

func widePadding(info unidata.Codepoint) string {
	if info.Width != unidata.WidthFullWidth && info.Width != unidata.WidthWide {
		return " "
	}
	return ""
}

func keysym(c rune) string {
	s, ok := unidata.KeySyms[c]
	if !ok {
		return "(none)"
	}
	return s[0]
}

func htmlEntity(c rune, cp uint32) string {
	html := unidata.Entities[c]
	if html == "" {
		html = fmt.Sprintf("#x%x", cp)
	}
	return "&" + html + ";"
}

func utf8Bytes(c rune) string {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, c)
	return fmt.Sprintf("% x", buf[:n])
}

func fmtChar(c rune, raw bool) string {
	if raw {
		return string(c)
	}

	// Display combining characters with ◌.
	if unicode.In(c, unicode.Mn, unicode.Mc, unicode.Me) {
		return "\u25cc" + string(c)
	}

	switch {
	case unicode.IsControl(c):
		switch {
		case c < 0x20: // C0; use "Control Pictures" block
			c += 0x2400
		case c == 0x7f: // DEL
			c = 0x2421
		// No control pictures for C1 or anything else, use "open box".
		default:
			c = 0x2423
		}
	// "Other, Format" category except the soft hyphen and spaces.
	case !unicode.IsPrint(c) && c != 0x00ad && !unicode.In(c, unicode.Zs):
		c = 0xfffd
	}

	return string(c)
}

// func (p *printer) PrintSorted(fp io.Writer, quiet, raw bool) {
// 	s := []unidata.Codepoint(*p)
// 	sort.Slice(s, func(i int, j int) bool { return s[i].Codepoint < s[j].Codepoint })
// 	p.Print(fp, quiet, raw)
// }
