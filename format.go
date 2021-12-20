package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"zgo.at/termtext"
	"zgo.at/uni/v2/unidata"
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
	fill  rune
}

type Format struct {
	format    string         // Format string: %(..)
	re        *regexp.Regexp // Cached regexp for format.
	cols      []column       // Columns we know about.
	lines     [][]string     // Processed lines, to be printed.
	autoalign []int          // Max line lengths for autoalign.
	ntrim     int            // Number of columns with "trim"
	json      bool           // Print as JSON.

	printHeader bool
}

func NewFormat(format string, asJSON, printHeader bool, knownCols ...string) (*Format, error) {
	var (
		reFindCols = regexp.MustCompile(`%\((.*?)(?: .+?)?\)`)
		f          = Format{format: format, printHeader: printHeader, json: asJSON}
	)
	for _, m := range reFindCols.FindAllString(format, -1) {
		err := f.processColumn(m)
		if err != nil {
			return nil, fmt.Errorf("-format flag: %w", err)
		}
	}

	f.autoalign = make([]int, len(f.cols))

	cols := make([]string, 0, len(f.cols))
	h := map[string]string{}
	for _, c := range f.cols {
		if !zstring.Contains(knownCols, c.name) {
			return nil, fmt.Errorf("-format flag: unknown placeholder: %q", c.name)
		}
		cols = append(cols, c.name)
		h[c.name] = c.name
	}

	h["emoji"] = ""
	h["char"] = ""
	h["tab"] = tabOrSpace()
	h["wide_padding"] = " "

	if printHeader && !asJSON {
		f.Line(h)
	}

	// TODO: is this actually faster than just .*?
	// TODO: don't really need to use regexp for this; can just scan for "%(name".
	f.re = regexp.MustCompile(`%\((` + strings.Join(cols, "|") + `)(?: .+?)?\)`)
	return &f, nil
}

func (f *Format) processColumn(line string) error {
	s := zstring.Fields(line[2:len(line)-1], " ") // name, flags
	name := s[0]
	col := column{name: name}

	if len(s) == 1 { // No flags
		f.cols = append(f.cols, col)
		return nil
	}

	for _, flag := range s[1:] {
		switch {
		default:
			return fmt.Errorf("unknown flag %q in %q", flag, line)
		case flag == "":
			continue
		case flag == "q":
			col.quote = true
		case flag[0] == 'f':
			if utf8.RuneCountInString(flag) != 3 {
				return fmt.Errorf("need exactly one character after f: for %q", line)
			}
			col.fill = []rune(flag[2:])[0]
		case flag == "t":
			f.ntrim++
			col.trim = true
		case flag[0] == 'l' || flag[0] == 'r':
			n := strings.Split(flag, ":")
			if len(n) != 2 {
				return fmt.Errorf("need width after : for %q", line)
			}

			col.align = alignLeft
			if flag[0] == 'r' {
				col.align = alignRight
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
					return fmt.Errorf(`width needs to be a number or "auto" in %q`, line)
				}
			}
		}
	}

	f.cols = append(f.cols, col)
	return nil
}

// Add a new line.
func (f *Format) Line(columns map[string]string) error {
	line := make([]string, len(f.cols))
	for i, c := range f.cols {
		line[i] = columns[c.name]
		if c.width == alignAuto {
			if l := termtext.Width(columns[c.name]); l > f.autoalign[i] {
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
	sort.Slice(f.lines, func(i, j int) bool {
		return f.lines[i][coli] < f.lines[j][coli]
	})
}

func (f *Format) SortNum(col string) {
	if len(f.cols) == 0 {
		return
	}

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

func (f *Format) printJSON(out io.Writer) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")

	out.Write([]byte("["))
	for i, l := range f.lines {
		m := make(map[string]string, len(f.cols))
		for j, c := range f.cols {
			if c.name == "wide_padding" || c.name == "tab" {
				continue
			}
			m[c.name] = l[j]
		}

		enc.Encode(m)
		out.Write(bytes.TrimSpace(buf.Bytes())) // Adds \n at end.
		buf.Reset()

		// j, _ := json.MarshalIndent(m, "", "\t")
		// out.Write(j)
		if i != len(f.lines)-1 {
			out.Write([]byte(", "))
		}
	}
	out.Write([]byte("]\n"))
}

func (f *Format) Print(out io.Writer) {
	if f.json {
		f.printJSON(out)
		return
	}

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
		if f.ntrim > 0 && termtext.Width(line) > termWidth {
			tooLongBy := termtext.Width(line) - termWidth
			var t = make([]int, len(f.cols))
			for i, text := range l {
				if f.cols[i].trim {
					t[i] = termtext.Width(text)
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
	if c.fill > 0 {
		if !f.printHeader || lineno > 0 {
			text = strings.ReplaceAll(text, " ", string(c.fill))
		}
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

var knownColumns = []string{"char", "wide_padding", "cpoint", "dec", "hex",
	"oct", "bin", "utf8", "utf16be", "utf16le", "html", "xml", "json", "keysym",
	"digraph", "name", "cat", "block", "plane", "width"}

func toLine(info unidata.Codepoint, raw bool) map[string]string {
	// TODO: would be better to include only the columns that are actually used.
	return map[string]string{
		"char":         info.Repr(raw),
		"wide_padding": widePadding(info),
		"cpoint":       info.FormatCodepoint(),
		"dec":          info.Format(10),
		"hex":          info.Format(16),
		"oct":          info.Format(8),
		"bin":          info.Format(2),
		"utf8":         info.UTF8(),
		"utf16be":      info.UTF16(true),
		"utf16le":      info.UTF16(false),
		"html":         info.HTMLEntity(),
		"xml":          info.XMLEntity(),
		"json":         info.JSON(),
		"keysym":       info.KeySym,
		"digraph":      info.Digraph,
		"name":         info.Name,
		"cat":          info.Category(),
		"block":        info.Block(),
		"plane":        info.Plane(),
		"width":        info.WidthName(),
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
//
// TODO: make this an option; -expandtab or some such.
func tabOrSpace() string {
	if os.Getenv("README") != "" { // Temporary hack for README generation.
		return "\t"
	}
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
