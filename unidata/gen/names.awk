# Format:
#    <codepoint>\t<description>
#    \t= alias
#    \t* comments/see also
#    \tx cross-reference
#    \t# comment(?)
#
# Example:
#
#   00B2 SUPERSCRIPT TWO
#       = squared
#       * other superscript digit characters: 2070-2079
#       x (superscript one - 00B9)
#       # <super> 0032
#
#   0027 APOSTROPHE
#       = apostrophe-quote (1.0)
#       = APL quote
#       * neutral (vertical) glyph with mixed usage
#       * 2019 is preferred for apostrophe
#       * preferred characters in English for paired quotation marks are 2018 & 2019
#       * 05F3 is preferred for geresh when writing Hebrew
#       x (modifier letter prime - 02B9)
#       x (modifier letter apostrophe - 02BC)
#       x (modifier letter vertical line - 02C8)
#       x (combining acute accent - 0301)
#       x (hebrew punctuation geresh - 05F3)
#       x (prime - 2032)
#       x (latin small letter saltillo - A78C)
#
# 0021	EXCLAMATION MARK
# 	= factorial
# 	= bang
# 	x (inverted exclamation mark - 00A1)
# 	x (latin letter retroflex click - 01C3)
# 	x (double exclamation mark - 203C)
# 	x (interrobang - 203D)
# 	x (warning sign - 26A0)
# 	x (heavy exclamation mark symbol - 2757)
# 	x (heavy exclamation mark ornament - 2762)
# 	x (medieval exclamation mark - 2E53)
# 	x (modifier letter raised exclamation mark - A71D)


#BEGIN        { FS = " *; *" }

# Skip everything before first control character, because the format is idiotic.
/0000/       { started = 1 }
!started     { next }
/^$/ || /^@/ { next }

# Alias and "formal alias".
/^\t[=%] / {
    # Skip controls for aliases, as they're always just the same as the actual
    # name.
    if (!(cp < 0x20 || cp == 0x7f || (cp >= 0x80 && cp <= 0x9f))) {
        # Stuff between (..) is rarely relevant, so just remove it.
        # TODO:
        # 0x2052
        # 0x2694
        # More?
        l = skip(1)
        gsub(" \\([a-zA-Z0-9. :°/,-]+\\)$", "", l)
        gsub("\"", "", l)
        split(l, arr, ", ")
        for (k in arr) aliases[++a] = arr[k]

        # Special cases
        if (cp == 0x2c) {
            # Not an alias: "the use as decimal or thousands separator is locale dependent"
            # TODO: report as bug?
            delete aliases[a--]
        }
        if (cp == 0xc6) {  # Æ
            delete aliases[a--]
        }
        if (cp == 0x00e6 && index(l, "(") > 0) {  # æ
            aliases[a - 1] = "ash"  # Don't include "ligature" line
            aliases[a]     = "æsc"
        }
        if (cp == 0x0153 && index(l, "(") > 0) {  # œ
            aliases[a - 1] = "eðel"
            aliases[a]     = "ethel"
        }
    }
    next
}

# Cross reference
/^\tx / {
    refs[++r] = gensub("\\((.+) - ([A-F0-9]{4,})\\)$", "\\2", "g", skip(1))
    next
}

# COMMENT_LINE
# /^\t\* /        { notes[++n]   = $0; next }

# VARIATION_LINE
# /^\t~ /         { x[++a] = $0; next }

# DECOMPOSITION
# /^\t: /         { x[++a] = $0; next }

# COMPAT_MAPPING
# /^\t# /         { comments[++c] = $0; next }


BEGIN {
    print("// Code generated by gen.zsh; DO NOT EDIT\n\npackage unidata\n")
    print("var names = map[rune]name{")
}

/^[^\t]/ {
    if (a + n + r == 0) {  # First item or nothing to print.
        cp = strtonum("0x" $1)
        next
    }

    printf("\t%s: {\n", prchr(cp))
    if (a > 0) {
        printf("\t\taliases: []string{")
        for (a in aliases) printf("`%s`,", aliases[a])
        print("},")
    }
    if (r > 0) {
        printf("\t\trefs: []rune{")
        for (r in refs) printf("0x%s,", refs[r])
        print("},")
    }

    # for (n in notes)    print "\t" notes[n]
    # for (r in comments) print "\t" comments[c]
    # printf("\n")
    printf("\t},\n")

    cp = strtonum("0x" $1)
    delete aliases; delete refs; delete notes; delete comments
    a = 0; r = 0; n = 0; c = 0
}

END {
    print("}")
}

function skip(n,    t) {
    n += 1
    for (i = n; i <= NF; i++) {
        if (i > n)
            t = t " "
        t = t $i
    }
    return t
}

function prchr(cp) {
    return sprintf("0x%02x", cp)

    if (cp < 0x20 || cp == 0x7f || (cp >= 0x80 && cp <= 0x9f))  # TODO: also other unprintable
        return sprintf("0x%02x", cp)
    return sprintf("'%c'", cp)
}
