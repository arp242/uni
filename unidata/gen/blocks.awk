BEGIN        { FS = " *; *" }
/^$/ || /^#/ { next }

{
    split($1, se, /\.\./)
    start = strtonum("0x" se[1])
    end   = strtonum("0x" se[2])
    name  = $2
    const = "Block" gensub(/[ -]/, "", "g", name)

    constants = constants "\t" const "\n"
    blocks    = blocks sprintf("\t%s: {[2]rune{0x%06X, 0x%06X}, \"%s\"},\n", const, start, end, name)
}

END {
    print("// Code generated by gen.zsh; DO NOT EDIT\n\npackage unidata\n")

    print("// Unicode blocks\nconst (\n" \
          "\tBlockUnknown = Block(iota)\n" \
          constants ")\n")

    print("// Blocks is a list of all Unicode blocks.\n" \
          "var Blocks = map[Block]struct {\n" \
              "\tRange [2]rune\n" \
              "\tName  string\n" \
          "}{\n" blocks "}")
}