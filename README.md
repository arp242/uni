`uni` queries the Unicode database from the commandline. It supports Unicode
15.1 (September 2023) and has good support for emojis.

There are four commands: `identify` codepoints in a string, `search` for
codepoints, `print` codepoints by class, block, or range, and `emoji` to find
emojis.

There are binaries on the [releases][release] page, and [packages][pkg] for a
number of platforms. You can also [run it in your browser][uni-wasm].

Compile from source with:

    $ git clone https://github.com/arp242/uni
    $ cd uni
    $ go build

which will give you a `uni` binary.

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
[release]: https://github.com/arp242/uni/releases
[pkg]: https://repology.org/project/uni/versions

Integrations
------------

- [dmenu][dmenu], [rofi][rofi], and [fzf][fzf] script at
  [`dmenu-uni`](/dmenu-uni). See the top of the script for some options you may
  want to frob with.

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

    $ uni identify €
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN (Currency_Symbol)

`i` is a shortcut for `identify`:

    $ uni i h€ý
         CPoint  Dec    UTF8        HTML       Name (Cat)
    'h'  U+0068  104    68          &#x68;     LATIN SMALL LETTER H (Lowercase_Letter)
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN (Currency_Symbol)
    'ý'  U+00FD  253    c3 bd       &yacute;   LATIN SMALL LETTER Y WITH ACUTE (Lowercase_Letter)

It reads from stdin:

    $ head -c2 README.markdown | uni i
         cpoint  dec    utf-8       html       name (cat)
    '['  U+005B  91     5b          &lsqb;     LEFT SQUARE BRACKET (Open_Punctuation)
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK (Other_Punctuation)

You can use `-compact` (or `-c`) to suppress the header, and `-format` (or `-f`)
to control the output format, for example you may want to generate a codepoint
to X11 keysym mapping:

    $ uni i -c -f '0x%(hex): "%(keysym)", // %(name)' h€ý
    0x68: "h", // LATIN SMALL LETTER H
    0x20ac: "EuroSign", // EURO SIGN
    0xfd: "yacute", // LATIN SMALL LETTER Y WITH ACUTE

See `uni help` for more details on the `-format` flag; this flag can also be
added to other commands.

### Search

Search description:

    $ uni search euro
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '₠'  U+20A0  8352   e2 82 a0    &#x20a0;   EURO-CURRENCY SIGN (Currency_Symbol)
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN (Currency_Symbol)
    '𐡷'  U+10877 67703  f0 90 a1 b7 &#x10877;  PALMYRENE LEFT-POINTING FLEURON (Other_Symbol)
    '𐡸'  U+10878 67704  f0 90 a1 b8 &#x10878;  PALMYRENE RIGHT-POINTING FLEURON (Other_Symbol)
    '𐫱'  U+10AF1 68337  f0 90 ab b1 &#x10af1;  MANICHAEAN PUNCTUATION FLEURON (Other_Punctuation)
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA (Other_Symbol)
    '🏤' U+1F3E4 127972 f0 9f 8f a4 &#x1f3e4;  EUROPEAN POST OFFICE (Other_Symbol)
    '🏰' U+1F3F0 127984 f0 9f 8f b0 &#x1f3f0;  EUROPEAN CASTLE (Other_Symbol)
    '💶' U+1F4B6 128182 f0 9f 92 b6 &#x1f4b6;  BANKNOTE WITH EURO SIGN (Other_Symbol)

The `s` command is a shortcut for `search`. Multiple words are matched
individually:

    $ uni s globe earth
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA (Other_Symbol)
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS (Other_Symbol)
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA (Other_Symbol)

Use shell quoting for more literal matches:

    $ uni s rightwards black arrow
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '➡'  U+27A1  10145  e2 9e a1    &#x27a1;   BLACK RIGHTWARDS ARROW (Other_Symbol)
    '➤'  U+27A4  10148  e2 9e a4    &#x27a4;   BLACK RIGHTWARDS ARROWHEAD (Other_Symbol)
    …

    $ uni s 'rightwards black arrow'
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '⮕'  U+2B95  11157  e2 ae 95    &#x2b95;   RIGHTWARDS BLACK ARROW (Other_Symbol)

Add `-or` or `-o` to combine the search terms with "OR" instead of "AND":

    $ uni s -o globe milky
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '🌌' U+1F30C 127756 f0 9f 8c 8c &#x1f30c;  MILKY WAY (Other_Symbol)
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA (Other_Symbol)
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS (Other_Symbol)
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA (Other_Symbol)
    '🌐' U+1F310 127760 f0 9f 8c 90 &#x1f310;  GLOBE WITH MERIDIANS (Other_Symbol)

### Print

Print specific codepoints or groups of codepoints:

    $ uni print U+2042
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM (Other_Punctuation)

Print a custom range; `U+2042`, `U2042`, and `2042` are all identical:

    $ uni print 2042..2044
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM (Other_Punctuation)
    '⁃'  U+2043  8259   e2 81 83    &hybull;   HYPHEN BULLET (Other_Punctuation)
    '⁄'  U+2044  8260   e2 81 84    &frasl;    FRACTION SLASH (Math_Symbol)

You can also use hex, octal, and binary numbers: `0x2024`, `0o20102`, or
`0b10000001000010`.

General category:

    $ uni p Po
    Showing category Po (Other_Punctuation)
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK (Other_Punctuation)
    …

Blocks:

    $ uni p arrows 'box drawing'
    Showing block Arrows
    Showing block Box Drawing
         CPoint  Dec    UTF8        HTML       Name (Cat)
    '←'  U+2190  8592   e2 86 90    &larr;     LEFTWARDS ARROW (Math_Symbol)
    '↑'  U+2191  8593   e2 86 91    &uarr;     UPWARDS ARROW (Math_Symbol)
    …

Print as table, and with a shorter name:

    $ uni p -as table box
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

    $ uni p -as table box -compact
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

    $ uni e cry
    	Name                     (Cldr)
    🥹	face holding back tears  (angry, cry, proud, resist, sad)
    😢	crying face              (sad, tear)
    😭	loudly crying face       (sad, sob, tear)
    😿	crying cat               (face, sad, tear)
    🔮	crystal ball             (fairy tale, fantasy, fortune, tool)

By default both the name and CLDR data are searched; the CLDR data is a list of
keywords for an emoji; prefix with `name:` or `n:` to search on the name only:

    $ uni e smile
    	Name                             (Cldr)
    😃	grinning face with big eyes      (mouth, open, smile)
    😄	grinning face with smiling eyes  (mouth, open, smile)
    …

    $ uni e name:smile
    	Name                (Cldr)
    😼	cat with wry smile  (face, ironic)

As you can see, the CLDR is pretty useful, as "smile" only gives one result as
most emojis use "smiling".

Prefix with `group:` to search by group:

    $ uni e group:hands
    	Name               (Cldr)
    👏	clapping hands     ()
    🙌	raising hands      (celebration, gesture, hooray, raised)
    🫶	heart hands        (love)
    👐	open hands         ()
    🤲	palms up together  (prayer)
    🤝	handshake          (agreement, meeting)
    🙏	folded hands       (ask, high 5, high five, please, pray, thanks)

Group and search can be combined, and `group:` can be abbreviated to `g:`:

    $ uni e g:cat-face grin
    	Name                            (Cldr)
    😺	grinning cat                    (face, mouth, open, smile)
    😸	grinning cat with smiling eyes  (face, smile)

Like with `search`, use `-or` to OR the parameters together instead of AND:

    $ uni e -or g:face-glasses g:face-hat
    	Name                          (Cldr)
    🤠	cowboy hat face               (cowgirl)
    🥳	partying face                 (celebration, hat, horn)
    🥸	disguised face                (glasses, incognito, nose)
    😎	smiling face with sunglasses  (bright, cool)
    🤓	nerd face                     (geek)
    🧐	face with monocle             (stuffy)

Apply skin tone modifiers with `-tone`:

    $ uni e -tone dark g:hands
    	Name                               (Cldr)
    👏🏿	clapping hands: dark skin tone     ()
    🙌🏿	raising hands: dark skin tone      (celebration, gesture, hooray, raised)
    🫶🏿	heart hands: dark skin tone        (love)
    👐🏿	open hands: dark skin tone         ()
    🤲🏿	palms up together: dark skin tone  (prayer)
    🤝	handshake                          (agreement, meeting)
    🙏🏿	folded hands: dark skin tone       (ask, high 5, high five, please, pray, thanks)

The "heart hands" may not show as it's very recent. The handshake emoji supports
setting individual skin tones per hand since Unicode 14, but this isn't
supported, mostly because I can't really really think a good CLI interface for
setting this without breaking compatibility (there are some other emojis too,
like "holding hands" and "kissing" where you can set both the gender and skin
tone of both sides individually). Maybe for uni v3 someday.

The default is to display only the gender-neutral "person", but this can be
changed with the `-gender` option:

    $ uni e -gender man g:person-gesture
    	Name              (Cldr)
    🙍‍♂️	man frowning      (gesture, person frowning)
    🙎‍♂️	man pouting       (gesture, person pouting)
    🙅‍♂️	man gesturing NO  (forbidden, gesture, hand, person gesturing NO, prohibited)
    🙆‍♂️	man gesturing OK  (gesture, hand, person gesturing OK)
    💁‍♂️	man tipping hand  (help, information, person tipping hand, sassy)
    🙋‍♂️	man raising hand  (gesture, happy, person raising hand, raised)
    🧏‍♂️	deaf man          (accessibility, deaf person, ear, hear)
    🙇‍♂️	man bowing        (apology, gesture, person bowing, sorry)
    🤦‍♂️	man facepalming   (disbelief, exasperation, person facepalming)
    🤷‍♂️	man shrugging     (doubt, ignorance, indifference, person shrugging)

Both `-tone` and `-gender` accept multiple values. `-gender women,man` will
display both the female and male variants, and `-tone light,dark` will display
both a light and dark skin tone; use `all` to display all skin tones or genders:

    $ uni e -tone light,dark -gender f,m shrug
    	Name                              (Cldr)
    🤷🏻‍♂️	man shrugging: light skin tone    (doubt, ignorance, indifference, person shrugging)
    🤷🏻‍♀️	woman shrugging: light skin tone  (doubt, ignorance, indifference, person shrugging)
    🤷🏿‍♂️	man shrugging: dark skin tone     (doubt, ignorance, indifference, person shrugging)
    🤷🏿‍♀️	woman shrugging: dark skin tone   (doubt, ignorance, indifference, person shrugging)

Like `print` and `identify`, you can use `-format`:

    $ uni e g:cat-face -c -format '%(name): %(emoji)'
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

    $ uni i -as json h€ý
    [{
    	"cat": "Lowercase_Letter",
    	"char": "h",
    	"cpoint": "U+0068",
    	"dec": "104",
    	"html": "&#x68;",
    	"name": "LATIN SMALL LETTER H",
    	"utf8": "68"
    }, {
    	"cat": "Currency_Symbol",
    	"char": "€",
    	"cpoint": "U+20AC",
    	"dec": "8364",
    	"html": "&euro;",
    	"name": "EURO SIGN",
    	"utf8": "e2 82 ac"
    }, {
    	"cat": "Lowercase_Letter",
    	"char": "ý",
    	"cpoint": "U+00FD",
    	"dec": "253",
    	"html": "&yacute;",
    	"name": "LATIN SMALL LETTER Y WITH ACUTE",
    	"utf8": "c3 bd"
    }]

All the columns listed in `-f` will be included; you can use `-f all` to include
all columns:

    $ uni i -as json -f all h€ý
    [{
    	"bin": "1101000",
    	"block": "Basic Latin",
    	"cat": "Lowercase_Letter",
    	"char": "h",
    	"cpoint": "U+0068",
    	"dec": "104",
    	"digraph": "h",
    	"hex": "68",
    	"html": "&#x68;",
    	"json": "\\u0068",
    	"keysym": "h",
    	"name": "LATIN SMALL LETTER H",
    	"oct": "150",
    	"plane": "Basic Multilingual Plane",
    	"props": "",
    	"script": "Latin",
    	"unicode": "1.1",
    	"utf16be": "00 68",
    	"utf16le": "68 00",
    	"utf8": "68",
    	"width": "neutral",
    	"xml": "&#x68;"
    }, {
    	"bin": "10000010101100",
    	"block": "Currency Symbols",
    	"cat": "Currency_Symbol",
    	"char": "€",
    	"cpoint": "U+20AC",
    	"dec": "8364",
    	"digraph": "=e",
    	"hex": "20ac",
    	"html": "&euro;",
    	"json": "\\u20ac",
    	"keysym": "EuroSign",
    	"name": "EURO SIGN",
    	"oct": "20254",
    	"plane": "Basic Multilingual Plane",
    	"props": "",
    	"script": "Common",
    	"unicode": "2.1",
    	"utf16be": "20 ac",
    	"utf16le": "ac 20",
    	"utf8": "e2 82 ac",
    	"width": "ambiguous",
    	"xml": "&#x20ac;"
    }, {
    	"bin": "11111101",
    	"block": "Latin-1 Supplement",
    	"cat": "Lowercase_Letter",
    	"char": "ý",
    	"cpoint": "U+00FD",
    	"dec": "253",
    	"digraph": "y'",
    	"hex": "fd",
    	"html": "&yacute;",
    	"json": "\\u00fd",
    	"keysym": "yacute",
    	"name": "LATIN SMALL LETTER Y WITH ACUTE",
    	"oct": "375",
    	"plane": "Basic Multilingual Plane",
    	"props": "",
    	"script": "Latin",
    	"unicode": "1.1",
    	"utf16be": "00 fd",
    	"utf16le": "fd 00",
    	"utf8": "c3 bd",
    	"width": "narrow",
    	"xml": "&#xfd;"
    }]

This also works for the `emoji` command:

    $ uni e -as json -f all 'kissing cat'
    [{
    	"cldr": "eye, face",
    	"cldr_full": "cat, eye, face, kiss, kissing cat",
    	"cpoint": "U+1F63D",
    	"emoji": "😽",
    	"group": "Smileys & Emotion",
    	"name": "kissing cat",
    	"subgroup": "cat-face"
    }]

All values are always a string, even numerical values. This makes things a bit
easier/consistent as JSON doesn't support hex literals and such. Use `jq` or
some other tool if you want to process the data further.


ChangeLog
---------

### v2.6.0 (2023-11-24)

- Update to Unicode 15.1.

- Add "script" property – also supported in the list and print commands:

      % uni identify -f '%(script l:auto) %(cpoint) %(name)' 'a Ω'
      Script CPoint Name
      Latin  U+0061 LATIN SMALL LETTER A
      Common U+0020 SPACE
      Greek  U+03A9 GREEK CAPITAL LETTER OMEGA

      % uni list scripts
      Scripts:
      Name                    Assigned
      Adlam                         83
      Ahom                          54
      Anatolian Hieroglyphs        582
      …
      
      % uni print 'script:linear a'
      Showing script Linear A
           CPoint  Dec    UTF8        HTML       Name (Cat)
      '𐘀'  U+10600 67072  f0 90 98 80 &#x10600;  LINEAR A SIGN AB001 (Other_Letter)
      '𐘁'  U+10601 67073  f0 90 98 81 &#x10601;  LINEAR A SIGN AB002 (Other_Letter)
      '𐘂'  U+10602 67074  f0 90 98 82 &#x10602;  LINEAR A SIGN AB003 (Other_Letter)
      …


- Add "unicode" property, which tells you in which Unicode version a codepoint
  was introduced:

      % uni identify -f '%(unicode l:auto) %(cpoint l:auto) %(name)' a𐘂🫁
      Unicode CPoint  Name
      1.1     U+0061  LATIN SMALL LETTER A
      7.0     U+10602 LINEAR A SIGN AB003
      13.0    U+1FAC1 LUNGS

- Show unprintable control characters as the open box (␣, U+2423) instead of the
  replacement character (�, U+FFFD). It already did that for C1 control
  characters, and U+FFFD looked more like a bug than intentional. The -raw/-r
  flag still overrides this.

- Always print Private Use characters as-is for %(char) instead of using U+FFFD
  replacement character. It's usually safe to print this, and having to use -raw
  is confusing.

- `ls` command is now an alias for `list.

### 2.5.1 (2022-05-09)

- Fix build on Go 1.17 and earlier.

### 2.5.0 (2022-05-03)

- Add support for properties; they can be displayed with `%(props)` in
  `-format`, and selected in `print` (e.g. `uni print dash`).

- Add `uni list` command, to list categories, blocks, and properties.

- Allow explicitly selecting a block, category, or property in `print` with
  `block:name` (`b:name`), `category:name` (`cat:name`, `c:name`), or
  `property:name` (`prop:name`, `p:name`).

  Also print an error if a string without prefix matched more than one group
  (i.e. `uni p dash` matches both the property `Dash` and category
  `Dash_Punctuation`).

- Add table layout with `-as table`. Also change `-json`/`-j` to `-as json` or
  `-as j`. The `-json` flag is still accepted as an alias for compatibility.

- Change `-q`/`-quiet` to `-c`/`-compact`; `-as json` will print as minified if
  given, and `-as table` will include less padding. `-q` is still accepted as an
  alias for compatibility.

- Don't use the Go stdlib `unicode` package; since this is a Unicode 13 database
  and some operations would fail on codepoints added in Unicode 14 due to the
  mismatch.

### v2.4.0 (2021-12-20)

- Update import path to `zgo.at/uni/v2`.

- Add `oct` and `bin` flags for `-f` to print a codepoint as octal or binary.

- Add `f` format flag to change the fill character with alignment; e.g.
  `%(bin r:auto f:0)` will print zeros on the left.

- Allow using just `o123` for an octal number (instead of `0o123`). We can't do
  this for binary and decimal numbers (since `b` and `d` are valid
  hexidecimals), but no reason not to do it for `o`.

### v2.3.0 (2021-10-05)

- Update to Unicode 14.0.

- UTF-16 and JSON are printed as lower case, just like UTF-8 was. Upper-case is
  used only for codepoints (i.e. U+00AC).

- `uni print` can now print from UTF-8 byte sequence; for example to print the €
  sign:

      uni p utf8:e282ac
      uni p 'utf8:e2 82 ac'
      uni p 'utf8:0xe2 0x82 0xac'

  Bytes can optionally be separated by any combination of `0x`, `-`, `_`, or spaces.

### v2.2.1 (2021-06-15)

- You can now use `uni p 0d40` to get U+28 by decimal.

  `uni print 40` interprets the `40` as hex instead of decimal, and there was no
  way to get a codepoint by decimal number. Since codepoints are much more more
  common than decimals, leaving off the `U+` and `U` is a useful shortcut I'd
  like to keep. AFAIK there isn't really a standard(-ish) was to explicitly
  indicate a number is a decimal, so this is probably the closest.

### v2.2.0 (2021-06-05)

- Make proper use of the `/v2` import path so that `go get` and `go install`
  work. (#26)

- Don't panic if `-f` doesn't contain any formatting characters.

### v2.1.0 (2021-03-30)

- Can now output as JSON with `-j` or `-json`.

- `-format all` is a special value to include all columns uni knows about. This
  is useful especially in combination with `-json`.

- Add `%(block)`, `%(plane)`, `%(width)`, `%(utf16be)`, `%(utf16le)`, and
  `%(json) to `-f`.

- Refactor the arp242.net/uni/unidata package to be more useful for other use
  cases. This isn't really relevant for `uni` users as such, but if you want to
  get information about codepoints or emojis then this package is a nice
  addition to the standard library's `unicode` package.

### v2.0.0 (2021-01-03)

This changes some flags, semantics, and defaults in **incompatible** ways, hence
the bump to 2.0. If you use the `dmenu-uni` script with dmenu or fzf, then
you'll need to update that to.

- Remove the `-group` flag in favour of `group:name` syntax; this is more
  flexible and will allow adding more query syntax later.

      uni emoji -group groupname,othergroup                  Old syntax
      uni emoji -group groupname,othergroup smile            Old syntax

      uni emoji -or group:groupname group:othergroup         New syntax
      uni emoji -or group:groupname group:othergroup smile   New syntax

      uni emoji -or g:groupname g:othergroup                 Can use shorter g: instead of group:

- Default for `-gender` is now `person` instead of `all`; including all genders
  by default isn't all that useful, and the gender-neutral "person" should be a
  fine default for most, just as the skin colour-neutral "yellow" is probably a
  fine default for most.

- Add new `-or`/`-o` flag. The default for `search` and `emoji` is to show
  everything where all query parameters match ("AND"); with this flag it shows
  everything where at least one parameter matches ("OR").

- Add new `-format`/`-f` flag to control which columns to output and column
  width. You can now also print X11 keysyms and Vim digraphs. See `uni help` for
  details.

- Include CLDR data for emojis, which is searched by default if you use `uni e
  <something>`. You can use `uni e name:x` to search for the name specifically.

- Show a short terse help when using just `uni`, and a more detailed help on
  `uni help`. I hate it when programs print 5 pages of text to my terminal when
  I didn't ask for it.

- Update Unicode data to 13.1.

- Add option to output to `$PAGER` with `-p` or `-pager`. This isn't done
  automatically (I don't really like it when programs throw me in a pager), but
  you can define a shell alias (`alias uni='uni -p'`) if you want it by default
  since flags can be both before or after the command.

### v1.1.1 (2020-05-31)

- Fix tests of v1.1.0, requested by a packager. No changes other than this.

### v1.1.0 (2020-03-17)

- Update Unicode data from 12.1 to 13.0.

- `print` command supports codepoints as hex (`0xff`), octal (`0o42`), and
  binary (`0b1001`).

- A few very small bugfixes.

### v1.0.0 (2019-12-12)

- Initial release

Development
-----------

Re-generate the Unicode data with `go generate unidata`. Files are cached in
`unidata/.cache`, so clear that if you want to update the files from remote.
This requires zsh and GNU awk (gawk).

Alternatives
------------

### CLI/TUI

- https://github.com/philpennock/character

  More or less similar to uni, but very different CLI, and has some additional
  features. Seems pretty good.

- https://github.com/sindresorhus/emoj

  Doesn't support emojis sequences (e.g. MAN SHRUGGING is PERSON SHRUGGING +
  MAN, FIREFIGHTER is PERSON + FIRE TRUCK, etc); quite slow for a CLI program
  (`emoj smiling` takes 1.8s on my system, sometimes a lot longer), search
  results are pretty bad (`shrug` returns unamused face, thinking face, eyes,
  confused face, neutral face, tears of joy, and expressionless face ... but not
  the shrugging emoji), not a fan of npm (has 1862 dependencies).

- https://github.com/Fingel/tuimoji

  Grouping could be better, doesn't support emojis sequences, only interactive
  TUI, feels kinda slow-ish especially when searching.

- https://github.com/pemistahl/chr

  Only deals with codepoints, not emojis.

### GUI

- gnome-characters

  Uses Gnome interface/window decorations and won't work well with other WMs,
  doesn't deal with emoji sequences, I don't like the grouping/ordering it uses,
  requires two clicks to copy a character.

- gucharmap

  Doesn't display emojis, just unicode blocks.

- KCharSelect

  Many KDE-specific dependencies (106M). Didn't try it.

- https://github.com/Mange/rofi-emoji and https://github.com/fdw/rofimoji

  Both are pretty similar to the dmenu/rofi integration of uni with some minor
  differences, and both seem to work well with no major issues.

- gtk3 emoji picker (Ctrl+; or Ctrl+. in gtk 3.93 or newer)

  Only works in GTK, doesn't work with `GTK_IM_MODULE=xim` (needed for compose
  key), for some reasons the emojis look ugly, doesn't display emojis sequences,
  doesn't have a tooltip or other text description about what the emoji actually
  is, the variation selector doesn't seem to work (never displays skin tone?),
  doesn't work in Firefox.

  This is so broken on my system that it seems that I'm missing something for
  this to work or something?

- https://github.com/rugk/awesome-emoji-picker

  Only works in Firefox; takes a tad too long to open; doesn't support skin
  tones.

### Didn't investigate (yet)

Some alternatives people have suggested that I haven't looked at; make an issue
or email me if you know of any others.

- https://github.com/cassidyjames/ideogram
- https://github.com/OzymandiasTheGreat/emoji-keyboard
- https://github.com/salty-horse/ibus-uniemoji
- https://fcitx-im.org/wiki/Unicode
- http://kassiopeia.juls.savba.sk/~garabik/software/unicode/ and https://github.com/garabik/unicode (same?)
- https://billposer.org/Software/unidesc.html
- https://github.com/NoraCodes/charpicker (rofi)
