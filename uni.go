//go:generate go run gen.go data.go

// Command uni prints Unicode information about characters.
package main // import "arp242.net/uni"

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"arp242.net/uni/isatty"
	"arp242.net/uni/terminal"
)

var (
	isTerm       = isatty.IsTerminal(os.Stdout.Fd())
	termWidth    = 0
	errFlag      = errors.New("")
	errNoMatches = errors.New("no matches")
)

func init() {
	if !isTerm {
		return
	}
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err == nil && w > 50 {
		termWidth = w
	}
}

func usage(err error) {
	out := os.Stdout
	e := 0
	if err != nil {
		out = os.Stderr
		e = 1
		if err != errFlag {
			_, _ = fmt.Fprintf(out, "%s: error: %v\n", os.Args[0], err)
		}
	}

	_, _ = fmt.Fprintf(out, `Usage: %s [-hrq] [identify | search | print]

Flags:
    -h      Show this help.
    -q      Quiet output; don't print header, "no matches", etc.
    -r      "Raw" unprocessed output; default is to display graphical variants
            for control characters and display ◌ (U+25CC) before combining
            characters. Note control characters may mangle the output.

Commands:
    identify [string string ... | file:#loc]
        Idenfity all the characters in the given strings. Pass a filename ending
        with :#loc to identify the characters at loc's byte (not character!)
        offset. This can be a range as :#start-end.

    search word [word ...]
        Search description for any of the words.

    print ident [ident ...]
        Print characters by codepoint, category, or block, or special name:

            Codepoint    U+2042
            Range        U+2042..U+2050
            Category     PunctuationOther
            Block        GeneralPunctuation
            all          Everything
            emoji        Alias for "Miscellaneous Symbols",
                         "Emoticons", and "Supplemental Symbols and
                         Pictographs" blocks

        Names are matched case insensitive. Spaces and commas are optional and
        can be replaced with an underscore. "Po", "po", "punction, OTHER",
        "Punctuation_other", and PunctuationOther are all identical.
`, os.Args[0])

	os.Exit(e)
}

func main() {
	var (
		//output string
		quiet bool
		help  bool
		raw   bool
	)
	//flag.StringVar(&output, "o", "human", "")
	flag.BoolVar(&quiet, "q", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&raw, "r", false, "")
	flag.Usage = func() { usage(errFlag) }
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage(errors.New("no command given"))
	}

	var err error
	switch strings.ToLower(args[0]) {
	default:
		usage(fmt.Errorf("unknown command: %q", args[0]))

	case "identify", "i":
		err = identify(getargs(args[1:], quiet), quiet, raw)
	case "search", "s":
		err = search(getargs(args[1:], quiet), quiet, raw)
	case "print", "p":
		err = print(getargs(args[1:], quiet), quiet, raw)
	}
	if err == errNoMatches && quiet {
		err = nil
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Use commandline args or stdin.
func getargs(args []string, quiet bool) []string {
	if len(args) > 0 {
		return args
	}

	if !quiet {
		fmt.Fprintf(os.Stderr, "uni: reading from stdin...\n")
	}
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("read stdin: %s", err))
	}
	return strings.Split(string(stdin), "\n")
}

func search(args []string, quiet, raw bool) error {
	var na []string
	for _, a := range args {
		if a != "" {
			na = append(na, a)
		}
	}
	args = na
	if len(args) == 0 {
		return errors.New("search: need search term")
	}

	var out printer
	words := make([]string, len(args))
	for i := range args {
		words[i] = strings.ToUpper(args[i])
	}
	for _, info := range unidata {
		m := 0
		for _, w := range words {
			if strings.Contains(info.name, w) {
				m++
			}
		}
		if m == len(words) {
			out = append(out, info)
		}
	}

	if len(out) == 0 {
		return errNoMatches
	}

	out.PrintSorted(os.Stdout, quiet, raw)
	return nil
}

// TODO: improve.
func canonCat(cat string) string {
	cat = strings.Replace(cat, " ", "", -1)
	cat = strings.Replace(cat, ",", "", -1)
	cat = strings.Replace(cat, "_", "", -1)
	cat = strings.ToLower(cat)
	return cat
}

func print(args []string, quiet, raw bool) error {
	var out printer

	// TODO: use this data:
	// https://unicode.org/emoji/charts-12.0/full-emoji-list.html
	rm := -1
	for i, a := range args {
		if canonCat(a) == "emoji" {
			rm = i
			break
		}
	}
	if rm > -1 {
		args = append(args[:rm], args[rm+1:]...)
		args = append(args[:rm], append(
			[]string{"Miscellaneous Symbols", "Emoticons", "Supplemental Symbols and Pictographs"},
			args[rm:]...)...)
	}

	for _, a := range args {
		canon := canonCat(a)

		// Print everything.
		if canon == "all" {
			for _, info := range unidata {
				out = append(out, info)
			}
			continue
		}

		// Category name.
		if cat, ok := catmap[canon]; ok {
			for _, info := range unidata {
				if info.cat == cat {
					out = append(out, info)
				}
			}
			continue
		}

		// Block.
		if bl, ok := blockmap[canon]; ok {
			for cp := blocks[bl][0]; cp <= blocks[bl][1]; cp++ {
				s, ok := unidata[fmt.Sprintf("%04X", cp)]
				if ok {
					out = append(out, s)
				}
			}
			continue
		}

		// U2042, U+2042
		// TODO: support 2042
		// TODO: support ranges
		if strings.HasPrefix(canon, "u") {
			canon = strings.TrimLeft(strings.TrimLeft(canon, "u"), "+")
			info, ok := unidata[canon]
			if !ok {
				return fmt.Errorf("unknown codepoint: %q", a)
			}
			out = append(out, info)
			continue
		}

		return fmt.Errorf("unknown identifier: %q", a)
	}

	out.PrintSorted(os.Stdout, quiet, raw)
	return nil
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

func identify(ins []string, quiet, raw bool) error {
	in := strings.Join(ins, "") // TODO: not needed. Loop instead.
	if strings.Contains(in, ":#") {
		in = fromFile(in)
	}

	if !utf8.ValidString(in) {
		_, _ = fmt.Fprintf(os.Stderr, "uni: WARNING: input string is not valid UTF-8\n")
	}

	var out printer
	for _, c := range in {
		find := fmt.Sprintf("%04X", c)
		m := false

		if name := inRange(c); name != "" {
			// TODO: fill in width, category
			out = append(out, char{0, 0, 0, name})
			m = true
		}

		if !m {
			for cp, info := range unidata {
				if cp == find {
					out = append(out, info)
					m = true
					break
				}
			}
			if !m {
				fmt.Printf("NO MATCH: %s\n", find)
			}
		}
	}

	out.Print(os.Stdout, quiet, raw)
	return nil
}

// The UnicodeData.txt file doesn't list every character; some are included as a
// range:
//
//   3400;<CJK Ideograph Extension A, First>;Lo;0;L;;;;;N;;;;;
//   4DB5;<CJK Ideograph Extension A, Last>;Lo;0;L;;;;;N;;;;;
func inRange(c rune) string {
	for i, r := range ranges {
		if c >= r[0] && c <= r[1] {
			return rangeNames[i]
		}
	}
	return ""
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
