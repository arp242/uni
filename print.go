package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"arp242.net/uni/unidata"
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
	reFormat = regexp.MustCompile(`%\((char|wide_padding|cpoint|hex|dec|utf8|html|name|cat)(?: .+?)?\)`)
)

type printer []unidata.Codepoint

func formatFlags(in string, flags []string, header bool) string {
	if len(flags) == 0 {
		return in
	}

	var (
		quote      bool
		trim       bool
		alignLeft  int
		alignRight int
	)
	for _, f := range flags {
		if f == "" {
			continue
		}
		if f == "q" {
			quote = true
			continue
		}
		if f == "t" {
			trim = true
			continue
		}
		if f[0] == 'l' || f[0] == 'r' {
			n := strings.Split(f, ":")
			if len(n) != 2 {
				return fmt.Sprintf("%%(ERROR format %q: need width after :)", in)
			}

			num, err := strconv.ParseInt(n[1], 10, 8)
			if err != nil {
				return fmt.Sprintf("%%(ERROR format %q: %s)", in, err)
			}

			if f[0] == 'l' {
				alignLeft = int(num)
			} else {
				alignRight = int(num)
			}

			continue
		}

		return fmt.Sprintf("%%(ERROR format %q: unknown flag %q", in, f)
	}

	if quote {
		if header {
			in = " " + in + " "
		} else {
			in = "'" + in + "'"
		}
	}
	if alignLeft > 0 {
		in = zstring.AlignLeft(in, alignLeft)
	}
	if alignRight > 0 {
		in = zstring.AlignRight(in, alignRight)
	}
	if trim {
		// TODO
		// if isTerm && termWidth > 0 && utf8.RuneCountInString(in) > termWidth-size {
		// 	name = name[:termWidth-size] + "…"
		// }
	}

	return in
}

func (p printer) Print(fp io.Writer, format string, quiet, raw bool) {
	if len(p) == 0 {
		return
	}

	if format == "" {
		format = "%(char q l:3)%(wide_padding) %(cpoint l:7) %(dec l:6) %(utf8 l:11) %(html l:10) %(name t) (%(cat t))"
	}

	var b strings.Builder
	if !quiet {
		b.WriteString(reFormat.ReplaceAllStringFunc(format, func(m string) string {
			s := zstring.Fields(m[2:len(m)-1], " ")
			switch s[0] {
			case "char":
				s[0] = ""
			case "wide_padding":
				s[0] = " "
			case "cpoint":
				s[0] = "cpoint"
			case "hex":
				s[0] = "hex"
			case "dec":
				s[0] = "dec"
			case "utf8":
				s[0] = "utf-8"
			case "html":
				s[0] = "html"
			case "name":
				s[0] = "name"
			case "cat":
				s[0] = "category"
			}
			return formatFlags(s[0], s[1:], true)
		}))
		b.WriteString("\n")
	}

	for _, c := range p {
		p.fmtChar(&b, c, format, raw)
	}
	fmt.Fprint(fp, b.String())
}

func (p *printer) PrintSorted(fp io.Writer, format string, quiet, raw bool) {
	s := []unidata.Codepoint(*p)
	sort.Slice(s, func(i int, j int) bool { return s[i].Codepoint < s[j].Codepoint })
	p.Print(fp, format, quiet, raw)
}

func (p printer) fmtChar(b *strings.Builder, info unidata.Codepoint, format string, raw bool) {
	c := rune(info.Codepoint)

	b.WriteString(reFormat.ReplaceAllStringFunc(format, func(m string) string {
		s := zstring.Fields(m[2:len(m)-1], " ")
		switch s[0] {
		case "char":
			s[0] = fmtChar(c, raw)
		case "wide_padding":
			if info.Width != unidata.WidthFullWidth && info.Width != unidata.WidthWide {
				s[0] = " "
			} else {
				s[0] = ""
			}
		case "cpoint":
			s[0] = fmt.Sprintf("U+%04X", info.Codepoint)
		case "hex":
			s[0] = fmt.Sprintf("%04X", info.Codepoint)
		case "dec":
			s[0] = strconv.FormatUint(uint64(info.Codepoint), 10)
		case "utf8":
			s[0] = p.utf8(c)
		case "html":
			s[0] = p.entity(c, info.Codepoint)
		case "name":
			s[0] = info.Name
		case "cat":
			s[0] = unidata.Catnames[info.Cat]
		}
		return formatFlags(s[0], s[1:], false)
	}))
	b.WriteString("\n")

	// size := 44
	// b.WriteString(fmt.Sprintf("'%v' ", fmtChar(c, raw)))
	// if info.Width != unidata.WidthFullWidth && info.Width != unidata.WidthWide {
	// 	size++
	// 	b.WriteString(" ")
	// }

	// name := fmt.Sprintf("%s (%s)", info.Name, unidata.Catnames[info.Cat])
	// if isTerm && termWidth > 0 && utf8.RuneCountInString(name) > termWidth-size {
	// 	name = name[:termWidth-size] + "…"
	// }
}

func (p printer) entity(c rune, cp uint32) string {
	html := unidata.Entities[c]
	if html == "" {
		html = fmt.Sprintf("#x%x", cp)
	}
	return "&" + html + ";"
}

func (p printer) utf8(r rune) string {
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, r)
	return fmt.Sprintf("% x", buf[:n])
}
