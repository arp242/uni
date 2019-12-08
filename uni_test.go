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
		{[]string{"a"}, "SMALL LETTER A"},

		// Make sure it uses the lower-case and short variant.
		{[]string{`"`}, "&quot;"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, "i", tt.in)
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
	}{
		{[]string{"asterism"}, "ASTERISM", 1},
		{[]string{"floral"}, "HEART", 3},
		{[]string{"floral", "bullet"}, "HEART", 2},
		{[]string{"rightwards arrow", "heavy"}, "HEAVY", 16},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, "s", tt.in)
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
		{[]string{"U+2042"}, "ASTERISM", 1},
		{[]string{"U+2042..U+2044"}, "ASTERISM", 3},
		{[]string{"U+3402"}, "'㐂'", 1},
		{[]string{"U+3402..U+3404"}, "<CJK Ideograph Extension A>", 3},
		{[]string{"OtherPunctuation"}, "ASTERISM", 588},
		{[]string{"Po"}, "ASTERISM", 588},
		{[]string{"GeneralPunctuation"}, "ASTERISM", 111},
		{[]string{"all"}, "ASTERISM", 32841},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, "p", tt.in)
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
		{[]string{"-groups", "hands"},
			[]string{"👏", "🙌", "👐", "🤲", "🤝", "🙏"}},
		{[]string{"-tone", "dark", "-groups", "hands"},
			[]string{"👏Z🏿", "🙌Z🏿", "👐Z🏿", "🤲Z🏿", "🤝", "🙏Z🏿"}},

		{[]string{"shrug"},
			[]string{"🤷", "🤷Z♂S", "🤷Z♀S"}},
		{[]string{"-gender", "m", "shrug"},
			[]string{"🤷Z♂S"}},
		{[]string{"-gender", "m", "-tone", "light", "shrug"},
			// TODO: one ZWJ too many; still works fine though.
			// 1F937 1F3FB 200D 2642 FE0F ; fully-qualified # 🤷🏻‍♂️ E4.0 man shrugging: light skin tone
			[]string{"🤷Z🏻Z♂S"}},

		{[]string{"farmer"},
			[]string{"🧑Z🌾", "👨Z🌾", "👩Z🌾"}},
		{[]string{"-gender", "f,m", "farmer"},
			[]string{"👩Z🌾", "👨Z🌾"}},
		{[]string{"-gender", "f", "-tone", "medium", "farmer"},
			// TODO: one ZWJ too many
			// 1F469 1F3FD 200D 1F33E ; fully-qualified # 👩🏽‍🌾 E4.0 woman farmer: medium skin tone
			[]string{"👩Z🏽Z🌾"}},

		{[]string{"-tone", "mediumlight", "bride"},
			// TODO: one ZWJ too many
			// 1F470 1F3FC ; fully-qualified # 👰🏼 E2.0 bride with veil: medium-light skin tone
			[]string{"👰Z🏼"}},

		// TODO: below all fail. Unicode is so inconsistent :-(

		// 1F575 FE0F ; fully-qualified # 🕵️ E2.0 detective
		//{[]string{"-gender", "p", "detective"},
		//	[]string{"🕵S"}},

		// 1F575 1F3FE ; fully-qualified # 🕵🏾 E2.0 detective: medium-dark skin tone
		// {[]string{"-gender", "p", "-tone", "mediumdark", "detective"},
		// 	[]string{"🕵🏾"}},

		// 1F575 FE0F 200D 2642 FE0F ; fully-qualified # 🕵️‍♂️ E4.0 man detective
		//{[]string{"-gender", "m", "detective"},
		//	[]string{"🕵️SZ♂️S"}},

		// 1F575 1F3FE 200D 2642 FE0F ; fully-qualified # 🕵🏾‍♂️ E4.0 man detective: medium-dark skin tone
		//{[]string{"-gender", "m", "-tone", "mediumdark", "detective"},
		//	[]string{"🕵🏾Z♂️S"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, "e", tt.in)
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

func setup(t *testing.T, cmd string, args []string) (fmt.Stringer, func()) {
	outbuf := new(bytes.Buffer)
	stdout = outbuf
	stderr = outbuf

	os.Args = append([]string{"testuni", "-q", cmd}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	exit = func(code int) {
		if code > 0 {
			t.Fatalf("os.Exit(%d) called\n%s", code, outbuf.String())
		}
	}

	return outbuf, func() {
		stdout = os.Stdout
		stderr = os.Stderr
		exit = os.Exit
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
