package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"arp242.net/uni/unidata"
	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

var termWidth = func() int {
	if !zli.IsTerminal(os.Stdout.Fd()) {
		return 0
	}

	w, _, err := zli.TerminalSize(os.Stdout.Fd())
	if err != nil || w < 50 {
		return 0
	}
	return w
}()

const (
	alignNone = iota
	alignLeft
	alignRight
	alignAuto = -1
)

type column struct {
	name  string
	width int
	align int
	trim  bool
	quote bool
}

type Format struct {
	format    string         // Format string: %(..)
	re        *regexp.Regexp // Cached regexp for format.
	cols      []column       // Columns we know about.
	lines     [][]string     // Processed lines, to be printed.
	autoalign []int          // Max line lengths for autoalign.
	ntrim     int            // Number of columns with "trim"

	printHeader bool
}

func NewFormat(format string, printHeader bool) Format {
	var (
		reFindCols = regexp.MustCompile(`%\((.*?)(?: .+?)?\)`)
		f          = Format{format: format, printHeader: printHeader}
	)
	for _, m := range reFindCols.FindAllString(format, -1) {
		err := f.processColumn(m)
		if err != nil {
			panic(err)
		}
	}

	f.autoalign = make([]int, len(f.cols))

	cols := make([]string, 0, len(f.cols))
	h := map[string]string{}
	for _, c := range f.cols {
		cols = append(cols, c.name)
		h[c.name] = c.name
	}

	h["emoji"] = ""
	h["char"] = ""
	h["tab"] = tabOrSpace()
	h["wide_padding"] = " "

	if printHeader {
		f.Line(h)
	}

	// TODO: is this actually faster than just .*?
	// TODO: don't really need to use regexp for this; can just scan for "%(name".
	f.re = regexp.MustCompile(`%\((` + strings.Join(cols, "|") + `)(?: .+?)?\)`)
	return f
}

func (f *Format) processColumn(line string) error {
	s := zstring.Fields(line[2:len(line)-1], " ") // name, flags
	name := s[0]
	col := column{name: name}

	if len(s) == 1 { // No flags
		f.cols = append(f.cols, col)
		return nil
	}

	for _, flag := range strings.Split(s[1], " ") {
		switch {
		default:
			return fmt.Errorf("unknown flag: %q", f)
		case flag == "":
			continue
		case flag == "q":
			col.quote = true
		case flag == "t":
			f.ntrim++
			col.trim = true
		case flag[0] == 'l' || flag[0] == 'r':
			n := strings.Split(flag, ":")
			if len(n) != 2 {
				return fmt.Errorf("need width after :")
			}

			switch flag[0] {
			case 'l':
				col.align = alignLeft
			case 'r':
				col.align = alignRight
			default:
				return fmt.Errorf("error")
			}

			if n[1] == "auto" {
				// TODO: allow setting a minimum/max width:
				//   %(col l:auto:5)     at least 5
				//   %(col l:auto:5:10)  5-10
				//   %(col l:auto:0:10)  at most 10
				col.width = alignAuto
			} else {
				var err error
				col.width, err = strconv.Atoi(n[1])
				if err != nil {
					return fmt.Errorf("%s", err)
				}
			}
		}
	}

	f.cols = append(f.cols, col)
	return nil
}

// Add a new line.
//
// TODO: check if colums match with f.cols
func (f *Format) Line(columns map[string]string) error {
	line := make([]string, len(f.cols))
	for i, c := range f.cols {
		line[i] = columns[c.name]
		if c.width == alignAuto {
			if l := zstring.TabWidth(columns[c.name]); l > f.autoalign[i] {
				f.autoalign[i] = l
			}
		}
	}
	f.lines = append(f.lines, line)
	return nil
}

// Sort by column.
func (f *Format) Sort(col string) {
	coli := 0
	for i, c := range f.cols {
		if c.name == col {
			coli = i
			break
		}
	}
	sort.Slice(f.lines, func(i, j int) bool { return f.lines[i][coli] < f.lines[j][coli] })
}

func (f *Format) SortNum(col string) {
	coli := 0
	for i, c := range f.cols {
		if c.name == col {
			coli = i
			break
		}
	}
	sort.Slice(f.lines, func(i, j int) bool {
		a, _ := strconv.Atoi(f.lines[i][coli])
		b, _ := strconv.Atoi(f.lines[j][coli])
		return a < b
	})
}

func (f *Format) Print(out io.Writer) {
	for lineno, l := range f.lines {
		line := f.format

		// TODO: we can probably make this a bit faster by getting the text
		// between the placeholders once; this is kind of slow for large sets
		// (1.5s for "uni p all"), partly because we do all of this twice.
		for i, text := range l {
			m := f.re.FindAllString(line, 1)
			line = strings.Replace(line, m[0], f.fmtPlaceholder(i, lineno, text, 0), 1)
		}

		// This line is too long and we want to trim: reformat the lot.
		// TODO: this can be a bit more efficient: we know the column widths and
		// text already, but this is easier.
		if f.ntrim > 0 && zstring.TabWidth(line) > termWidth {
			tooLongBy := zstring.TabWidth(line) - termWidth
			var t = make([]int, len(f.cols))
			for i, text := range l {
				if f.cols[i].trim {
					t[i] = zstring.TabWidth(text)
				}
			}
			trim := nratio(tooLongBy, t...)

			line = f.format
			for i, text := range l {
				m := f.re.FindAllString(line, 1)
				line = strings.Replace(line, m[0], f.fmtPlaceholder(i, lineno, text, trim[i]-1), 1)
			}
		}

		out.Write(append([]byte(line), '\n'))
	}
}

// nratio subtracts "sub" from all the numbers in "nums" proportionally. That
// is, the total subtraction over all the numbers is equal to "sub", but smaller
// numbers get subtracted less.
//
// There's probably some math concept for this, but I don't know what that is
// and I can't really find it...
func nratio(sub int, nums ...int) []int {
	var total int
	for _, n := range nums {
		total += n
	}
	for i := range nums {
		ratio := float64(nums[i]) / float64(total)
		nums[i] -= int(math.Ceil(float64(sub) * ratio))
	}
	// TODO: checked if we substracted more than intended, and add that back.
	return nums
}

func (f *Format) fmtPlaceholder(i, lineno int, text string, applyTrim int) string {
	c := f.cols[i]

	if c.quote {
		if f.printHeader && lineno == 0 {
			text = " " + text + "  " // TODO: why two spaces?
		} else {
			text = "'" + text + "'"
		}
	}

	w := c.width
	if w == alignAuto {
		w = f.autoalign[i]
	}
	switch c.align {
	case alignLeft:
		text = zstring.AlignLeft(text, w)
	case alignRight:
		text = zstring.AlignRight(text, w)
	}

	if c.trim && applyTrim > 0 {
		text = zstring.ElideLeft(text, applyTrim)
	}

	return text
}

func (f *Format) String() string {
	b := new(strings.Builder)
	f.Print(b)
	return b.String()
}

func toLine(info unidata.Codepoint, raw bool) map[string]string {
	c := rune(info.Codepoint)
	return map[string]string{
		"char":         fmtChar(c, raw),
		"wide_padding": widePadding(info),
		"cpoint":       fmt.Sprintf("U+%04X", info.Codepoint),
		"dec":          strconv.FormatUint(uint64(info.Codepoint), 10),
		"utf8":         utf8Bytes(c),
		"html":         htmlEntity(c, info.Codepoint),
		"xml":          fmt.Sprintf("#x%x", info.Codepoint),
		"keysym":       keysym(c),
		"digraph":      unidata.Digraphs[c],
		"name":         info.Name,
		"cat":          unidata.Catnames[info.Cat],
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
	return "\t"
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
	return strings.Join(s, " ")
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
