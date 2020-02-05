// Command uni prints Unicode information about characters.
package main // import "arp242.net/uni"

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"arp242.net/uni/isatty"
	"arp242.net/uni/unidata"
)

var (
	errNoMatches = errors.New("no matches")
)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
	exit             = os.Exit
)

const usage = `Usage: uni [-hrq] [help | identify | search | print | emoji]

Flags:
    -q      Quiet output; don't print header, "no matches", etc.
    -r      "Raw" output instead of displaying graphical variants for control
            characters and ◌ (U+25CC) before combining characters.

Commands:
    identify [string string ...]
        Idenfity all the characters in the given strings.

    search [word word ...]
        Search description for any of the words.

    print [ident ident ...]
        Print characters by codepoint, category, or block:

            Codepoints             U+2042, U+2042..U+2050, 0x20, 0o40, 0b100000
            Categories and Blocks  OtherPunctuation, Po, GeneralPunctuation
            all                    Everything

        Names are matched case insensitive; spaces and commas are optional and
        can be replaced with an underscore. "Po", "po", "punction, OTHER",
        "Punctuation_other", and PunctuationOther are all identical.

    emoji [-tone tone,..] [-gender gender,..] [-groups word] [word word ...]
        Search emojis. The special keyword "all" prints all emojis.

        -group   comma-separated list of group and/or subgroup names.
        -tone    comma-separated list of light, mediumlight, medium,
                 mediumdark, dark. Default is to include none.
        -gender  comma-separated list of person, man, or woman.
                 Default is to include all.

        Note: output may contain unprintable character (U+200D and U+FE0F) which
        may not survive a select and copy operation from text-based applications
        such as terminals. It's recommended to copy to the clipboard directly
        with e.g. xclip.
`

func main() {
	var (
		quiet bool
		help  bool
		raw   bool
	)
	flag.BoolVar(&quiet, "q", false, "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&raw, "r", false, "")
	flag.Usage = func() { fmt.Fprint(stdout, usage) }
	flag.Parse()

	if help {
		flag.Usage()
		exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		die("no command given")
		return
	}

	var err error
	switch strings.ToLower(args[0]) {
	default:
		die("unknown command: %q", args[0])
	case "help", "h":
		flag.Usage()
		exit(0)
	case "identify", "i":
		err = identify(getargs(args[1:], quiet), quiet, raw)
	case "search", "s":
		err = search(getargs(args[1:], quiet), quiet, raw)
	case "print", "p":
		err = print(getargs(args[1:], quiet), quiet, raw)
	case "emoji", "e":
		err = emoji(getargs(args[1:], quiet), quiet, raw)
	}

	if err != nil {
		if !(err == errNoMatches && quiet) {
			fmt.Fprintf(stderr, "%s\n", err)
		}
		exit(1)
	}
}

func die(f string, a ...interface{}) {
	fmt.Fprintf(stderr, "uni: "+f+"\n", a...)
	exit(1)
}

// Use commandline args or stdin.
func getargs(args []string, quiet bool) []string {
	if len(args) > 0 {
		return args
	}

	interactive := isatty.IsTerminal(os.Stdin.Fd())

	// Print message so people aren't left waiting when typing "uni print". We
	// don't print a newline and a \r later on, so you don't see it in actual
	// pipe usage, just when it would "hang" uni.
	if !quiet && interactive {
		fmt.Fprintf(stderr, "uni: reading from stdin...")
		os.Stderr.Sync()
	}
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Errorf("read stdin: %s", err))
	}
	if !quiet && interactive {
		fmt.Fprintf(stderr, "\r")
	}

	var words []string
	for _, l := range strings.Split(strings.TrimRight(string(stdin), "\n"), "\n") {
		words = append(words, strings.Split(l, " ")...)
	}
	return words
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
	for _, info := range unidata.Codepoints {
		m := 0
		for _, w := range words {
			if strings.Contains(info.Name, w) {
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

	out.PrintSorted(stdout, quiet, raw)
	return nil
}

func emoji(args []string, quiet, raw bool) error {
	var (
		tone, gender, group string
		subargs             []string
		skip                bool
	)

	needarg := func(i int) error {
		if i+2 > len(args) || len(args[i+1]) == 0 || args[i+1][0] == '-' {
			return fmt.Errorf("argument required for %s", args[i])
		}
		return nil
	}
	for i := range args {
		if skip {
			skip = false
			continue
		}
		switch args[i] {
		case "-t", "-tone", "-tones":
			if err := needarg(i); err != nil {
				return err
			}
			tone += args[i+1]
			skip = true
		case "-gender", "-genders":
			if err := needarg(i); err != nil {
				return err
			}
			gender += args[i+1]
			skip = true
		case "-g", "-group", "-groups":
			if err := needarg(i); err != nil {
				return err
			}
			group += args[i+1]
			skip = true
		default:
			if len(args[i]) > 0 && args[i][0] == '-' {
				return fmt.Errorf("unknown option: %s", args[i])
			}
			subargs = append(subargs, args[i])
		}
	}

	groups := parseEmojiGroups(group)

	if len(subargs) == 0 && len(groups) > 0 {
		subargs = []string{""} // Imply all
	}

	out := [][]string{}
	cols := []int{4, 0, 0}
	for _, a := range subargs {
		a = strings.ToLower(a)
		switch a {
		case "all":
			a = ""
		}

		for _, e := range unidata.Emojis {
			found := false
			for _, g := range groups {
				if strings.Contains(strings.ToLower(e.Group), g) ||
					strings.Contains(strings.ToLower(e.Subgroup), g) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
			if !strings.Contains(strings.ToLower(e.Name), a) {
				continue
			}

			for _, ee := range applyGender(applyTone(e, tone), gender) {
				out = append(out, []string{ee.String(), ee.Name, ee.Group, ee.Subgroup})
				if l := utf8.RuneCountInString(ee.Name); l > cols[1] {
					cols[1] = l
				}
				if l := utf8.RuneCountInString(ee.Group); l > cols[2] {
					cols[2] = l
				}
			}
		}
	}

	// TODO: not always correctly aligned as some emojis are double-width and
	// some are not. As far as I can tell, there is no good way to predict this
	// as it will depend on the font. Unicode recommends "emoji presentation
	// sequences behave as though they were East Asian Wide", but that's too
	// simplistic too.
	for _, o := range out {
		for i, c := range o {
			switch i {
			case 0:
				fmt.Fprint(stdout, c+" ")
			case 3: // Last column
				fmt.Fprint(stdout, c)
			default:
				fmt.Fprint(stdout, fill(c, cols[i]+2))
			}
		}
		fmt.Fprintln(stdout, "")
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

func applyTone(e unidata.Emoji, t string) []unidata.Emoji {
	if !e.SkinTones || t == "" {
		return []unidata.Emoji{e}
	}

	if t == "all" {
		t = "none,light,mediumlight,medium,mediumdark,dark"
	}

	tones := strings.Split(t, ",")
	emojis := make([]unidata.Emoji, len(tones))
	// Skintone always comes after the base emoji:
	//   1F937 1F3FD                   # 🤷🏽 E4.0 person shrugging: medium skin tone
	//   1F937 1F3FB 200D 2642 FE0F    # 🤷🏻‍♂️ E4.0 man shrugging: light skin tone
	//   1F9D1 200D 1F692              # 🧑‍🚒 E12.1 firefighter
	//   1F9D1 1F3FB 200D 1F692        # 🧑🏻‍🚒 E12.1 firefighter: light skin tone
	//   1F469 200D 1F692              # 👩‍🚒 E4.0 woman firefighter
	//   1F469 1F3FB 200D 1F692        # 👩🏻‍🚒 E4.0 woman firefighter: light skin tone
	for i, t := range tones {
		tcp, ok := tonemap[t]
		if !ok {
			die("invalid skin tone: %q", t)
		}

		emojis[i] = unidata.Emoji{
			Codepoints: e.Codepoints,
			Name:       e.Name,
			Group:      e.Group,
			Subgroup:   e.Subgroup,
			SkinTones:  e.SkinTones,
			Genders:    e.Genders,
		}

		if tcp > 0 {
			emojis[i].Codepoints = append(append([]uint32{e.Codepoints[0]}, tcp), e.Codepoints[1:]...)
			l := len(emojis[i].Codepoints) - 1
			if emojis[i].Codepoints[l] == 0xfe0f {
				emojis[i].Codepoints = emojis[i].Codepoints[:l]
			}

			emojis[i].Name += fmt.Sprintf(": %s skin tone", t)
		}

	}

	return emojis
}

func applyGender(emojis []unidata.Emoji, gender string) []unidata.Emoji {
	genders := []string{"person", "man", "woman"}
	if gender != "" {
		genders = strings.Split(gender, ",")
		for i, g := range genders {
			switch g {
			case "person", "p", "people":
				g = "person"
			case "man", "men", "m", "male":
				g = "man"
			case "woman", "women", "w", "female", "f":
				g = "woman"
			default:
				die("invalid gender : %q", g)
			}
			genders[i] = g
		}
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
			die("doesn't match any emoji group or subgroup: %q", g)
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

	out.PrintSorted(stdout, quiet, raw)
	return nil
}

func identify(ins []string, quiet, raw bool) error {
	in := strings.Join(ins, "")
	if !utf8.ValidString(in) {
		fmt.Fprintf(stderr, "uni: WARNING: input string is not valid UTF-8\n")
	}

	var out printer
	for _, c := range in {
		info, ok := unidata.FindCodepoint(c)
		if !ok {
			return fmt.Errorf("unknown codepoint: U+%.4X", c) // Should never happen.
		}

		out = append(out, info)
	}

	out.Print(stdout, quiet, raw)
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
