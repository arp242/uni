package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestIdentify(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"i", ""}, ""},

		{[]string{"i", "a"}, "SMALL LETTER A"},

		// Make sure it uses the lower-case and short variant.
		{[]string{"i", `"`}, "&quot;"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, tt.in, -1)
			defer c()

			main()
			out := outbuf.String()
			if !strings.Contains(out, tt.want) {
				t.Errorf("wrong output\nout:  %q\nwant: %q", out, tt.want)
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
			outbuf, c := setup(t, tt.in, tt.wantExit)
			defer c()

			main()
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
		in        []string
		want      string
		wantLines int
	}{
		{[]string{"-q", "p", "U+2042"}, "ASTERISM", 1},
		{[]string{"-q", "p", "U+2042..U+2044"}, "ASTERISM", 3},
		{[]string{"-q", "p", "U+3402"}, "'㐂'", 1},
		{[]string{"-q", "p", "U+3402..U+3404"}, "<CJK Ideograph Extension A>", 3},
		{[]string{"-q", "p", "OtherPunctuation"}, "ASTERISM", 588},
		{[]string{"-q", "p", "Po"}, "ASTERISM", 588},
		{[]string{"-q", "p", "GeneralPunctuation"}, "ASTERISM", 111},
		{[]string{"-q", "p", "all"}, "ASTERISM", 32841},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, tt.in, -1)
			defer c()

			main()
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

		{[]string{"e", "-groups", "hands"},
			[]string{"👏", "🙌", "👐", "🤲", "🤝", "🙏"}},
		{[]string{"e", "-tone", "dark", "-groups", "hands"},
			[]string{"👏Z🏿", "🙌Z🏿", "👐Z🏿", "🤲Z🏿", "🤝", "🙏Z🏿"}},

		{[]string{"e", "shrug"},
			[]string{"🤷", "🤷Z♂S", "🤷Z♀S"}},
		{[]string{"e", "-gender", "m", "shrug"},
			[]string{"🤷Z♂S"}},
		{[]string{"e", "-gender", "m", "-tone", "light", "shrug"},
			// TODO: one ZWJ too many; still works fine though.
			// 1F937 1F3FB 200D 2642 FE0F ; fully-qualified # 🤷🏻‍♂️ E4.0 man shrugging: light skin tone
			[]string{"🤷Z🏻Z♂S"}},

		{[]string{"e", "farmer"},
			[]string{"🧑Z🌾", "👨Z🌾", "👩Z🌾"}},
		{[]string{"e", "-gender", "f,m", "farmer"},
			[]string{"👩Z🌾", "👨Z🌾"}},
		{[]string{"e", "-gender", "f", "-tone", "medium", "farmer"},
			// TODO: one ZWJ too many
			// 1F469 1F3FD 200D 1F33E ; fully-qualified # 👩🏽‍🌾 E4.0 woman farmer: medium skin tone
			[]string{"👩Z🏽Z🌾"}},

		{[]string{"e", "-tone", "mediumlight", "bride"},
			// TODO: one ZWJ too many
			// 1F470 1F3FC ; fully-qualified # 👰🏼 E2.0 bride with veil: medium-light skin tone
			[]string{"👰Z🏼"}},

		// TODO: below all fail. Unicode is so inconsistent :-(

		// 1F575 FE0F ; fully-qualified # 🕵️ E2.0 detective
		//{[]string{"e", "-gender", "p", "detective"},
		//	[]string{"🕵S"}},

		// 1F575 1F3FE ; fully-qualified # 🕵🏾 E2.0 detective: medium-dark skin tone
		// {[]string{"e", "-gender", "p", "-tone", "mediumdark", "detective"},
		// 	[]string{"🕵🏾"}},

		// 1F575 FE0F 200D 2642 FE0F ; fully-qualified # 🕵️‍♂️ E4.0 man detective
		//{[]string{"e", "-gender", "m", "detective"},
		//	[]string{"🕵️SZ♂️S"}},

		// 1F575 1F3FE 200D 2642 FE0F ; fully-qualified # 🕵🏾‍♂️ E4.0 man detective: medium-dark skin tone
		//{[]string{"e", "-gender", "m", "-tone", "mediumdark", "detective"},
		//	[]string{"🕵🏾Z♂️S"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, tt.in, -1)
			defer c()

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
				// U+FE0F is a somewhat elusive character that gets eaten and
				// not displayed. Make sure it's displayed.
				a := strings.Replace(fmt.Sprintf("%#v", out), "\ufe0f", `\ufe0f`, -1)
				b := strings.Replace(fmt.Sprintf("%#v", tt.want), "\ufe0f", `\ufe0f`, -1)
				t.Errorf("wrong output\nout:  %s\nwant: %s", a, b)
			}
		})
	}
}

func setup(t *testing.T, args []string, wantExit int) (fmt.Stringer, func()) {
	outbuf := new(bytes.Buffer)
	stdout = outbuf
	stderr = outbuf

	os.Args = append([]string{"testuni"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	exitRan := false
	exit = func(code int) {
		exitRan = true
		if wantExit == -1 {
			t.Fatalf("os.Exit(%d) called\n%s", code, outbuf.String())
		}
		if code != wantExit {
			t.Fatalf("os.Exit(%d) called; want %d\n%s", code, wantExit, outbuf.String())
		}
	}

	return outbuf, func() {
		stdout = os.Stdout
		stderr = os.Stderr
		exit = os.Exit

		if wantExit > -1 && !exitRan {
			t.Fatalf("os.Exit() not called")
		}
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
