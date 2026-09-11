`uni` queries the Unicode database from the commandline. It supports Unicode
15.1 (September 2023) and has good support for emojis.

There are four commands: `identify` codepoints in a string, `search` for
codepoints, `print` codepoints by class, block, or range, and `emoji` to find
emojis.

There are binaries on the [releases] page, and [packages] for a number of
platforms. You can also [run it in your browser][uni-wasm].

Compile from source with:

    % go install zgo.at/uni/v2@latest

which will give you a `uni` binary in `~/go/bin`.

README index:
- [Integrations](#integrations)
- [Usage](#usage)
  - [Identify](#identify)
  - [Search](#search)
  - [Print](#identify)
  - [Emoji](#emoji)
  - [JSON](#json)
- [ChangeLog](#changelog)
- [Development](#development)
- [Alternatives](#alternatives)

[uni-wasm]: https://arp242.github.io/uni-wasm/
[releases]: https://github.com/arp242/uni/releases
[packages]: https://repology.org/project/uni/versions

Integrations
------------

- [dmenu], [rofi], and [fzf] script at [`dmenu-uni`](/dmenu-uni). See the top of
  the script for some options you may want to frob with.

- For a Vim command see [`uni.vim`](/uni.vim); just copy/paste it in your vimrc.

[dmenu]: http://tools.suckless.org/dmenu
[rofi]: https://github.com/davatorium/rofi
[fzf]: https://github.com/junegunn/fzf

Usage
-----
*Note: the alignment is slightly off for some entries due to the way GitHub
renders wide characters; in terminals it should be aligned correctly.*

### Identify

Identify characters in a string, as a kind of a unicode-aware `hexdump`:

    % uni identify €
                 Dec    UTF8        HTML       Name
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN

`i` is a shortcut for `identify`:

    % uni i h€ý
                 Dec    UTF8        HTML       Name
    'h'  U+0068  104    68          &#x68;     LATIN SMALL LETTER H
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN
    'ý'  U+00FD  253    c3 bd       &yacute;   LATIN SMALL LETTER Y WITH ACUTE

It reads from stdin:

     % head -c5 README.md | uni i
          CPoint  Dec    UTF8        HTML       Name
     '`'  U+0060  96     60          &grave;    GRAVE ACCENT [backtick, backquote]
     'u'  U+0075  117    75          &#x75;     LATIN SMALL LETTER U
     'n'  U+006E  110    6e          &#x6e;     LATIN SMALL LETTER N
     'i'  U+0069  105    69          &#x69;     LATIN SMALL LETTER I
     '`'  U+0060  96     60          &grave;    GRAVE ACCENT [backtick, backquote]

    % echo 'U+1234 U+1111' | uni p
         CPoint  Dec    UTF8        HTML       Name
    'ᄑ' U+1111  4369   e1 84 91    &#x1111;   HANGUL CHOSEONG PHIEUPH [P]
    'ሴ'  U+1234  4660   e1 88 b4    &#x1234;   ETHIOPIC SYLLABLE SEE

You can use `-compact` (or `-c`) to suppress the header, and `-format` (or `-f`)
to control the output format:

    % uni i -f '%unicode %name' a€🧟
    Unicode Name
    1.1     LATIN SMALL LETTER A
    2.1     EURO SIGN
    10.0    ZOMBIE

If the format string starts with `+` it will automatically be prepended with the
character, codepoint, and name:

    % uni i -f +%unicode a€🧟
                 Name                 Unicode
    'a'  U+0061  LATIN SMALL LETTER A 1.1
    '€'  U+20AC  EURO SIGN            2.1
    '🧟' U+1F9DF ZOMBIE               10.0

You can add more advanced options with `%(name flags)`; for example to generate
an aligned codepoint to X11 keysym mapping:

    % uni i -c -f '0x%(hex l:auto f:0): %(keysym l:auto q:":",) // %name' h€ý
    0x6800: "h",        // LATIN SMALL LETTER H
    0x20ac: "EuroSign", // EURO SIGN
    0xfd00: "yacute",   // LATIN SMALL LETTER Y WITH ACUTE

See `uni help` for more details on the `-format` flag; this flag can also be
added to other commands.

### Search

Search description:

    % uni search euro
                 Dec    UTF8        HTML       Name
    '₠'  U+20A0  8352   e2 82 a0    &#x20a0;   EURO-CURRENCY SIGN
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN
    '𐡷'  U+10877 67703  f0 90 a1 b7 &#x10877;  PALMYRENE LEFT-POINTING FLEURON
    '𐡸'  U+10878 67704  f0 90 a1 b8 &#x10878;  PALMYRENE RIGHT-POINTING FLEURON
    '𐫱'  U+10AF1 68337  f0 90 ab b1 &#x10af1;  MANICHAEAN PUNCTUATION FLEURON
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
    '🏤' U+1F3E4 127972 f0 9f 8f a4 &#x1f3e4;  EUROPEAN POST OFFICE
    '🏰' U+1F3F0 127984 f0 9f 8f b0 &#x1f3f0;  EUROPEAN CASTLE
    '💶' U+1F4B6 128182 f0 9f 92 b6 &#x1f4b6;  BANKNOTE WITH EURO SIGN

The `s` command is a shortcut for `search`. Multiple words are matched
individually:

    % uni s globe earth
                 Dec    UTF8        HTML       Name
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA

Use shell quoting for more literal matches:

    % uni s rightwards black arrow
                 Dec    UTF8        HTML       Name
    '➡'  U+27A1  10145  e2 9e a1    &#x27a1;   BLACK RIGHTWARDS ARROW
    '➤'  U+27A4  10148  e2 9e a4    &#x27a4;   BLACK RIGHTWARDS ARROWHEAD
    …

    % uni s 'rightwards black arrow'
                 Dec    UTF8        HTML       Name
    '⮕'  U+2B95  11157  e2 ae 95    &#x2b95;   RIGHTWARDS BLACK ARROW

Add `-or` or `-o` to combine the search terms with "OR" instead of "AND":

    % uni s -o globe milky
                 Dec    UTF8        HTML       Name
    '🌌' U+1F30C 127756 f0 9f 8c 8c &#x1f30c;  MILKY WAY
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA
    '🌐' U+1F310 127760 f0 9f 8c 90 &#x1f310;  GLOBE WITH MERIDIANS

### Print

Print specific codepoints or groups of codepoints:

    % uni print U+2042
                 Dec    UTF8        HTML       Name
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM

Print a custom range; `U+2042`, `U2042`, and `2042` are all identical:

    % uni print 2042..2044
                 Dec    UTF8        HTML       Name
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM
    '⁃'  U+2043  8259   e2 81 83    &hybull;   HYPHEN BULLET
    '⁄'  U+2044  8260   e2 81 84    &frasl;    FRACTION SLASH [solidus]

You can also use hex, octal, and binary numbers: `0x2024`, `0o20102`, or
`0b10000001000010`.

General category:

    % uni p Po
    Showing category Po (Other_Punctuation)
                 Dec    UTF8        HTML       Name
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK [factorial, bang]
    …

Blocks:

    % uni p arrows 'box drawing'
    Showing block Arrows
    Showing block Box Drawing
                 Dec    UTF8        HTML       Name
    '←'  U+2190  8592   e2 86 90    &larr;     LEFTWARDS ARROW
    '↑'  U+2191  8593   e2 86 91    &uarr;     UPWARDS ARROW
    …

Print as table, and with a shorter name:

    % uni p -as table box
    Showing block Box Drawing
             0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
           ┌────────────────────────────────────────────────────────────────
    U+250x │ ─   ━   │   ┃   ┄   ┅   ┆   ┇   ┈   ┉   ┊   ┋   ┌   ┍   ┎   ┏
           │
    U+251x │ ┐   ┑   ┒   ┓   └   ┕   ┖   ┗   ┘   ┙   ┚   ┛   ├   ┝   ┞   ┟
           │
    U+252x │ ┠   ┡   ┢   ┣   ┤   ┥   ┦   ┧   ┨   ┩   ┪   ┫   ┬   ┭   ┮   ┯
           │
    U+253x │ ┰   ┱   ┲   ┳   ┴   ┵   ┶   ┷   ┸   ┹   ┺   ┻   ┼   ┽   ┾   ┿
           │
    U+254x │ ╀   ╁   ╂   ╃   ╄   ╅   ╆   ╇   ╈   ╉   ╊   ╋   ╌   ╍   ╎   ╏
           │
    U+255x │ ═   ║   ╒   ╓   ╔   ╕   ╖   ╗   ╘   ╙   ╚   ╛   ╜   ╝   ╞   ╟
           │
    U+256x │ ╠   ╡   ╢   ╣   ╤   ╥   ╦   ╧   ╨   ╩   ╪   ╫   ╬   ╭   ╮   ╯
           │
    U+257x │ ╰   ╱   ╲   ╳   ╴   ╵   ╶   ╷   ╸   ╹   ╺   ╻   ╼   ╽   ╾   ╿
           │

Or more compact table:

    % uni p -as table box -compact
             0   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
           ┌────────────────────────────────────────────────────────────────
    U+250x │ ─   ━   │   ┃   ┄   ┅   ┆   ┇   ┈   ┉   ┊   ┋   ┌   ┍   ┎   ┏
    U+251x │ ┐   ┑   ┒   ┓   └   ┕   ┖   ┗   ┘   ┙   ┚   ┛   ├   ┝   ┞   ┟
    U+252x │ ┠   ┡   ┢   ┣   ┤   ┥   ┦   ┧   ┨   ┩   ┪   ┫   ┬   ┭   ┮   ┯
    U+253x │ ┰   ┱   ┲   ┳   ┴   ┵   ┶   ┷   ┸   ┹   ┺   ┻   ┼   ┽   ┾   ┿
    U+254x │ ╀   ╁   ╂   ╃   ╄   ╅   ╆   ╇   ╈   ╉   ╊   ╋   ╌   ╍   ╎   ╏
    U+255x │ ═   ║   ╒   ╓   ╔   ╕   ╖   ╗   ╘   ╙   ╚   ╛   ╜   ╝   ╞   ╟
    U+256x │ ╠   ╡   ╢   ╣   ╤   ╥   ╦   ╧   ╨   ╩   ╪   ╫   ╬   ╭   ╮   ╯
    U+257x │ ╰   ╱   ╲   ╳   ╴   ╵   ╶   ╷   ╸   ╹   ╺   ╻   ╼   ╽   ╾   ╿

### Emoji
The `emoji` command (shortcut: `e`) is is the real reason I wrote this:

    % uni e cry
    	Name                      CLDR
    🥹	face holding back tears  [admiration, aww, cry, embarrassed, feelings, grateful, gratitude, joy, please, proud, resist, sad]
    😢	crying face              [awful, feels, miss, sad, tear, triste, unhappy]
    😭	loudly crying face       [bawling, sad, sob, tear, tears, unhappy]
    😿	crying cat               [animal, face, sad, tear]
    🔮	crystal ball             [fairy, fairytale, fantasy, fortune, future, magic, tale, tool]

By default both the name and CLDR data are searched; the CLDR data is a list of
keywords for an emoji; prefix with `name:` or `n:` to search on the name only:

    % uni e smile
    	Name                              CLDR
    😀	grinning face                    [cheerful, cheery, happy, laugh, nice, smile, smiling, teeth]
    😃	grinning face with big eyes      [awesome, happy, mouth, open, smile, smiling, teeth, yay]
    …

    % uni e name:smile
    	Name                 CLDR
    😼	cat with wry smile  [animal, face, ironic]

As you can see, the CLDR is pretty useful, as "smile" only gives one result as
most emojis use "smiling".

Prefix with `group:` to search by group:

    % uni e group:hands
    	Name                CLDR
    👏	clapping hands     [applause, approval, awesome, congrats, congratulations, excited, good, great, homie, job, nice, prayed, well, yay]
    🙌	raising hands      [celebration, gesture, hooray, praise, raised]
    🫶	heart hands        [<3, love, you]
    👐	open hands         [hug, jazz, swerve]
    🤲	palms up together  [cupped, dua, hands, pray, prayer, wish]
    🤝	handshake          [agreement, deal, meeting]
    🙏	folded hands       [appreciate, ask, beg, blessed, bow, cmon, five, gesture, high, please, pray, thanks, thx]

Group and search can be combined, and `group:` can be abbreviated to `g:`:

    % uni e g:cat-face grin
    	Name                             CLDR
    😺	grinning cat                    [animal, face, mouth, open, smile, smiling]
    😸	grinning cat with smiling eyes  [animal, face, smile]

Like with `search`, use `-or` to OR the parameters together instead of AND:

    % uni e -or g:face-glasses g:face-hat
    	Name                           CLDR
    🤠	cowboy hat face               [cowgirl]
    🥳	partying face                 [bday, birthday, celebrate, celebration, excited, happy, hat, hooray, horn]
    🥸	disguised face                [eyebrow, glasses, incognito, moustache, mustache, nose, person, spy, tache, tash]
    😎	smiling face with sunglasses  [awesome, beach, bright, bro, chilling, cool, rad, relaxed, shades, slay, smile, style, swag, win]
    🤓	nerd face                     [brainy, clever, expert, geek, gifted, glasses, intelligent, smart]
    🧐	face with monocle             [classy, fancy, rich, stuffy, wealthy]

Apply skin tone modifiers with `-tone`:

    % uni e -tone dark g:hands
    	Name                                CLDR
    👏🏿	clapping hands: dark skin tone     [applause, approval, awesome, congrats, congratulations, excited, good, great, homie, job, nice, prayed, well, yay]
    🙌🏿	raising hands: dark skin tone      [celebration, gesture, hooray, praise, raised]
    🫶🏿	heart hands: dark skin tone        [<3, love, you]
    👐🏿	open hands: dark skin tone         [hug, jazz, swerve]
    🤲🏿	palms up together: dark skin tone  [cupped, dua, hands, pray, prayer, wish]
    🤝🏿	handshake: dark skin tone          [agreement, deal, meeting]
    🙏🏿	folded hands: dark skin tone       [appreciate, ask, beg, blessed, bow, cmon, five, gesture, high, please, pray, thanks, thx]

The handshake emoji supports setting individual skin tones per hand since
Unicode 14, but this isn't supported, mostly because I can't really really think
a good CLI interface for setting this without breaking compatibility (there are
some other emojis too, like "holding hands" and "kissing" where you can set both
the gender and skin tone of both sides individually). Maybe for uni v3 someday.

The default is to display only the gender-neutral "person", but this can be
changed with the `-gender` option:

    % uni e -gender man g:person-gesture
    	Name               CLDR
    🙍‍♂️	man frowning      [annoyed, disappointed, disgruntled, disturbed, frustrated, gesture, irritated, person, upset]
    🙎‍♂️	man pouting       [disappointed, downtrodden, frown, grimace, person, scowl, sulk, upset, whine]
    🙅‍♂️	man gesturing NO  [forbidden, gesture, hand, not, person, prohibit]
    🙆‍♂️	man gesturing OK  [exercise, gesture, hand, omg, person]
    💁‍♂️	man tipping hand  [fetch, flick, flip, gossip, person, sarcasm, sarcastic, sassy, seriously, whatever]
    🙋‍♂️	man raising hand  [gesture, here, know, me, person, pick, question, raise]
    🧏‍♂️	deaf man          [accessibility, ear, gesture, hear, person]
    🙇‍♂️	man bowing        [apology, ask, beg, favor, forgive, gesture, meditate, meditation, person, pity, regret, sorry]
    🤦‍♂️	man facepalming   [again, bewilder, disbelief, exasperation, no, not, oh, omg, person, shock, smh]
    🤷‍♂️	man shrugging     [doubt, dunno, guess, idk, ignorance, indifference, knows, maybe, person, whatever, who]

Both `-tone` and `-gender` accept multiple values. `-gender women,man` will
display both the female and male variants, and `-tone light,dark` will display
both a light and dark skin tone; use `all` to display all skin tones or genders:

    % uni e -tone light,dark -gender f,m shrug
    	Name                               CLDR
    🤷🏻‍♂️	man shrugging: light skin tone    [doubt, dunno, guess, idk, ignorance, indifference, knows, maybe, person, whatever, who]
    🤷🏻‍♀️	woman shrugging: light skin tone  [doubt, dunno, guess, idk, ignorance, indifference, knows, maybe, person, whatever, who]
    🤷🏿‍♂️	man shrugging: dark skin tone     [doubt, dunno, guess, idk, ignorance, indifference, knows, maybe, person, whatever, who]
    🤷🏿‍♀️	woman shrugging: dark skin tone   [doubt, dunno, guess, idk, ignorance, indifference, knows, maybe, person, whatever, who]

Like `print` and `identify`, you can use `-format`:

    % uni e g:cat-face -c -format '%(name): %(emoji)'
    grinning cat: 😺
    grinning cat with smiling eyes: 😸
    cat with tears of joy: 😹
    smiling cat with heart-eyes: 😻
    cat with wry smile: 😼
    kissing cat: 😽
    weary cat: 🙀
    crying cat: 😿
    pouting cat: 😾

See `uni help` for more details on the `-format` flag.

### JSON

With `-as json` or `-as j` you can output the data as JSON:

    % uni i -as json h€ý
    [{
    	"aliases": "",
    	"char":    "h",
    	"cpoint":  "U+0068",
    	"dec":     "104",
    	"html":    "&#x68;",
    	"name":    "LATIN SMALL LETTER H",
    	"utf8":    "68"
    }, {
    	"aliases": "",
    	"char":    "€",
    	"cpoint":  "U+20AC",
    	"dec":     "8364",
    	"html":    "&euro;",
    	"name":    "EURO SIGN",
    	"utf8":    "e2 82 ac"
    }, {
    	"aliases": "",
    	"char":    "ý",
    	"cpoint":  "U+00FD",
    	"dec":     "253",
    	"html":    "&yacute;",
    	"name":    "LATIN SMALL LETTER Y WITH ACUTE",
    	"utf8":    "c3 bd"
    }]

All the columns listed in `-f` will be included; you can use `-f all` to include
all columns:

    % uni i -as json -f all h€ý
    [{
    	"aliases": "",
    	"bin":     "1101000",
    	"block":   "Basic Latin",
    	"cat":     "Lowercase_Letter",
    	"cells":   "1",
    	"char":    "h",
    	"cpoint":  "U+0068",
    	"dec":     "104",
    	"digraph": "h",
    	"hex":     "68",
    	"html":    "&#x68;",
    	"json":    "\\u0068",
    	"keysym":  "h",
    	"name":    "LATIN SMALL LETTER H",
    	"oct":     "150",
    	"plane":   "Basic Multilingual Plane",
    	"props":   "",
    	"refs":    "U+04BB, U+210E",
    	"script":  "Latin",
    	"unicode": "1.1",
    	"utf16be": "00 68",
    	"utf16le": "68 00",
    	"utf8":    "68",
    	"width":   "narrow",
    	"xml":     "&#x68;"
    }, {
    	"aliases": "",
    	"bin":     "10000010101100",
    	"block":   "Currency Symbols",
    	"cat":     "Currency_Symbol",
    	"cells":   "1",
    	"char":    "€",
    	"cpoint":  "U+20AC",
    	"dec":     "8364",
    	"digraph": "=e",
    	"hex":     "20ac",
    	"html":    "&euro;",
    	"json":    "\\u20ac",
    	"keysym":  "EuroSign",
    	"name":    "EURO SIGN",
    	"oct":     "20254",
    	"plane":   "Basic Multilingual Plane",
    	"props":   "",
    	"refs":    "U+20A0",
    	"script":  "Common",
    	"unicode": "2.1",
    	"utf16be": "20 ac",
    	"utf16le": "ac 20",
    	"utf8":    "e2 82 ac",
    	"width":   "ambiguous",
    	"xml":     "&#x20ac;"
    }, {
    	"aliases": "",
    	"bin":     "11111101",
    	"block":   "Latin-1 Supplement",
    	"cat":     "Lowercase_Letter",
    	"cells":   "1",
    	"char":    "ý",
    	"cpoint":  "U+00FD",
    	"dec":     "253",
    	"digraph": "y'",
    	"hex":     "fd",
    	"html":    "&yacute;",
    	"json":    "\\u00fd",
    	"keysym":  "yacute",
    	"name":    "LATIN SMALL LETTER Y WITH ACUTE",
    	"oct":     "375",
    	"plane":   "Basic Multilingual Plane",
    	"props":   "",
    	"refs":    "",
    	"script":  "Latin",
    	"unicode": "1.1",
    	"utf16be": "00 fd",
    	"utf16le": "fd 00",
    	"utf8":    "c3 bd",
    	"width":   "neutral",
    	"xml":     "&#xfd;"
    }]

This also works for the `emoji` command:

    % uni e -as json -f all 'kissing cat'
    [{
    	"cldr":      "animal, closed, eye, eyes, face",
    	"cldr_full": "animal, cat, closed, eye, eyes, face, kiss, kissing",
    	"cpoint":    "U+1F63D",
    	"emoji":     "😽",
    	"group":     "Smileys & Emotion",
    	"name":      "kissing cat",
    	"subgroup":  "cat-face"
    }]

All values are always a string, even numerical values. This makes things a bit
easier/consistent as JSON doesn't support hex literals and such. Use `jq` or
some other tool if you want to process the data further.

ChangeLog
---------
Moved to [CHANGELOG.md](/CHANGELOG.md).

Development
-----------
Re-generate the Unicode data with `go generate unidata`. Files are cached in
`unidata/.cache`, so clear that if you want to update the files from remote.
This requires zsh and GNU awk (gawk).
