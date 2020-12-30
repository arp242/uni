// Command uni prints Unicode information about characters.
package main // import "arp242.net/uni"

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"arp242.net/uni/unidata"
	"zgo.at/zli"
	"zgo.at/zstd/zstring"
)

var (
	errNoMatches = errors.New("no matches")
	version      = "git"
)

var usageShort = zli.Usage(zli.UsageHeaders|zli.UsageProgram|zli.UsageTrim, `
Usage: %(prog) [help | version | identify | search | print | emoji] [-qrpao] 

Flags:
    -q, -quiet     Quiet output.
    -r, -raw       Don't use graphical variants or add combining characters.
    -p, -pager     Output to $PAGER.
    -o, -or        Use "or" when searching instead of "and".

Commands:
    identify       Idenfity all the characters in the given strings.
    search         Search description for any of the words.
    print          Print characters by codepoint, category, or block.
    emoji          Search emojis.

Use "%(prog) help" for a more detailed help.
`)

var usage = zli.Usage(zli.UsageHeaders|zli.UsageProgram|zli.UsageTrim, `
Usage: %(prog) [help | version | identify | search | print | emoji] [-qrpao]

Flags:
    -q, -quiet     Quiet output; don't print header, "no matches", etc.
    -r, -raw       Don't use graphical variants for control characters and
                   don't add ◌ (U+25CC) before combining characters.
    -p, -pager     Output to $PAGER.
    -o, -or        Use "or" when searching: only match if all parameters match,
                   instead of anything where at least one matches.

Commands:
    identify [string..]    Idenfity all the characters in the given strings.

    search   [query..]     Search description for any of the words.

    print    [ident..]     Print characters by codepoint, category, or block.

                           Codepoints             U+2042, 0x20, 0o40, 0b100000
                           Range                  U+2042..U+2050, 0o101..0x5a
                           Categories and Blocks  OtherPunctuation, Po, GeneralPunctuation
                           all                    Everything

    emoji [..] [query..]   Search emojis.

                           The query parameters are matched on the emoji name.
                           Parameters prefixed with "group:" or "g:" are matched
                           on the emoji group name. Use "all" to show all
                           emojis. The query parameters are AND'd together.

                           Modifier flags, both accept a comma-separated list:
                               -tone    light, mediumlight, medium, mediumdark, dark
                               -gender  person, man, woman.
                           Default is to include no skin tones and the "person" gender.

                           Note: emojis may not be accurately copied by select &
                           copy in terminals. It's recommended to copy to the
                           clipboard directly with e.g. xclip.
`)

func main() {
	flag := zli.NewFlags(os.Args)
	var (
		quietF   = flag.Bool(false, "q", "quiet")
		help     = flag.Bool(false, "h", "help")
		versionF = flag.Bool(false, "v", "version")
		rawF     = flag.Bool(false, "r", "raw")
		pager    = flag.Bool(false, "p", "pager")
		or       = flag.Bool(false, "o", "or")
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

	cmd := strings.ToLower(flag.Shift())
	if cmd == "" {
		fmt.Fprint(zli.Stdout, usageShort)
		return
	}
	if cmd == "h" || cmd == "help" {
		fmt.Fprint(zli.Stdout, usage)
		return
	}
	if cmd == "v" || cmd == "version" {
		fmt.Println(version)
		return
	}

	quiet := quietF.Set()
	raw := rawF.Set()
	args := flag.Args
	args, err = zli.InputOrArgs(args, " \t\n", quiet)
	zli.F(err)

	switch cmd {
	default:
		zli.Fatalf("unknown command")
	case "identify", "i":
		err = identify(args, quiet, raw)
	case "search", "s":
		err = search(args, quiet, raw, or.Bool())
	case "print", "p":
		err = print(args, quiet, raw)
	case "emoji", "e":
		err = emoji(args, quiet, raw, or.Bool(), parseToneFlag(tone.String()), parseGenderFlag(gender.String()))
	}
	if err != nil {
		if !(err == errNoMatches && quiet) {
			zli.Fatalf(err)
		}
		zli.Exit(1)
	}
}

func parseToneFlag(tone string) []string {
	if tone == "" {
		return nil
	}

	var allTones = []string{"none", "light", "mediumlight", "medium", "mediumdark", "dark"}
	if tone == "all" {
		tone = strings.Join(allTones, ",")
	}

	var tones []string
	for _, t := range zstring.Fields(tone, ",") {
		if !zstring.Contains(allTones, t) {
			zli.Fatalf("invalid skin tone: %q", tone)
		}
		tones = append(tones, t)
	}

	return tones
}

func parseGenderFlag(gender string) []string {
	if gender == "" {
		return nil
	}

	if gender == "all" {
		gender = "person,man,woman"
	}

	var genders []string
	for _, g := range zstring.Fields(gender, ",") {
		switch g {
		case "person", "p", "people":
			g = "person"
		case "man", "men", "m", "male":
			g = "man"
		case "woman", "women", "w", "female", "f":
			g = "woman"
		default:
			zli.Fatalf("invalid gender: %q", gender)
		}
		genders = append(genders, g)
	}

	return genders
}

func search(args []string, quiet, raw, or bool) error {
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

	var out printer
	for _, info := range unidata.Codepoints {
		m := 0
		for _, a := range args {
			if strings.Contains(info.Name, a) {
				if or {
					out = append(out, info)
					break
				}
				m++
			}
		}
		if !or && m == len(args) {
			out = append(out, info)
		}
	}

	if len(out) == 0 {
		return errNoMatches
	}

	out.PrintSorted(zli.Stdout, quiet, raw)
	return nil
}

func emoji(args []string, quiet, raw, or bool, tones, genders []string) error {
	type matchArg struct {
		group bool
		text  string
	}
	var matchArgs = make([]matchArg, 0, len(args))
	for _, a := range args {
		a := strings.ToLower(a)
		group := strings.HasPrefix(a, "g:") || strings.HasPrefix(a, "group:")
		if group {
			a = strings.TrimPrefix(strings.TrimPrefix(a, "group:"), "g:")
		}
		matchArgs = append(matchArgs, matchArg{text: a, group: group})
	}

	all := zstring.Contains(args, "all")
	if all {
		args = []string{"all"}
	}

	out := make([]unidata.Emoji, 0, 16)
	for _, e := range unidata.Emojis {
		m := 0
		for _, a := range matchArgs {
			var match bool
			switch {
			case a.group:
				match = strings.Contains(strings.ToLower(e.Group), a.text) || strings.Contains(strings.ToLower(e.Subgroup), a.text)
			default:
				match = strings.Contains(strings.ToLower(e.Name), a.text)
			}
			if match {
				if or {
					out = append(out, e)
					break
				}
				m++
			}
		}
		if all || (!or && m == len(matchArgs)) {
			out = append(out, applyGenders(applyTones(e, tones), genders)...)
		}
	}

	if len(out) == 0 {
		return errNoMatches
	}

	colSize := []int{0, 0}
	for _, e := range out {
		if len(e.Name) > colSize[0] {
			colSize[0] = len(e.Name)
		}
		if len(e.Group) > colSize[1] {
			colSize[1] = len(e.Group)
		}
	}
	colSize[0] += 2
	colSize[1] += 2

	// Alignment with spaces is tricky, as some emojis are double-width and some
	// are not. As far as I can tell, there is no good way to predict this as it
	// will depend on the font. Unicode recommends "emoji presentation sequences
	// behave as though they were East Asian Wide", but that's too simplistic
	// too. So use a tab character for this, which aligns correctly, even though
	// it adds some unnecessary whitespace.
	//
	// Don't do this when piping, since dmenu doesn't display tabs well :-/ This
	// seems like a problem in Xft as near as I can determine.
	tab := "\t"
	if !zli.IsTerminal(os.Stdout.Fd()) {
		tab = "    "
	}

	if !quiet {
		fmt.Fprintf(zli.Stdout, "  %s%s%s%s\n", tab,
			zstring.AlignLeft("name", colSize[0]),
			zstring.AlignLeft("group", colSize[1]),
			"subgroup")
	}

	for _, e := range out {
		fmt.Fprintf(zli.Stdout, "%s%s%s%s%s\n", e.String(), tab,
			zstring.AlignLeft(e.Name, colSize[0]),
			zstring.AlignLeft(e.Group, colSize[1]),
			e.Subgroup)
	}
	return nil
}

var tonemap = map[string]uint32{
	"none":        0,
	"light":       0x1f3fb,
	"mediumlight": 0x1f3fc,
	"medium":      0x1f3fd,
	"mediumdark":  0x1f3fe,
	"dark":        0x1f3ff,
}

// Skintone always comes after the base emoji and doesn't required a ZWJ.
func applyTones(e unidata.Emoji, tones []string) []unidata.Emoji {
	if !e.SkinTones || len(tones) == 0 {
		return []unidata.Emoji{e}
	}

	emojis := make([]unidata.Emoji, len(tones))
	for i, t := range tones {
		emojis[i] = e // This makes a copy, but beware of directly modifying lists as they're pointers.

		if tcp := tonemap[t]; tcp > 0 {
			emojis[i].Name += fmt.Sprintf(": %s skin tone", t)
			emojis[i].Codepoints = append(append([]uint32{e.Codepoints[0]}, tcp), e.Codepoints[1:]...)
			l := len(emojis[i].Codepoints) - 1
			if emojis[i].Codepoints[l] == 0xfe0f {
				emojis[i].Codepoints = emojis[i].Codepoints[:l]
			}
		}
	}

	return emojis
}

func applyGenders(emojis []unidata.Emoji, genders []string) []unidata.Emoji {
	if len(genders) == 0 {
		return emojis
	}

	var ret []unidata.Emoji
	for _, e := range emojis {
		if e.Genders == unidata.GenderNone {
			ret = append(ret, e)
			continue
		}

		for _, g := range genders {
			ee := e
			ee.Codepoints = make([]uint32, len(ee.Codepoints))
			copy(ee.Codepoints, e.Codepoints)

			if e.Genders == unidata.GenderSign {
				// Append male or female sign
				//   1F937 1F3FD                   # 🤷🏽 E4.0 person shrugging: medium skin tone
				//   1F937 1F3FB 200D 2642 FE0F    # 🤷🏻‍♂️ E4.0 man shrugging: light skin tone
				switch g {
				case "woman":
					ee.Name = strings.Replace(ee.Name, "person", "woman", 1)
					ee.Codepoints = append(ee.Codepoints, []uint32{0x2640, 0xfe0f}...)
				case "man":
					ee.Name = strings.Replace(ee.Name, "person", "man", 1)
					ee.Codepoints = append(ee.Codepoints, []uint32{0x2642, 0xfe0f}...)
				}
			} else if e.Genders == unidata.GenderRole {
				// Replace first "person" with "man" or "woman".
				//   1F9D1 200D 1F692              # 🧑‍🚒 E12.1 firefighter
				//   1F9D1 1F3FB 200D 1F692        # 🧑🏻‍🚒 E12.1 firefighter: light skin tone
				//   1F469 200D 1F692              # 👩‍🚒 E4.0 woman firefighter
				//   1F469 1F3FB 200D 1F692        # 👩🏻‍🚒 E4.0 woman firefighter: light skin tone
				switch g {
				case "woman":
					ee.Name = "woman " + ee.Name
					ee.Codepoints = append([]uint32{0x1f469}, ee.Codepoints[1:]...)
				case "man":
					ee.Name = "man " + ee.Name
					ee.Codepoints = append([]uint32{0x1f468}, ee.Codepoints[1:]...)
				}
			}

			ret = append(ret, ee)
		}
	}

	return ret
}

func parseEmojiGroups(group string) []string {
	groups := strings.Split(strings.ToLower(group), ",")
	for _, g := range groups {
		found := false
	outer:
		for eg, subs := range unidata.EmojiSubgroups {
			if strings.Contains(strings.ToLower(eg), g) {
				found = true
				break
			}

			for _, sg := range subs {
				if strings.Contains(strings.ToLower(sg), g) {
					found = true
					break outer
				}
			}
		}
		if !found {
			zli.Fatalf("doesn't match any emoji group or subgroup: %q", g)
		}
	}

	return groups
}

func print(args []string, quiet, raw bool) error {
	var out printer

	for _, a := range args {
		canon := unidata.CanonicalCategory(a)

		// Print everything.
		if canon == "all" {
			for _, info := range unidata.Codepoints {
				out = append(out, info)
			}
			continue
		}

		// Category name.
		if cat, ok := unidata.Catmap[canon]; ok {
			for _, info := range unidata.Codepoints {
				if info.Cat == cat {
					out = append(out, info)
				}
			}
			continue
		}

		// Block.
		if bl, ok := unidata.Blockmap[canon]; ok {
			for cp := unidata.Blocks[bl][0]; cp <= unidata.Blocks[bl][1]; cp++ {
				s, ok := unidata.Codepoints[fmt.Sprintf("%04X", cp)]
				if ok {
					out = append(out, s)
				}
			}
			continue
		}

		// U2042, U+2042, U+2042..U+2050, 2042..2050, 0x2041, etc.
		var s []string
		switch {
		case strings.Contains(canon, ".."):
			s = strings.SplitN(canon, "..", 2)
		case strings.Contains(canon, "-"):
			s = strings.SplitN(canon, "-", 2)
		default:
			s = []string{canon, canon}
		}
		start, err1 := unidata.ToCodepoint(s[0])
		end, err2 := unidata.ToCodepoint(s[1])
		if len(s) != 2 || err1 != nil || err2 != nil {
			return fmt.Errorf("unknown identifier: %q", a)
		}
		for i := start; i <= end; i++ {
			info, _ := unidata.FindCodepoint(rune(i))
			out = append(out, info)
		}
	}

	out.PrintSorted(zli.Stdout, quiet, raw)
	return nil
}

func identify(ins []string, quiet, raw bool) error {
	in := strings.Join(ins, "")
	if !utf8.ValidString(in) {
		fmt.Fprintf(zli.Stderr, "uni: WARNING: input string is not valid UTF-8\n")
	}

	var out printer
	for _, c := range in {
		info, ok := unidata.FindCodepoint(c)
		if !ok {
			return fmt.Errorf("unknown codepoint: U+%.4X", c) // Should never happen.
		}
		out = append(out, info)
	}
	out.Print(zli.Stdout, quiet, raw)
	return nil
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
