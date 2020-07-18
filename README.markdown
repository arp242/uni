[![Build Status](https://travis-ci.org/arp242/uni.svg?branch=master)](https://travis-ci.org/arp242/uni)
[![codecov](https://codecov.io/gh/arp242/uni/branch/master/graph/badge.svg)](https://codecov.io/gh/arp242/uni)

`uni` queries the Unicode database from the commandline. It supports Unicode 13
(March 2020) including almost-complete<sup>\*</sup> Emoji support.

There are four commands: `identify` codepoints in a string, `search` for
codepoints, `print` codepoints by class, block, or range, and `emoji` to find
emojis.

There are binaries on the [releases][release] page, or compile from source with
`go get arp242.net/uni`, which will put the binary at `~/go/bin/uni`. You can
also [run it in your browser][uni-wasm].<br>
Packages:
[Arch Linux](https://aur.archlinux.org/packages/uni/) ·
[FreeBSD](https://www.freshports.org/textproc/uni) ·
[Homebrew](https://formulae.brew.sh/formula/uni) ·
[Void Linux](https://github.com/void-linux/void-packages/tree/master/srcpkgs/uni)

Readme Index:
[Integrations](#integrations) ·
[Usage](#usage) ·
[ChangeLog](#changelog) ·
[Development](#development) ·
[Alternatives](#alternatives)

<sup>\*: the part that doesn't work is the complex "family emojis", which can
consist of a bunch of codepoints like "family of man, woman, boy, girl" or "man
kissing man"; supporting that would complicate the CLI quite a lot and is IMO
not worth it.</sup>

[uni-wasm]: https://arp242.github.io/uni-wasm/
[release]: https://github.com/arp242/uni/releases

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

Identify a character:

    $ uni identify €
         cpoint  dec    utf-8       html       name
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN (Currency_Symbol)

Or a string; `i` is a shortcut for `identify`:

    $ uni i h€ý
         cpoint  dec    utf-8       html       name
    'h'  U+0068  104    68          &#x68;     LATIN SMALL LETTER H (Lowercase_Letter)
    '€'  U+20AC  8364   e2 82 ac    &euro;     EURO SIGN (Currency_Symbol)
    'ý'  U+00FD  253    c3 bd       &yacute;   LATIN SMALL LETTER Y WITH ACUTE (Lowercase_Letter)

It reads from stdin:

    $ head -c2 README.markdown | uni i
         cpoint  dec    utf-8       html       name
    '['  U+005B  91     5b          &lsqb;     LEFT SQUARE BRACKET (Open_Punctuation)
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK (Other_Punctuation)

Search description:

    $ uni search euro
         cpoint  dec    utf-8       html       name
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
         cpoint  dec    utf-8       html       name
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA (Other_Symbol)
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS (Other_Symbol)
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA (Other_Symbol)

Use standard shell quoting for more literal matches:

    $ uni s rightwards black arrow
         cpoint  dec    utf-8       html       name
    '➡'  U+27A1  10145  e2 9e a1    &#x27a1;   BLACK RIGHTWARDS ARROW (Other_Symbol)
    '➤'  U+27A4  10148  e2 9e a4    &#x27a4;   BLACK RIGHTWARDS ARROWHEAD (Other_Symbol)
    [..]

    $ uni s 'rightwards black arrow'
         cpoint  dec    utf-8       html       name
    '⮕'  U+2B95  11157  e2 ae 95    &#x2b95;   RIGHTWARDS BLACK ARROW (Other_Symbol)

The `print` command (shortcut `p`) can be used to print specific codepoints or
groups of codepoints:

    $ uni print U+2042
         cpoint  dec    utf-8       html       name
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM (Other_Punctuation)

Print a custom range; `U+2042`, `U2042`, and `2042` are all identical:

    $ uni print 2042..2044
         cpoint  dec    utf-8       html       name
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM (Other_Punctuation)
    '⁃'  U+2043  8259   e2 81 83    &hybull;   HYPHEN BULLET (Other_Punctuation)
    '⁄'  U+2044  8260   e2 81 84    &frasl;    FRACTION SLASH (Math_Symbol)

You can also use hex, octal, and binary numbers: `0x2024`, `0o20102`, or
`0b10000001000010`.

General category:

    $ uni p Po
         cpoint  dec    utf-8       html       name
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK (Other_Punctuation)
    '"'  U+0022  34     22          &quot;     QUOTATION MARK (Other_Punctuation)
    [..]

Blocks:

    $ uni p arrows 'box drawing'
         cpoint  dec    utf-8       html       name
    '←'  U+2190  8592   e2 86 90    &larr;     LEFTWARDS ARROW (Math_Symbol)
    '↑'  U+2191  8593   e2 86 91    &uarr;     UPWARDS ARROW (Math_Symbol)
    [..]
    '─'  U+2500  9472   e2 94 80    &boxh;     BOX DRAWINGS LIGHT HORIZONTAL (Other_Symbol)
    '━'  U+2501  9473   e2 94 81    &#x2501;   BOX DRAWINGS HEAVY HORIZONTAL (Other_Symbol)
    [..]

And finally, there is the `emoji` command (shortcut: `e`), which is the real
reason I wrote this:

    $ uni e cry
    😢 crying face         Smileys & Emotion  face-concerned
    😭 loudly crying face  Smileys & Emotion  face-concerned
    😿 crying cat          Smileys & Emotion  cat-face
    🔮 crystal ball        Activities         game

Filter by group:

    $ uni e group:hands
    🤲 palms up together  People & Body  hands
    🤝 handshake          People & Body  hands
    👏 clapping hands     People & Body  hands
    🙏 folded hands       People & Body  hands
    👐 open hands         People & Body  hands
    🙌 raising hands      People & Body  hands

Group and search can be combined, and `group:` can be abbreviated to `g:`:

    $ uni e g:cat-face grin
    😺 grinning cat                    Smileys & Emotion  cat-face
    😸 grinning cat with smiling eyes  Smileys & Emotion  cat-face

Apply skin tone modifiers with `-tone`:

    $ uni e -tone dark -groups hands
    🤲🏿 palms up together  People & Body  hands
    🤝 handshake          People & Body  hands    [doesn't support skin tone; it's displayed correct]
    👏🏿 clapping hands     People & Body  hands
    🙏🏿 folded hands       People & Body  hands
    👐🏿 open hands         People & Body  hands
    🙌🏿 raising hands      People & Body  hands

The default is to display only the gender-neutral "person", but this can be
changed with the `-gender` option:

    $ uni e -gender man -groups person-gesture
    🙍‍♂️ man frowning      People & Body  person-gesture
    🙎‍♂️ man pouting       People & Body  person-gesture
    🙅‍♂️ man gesturing NO  People & Body  person-gesture
    🙆‍♂️ man gesturing OK  People & Body  person-gesture
    💁‍♂️ man tipping hand  People & Body  person-gesture
    🙋‍♂️ man raising hand  People & Body  person-gesture
    🧏‍♂️ deaf man          People & Body  person-gesture
    🙇‍♂️ man bowing        People & Body  person-gesture
    🤦‍♂️ man facepalming   People & Body  person-gesture
    🤷‍♂️ man shrugging     People & Body  person-gesture

Both `-tone` and `-gender` accept multiple values. `-gender women,man` will
display both the female and male variants (in that order), and `-tone
light,dark` will display both a light and dark skin tone; use `all` to display
all skin tones or genders:

    $ uni e -tone light,dark -gender f,m shrug
    🤷🏻‍♀️ woman shrugging: light skin tone  People & Body  person-gesture
    🤷🏻‍♂️ man shrugging: light skin tone    People & Body  person-gesture
    🤷🏿‍♀️ woman shrugging: dark skin tone   People & Body  person-gesture
    🤷🏿‍♂️ man shrugging: dark skin tone     People & Body  person-gesture


ChangeLog
---------

### master branch

- Replace `uni emoji -group group1,group2` syntax with `uni emoji group:group1
  group:group2`. This is more flexible and will allow adding more query syntax
  later.

- Default for `-gender` is now `-person` instead of `all`.

- Output to `$PAGER` for interactive terminals by default. Use `-no-pager` or
  `-p` to disable.

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

Alternatives
------------

### CLI/TUI

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

- https://github.com/philpennock/character

  More or less similar to uni, but very different CLI.

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
