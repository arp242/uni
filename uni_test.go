package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
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
		in        []string
		want      string
		wantLines int
	}{
		{[]string{"hands"}, "clapping", 6},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			outbuf, c := setup(t, "e", tt.in)
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
