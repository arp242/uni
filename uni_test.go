package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"zgo.at/zli"
	"zgo.at/zstd/ztest"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		in   []string
		want string
	}{
		{[]string{"xxx"}, "uni: unknown command"},
		//{[]string{""}, "uni: unknown command"},
		//{[]string{}, "Show this help"},
		{[]string{"e", "-tone"}, "testuni: -tone: needs an argument"},
		{[]string{"e", "-tone", "g"}, `testuni: invalid skin tone: "g"`},
		{[]string{"e", "-x"}, `testuni: unknown flag: "-x"`},
		{[]string{"e", "-tone", "xx"}, "invalid skin"},
		{[]string{"e", "-gender", "xx"}, "invalid gender"},
		{[]string{"e", "-g", "xxsxxxx"}, `invalid gender: "xxsxxxx"`},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, "_"), func(t *testing.T) {
			exit, _, out := zli.Test(t)
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

		{[]string{"i", "\u0600"}, "␣"},                                 // ARABIC NUMBER SIGN
		{[]string{"i", "\u200d"}, "␣"},                                 // ZERO WIDTH JOINER
		{[]string{"i", "\u200e"}, "␣"},                                 // LEFT-TO-RIGHT MARK
		{[]string{"i", "\U000E007E"}, "␣"},                             // TAG TILDE
		{[]string{"i", "\uE000"}, "'\uE000'"},                          // <Private Use> (First)
		{[]string{"i", "\uE001"}, "'\uE001'"},                          // <Private Use>
		{[]string{"i", "\uF8FF"}, "'\uF8FF'"},                          // <Private Use> (Last)
		{[]string{"i", "\U000F0000"}, "'\U000F0000'"},                  // <Plane 15 Private Use> (First)
		{[]string{"i", "\U000F0001"}, "'\U000F0001'"},                  // <Plane 15 Private Use>
		{[]string{"i", "\U000FFFFD"}, "'\U000FFFFD'"},                  // <Plane 15 Private Use> (Last)
		{[]string{"i", "\U00100000"}, "'\U00100000'"},                  // <Plane 16 Private Use> (First)
		{[]string{"i", "\U00100001"}, "'\U00100001'"},                  // <Plane 16 Private Use>
		{[]string{"i", "\U0010FFFD"}, "'\U0010FFFD'"},                  // <Plane 16 Private Use> (Last)
		{[]string{"i", "\uE000"}, "<Private Use, First>"},              // <Private Use> (First)
		{[]string{"i", "\uE001"}, "<Private Use>"},                     // <Private Use>
		{[]string{"i", "\uF8FF"}, "<Private Use, Last>"},               // <Private Use> (Last)
		{[]string{"i", "\U000F0000"}, "<Plane 15 Private Use, First>"}, // <Plane 15 Private Use> (First)
		{[]string{"i", "\U000F0001"}, "<Plane 15 Private Use>"},        // <Plane 15 Private Use>
		{[]string{"i", "\U000FFFFD"}, "<Plane 15 Private Use, Last>"},  // <Plane 15 Private Use> (Last)
		{[]string{"i", "\U00100000"}, "<Plane 16 Private Use, First>"}, // <Plane 16 Private Use> (First)
		{[]string{"i", "\U00100001"}, "<Plane 16 Private Use>"},        // <Plane 16 Private Use>
		{[]string{"i", "\U0010FFFD"}, "<Plane 16 Private Use, Last>"},  // <Plane 16 Private Use> (Last)
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, "_"), func(t *testing.T) {
			exit, _, out := zli.Test(t)
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

		{[]string{"-qo", "s", "floral", "bullet"}, "WHITE BULLET", 15, -1},

		{[]string{"s", "nomatch_nomatch"}, "no matches", 1, 1},
		{[]string{"-q", "s", "nomatch_nomatch"}, "", 0, 1},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, "_"), func(t *testing.T) {
			exit, _, outbuf := zli.Test(t)
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
		{[]string{"p", "9999999999"}, `out of range: "9999999999"`, 1, 1},
		{[]string{"p", "99999"}, `CODEPOINT NOT IN UNICODE`, 2, -1},

		{[]string{"-q", "p", "0d90"}, "CAPITAL LETTER Z", 1, -1},
		{[]string{"-q", "p", "0D90"}, "CAPITAL LETTER Z", 1, -1},

		{[]string{"p", ""}, `invalid codepoint: not a number or codepoint: ""`, 1, 1},
		{[]string{"p", "nonsense"}, `invalid codepoint: not a number or codepoint: "nonsense"`, 1, 1},
		{[]string{"p", "2042..xxx"}, `invalid codepoint: not a number or codepoint: "xxx"`, 1, 1},
		{[]string{"p", "xxx..xxx"}, `invalid codepoint: not a number or codepoint: "xxx"`, 1, 1},
		{[]string{"p", "xxx..xxx"}, `invalid codepoint: not a number or codepoint: "xxx"`, 1, 1},

		{[]string{"-q", "p", "U+3402"}, "'㐂'", 1, -1},
		{[]string{"-q", "p", "U+3402..U+3404"}, "<CJK Ideograph Extension A>", 3, -1},
		{[]string{"-q", "p", "OtherPunctuation"}, "ASTERISM", 628, -1},
		{[]string{"-q", "p", "Po"}, "ASTERISM", 628, -1},
		{[]string{"-q", "p", "GeneralPunctuation"}, "ASTERISM", 111, -1},
		{[]string{"-q", "p", "all"}, "ASTERISM", 34931, -1},

		{[]string{"-q", "-r", "p", "U9"}, "'\t'", 1, -1},

		// UTF-8
		{[]string{"-q", "p", "utf8:75"}, "'u'", 1, -1},
		{[]string{"-q", "p", "UTF8:75"}, "'u'", 1, -1},
		{[]string{"-q", "p", "utf8:e282ac"}, "'€'", 1, -1},
		{[]string{"-q", "p", "utf8:e2 82 ac"}, "'€'", 1, -1},
		{[]string{"-q", "p", "utf8:0xe20x820xac"}, "'€'", 1, -1},
		{[]string{"-q", "p", "utf8:0xE2 0x82 0xAC"}, "'€'", 1, -1},
		// Issue #46
		{[]string{"-q", "p", "utf8:ef bf bd"}, "U+FFFD", 1, -1},

		// Surrogates
		{[]string{"-q", "p", "U+DC00"}, "␣", 1, -1},                                       // <Low Surrogate> (First)
		{[]string{"-q", "p", "U+DC01"}, "␣", 1, -1},                                       // <Low Surrogate>
		{[]string{"-q", "p", "U+DFFF"}, "␣", 1, -1},                                       // <Low Surrogate> (Last)
		{[]string{"-q", "p", "U+D800"}, "␣", 1, -1},                                       // <Non Private Use High Surrogate> (First)
		{[]string{"-q", "p", "U+D801"}, "␣", 1, -1},                                       // <Non Private Use High Surrogate>
		{[]string{"-q", "p", "U+DB7F"}, "␣", 1, -1},                                       // <Non Private Use High Surrogate> (Last)
		{[]string{"-q", "p", "U+DB80"}, "␣", 1, -1},                                       // <Private Use High Surrogate> (First)
		{[]string{"-q", "p", "U+DB81"}, "␣", 1, -1},                                       // <Private Use High Surrogate>
		{[]string{"-q", "p", "U+DBFF"}, "␣", 1, -1},                                       // <Private Use High Surrogate> (Last)
		{[]string{"-q", "p", "U+DC00"}, "<Low Surrogate, First>", 1, -1},                  // <Low Surrogate> (First)
		{[]string{"-q", "p", "U+DC01"}, "<Low Surrogate>", 1, -1},                         // <Low Surrogate>
		{[]string{"-q", "p", "U+DFFF"}, "<Low Surrogate, Last>", 1, -1},                   // <Low Surrogate> (Last)
		{[]string{"-q", "p", "U+D800"}, "<Non Private Use High Surrogate, First>", 1, -1}, // <Non Private Use High Surrogate> (First)
		{[]string{"-q", "p", "U+D801"}, "<Non Private Use High Surrogate>", 1, -1},        // <Non Private Use High Surrogate>
		{[]string{"-q", "p", "U+DB7F"}, "<Non Private Use High Surrogate, Last>", 1, -1},  // <Non Private Use High Surrogate> (Last)
		{[]string{"-q", "p", "U+DB80"}, "<Private Use High Surrogate, First>", 1, -1},     // <Private Use High Surrogate> (First)
		{[]string{"-q", "p", "U+DB81"}, "<Private Use High Surrogate>", 1, -1},            // <Private Use High Surrogate>
		{[]string{"-q", "p", "U+DBFF"}, "<Private Use High Surrogate, Last>", 1, -1},      // <Private Use High Surrogate> (Last)
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, "_"), func(t *testing.T) {
			exit, _, outbuf := zli.Test(t)
			os.Args = append([]string{"testuni"}, tt.in...)

			func() {
				defer exit.Recover()
				main()
			}()

			if int(*exit) != tt.wantExit {
				t.Fatalf("exit %d: %s", *exit, outbuf.String())
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

		{[]string{"e", "-q", "group:hands"},
			[]string{"👏", "🙌", "🫶", "👐", "🤲", "🤝", "🙏"}},
		{[]string{"e", "-q", "-tone", "dark", "g:hands"},
			[]string{"👏🏿", "🙌🏿", "🫶🏿", "👐🏿", "🤲🏿", "🤝🏿", "🙏🏿"}},

		{[]string{"e", "-q", "shrug"},
			[]string{"🤷"}},
		{[]string{"e", "-q", "shrug", "-gender", "all"},
			[]string{"🤷", "🤷Z♂S", "🤷Z♀S"}},
		{[]string{"e", "-q", "-gender", "m", "shrug"},
			[]string{"🤷Z♂S"}},
		{[]string{"e", "-q", "-gender", "m", "-tone", "light", "shrug"},
			[]string{"🤷🏻Z♂S"}},

		{[]string{"e", "-q", "farmer"},
			[]string{"🧑Z🌾"}},
		{[]string{"e", "-q", "farmer", "-gender", "all"},
			[]string{"🧑Z🌾", "👨Z🌾", "👩Z🌾"}},
		{[]string{"e", "-q", "-gender", "f,m", "farmer"},
			[]string{"👨Z🌾", "👩Z🌾"}},
		{[]string{"e", "-q", "-gender", "f", "-t", "medium", "farmer"},
			[]string{"👩🏽Z🌾"}},

		{[]string{"e", "-q", "-gender", "p", "detective"},
			[]string{"🕵S"}},
		{[]string{"e", "-q", "-gender", "p", "-tone", "mediumdark", "detective"},
			[]string{"🕵🏾"}},
		{[]string{"e", "-q", "-gender", "m", "detective"},
			[]string{"🕵SZ♂S"}},
		{[]string{"e", "-q", "-gender", "m", "-tone", "mediumdark", "detective"},
			[]string{"🕵🏾Z♂S"}},

		{[]string{"e", "-qo", "zimbabwe", "#", "england"},
			[]string{"#S⃣", "🇿🇼", "🏴󠁧󠁢󠁥󠁮󠁧󠁿"}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.in, "_"), func(t *testing.T) {
			_, _, outbuf := zli.Test(t)
			os.Args = append([]string{"testuni"}, tt.in...)

			main()

			var out []string
			for _, line := range strings.Split(strings.TrimSpace(outbuf.String()), "\n") {
				line = strings.ReplaceAll(line, "\t", " ")
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

func TestJSON(t *testing.T) {
	_, _, outbuf := zli.Test(t)

	os.Args = append([]string{"testuni"}, "i", "€", "-f=all", "-j")
	main()

	want := ` [{
	"bin": "10000010101100",
	"block": "Currency Symbols",
	"cat": "Currency_Symbol",
	"char": "€",
	"cpoint": "U+20AC",
	"dec": "8364",
	"digraph": "=e",
	"hex": "20ac",
	"html": "&euro;",
	"json": "\\u20ac",
	"keysym": "EuroSign",
	"name": "EURO SIGN",
	"oct": "20254",
	"plane": "Basic Multilingual Plane",
	"props": "",
	"script": "Common",
	"unicode": "2.1",
	"utf16be": "20 ac",
	"utf16le": "ac 20",
	"utf8": "e2 82 ac",
	"width": "ambiguous",
	"xml": "&#x20ac;"
}]
`
	got := outbuf.String()

	if d := ztest.Diff(want, got); d != "" {
		t.Error(d)
	}
}

func BenchmarkUni(b *testing.B) {
	zli.Stdout = new(bytes.Buffer)

	b.Run("print one", func(b *testing.B) {
		os.Args = []string{"uni", "p", "20ac"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})
	b.Run("print three", func(b *testing.B) {
		os.Args = []string{"uni", "p", "20ac", "20ad", "20ae"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})

	b.Run("print all", func(b *testing.B) {
		os.Args = []string{"uni", "p", "all"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})
	b.Run("print all json", func(b *testing.B) {
		os.Args = []string{"uni", "p", "all", "-j"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})

	b.Run("print all columns", func(b *testing.B) {
		os.Args = []string{"uni", "p", "all", "-f", "all"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})
	b.Run("print all columns json", func(b *testing.B) {
		os.Args = []string{"uni", "p", "all", "-f", "all", "-j"}
		for n := 0; n < b.N; n++ {
			main()
		}
	})
}
