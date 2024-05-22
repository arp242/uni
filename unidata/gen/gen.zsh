#!/usr/bin/env zsh
[ "${ZSH_VERSION:-}" = "" ] && echo >&2 "Only works with zsh" && exit 1
setopt no_unset pipefail
cd $0:P:h:h

# TODO: add "confusable" information from
# https://www.unicode.org/Public/security/13.0.0/
# https://www.unicode.org/reports/tr39/

use_cache=0
use_beta=0
for a in $argv; do
	case $a in
		cache) use_cache=1; shift ;;
		beta)  use_beta=1;  shift ;;
	esac
done

get() {
	if [[ $use_beta = 1 && $1 =~ '^https://www.unicode.org/Public/UCD/latest/' ]] then
		1=${1/UCD\/latest/draft\/UCD}
	fi

	if [[ $use_cache = 1 && -f .cache/$1:t ]] then
		print "Using cache at .cache/$1:t"
		return
	fi
	print "Fetching $1"
	curl -sL $1 >.cache/$1:t
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
get 'https://www.unicode.org/Public/UCD/latest/ucd/NamesList.txt'
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
[[ $1 =~ "all|names?"      ]] && mk names      '.cache/NamesList.txt'
[[ $1 =~ "all|scripts?"    ]] && mk scripts    '.cache/Scripts.txt'
[[ $1 =~ "all|emojis?"     ]] && mkgo emojis   '.cache/emoji-test.txt' '.cache/en.xml'

exit 0
