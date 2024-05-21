#!/usr/bin/env zsh
[ "${ZSH_VERSION:-}" = "" ] && echo >&2 "Only works with zsh" && exit 1
setopt no_unset pipefail
cd $0:P:h:h

# TODO: add "confusable" information from
# https://www.unicode.org/Public/security/13.0.0/
# https://www.unicode.org/reports/tr39/

# TODO: add "alias" information from
# https://www.unicode.org/Public/UCD/latest/ucd/NamesList.txt
# https://www.unicode.org/Public/UCD/latest/ucd/NamesList.html
# https://www.unicode.org/versions/Unicode14.0.0/ch24.pdf
#
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

get() {
	if [[ ! -f .cache/$1:t ]]; then
		print "Fetching $1"
		curl -sL $1 >.cache/$1:t
	fi
}
mk() {
	local go=gen_$1.go
	print "Generating $go"

	if [[ ${PRINT:-} -eq 1 ]]; then
		gawk -f gen/$1.awk $argv[2,-1] || exit $?
		return 0
	fi

	gawk -f gen/$1.awk $argv[2,-1] >$go || exit $?
	err=$(gofmt -w $go 2>&1)
	if [[ $? -ne 0 ]]; then
		for line in ${(ps:\n:)err}; \
			printf "%s\n\t%s\n\n" "$line" "$(head -n ${${(s/:/)line}[2]} $go | tail -n1)"
		exit 1
	fi
}
mkgo() {
	local go=gen_$1.go
	print "Generating $go"

	if [[ ${PRINT:-} -eq 1 ]]; then
		go run gen/$1.go $argv[2,-1] || exit $?
		return 0
	fi

	go run gen/$1.go $argv[2,-1] >$go || exit $?
	err=$(gofmt -w $go 2>&1)
	if [[ $? -ne 0 ]]; then
		for line in ${(ps:\n:)err}; \
			printf "%s\n\t%s\n\n" "$line" "$(head -n ${${(s/:/)line}[2]} $go | tail -n1)"
		exit 1
	fi
}

# go run ./gen/emojis2.go .cache/emoji-test.txt .cache/en.xml |gofmt >!x

mkdir -p .cache
get 'https://www.unicode.org/Public/UCD/latest/ucd/Blocks.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/DerivedAge.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/EastAsianWidth.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/PropList.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/PropertyValueAliases.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/Scripts.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt'
get 'https://www.unicode.org/Public/emoji/latest/emoji-test.txt'
get 'https://html.spec.whatwg.org/entities.json'
get 'https://gitlab.freedesktop.org/xorg/proto/xorgproto/-/raw/master/include/X11/keysymdef.h'
get 'https://tools.ietf.org/rfc/rfc1345.txt'
get 'https://raw.githubusercontent.com/unicode-org/cldr/master/common/annotations/en.xml'


1=${1:-all}
[[ $1 =~ "all|props?"      ]] && mk props      '.cache/PropList.txt'
[[ $1 =~ "all|blocks?"     ]] && mk blocks     '.cache/Blocks.txt'
[[ $1 =~ "all|cats?"       ]] && mk cats       '.cache/PropertyValueAliases.txt'
[[ $1 =~ "all|codepoints?" ]] && mk codepoints '.cache/UnicodeData.txt'
[[ $1 =~ "all|scripts?"    ]] && mk scripts    '.cache/Scripts.txt'
[[ $1 =~ "all|emojis?"     ]] && mkgo emojis   '.cache/emoji-test.txt' '.cache/en.xml'

exit 0
