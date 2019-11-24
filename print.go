package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"arp242.net/uni/isatty"
	"arp242.net/uni/terminal"
)

var (
	isTerm    = isatty.IsTerminal(os.Stdout.Fd())
	termWidth = func() int {
		if !isTerm {
			return 0
		}

		w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil || w < 50 {
			return 0
		}
		return w
	}()
)

type printer []char

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

// TODO: add option to choose sorting order.
func (p *printer) PrintSorted(fp io.Writer, quiet, raw bool) {
	s := []char(*p)
	sort.Slice(s, func(i int, j int) bool { return s[i].codepoint < s[j].codepoint })
	p.Print(fp, quiet, raw)
}

func (p printer) fmtChar(b *strings.Builder, info char, raw bool) {
	c := rune(info.codepoint)

	size := 44
	b.WriteString(fmt.Sprintf("'%v' ", fmtChar(c, raw)))
	if info.width != widthFullWidth && info.width != widthWide {
		size++
		b.WriteString(" ")
	}

	name := fmt.Sprintf("%s (%s)", info.name, catnames[info.cat])
	if isTerm && termWidth > 0 && len(name) > termWidth-size {
		name = name[:termWidth-size] + "…"
	}

	b.WriteString(fmt.Sprintf("U+%s %s %s %s %s\n",
		fill(fmt.Sprintf("%04X", info.codepoint), 5),
		fill(strconv.FormatUint(uint64(info.codepoint), 10), 6),
		fill(p.utf8(c), 11),
		fill(p.entity(c, info.codepoint), 10),
		name))
}

func (p printer) entity(c rune, cp uint32) string {
	html := entities[c]
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

func fill(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
