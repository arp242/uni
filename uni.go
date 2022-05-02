//go:generate go run gen.go

// Command uni prints Unicode information about characters.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"zgo.at/uni/v2/unidata"
	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

var (
	errNoMatches = errors.New("no matches")
	version      = "git"
)

var usageShort = zli.Usage(zli.UsageHeaders|zli.UsageProgram|zli.UsageTrim, `
Usage: %(prog) [command] [flags]

uni queries the unicode database. https://github.com/arp242/uni

Flags:
    -q, -quiet     Quiet output.
    -r, -raw       Don't use graphical variants or add combining characters.
    -p, -pager     Output to $PAGER.
    -o, -or        Use "or" when searching instead of "and".
    -f, -format    Output format.
    -j, -json      Output as JSON.

Commands:
    identify       Identify all the characters in the given strings.
    search         Search description for any of the words.
    print          Print characters by codepoint, category, or block.
    emoji          Search emojis.

Use "%(prog) help" for a more detailed help.
`)

var usage = zli.Usage(zli.UsageHeaders|zli.UsageProgram|zli.UsageTrim, `
Usage: %(prog) [command] [flags] 

uni queries the unicode database. https://github.com/arp242/uni

Flags:
    Flags can appear anywhere; "uni search euro -q" and "uni -q search euro"
    are identical. Use "uni search -- -q" if you want to search for "-q". This
    also applies to flags specific to a command (e.g. "-gender" for "emoji").

    -q, -quiet     Quiet output; don't print header, "no matches", etc.
    -r, -raw       Don't use graphical variants for control characters and
                   don't add ◌ (U+25CC) before combining characters.
    -p, -pager     Output to $PAGER.
    -o, -or        Use "or" when searching: print if at least one parameter
                   matches, instead of only when all parameters match.
    -f, -format    Output format; see Format section below for details.
    -j, -json      Output as JSON; the columns listed in -format are included,
                   ignoring formatting flags. Use "-format all" to include all
                   columns.

Commands:
    identify [text]  Identify all the characters in the given strings.

    search [query]   Search description for any of the words.

    print [query]    Print characters by codepoint, category, or block.

                       Codepoints             U+20, U20, 0x20, 0d32 (decimal),
                                              0o40, 0b100000
                       Range                  U+2042..U+2050, 0o101..0x5a
                       Categories and Blocks  OtherPunctuation, Po,
                                              GeneralPunctuation
                       UTF-8                  UTF-8 byte sequence, optionally
                                              separated by any combination of
                                              '0x', '-', '_', or spaces:
                                                utf8:e282ac
                                                utf8:0xe20x820xac
                                                'utf8:e2 82 ac'
                                                'utf8:0xe2 0x82 0xac'
                       all                    Everything

    emoji [query]    Search emojis.

                     The query is matched on the emoji name and CLDR data.

                     The CLDR data is a list of keywords. For example 🙏
                     (folded hands) contains "ask, high 5, high five, please,
                     pray, thanks", which represents the various scenarios in
                     which it's used.

                     You can use <prefix>:query to search in specific fields:

                       group: g:    Group and subgroup
                       name:  n:    Emoji name
                       cldr:  c:    CLDR data

                     The query parameters are AND'd together, so this:

                       uni emoji smiling g:cat-face

                     Will match everything in the cat-face group with smiling
                     in the name. Use the -or flag to change this to "cat-face
                     group OR smiling in the name".

                     Use "all" to show all emojis.

                     Modifier flags, both accept a comma-separated list:

                        -g, -gender   Set the gender:
                                        p, person, people
                                        m, man, men, male
                                        f, female, w, woman, women

                        -t, -tone     Set the skin tone modifier:
                                        n,  none
                                        l,  light
                                        ml, mediumlight, medium-light
                                        m,  medium
                                        md, mediumdark, medium-dark
                                        d,  dark

                     Use "all" to include all combinations; the default is to
                     include no skin tones and the "person" gender.

                     Note: emojis may not be accurately copied by select & copy
                     in terminals. It's recommended to copy to the clipboard
                     directly with e.g. xclip.

Format:
    You can use the -format or -f flag to control what to print; placeholders
    are in the form of %(name) or %(name flags), where "name" is a column name
    and "flags" are some flags to control how it's printed.

    The special value "all" includes all columns; this is useful especially
    with -json if you want to get all information uni knows about a codepoint
    or emoji.

    Flags:
        %(name l:5)     Left-align and pad with 5 spaces
        %(name l:auto)  Left-align and pad to the longest value
        %(name r:5)     Right-align and pad with 5 spaces
        %(name q)       Quote with single quotes, excluding any padding
        %(name t)       Trim this column if it's longer than the screen width
        %(name f:C)     Fill this column with character C; especially useful for
                        numbers: %(bin r:auto f:0)

    Placeholders that work for all commands:
        %(tab)           A literal tab when outputting to a terminal, or four
                         spaces if not; this helps with aligning emojis in
                         terminals, but some tools like dmenu don't work well
                         with tabs.

    Placeholders for identify, search, and print:
        %(char)          The literal character          ✓
        %(cpoint)        As codepoint                   U+2713
        %(hex)           As hex                         2713
        %(oct)           As octal                       23423
        %(bin)           As binary (little-endian)      10011100010011
        %(dec)           As decimal                     10003
        %(utf8)          As UTF-8                       e2 9c 93
        %(utf16le)       As UTF-16 LE (Windows)         13 27
        %(utf16be)       As UTF-16 BE                   27 13
        %(html)          HTML entity                    &check;
        %(xml)           XML entity                     &#x2713;
        %(json)          JSON escape                    \u2713
        %(keysym)        X11 keysym; can be blank       checkmark
        %(digraph)       Vim Digraph; can be blank      OK
        %(name)          Code point name                CHECK MARK
        %(cat)           Category name                  Other_Symbol
        %(block)         Block name                     Dingbats
        %(props)         Properties, separated by ,     Pattern Syntax
        %(plane)         Plane name                     Basic Multilingual Plane
        %(width)         Character width                Narrow
        %(wide_padding)  Blank for wide characters,
                         space otherwise; for alignment

        The default is:
        `+defaultFormat+`

    Placeholders for emoji:

        %(emoji)       The emoji itself                 🧑‍🚒
        %(name)        Emoji name                       firefighter
        %(group)       Emoji group                      People & Body
        %(subgroup)    Emoji subgroup                   person-role
        %(cpoint)      Codepoints                       U+1F9D1 U+200D U+1F692
        %(cldr)        CLDR data, w/o duplicating name  firetruck
        %(cldr_full)   Full CLDR data                   firefighter, firetruck

        The default is:
        `+defaultEmojiFormat+`
`)

const (
	defaultFormat = "%(char q l:3)%(wide_padding) %(cpoint l:7) %(dec l:6) %(utf8 l:11) %(html l:10) %(name t) (%(cat t))"
	allFormat     = "%(char q l:3)%(wide_padding) %(cpoint l:auto) %(width l:auto) %(dec l:auto) %(hex l:auto)" +
		" %(oct l:auto) %(bin l:auto)" +
		" %(utf8 l:auto) %(utf16le l:auto) %(utf16be l:auto) %(html l:auto) %(xml l:auto) %(json l:auto)" +
		" %(keysym l:auto) %(digraph l:auto) %(name l:auto) %(plane l:auto) %(cat l:auto) %(block l:auto)" +
		" %(props l:auto)"

	defaultEmojiFormat = "%(emoji)%(tab)%(name l:auto)  (%(cldr t))"
	allEmojiFormat     = "%(emoji)%(tab)%(name l:auto) %(group l:auto) %(subgroup l:auto) %(cpoint l:auto) %(cldr l:auto) %(cldr_full)"
)

func main() {
	flag := zli.NewFlags(os.Args)
	var (
		quietF   = flag.Bool(false, "q", "quiet")
		help     = flag.Bool(false, "h", "help")
		versionF = flag.Bool(false, "v", "version")
		rawF     = flag.Bool(false, "r", "raw")
		pager    = flag.Bool(false, "p", "pager")
		or       = flag.Bool(false, "o", "or")
		formatF  = flag.String("", "format", "f")
		jsonF    = flag.Bool(false, "json", "j")
		tone     = flag.String("", "t", "tone", "tones")
		gender   = flag.String("person", "g", "gender", "genders")
	)
	err := flag.Parse()
	zli.F(err)

	if versionF.Set() {
		fmt.Println(version)
		return
	}

	if pager.Set() {
		defer zli.PagerStdout()()
	}

	if help.Set() {
		fmt.Fprint(zli.Stdout, usage)
		return
	}

	cmd := flag.ShiftCommand("identify", "print", "search", "emoji", "help", "version")
	switch cmd { // These commands don't read from stdin.
	case zli.CommandNoneGiven:
		fmt.Fprint(zli.Stdout, usageShort)
		return
	case zli.CommandUnknown:
		zli.Fatalf("unknown command")
		return
	case "help":
		fmt.Fprint(zli.Stdout, usage)
		return
	case "version":
		fmt.Println(version)
		return
	}

	quiet := quietF.Set()
	raw := rawF.Set()
	args := flag.Args
	args, err = zli.InputOrArgs(args, " \t\n", quiet)
	zli.F(err)

	format := formatF.String()
	if !formatF.Set() {
		format = defaultFormat
		if cmd == "emoji" {
			format = defaultEmojiFormat
		}
	}

	if formatF.String() == "all" {
		format = allFormat
		if cmd == "emoji" {
			format = allEmojiFormat
		}
	}

	switch cmd {
	case "identify":
		err = identify(args, format, quiet, raw, jsonF.Bool())
	case "search":
		err = search(args, format, quiet, raw, jsonF.Bool(), or.Bool())
	case "print":
		err = print(args, format, quiet, raw, jsonF.Bool())
	case "emoji":
		err = emoji(args, format, quiet, raw, jsonF.Bool(), or.Bool(),
			parseToneFlag(tone.String()), parseGenderFlag(gender.String()))
	}
	if err != nil {
		if !(err == errNoMatches && quiet) {
			zli.Fatalf(err)
		}
		zli.Exit(1)
	}
}

func parseToneFlag(tone string) unidata.EmojiModifier {
	if tone == "" {
		return 0
	}
	if tone == "all" {
		tone = "none,light,mediumlight,medium,mediumdark,dark"
	}

	var m unidata.EmojiModifier
	for _, t := range zstring.Fields(tone, ",") {
		switch t {
		case "none", "n":
			m |= unidata.ModNone
		case "l", "light":
			m |= unidata.ModLight
		case "ml", "mediumlight", "medium-light", "medium_light":
			m |= unidata.ModMediumLight
		case "m", "medium":
			m |= unidata.ModMedium
		case "md", "mediumdark", "medium-dark", "medium_dark":
			m |= unidata.ModMediumDark
		case "d", "dark":
			m |= unidata.ModDark
		default:
			zli.Fatalf("invalid skin tone: %q", tone)
		}
	}
	return m
}

func parseGenderFlag(gender string) unidata.EmojiModifier {
	if gender == "" {
		return 0
	}
	if gender == "all" {
		gender = "person,man,woman"
	}

	var m unidata.EmojiModifier
	for _, g := range zstring.Fields(gender, ",") {
		switch g {
		case "person", "p", "people":
			m |= unidata.ModPerson
		case "man", "men", "m", "male":
			m |= unidata.ModMale
		case "woman", "women", "w", "female", "f":
			m |= unidata.ModFemale
		default:
			zli.Fatalf("invalid gender: %q", gender)
		}
	}
	return m
}

func identify(ins []string, format string, quiet, raw, asJSON bool) error {
	in := strings.Join(ins, "")
	if !utf8.ValidString(in) {
		fmt.Fprintf(zli.Stderr, "uni: WARNING: input string is not valid UTF-8\n")
	}

	f, err := NewFormat(format, asJSON, !quiet, knownColumns...)
	if err != nil {
		return err
	}
	for _, c := range in {
		info, ok := unidata.Find(c)
		if !ok {
			return fmt.Errorf("unknown codepoint: U+%.4X", c) // Should never happen.
		}

		f.Line(toLine(info, raw))
	}
	f.Print(zli.Stdout)
	return nil
}

func search(args []string, format string, quiet, raw, asJSON, or bool) error {
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

	for i := range args {
		args[i] = strings.ToUpper(args[i])
	}

	found := false
	f, err := NewFormat(format, asJSON, !quiet, knownColumns...)
	if err != nil {
		return err
	}
	for _, info := range unidata.Codepoints {
		m := 0
		for _, a := range args {
			if strings.Contains(info.Name(), a) {
				if or {
					found = true
					f.Line(toLine(info, raw))
					break
				}
				m++
			}
		}
		if !or && m == len(args) {
			found = true
			f.Line(toLine(info, raw))
		}
	}

	if !found {
		return errNoMatches
	}
	f.SortNum("dec")
	f.Print(zli.Stdout)
	return nil
}

var utfClean = strings.NewReplacer("0x", "", " ", "", "_", "", "-", "")

func print(args []string, format string, quiet, raw, asJSON bool) error {
	f, err := NewFormat(format, asJSON, !quiet, knownColumns...)
	if err != nil {
		return err
	}
	for _, a := range args {
		a = strings.ToLower(a)

		// UTF-8
		if strings.HasPrefix(a, "utf8:") {
			a = a[5:]

			seq := utfClean.Replace(a)
			if len(seq)%2 == 1 {
				seq = "0" + seq
			}

			byt := make([]byte, 0, len(seq)/2)
			for i := 0; len(seq) > i; i += 2 {
				b, err := strconv.ParseUint(seq[i:i+2], 16, 8)
				if err != nil {
					return fmt.Errorf("invalid UTF-8 sequence %q: %q is not a hex number",
						a, seq[i:i+2])
				}
				byt = append(byt, byte(b))
			}

			r, s := utf8.DecodeRune(byt)
			if r == utf8.RuneError {
				return fmt.Errorf("invalid UTF-8 sequence: %q", a)
			}
			if s != len(byt) {
				return fmt.Errorf("multiple characters in sequence %q", a)
			}

			f.Line(toLine(unidata.Codepoints[r], raw))
			continue
		}

		// Print everything.
		if strings.ToLower(a) == "all" {
			for _, info := range unidata.Codepoints {
				f.Line(toLine(info, raw))
			}
			continue
		}

		// Category name.
		if cat, ok := unidata.FindCategory(a); ok {
			for _, info := range unidata.Codepoints {
				if info.Category() == cat {
					f.Line(toLine(info, raw))
				}
			}
			continue
		}

		// Block.
		if bl, ok := unidata.FindBlock(a); ok {
			for cp := unidata.Blocks[bl].Range[0]; cp <= unidata.Blocks[bl].Range[1]; cp++ {
				s, ok := unidata.Codepoints[cp]
				if ok {
					f.Line(toLine(s, raw))
				}
			}
			continue
		}

		// Properties
		if p, ok := unidata.FindProperty(a); ok {
			for _, pp := range unidata.Properties[p].Ranges {
				for cp := pp[0]; cp <= pp[1]; cp++ {
					s, ok := unidata.Codepoints[cp]
					if ok {
						f.Line(toLine(s, raw))
					}
				}
			}
			continue
		}

		// U2042, U+2042, U+2042..U+2050, 2042..2050, 2042-2050, 0x2041, etc.
		var s []string
		switch {
		case strings.Contains(a, ".."):
			s = strings.SplitN(a, "..", 2)
		case strings.Contains(a, "-"):
			s = strings.SplitN(a, "-", 2)
		default:
			s = []string{a, a}
		}

		start, err := unidata.FromString(s[0])
		if err != nil {
			return fmt.Errorf("invalid codepoint: %s", errors.Unwrap(err))
		}
		end, err := unidata.FromString(s[1])
		if err != nil {
			return fmt.Errorf("invalid codepoint: %s", errors.Unwrap(err))
		}

		for i := start.Codepoint; i <= end.Codepoint; i++ {
			info, _ := unidata.Find(i)
			f.Line(toLine(info, raw))
		}
	}
	f.SortNum("dec")
	f.Print(zli.Stdout)
	return nil
}

func emoji(args []string, format string, quiet, raw, asJSON, or bool, tones, genders unidata.EmojiModifier) error {
	type matchArg struct {
		group bool
		name  bool
		text  string
	}
	var (
		all       = zstring.Contains(args, "all")
		matchArgs = make([]matchArg, 0, len(args))
	)
	if !all {
		for _, a := range args {
			a := strings.ToLower(a)
			group := strings.HasPrefix(a, "g:") || strings.HasPrefix(a, "group:")
			if group {
				a = strings.TrimPrefix(strings.TrimPrefix(a, "group:"), "g:")
			}
			name := strings.HasPrefix(a, "n:") || strings.HasPrefix(a, "name:")
			if name {
				a = strings.TrimPrefix(strings.TrimPrefix(a, "name:"), "n:")
			}
			matchArgs = append(matchArgs, matchArg{text: a, group: group, name: name})
		}
	}

	out := make([]unidata.Emoji, 0, 16)
	for _, e := range unidata.Emojis {
		m := 0
		for _, a := range matchArgs {
			var match bool
			switch {
			case a.group:
				match = strings.Contains(strings.ToLower(e.Group().String()), a.text) ||
					strings.Contains(strings.ToLower(e.Subgroup().String()), a.text)
			case a.name:
				match = strings.Contains(strings.ToLower(e.Name), a.text)
			default:
				match = strings.Contains(strings.ToLower(e.Name), a.text) ||
					zstring.Contains(e.CLDR, a.text)
			}
			if match {
				m++
				if or {
					out = append(out, e)
					break
				}
			}
		}
		if all || (!or && m == len(matchArgs)) {
			out = append(out, applyGenders(applyTones(e, tones), genders)...)
		}
	}

	if len(out) == 0 {
		return errNoMatches
	}

	f, err := NewFormat(format, asJSON, !quiet, "emoji", "name", "group", "subgroup",
		"tab", "cldr", "cldr_full", "cpoint")
	if err != nil {
		return err
	}
	for _, e := range out {
		f.Line(map[string]string{
			"emoji":    e.String(),
			"name":     e.Name,
			"group":    e.Group().String(),
			"subgroup": e.Subgroup().String(),
			"tab":      tabOrSpace(),
			"cldr": func() string {
				// Remove words that duplicate what's already in the name; it's
				// kind of pointless.
				cldr := make([]string, 0, len(e.CLDR))
				for _, c := range e.CLDR {
					if !strings.Contains(e.Name, c) {
						cldr = append(cldr, c)
					}
				}
				return strings.Join(cldr, ", ")
			}(),
			"cldr_full": strings.Join(e.CLDR, ", "),
			"cpoint": func() string {
				cp := make([]string, 0, len(e.Codepoints))
				for _, c := range e.String() { // String() inserts ZWJ and whatnot
					cp = append(cp, fmt.Sprintf("U+%04X", c))
				}
				return strings.Join(cp, " ")
			}(),
		})
	}
	f.Print(zli.Stdout)
	return nil
}

func applyAll(e unidata.Emoji, mod unidata.EmojiModifier) []unidata.Emoji {
	emojis := make([]unidata.Emoji, 0, 1)
	i := unidata.EmojiModifier(1)
	for i <= unidata.ModDark {
		if mod&i != 0 {
			emojis = append(emojis, e.With(i))
		}
		i <<= 1
	}
	return emojis
}

func applyTones(e unidata.Emoji, mod unidata.EmojiModifier) []unidata.Emoji {
	if !e.Skintones() || mod == 0 {
		return []unidata.Emoji{e}
	}
	return applyAll(e, mod)
}

func applyGenders(emojis []unidata.Emoji, mod unidata.EmojiModifier) []unidata.Emoji {
	if mod == 0 {
		return emojis
	}

	var ret []unidata.Emoji
	for _, e := range emojis {
		if !e.Genders() {
			ret = append(ret, e)
			continue
		}
		ret = append(ret, applyAll(e, mod)...)
	}
	return ret
}
