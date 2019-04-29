[![This project is considered experimental](https://img.shields.io/badge/Status-experimental-red.svg)](https://arp242.net/status/experimental)

`uni` queries the Unicode database from the commandline.

There are three commands: `identify` to print Unicode information about a
string, `search` to search for codepoints, and `print` to print groups of
Unicode classes.

Install it with `go get arp242.net/uni`, which will put the binary at
`~/go/bin/uni`.

Usage
-----

Identify a character:

	$ uni identify €
	     cpoint  dec    utf-8      html       name
	'€'  U+20AC  8364   0xe282ac   &euro;     EURO SIGN

Or an entire string. `i` is a shortcut for `identify`:

	$ uni i h€łłø
	     cpoint  dec    utf-8       html       name
	'h'  U+0068  104    68          &#x68;     LATIN SMALL LETTER H
	'€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN
	'ł'  U+0142  322    c5 82       &lstrok;   LATIN SMALL LETTER L WITH STROKE
	'ł'  U+0142  322    c5 82       &lstrok;   LATIN SMALL LETTER L WITH STROKE
	'ø'  U+00F8  248    c3 b8       &oslash;   LATIN SMALL LETTER O WITH STROKE

Identify byte offset from a file (useful for editor integration):

	$ uni i 'README.markdown:#0'
	      cpoint  dec    utf-8       html       name
	 '`'  U+0060  96     60          &grave;    GRAVE ACCENT

Or a range from a file:

	$ uni i 'README.markdown:#0-4'
	     cpoint  dec    utf-8       html       name
	'`'  U+0060  96     60          &grave;    GRAVE ACCENT
	'u'  U+0075  117    75          &#x75;     LATIN SMALL LETTER U
	'n'  U+006E  110    6e          &#x6e;     LATIN SMALL LETTER N
	'i'  U+0069  105    69          &#x69;     LATIN SMALL LETTER I
	'`'  U+0060  96     60          &grave;    GRAVE ACCENT

Note that these are **byte** offsets, not *character* offsets:

	$ uni i 'README.markdown:#130'
	uni: WARNING: input string is not valid UTF-8
	     cpoint  dec    utf-8       html       name
	'�'  U+FFFD  65533  ef bf bd    &#xfffd;   REPLACEMENT CHARACTER

	$ uni i 'README.markdown:#130-132'
	     cpoint  dec    utf-8       html       name
	'€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN

Search description:

	$ uni search euro
	     cpoint  dec    utf-8      html       name
	'₠'  U+20A0  8352   e2 82 a0    &#x20a0;   EURO-CURRENCY SIGN
	'€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN
	'𐡷'  U+10877 67703  f0 90 a1 b7 &#x10877;  PALMYRENE LEFT-POINTING FLEURON
	'𐡸'  U+10878 67704  f0 90 a1 b8 &#x10878;  PALMYRENE RIGHT-POINTING FLEURON
	'𐫱'  U+10AF1 68337  f0 90 ab b1 &#x10af1;  MANICHAEAN PUNCTUATION FLEURON
	'🌍'  U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
	'🏤'  U+1F3E4 127972 f0 9f 8f a4 &#x1f3e4;  EUROPEAN POST OFFICE
	'🏰'  U+1F3F0 127984 f0 9f 8f b0 &#x1f3f0;  EUROPEAN CASTLE
	'💶'  U+1F4B6 128182 f0 9f 92 b6 &#x1f4b6;  BANKNOTE WITH EURO SIGN

The `s` command is a shortcut for `search`. Multiple words are matched
individually:

	$ uni s earth globe
	      cpoint  dec    utf-8       html       name
	'🌍'  U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
	'🌎'  U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS
	'🌏'  U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA

	$ uni s globe earth
	      cpoint  dec    utf-8       html       name
	'🌍'  U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
	'🌎'  U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS
	'🌏'  U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA

Use standard shell quoting for more literal matches:

	$ uni s rightwards black arrow
	     cpoint  dec    utf-8       html       name
	'➡'  U+27A1  10145  e2 9e a1    &#x27a1;   BLACK RIGHTWARDS ARROW
	'➤'  U+27A4  10148  e2 9e a4    &#x27a4;   BLACK RIGHTWARDS ARROWHEAD
	'➥'  U+27A5  10149  e2 9e a5    &#x27a5;   HEAVY BLACK CURVED DOWNWARDS AND RIGHTWARDS ARROW
	[..]

	$ uni s 'rightwards black arrow'
	     cpoint  dec    utf-8       html       name
	'⮕'  U+2B95  11157  e2 ae 95    &#x2b95;   RIGHTWARDS BLACK ARROW

The `print` command (shortcut `p`) can be used to print groups of characters:

	$ uni print printable
	[..]

	$ uni p emoji
	[..]

