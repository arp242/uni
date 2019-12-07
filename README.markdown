`uni` queries the Unicode database from the commandline.

There are four commands: `identify` to print Unicode information about a string,
`search` to search for codepoints, `print` to print groups of Unicode classes,
and `emoji` to find emojis.

Install it with `go get arp242.net/uni`, which will put the binary at
`~/go/bin/uni`.

Integrations
------------

There is a [dmenu][dmenu] example script at [`dmenu-uni`](dmenu-uni). This can
also be used with [rofi][rofi] or similar programs. Note that dmenu will *crash*
when using a colour emoji font (such as Noto), this is [a bug in Xft][xft].

You can add a `:UnicodeName` command to Vim with:

    command! UnicodeName echo
            \ system('uni -q i',
            \      [strcharpart(strpart(getline('.'), col('.') - 1), 0, 1)]
            \ )[:-2]

[dmenu]: http://tools.suckless.org/dmenu
[rofi]: https://github.com/davatorium/rofi
[xft]: https://bugs.freedesktop.org/show_bug.cgi?id=107534

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

It reads from stdin:

    $ head -c5 README.markdown | uni i
    uni: reading from stdin...
         cpoint  dec    utf-8       html       name
    '`'  U+0060  96     60          &grave;    GRAVE ACCENT (Modifier_Symbol)
    'u'  U+0075  117    75          &#x75;     LATIN SMALL LETTER U (Lowercase_Letter)
    'n'  U+006E  110    6e          &#x6e;     LATIN SMALL LETTER N (Lowercase_Letter)
    'i'  U+0069  105    69          &#x69;     LATIN SMALL LETTER I (Lowercase_Letter)
    '`'  U+0060  96     60          &grave;    GRAVE ACCENT (Modifier_Symbol)

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

    $ uni s earth globe
         cpoint  dec    utf-8       html       name
    '🌍' U+1F30D 127757 f0 9f 8c 8d &#x1f30d;  EARTH GLOBE EUROPE-AFRICA (Other_Symbol)
    '🌎' U+1F30E 127758 f0 9f 8c 8e &#x1f30e;  EARTH GLOBE AMERICAS (Other_Symbol)
    '🌏' U+1F30F 127759 f0 9f 8c 8f &#x1f30f;  EARTH GLOBE ASIA-AUSTRALIA (Other_Symbol)

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

The `print` command (shortcut `p`) can be used to print specific codepoints or
groups of codepoints:

    $ uni print U+2042
         cpoint  dec    utf-8       html       name
    '⁂'  U+2042  8258   e2 81 82    &#x2042;   ASTERISM (Other_Punctuation)

General category:

    $ uni p Po
         cpoint  dec    utf-8       html       name
    '!'  U+0021  33     21          &excl;     EXCLAMATION MARK (Other_Punctuation)
    '"'  U+0022  34     22          &quot;     QUOTATION MARK (Other_Punctuation)
    '#'  U+0023  35     23          &num;      NUMBER SIGN (Other_Punctuation)
    [..]

Blocks:

    $ uni p arrows 'box drawing'
         cpoint  dec    utf-8       html       name
    '←'  U+2190  8592   e2 86 90    &larr;     LEFTWARDS ARROW (Math_Symbol)
    '↑'  U+2191  8593   e2 86 91    &uarr;     UPWARDS ARROW (Math_Symbol)
    '→'  U+2192  8594   e2 86 92    &rarr;     RIGHTWARDS ARROW (Math_Symbol)
    '↓'  U+2193  8595   e2 86 93    &darr;     DOWNWARDS ARROW (Math_Symbol)
    [..]
    '─'  U+2500  9472   e2 94 80    &boxh;     BOX DRAWINGS LIGHT HORIZONTAL (Other_Symbol)
    '━'  U+2501  9473   e2 94 81    &#x2501;   BOX DRAWINGS HEAVY HORIZONTAL (Other_Symbol)
    '│'  U+2502  9474   e2 94 82    &boxv;     BOX DRAWINGS LIGHT VERTICAL (Other_Symbol)
    '┃'  U+2503  9475   e2 94 83    &#x2503;   BOX DRAWINGS HEAVY VERTICAL (Other_Symbol)
    [..]

And finally, there is the `emoji` command (shortcut: `e`), which is the real
reason I wrote this:

    $ uni e firefighter
    👨‍🚒 man firefighter
    👨‍🏻‍🚒 man firefighter: light skin tone
    👨‍🏼‍🚒 man firefighter: medium-light skin tone
    👨‍🏽‍🚒 man firefighter: medium skin tone
    👨‍🏾‍🚒 man firefighter: medium-dark skin tone
    👨‍🏿‍🚒 man firefighter: dark skin tone
    👩‍🚒 woman firefighter
    👩‍🏻‍🚒 woman firefighter: light skin tone
    👩‍🏼‍🚒 woman firefighter: medium-light skin tone
    👩‍🏽‍🚒 woman firefighter: medium skin tone
    👩‍🏾‍🚒 woman firefighter: medium-dark skin tone
    👩‍🏿‍🚒 woman firefighter: dark skin tone
