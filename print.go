package main

import (
	"fmt"
	"io"
	"os"
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
)

type printer []unidata.Codepoint

func (p printer) Print(fp io.Writer, quiet, raw bool) {
	if len(p) == 0 {
		return
	}

	var b strings.Builder
	if !quiet {
		b.WriteString("     cpoint  dec    utf-8       html       name\n")
	}

	for _, c := range p {
		p.fmtChar(&b, c, raw)
	}
	fmt.Fprint(fp, b.String())
}

func (p *printer) PrintSorted(fp io.Writer, quiet, raw bool) {
	s := []unidata.Codepoint(*p)
	sort.Slice(s, func(i int, j int) bool { return s[i].Codepoint < s[j].Codepoint })
	p.Print(fp, quiet, raw)
}

func (p printer) fmtChar(b *strings.Builder, info unidata.Codepoint, raw bool) {
	c := rune(info.Codepoint)

	size := 44
	b.WriteString(fmt.Sprintf("'%v' ", fmtChar(c, raw)))
	if info.Width != unidata.WidthFullWidth && info.Width != unidata.WidthWide {
		size++
		b.WriteString(" ")
	}

	name := fmt.Sprintf("%s (%s)", info.Name, unidata.Catnames[info.Cat])
	if isTerm && termWidth > 0 && utf8.RuneCountInString(name) > termWidth-size {
		name = name[:termWidth-size] + "…"
	}

	b.WriteString(fmt.Sprintf("U+%s %s %s %s %s\n",
		zstring.AlignLeft(fmt.Sprintf("%04X", info.Codepoint), 5),
		zstring.AlignLeft(strconv.FormatUint(uint64(info.Codepoint), 10), 6),
		zstring.AlignLeft(p.utf8(c), 11),
		zstring.AlignLeft(p.entity(c, info.Codepoint), 10),
		name))
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
