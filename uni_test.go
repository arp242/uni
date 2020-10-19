package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"zgo.at/zli"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"xxx"}, "uni: unknown command"},
		//{[]string{""}, "uni: unknown command"},
		//{[]string{}, "Show this help"},
		{[]string{"e", "-t"}, "testuni: -t: needs an argument"},
		{[]string{"e", "-t", "-g"}, `testuni: invalid skin tone: "-g"`},
		{[]string{"e", "-x"}, `testuni: unknown flag: "-x"`},
		{[]string{"e", "-t", "xx"}, "invalid skin"},
		{[]string{"e", "-gender", "xx"}, "invalid gender"},
		{[]string{"e", "-g", "xxsxxxx"}, `invalid gender: "xxsxxxx"`},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			exit, _, out, reset := zli.Test()
			defer reset()
			os.Args = append([]string{"testuni"}, tt.in...)

			func() {
				defer exit.Recover()
				main()
			}()

			if !strings.Contains(out.String(), tt.want) {
				t.Errorf("wrong output\nout:  %q\nwant: %q", out.String(), tt.want)
			}
			if *exit != 1 {
				t.Errorf("wrong exit: %d", *exit)
			}
		})
	}

	t.Run("usage", func(t *testing.T) {
		if strings.Contains(usage, "\t") {
			t.Errorf("usage text contains tabs")
		}
	})
}

func TestIdentify(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"i", ""}, ""},
		{[]string{"i", "a"}, "SMALL LETTER A"},
		{[]string{"i", `"`}, "&quot;"}, // Make sure it uses the lower-case and short variant.
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			exit, _, out, reset := zli.Test()
			defer reset()
			os.Args = append([]string{"testuni"}, tt.in...)

			func() {
				defer exit.Recover()
				main()
			}()

			if !strings.Contains(out.String(), tt.want) {
				t.Errorf("wrong output\nout:  %q\nwant: %q", out.String(), tt.want)
			}
			if *exit != -1 {
				t.Errorf("wrong exit: %d", *exit)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		in        []string
		want      string
		wantLines int
		wantExit  int
	}{
		{[]string{"s", ""}, "need search term", 1, 1},

		{[]string{"-q", "s", "asterism"}, "ASTERISM", 1, -1},
		{[]string{"-q", "s", "floral"}, "HEART", 3, -1},
		{[]string{"-q", "s", "floral", "bullet"}, "HEART", 2, -1},
		{[]string{"-q", "s", "rightwards arrow", "heavy"}, "HEAVY", 16, -1},

		{[]string{"s", "nomatch_nomatch"}, "no matches", 1, 1},
		{[]string{"-q", "s", "nomatch_nomatch"}, "", 0, 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			exit, _, outbuf, reset := zli.Test()
			defer reset()
			os.Args = append([]string{"testuni"}, tt.in...)

			func() {
				defer exit.Recover()
				main()
			}()
			if int(*exit) != tt.wantExit {
				t.Fatalf("wrong exit: %d", *exit)
			}

			out := outbuf.String()
			if lines := strings.Count(out, "\n"); lines != tt.wantLines {
				t.Errorf("wrong # of lines\nout:  %d\nwant: %d", lines, tt.wantLines)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("wrong output\nout:  %q\nwant: %q", out, tt.want)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	tests := []struct {
		in                  []string
		want                string
		wantLines, wantExit int
	}{
		{[]string{"-q", "p", "U+2042"}, "ASTERISM", 1, -1},
		{[]string{"-q", "p", "2042"}, "ASTERISM", 1, -1},
		{[]string{"-q", "p", "U+2042..U+2044"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "2042..2044"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "U+2042..2044"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "2042..U+2044"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "0x2042..0o20104"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "0b10000001000010..0o20104"}, "ASTERISM", 3, -1},
		{[]string{"-q", "p", "0b10000001000010-0o20104"}, "ASTERISM", 3, -1},
		{[]string{"p", "9999999999"}, `CODEPOINT NOT IN UNICODE`, 2, -1},

		{[]string{"p", ""}, `unknown identifier: ""`, 1, 1},
		{[]string{"p", "nonsense"}, `unknown identifier: "nonsense"`, 1, 1},
		{[]string{"p", "2042..xxx"}, `unknown identifier: "2042..xxx"`, 1, 1},
		{[]string{"p", "xxx..xxx"}, `unknown identifier: "xxx..xxx"`, 1, 1},
		{[]string{"p", "xxx..xxx"}, `unknown identifier: "xxx..xxx"`, 1, 1},

		{[]string{"-q", "p", "U+3402"}, "'㐂'", 1, -1},
		{[]string{"-q", "p", "U+3402..U+3404"}, "<CJK Ideograph Extension A>", 3, -1},
		{[]string{"-q", "p", "OtherPunctuation"}, "ASTERISM", 593, -1},
		{[]string{"-q", "p", "Po"}, "ASTERISM", 593, -1},
		{[]string{"-q", "p", "GeneralPunctuation"}, "ASTERISM", 111, -1},
		{[]string{"-q", "p", "all"}, "ASTERISM", 33797, -1},

		{[]string{"-q", "-r", "p", "U9"}, "'\t'", 1, -1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			exit, _, outbuf, reset := zli.Test()
			defer reset()
			os.Args = append([]string{"testuni"}, tt.in...)

			func() {
				defer exit.Recover()
				main()
			}()

			if int(*exit) != tt.wantExit {
				t.Fatalf("wrong exit: %d", *exit)
			}

			out := outbuf.String()
			if lines := strings.Count(out, "\n"); lines != tt.wantLines {
				t.Errorf("wrong # of lines\nout:  %d\nwant: %d", lines, tt.wantLines)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("wrong output\nout:  %q\nwant: %q", out, tt.want)
			}
		})
	}
}

func TestEmoji(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		//{[]string{"e", "all"},
		//[]string{}},

		//{[]string{"e", "-groups", "person", "all"},
		//[]string{}},

		{[]string{"e", "group:hands"},
			[]string{"👏", "🙌", "👐", "🤲", "🤝", "🙏"}},
		{[]string{"e", "-tone", "dark", "g:hands"},
			[]string{"👏🏿", "🙌🏿", "👐🏿", "🤲🏿", "🤝", "🙏🏿"}},

		{[]string{"e", "shrug"},
			[]string{"🤷"}},
		{[]string{"e", "shrug", "-gender", "all"},
			[]string{"🤷", "🤷Z♂S", "🤷Z♀S"}},
		{[]string{"e", "-gender", "m", "shrug"},
			[]string{"🤷Z♂S"}},
		{[]string{"e", "-gender", "m", "-tone", "light", "shrug"},
			[]string{"🤷🏻Z♂S"}},

		{[]string{"e", "farmer"},
			[]string{"🧑Z🌾"}},
		{[]string{"e", "farmer", "-gender", "all"},
			[]string{"🧑Z🌾", "👨Z🌾", "👩Z🌾"}},
		{[]string{"e", "-gender", "f,m", "farmer"},
			[]string{"👩Z🌾", "👨Z🌾"}},
		{[]string{"e", "-gender", "f", "-tone", "medium", "farmer"},
			[]string{"👩🏽Z🌾"}},

		{[]string{"e", "-gender", "p", "detective"},
			[]string{"🕵S"}},
		{[]string{"e", "-gender", "p", "-tone", "mediumdark", "detective"},
			[]string{"🕵🏾"}},
		{[]string{"e", "-gender", "m", "detective"},
			[]string{"🕵SZ♂S"}},
		{[]string{"e", "-gender", "m", "-tone", "mediumdark", "detective"},
			[]string{"🕵🏾Z♂S"}},

		{[]string{"e", "zimbabwe", "#", "england"},
			[]string{"🇿🇼", "#S⃣", "🏴󠁧󠁢󠁥󠁮󠁧󠁿"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			_, _, outbuf, reset := zli.Test()
			defer reset()
			os.Args = append([]string{"testuni"}, tt.in...)

			main()

			var out []string
			for _, line := range strings.Split(strings.TrimSpace(outbuf.String()), "\n") {
				out = append(out, strings.Split(line, " ")[0])
			}

			for i := range tt.want {
				tt.want[i] = strings.Replace(tt.want[i], "Z", "\u200d", -1)
				tt.want[i] = strings.Replace(tt.want[i], "S", "\ufe0f", -1)
			}

			if !reflect.DeepEqual(out, tt.want) {
				a := strings.ReplaceAll(fmt.Sprintf("%#v", out), "\ufe0f", `\ufe0f`)
				b := strings.ReplaceAll(fmt.Sprintf("%#v", tt.want), "\ufe0f", `\ufe0f`)
				t.Errorf("wrong output\nout:  %s\nwant: %s", a, b)
			}
		})
	}
}

func TestAllEmoji(t *testing.T) {
	exit, _, outbuf, reset := zli.Test()
	defer reset()
	os.Args = append([]string{"testuni"}, []string{"e", "-gender", "all", "-tone", "all", "all"}...)

	func() {
		defer exit.Recover()
		main()
	}()

	// grep -v '^#' unidata/.cache/emoji-test.txt |
	//     grep fully-qualified |
	//     grep -Ev '(holding hands|kiss:|couple with heart).*tone' |
	//     grep -Eo '# .+? E[0-9]' |
	//     cut -d ' ' -f2 >| testdata/emojis
	//
	// double tones: 70
	// family: 145
	w, err := ioutil.ReadFile("./testdata/emojis")
	if err != nil {
		t.Fatal(err)
	}
	wantEmojis := strings.Split(strings.TrimSpace(string(w)), "\n")

	out := strings.Split(strings.TrimRight(outbuf.String(), "\n"), "\n")
	outEmojis := make([]string, len(out))
	for i := range out {
		outEmojis[i] = out[i][:strings.Index(out[i], " ")]
	}

	if len(outEmojis) != len(wantEmojis) {
		t.Errorf("different length: want %d, got %d", len(wantEmojis), len(outEmojis))
	}

	// Still some \ufe0f issues
	t.Skip()

	// TODO: this shouldnt; be needed
	sort.Strings(wantEmojis)
	sort.Strings(outEmojis)

	for i := range wantEmojis {
		wantEmojis[i] = strings.ReplaceAll(wantEmojis[i], "\ufe0f", `\ufe0f`)
	}
	for i := range outEmojis {
		outEmojis[i] = strings.ReplaceAll(outEmojis[i], "\ufe0f", `\ufe0f`)
	}
	if !reflect.DeepEqual(outEmojis, wantEmojis) {
		t.Errorf("emoji lists not equal\nout:  %v\nwant: %v", outEmojis, wantEmojis)
	}
	//if d := ztest.Diff(outEmojis, wantEmojis); d != "" {
	//	t.Error(d)
	//}
}
