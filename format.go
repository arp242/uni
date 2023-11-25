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

	"github.com/mattn/go-runewidth"
	"zgo.at/termtext"
	"zgo.at/uni/v2/unidata"
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
	as        printAs        // How to print (list, table, json)
	re        *regexp.Regexp // Cached regexp for format.
	cols      []column       // Columns we know about.
	colNames  []string
	lines     [][]string // Processed lines, to be printed.
	autoalign []int      // Max line lengths for autoalign.
	ntrim     int        // Number of columns with "trim"

	tblData []unidata.Codepoint
}

var reFindCols = regexp.MustCompile(`%\((.*?)(?: .+?)?\)`)

type printAs uint8

const (
	// Keep compact regular + 1
	printAsList = printAs(iota)
	printAsListCompact
	printAsJSON
	printAsJSONCompact
	printAsTable
	printAsTableCompact
)

func header(h string) string {
	h = strings.ReplaceAll(h, "-", " ")
	switch h {
	case "utf8", "utf16", "utf16le", "utf16be", "html", "xml", "json":
		return strings.ToUpper(h)
	case "cpoint":
		return "CPoint"
	default:
		return zstring.UpperFirst(h)
	}
}

func NewFormat(format string, as printAs, knownCols ...string) (*Format, error) {
	f := Format{format: format, as: as}

	if as == printAsTable || as == printAsTableCompact {
		// Don't need all the rest of the logic.
		return &f, nil
	}

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

		if f.json() {
			h[c.name] = c.name
		} else {
			h[c.name] = header(c.name)
		}
	}

	h["emoji"] = ""
	h["char"] = ""
	h["tab"] = tabOrSpace()
	h["wide_padding"] = " "

	if as == printAsList {
		f.Line(h)
	}

	// TODO: is this actually faster than just .*?
	// TODO: don't really need to use regexp for this; can just scan for "%(name".
	f.re = regexp.MustCompile(`%\((` + strings.Join(cols, "|") + `)(?: .+?)?\)`)
	return &f, nil
}

func (f *Format) tbl() bool  { return f.as == printAsTable || f.as == printAsTableCompact }
func (f *Format) json() bool { return f.as == printAsJSON || f.as == printAsJSONCompact }

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
	if f.tbl() { // Don't need to do anything.
		return nil
	}

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
	if f.as != printAsJSONCompact {
		enc.SetIndent("", "\t")
	}

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

		if i != len(f.lines)-1 {
			out.Write([]byte(", "))
		}
	}
	out.Write([]byte("]\n"))
}

func (f *Format) printTbl(out io.Writer) {
	sort.Slice(f.tblData, func(i, j int) bool {
		return f.tblData[i].Codepoint < f.tblData[j].Codepoint
	})

	var (
		wide   = false
		tblMap = make(map[rune]string)
		head   = 4
	)
	for _, c := range f.tblData {
		if w := c.Width(); w == unidata.WidthFullWidth || w == unidata.WidthWide || w == unidata.WidthAmbiguous {
			wide = true
		}
		tblMap[c.Codepoint] = c.Display()
		if c.Codepoint > 0xffff {
			head = 5
		}
	}

	h := strings.Repeat(" ", head)
	if wide {
		fmt.Println(h + "     0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F")
		fmt.Println(h + "   ┌" + strings.Repeat("─", 64))
	} else {
		fmt.Println(h + "     0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F")
		fmt.Println(h + "   ┌" + strings.Repeat("─", 48))
	}

	start, end := f.tblData[0].Codepoint, f.tblData[len(f.tblData)-1].Codepoint
	start -= start % 16 /// Make sure we start at column 0

	var (
		row   = ""
		blank = 0
		didel = false
	)
	for i := start; i <= end; i++ {
		cp := fmt.Sprintf("U+%0"+strconv.Itoa(head)+"X", i)
		char, ok := tblMap[i]
		if _, has := unidata.Codepoints[i]; !has { /// Not assigned
			if isTerm {
				char = zli.Colorize(" ", zli.Color256(254).Bg())
			} else {
				char = "·"
			}
			blank++
		} else if !ok { /// Not in selection.
			blank++
			char = zli.Colorize("·", zli.Color256(249))
		}

		if strings.HasSuffix(cp, "0") {
			row += fmt.Sprintf("%sx │", strings.TrimSuffix(cp, "0"))
		}

		// Need to add space for alignment if some codepoints are wide but this
		// one isn't.
		row += fmt.Sprintf(" %s ", char)
		if wide && runewidth.RuneWidth([]rune(char)[0]) == 1 {
			row += " "
		}
		if strings.HasSuffix(cp, "F") || i == end {
			if blank < 16 {
				fmt.Print(row)
				fmt.Print("\n")
				if f.as != printAsTableCompact {
					fmt.Println(h + "   │")
				}
				didel = false
			} else if !didel {
				fmt.Println(h + " … │")
				didel = true
			}

			row = ""
			blank = 0
		}
	}
}

func (f *Format) Print(out io.Writer) {
	if f.json() {
		f.printJSON(out)
		return
	}
	if f.tbl() {
		f.printTbl(out)
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
		w := termtext.Width(line)
		if f.ntrim > 0 && w > termWidth {
			tooLongBy := w - termWidth
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
		if f.as != printAsListCompact && lineno == 0 {
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
		if f.as == printAsListCompact || lineno > 0 {
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
	"digraph", "name", "cat", "block", "plane", "width", "props", "script", "unicode"}

func (f *Format) toLine(info unidata.Codepoint, raw bool) map[string]string {
	if f.tbl() {
		f.tblData = append(f.tblData, info)
		return nil
	}

	if len(f.cols) == len(knownColumns) { // Optimize printing all columns.
		return map[string]string{
			"char":         map[bool]string{false: info.Display(), true: string(info.Codepoint)}[raw],
			"wide_padding": widePadding(info),
			"cpoint":       info.FormatCodepoint(),
			"dec":          info.Format(10),
			"hex":          info.Format(16),
			"oct":          info.Format(8),
			"bin":          info.Format(2),
			"utf8":         fmt.Sprintf("% x", info.UTF8()),
			"utf16be":      fmt.Sprintf("% x", info.UTF16(true)),
			"utf16le":      fmt.Sprintf("% x", info.UTF16(false)),
			"html":         info.HTML(),
			"xml":          info.XML(),
			"json":         info.JSON(),
			"keysym":       info.KeySym(),
			"digraph":      info.Digraph(),
			"name":         info.Name(),
			"cat":          info.Category().String(),
			"block":        info.Block().String(),
			"plane":        info.Plane().String(),
			"width":        info.Width().String(),
			"props":        info.Properties().String(),
			"script":       info.Script().String(),
			"unicode":      info.Unicode().String(),
		}
	}

	if f.colNames == nil {
		f.colNames = make([]string, len(f.cols))
		for _, c := range f.cols {
			f.colNames = append(f.colNames, c.name)
		}
	}

	cols := make(map[string]string)
	if zstring.Contains(f.colNames, "char") {
		cols["char"] = map[bool]string{false: info.Display(), true: string(info.Codepoint)}[raw]
	}
	if zstring.Contains(f.colNames, "wide_padding") {
		cols["wide_padding"] = widePadding(info)
	}
	if zstring.Contains(f.colNames, "cpoint") {
		cols["cpoint"] = info.FormatCodepoint()
	}
	if zstring.Contains(f.colNames, "dec") {
		cols["dec"] = info.Format(10)
	}
	if zstring.Contains(f.colNames, "hex") {
		cols["hex"] = info.Format(16)
	}
	if zstring.Contains(f.colNames, "oct") {
		cols["oct"] = info.Format(8)
	}
	if zstring.Contains(f.colNames, "bin") {
		cols["bin"] = info.Format(2)
	}
	if zstring.Contains(f.colNames, "utf8") {
		cols["utf8"] = fmt.Sprintf("% x", info.UTF8())
	}
	if zstring.Contains(f.colNames, "utf16be") {
		cols["utf16be"] = fmt.Sprintf("% x", info.UTF16(true))
	}
	if zstring.Contains(f.colNames, "utf16le") {
		cols["utf16le"] = fmt.Sprintf("% x", info.UTF16(false))
	}
	if zstring.Contains(f.colNames, "html") {
		cols["html"] = info.HTML()
	}
	if zstring.Contains(f.colNames, "xml") {
		cols["xml"] = info.XML()
	}
	if zstring.Contains(f.colNames, "json") {
		cols["json"] = info.JSON()
	}
	if zstring.Contains(f.colNames, "keysym") {
		cols["keysym"] = info.KeySym()
	}
	if zstring.Contains(f.colNames, "digraph") {
		cols["digraph"] = info.Digraph()
	}
	if zstring.Contains(f.colNames, "name") {
		cols["name"] = info.Name()
	}
	if zstring.Contains(f.colNames, "cat") {
		cols["cat"] = info.Category().String()
	}
	if zstring.Contains(f.colNames, "block") {
		cols["block"] = info.Block().String()
	}
	if zstring.Contains(f.colNames, "script") {
		cols["script"] = info.Script().String()
	}
	if zstring.Contains(f.colNames, "plane") {
		cols["plane"] = info.Plane().String()
	}
	if zstring.Contains(f.colNames, "width") {
		cols["width"] = info.Width().String()
	}
	if zstring.Contains(f.colNames, "props") {
		cols["props"] = info.Properties().String()
	}
	if zstring.Contains(f.colNames, "unicode") {
		cols["unicode"] = info.Unicode().String()
	}
	return cols
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
	if isTerm {
		return "\t"
	}
	return "    "
}

func widePadding(info unidata.Codepoint) string {
	if info.Width() != unidata.WidthFullWidth && info.Width() != unidata.WidthWide {
		return " "
	}
	return ""
}
