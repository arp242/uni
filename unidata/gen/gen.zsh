#!/usr/bin/env zsh
[ "${ZSH_VERSION:-}" = "" ] && echo >&2 "Only works with zsh" && exit 1
setopt no_unset pipefail
cd $0:P:h:h

# TODO: add casefolding
# https://unicode.org/Public/13.0.0/ucd/CaseFolding.txt
# CaseFold []rune
# 
# TODO: add properties:
# https://unicode.org/Public/13.0.0/ucd/PropList.txt
# "uni p dash" should print all dashes.
# 
# 
# TODO: add "confusable" information from
# https://www.unicode.org/Public/idna/13.0.0/
# and/or
# https://www.unicode.org/Public/security/13.0.0/
# 
# 
# TODO: add "alias" information from
# https://unicode.org/Public/13.0.0/ucd/NamesList.txt
# This is generated from other sources, but I can't really find where it gts
# that "x (modifier letter prime - 02B9)" from.
# 
# 0027 APOSTROPHE
#     = apostrophe-quote (1.0)
#     = APL quote
#     * neutral (vertical) glyph with mixed usage
#     * 2019 is preferred for apostrophe
#     * preferred characters in English for paired quotation marks are 2018 & 2019
#     * 05F3 is preferred for geresh when writing Hebrew
#     x (modifier letter prime - 02B9)
#     x (modifier letter apostrophe - 02BC)
#     x (modifier letter vertical line - 02C8)
#     x (combining acute accent - 0301)
#     x (hebrew punctuation geresh - 05F3)
#     x (prime - 2032)
#     x (latin small letter saltillo - A78C)

get() { [[ -f .cache/$1:t ]] || curl -sL $1 >.cache/$1:t; }
mk() {
	local go=gen_$1.go
	gawk -f gen/$1.awk $2 >$go || exit $?
	err=$(gofmt -w $go 2>&1)
	if [[ $? -ne 0 ]]; then
		for line in ${(ps:\n:)err}; \
			printf "%s\n\t%s\n\n" "$line" "$(head -n ${${(s/:/)line}[2]} $go | tail -n1)"
		exit 1
	fi
}

mkdir -p .cache
get 'https://www.unicode.org/Public/UCD/latest/ucd/Blocks.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/PropertyValueAliases.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt'
get 'https://www.unicode.org/Public/UCD/latest/ucd/EastAsianWidth.txt'
get 'https://www.unicode.org/Public/emoji/14.0/emoji-test.txt'
get 'https://html.spec.whatwg.org/entities.json'
get 'https://gitlab.freedesktop.org/xorg/proto/xorgproto/-/raw/master/include/X11/keysymdef.h'
get 'https://tools.ietf.org/rfc/rfc1345.txt'
get 'https://raw.githubusercontent.com/unicode-org/cldr/master/common/annotations/en.xml'

# export LC_ALL=C
mk blocks     '.cache/Blocks.txt'
mk cats       '.cache/PropertyValueAliases.txt'
mk emojis     '.cache/emoji-test.txt'
mk codepoints '.cache/UnicodeData.txt'
