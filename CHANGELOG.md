### Unreleased

- Update to Unicode 17.0.

- Fix sorting of `print` and `search` with a custom `-format` flag which didn't
  include `%(dec)`.

- Sometimes the would be printed on the second line instead of the first when
  printing everything or control characters.

- Recognize `c:[cat-name]` to print a category.

### v2.8.0 (2024-09-11)

- Update to Unicode 16.0.

### v2.7.0 (2024-05-22)

- Improve `-format` flag:

  - Add `%name` as an alias for `%(name l:auto)`; this is a lot less typing and
    requires less shell quoting, and >90% of the time this is what you want.

  - Automatically prepend character, codepoint, and name if the format flag
    starts with `+`; for example:

        % uni identify -f +'%unicode %plane' a
                     Name                 Unicode Plane
        'a'  U+0061  LATIN SMALL LETTER A 1.1     Basic Multilingual Plane

  This should make quickly printing some property a lot quicker.

- Align and colourize JSON output.

- Update CLDR information, adding significantly more aliases for emojis.

- Add `cells` column, which returns how many cells a codepoint will display at
  (0, 1, or 2).

- Add `aliases` column, which lists the alias names. Also add this to the
  default output:

      % uni s factorial
           CPoint  Dec    UTF8        HTML       Name  Aliases
      '!'  U+0021  33     21          &excl;     EXCLAMATION MARK [factorial, bang]

- Add `refs` columns, which references other related/similar codepoints:

      % uni p -q U+46 -f '%(name): %(refs)'
      LATIN CAPITAL LETTER F: U+2109, U+2131, U+2132

      % uni p -q U+46 -f '%(refs)' | uni p
           CPoint  Dec    UTF8        HTML       Name  Aliases
      '℉'  U+2109  8457   e2 84 89    &#x2109;   DEGREE FAHRENHEIT
      'ℱ'  U+2131  8497   e2 84 b1    &Fscr;     SCRIPT CAPITAL F [Fourier transform]
      'Ⅎ'  U+2132  8498   e2 84 b2    &#x2132;   TURNED CAPITAL F [Claudian digamma inversum]

- Allow arguments to `print`to start or end with a comma or slash. This comes up
  when copy/pasting some list of codepoints from another source; there's no real
  reason to error out on this.

- Allow listing unicode versions with `uni list unicode` and planes with `uni
  list planes`.

- `uni list` without arguments errors, instead of listing all.

- Add `h` format flag to not print the header for this column.

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
