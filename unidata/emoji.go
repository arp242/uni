package unidata

import (
	"strings"
)

// Emoji is an emoji sequence.
type (
	Emoji struct {
		Codepoints []rune        // Codepoints in this emoji
		Name       string        // Emoji name
		group      EmojiGroup    // Emoji group
		subgroup   EmojiSubgroup // Emoji subgroup
		CLDR       []string      // CLDR names
		skinTones  bool          // Supports skintones?
		gender     int           // Supports setting gender?
	}
	EmojiGroup    uint8  // Emoji group.
	EmojiSubgroup uint16 // Emoji subgroup.

	// EmojiGenderType   uint8
	// EmojiSkintoneType uint8
)

func (e EmojiGroup) String() string    { return EmojiGroups[e].Name }
func (e EmojiSubgroup) String() string { return EmojiSubgroups[e].Name }

func (e Emoji) Group() EmojiGroup       { return e.group }
func (e Emoji) Subgroup() EmojiSubgroup { return e.subgroup }
func (e Emoji) Skintones() bool         { return e.skinTones }
func (e Emoji) Genders() bool           { return e.gender > 0 }

func (e Emoji) String() string {
	if len(e.Codepoints) == 0 { // Should never happen.
		return ""
	}

	var c string

	// Flags
	// 1F1FF 1F1FC                                 # 🇿🇼 E2.0 flag: Zimbabwe
	// 1F3F4 E0067 E0062 E0065 E006E E0067 E007F   # 🏴󠁧󠁢󠁥󠁮󠁧󠁿 E5.0 flag: England
	if (e.Codepoints[0] >= 0x1f1e6 && e.Codepoints[0] <= 0x1f1ff) ||
		(len(e.Codepoints) > 1 && e.Codepoints[1] == 0xe0067) {
		for _, cp := range e.Codepoints {
			c += string(rune(cp))
		}
		return c
	}

	for i, cp := range e.Codepoints {
		c += string(rune(cp))

		// Don't add ZWJ as last item.
		if i == len(e.Codepoints)-1 {
			continue
		}

		switch e.Codepoints[i+1] {
		// Never add ZWJ before variation selector or skin tone.
		case 0xfe0f, 0x1f3fb, 0x1f3fc, 0x1f3fd, 0x1f3fe, 0x1f3ff:
			continue
		// Keycap: join with 0xfe0f
		case 0x20e3:
			continue
		}

		c += "\u200d"
	}
	return c
}

// Emoji genders types
const (
	genderNone = 0
	genderSign = 1
	genderRole = 2
)

// EmojiModifier is a modifier to apply to an emoji to change the gender(s) or
// skintone(s).
type EmojiModifier uint16

// EmojiModifier values.
const (
	ModPerson      = EmojiModifier(1 << iota) // Gender-neutral "person".
	ModMale                                   // Male
	ModFemale                                 // Female
	ModNone                                   // No skin tone
	ModLight                                  // Light skin tone
	ModMediumLight                            // Medium light skin tone
	ModMedium                                 // Medium skin tone
	ModMediumDark                             // Mediun dark skin tone
	ModDark                                   // Dark skin tone
)

func isEmoji(e Emoji, want ...rune) bool {
	if len(e.Codepoints) != len(want) {
		return false
	}
	for i := range e.Codepoints {
		if e.Codepoints[i] != want[i] {
			return false
		}
	}
	return true
}

// With returns a copy of this emoji with the given modifiers.
func (e Emoji) With(mod EmojiModifier, selmod ...EmojiModifier) Emoji {
	// Make explicit copy of the codepoints; as this is a slice/pointer and we
	// don't want to modify the original.
	orig := e.Codepoints
	e.Codepoints = make([]rune, len(orig))
	copy(e.Codepoints, orig)

	// Handshake supports setting the skintone individually for the left and
	// right side:
	//
	//   🤝
	//   1F91D                          handshake
	//   🤝    🏻
	//   1F91D 1F3FB                    handshake: light skin tone
	//   🫱     🏼         🫲     🏽
	//   1FAF1 1F3FC 200D 1FAF2 1F3FD   handshake: medium-light skin tone, medium skin tone
	if len(selmod) > 0 && isEmoji(e, 0x1F91D) {
		e.Codepoints = []rune{
			0x1FAF1, tonemap[mod&^(ModPerson|ModMale|ModFemale)],
			0x200D,
			0x1FAF2, tonemap[selmod[0]&^(ModPerson|ModMale|ModFemale)],
		}
		return e
	}

	// The holding hands, kissing, and couple with heart supports skintone and
	// gender for left & right side; this works in a bit of an odd way:
	//
	//   👫
	//   1F46B         woman and man holding hands
	//   👬
	//   1F46C         men holding hands
	//   👭
	//   1F46D         women holding hands
	//   👬    🏻
	//   1F46C 1F3FB   men holding hands: light skin tone
	//
	// But to set the skintone or gender invididually (or use gender-neutral
	// people) expand the 1F46{B,C,D}:
	//
	//   🧑         🤝         🧑
	//   1F9D1 200D 1F91D 200D 1F9D1                  people holding hands
	//   👨    🏿         🤝         👨    🏽
	//   1F468 1F3FF 200D 1F91D 200D 1F468 1F3FD      men holding hands: dark skin tone, medium skin tone
	//   👩    🏿         🤝         👨    🏻
	//   1F469 1F3FF 200D 1F91D 200D 1F468 1F3FB      woman and man holding hands: dark skin tone, light skin tone
	if isEmoji(e, 0x1F46B) || isEmoji(e, 0x1F46C) || isEmoji(e, 0x1F46D) {
		// TODO
		switch {
		// Special case to expand to "person"
		case len(selmod) == 0 && mod&ModPerson != 0:
		case len(selmod) > 1:
		}
	}

	// For kissing it's similar:
	//
	//   💏
	//   1F48F            kiss
	//   💏    🏻
	//   1F48F 1F3FB      kiss: light skin tone
	//
	// Expands to these codepoint poems if you want to set a gender or skin tone
	// individually:
	//
	//   👨    🏻         ❤️              💋         👨    🏼
	//   1F468 1F3FB 200D 2764 FE0F 200D 1F48B 200D 1F468 1F3FC    kiss: man, man, light skin tone, medium-light skin tone
	//   👩         ❤️              💋         👨
	//   1F469 200D 2764 FE0F 200D 1F48B 200D 1F468                kiss: woman, man
	//   👩    🏼         ❤️              💋         👨    🏽
	//   1F469 1F3FC 200D 2764 FE0F 200D 1F48B 200D 1F468 1F3FD    kiss: woman, man, medium-light skin tone, medium skin tone
	//if isEmoji(e, 0x1F48F) {
	//	// TODO
	//}

	// And finally with a heart:
	//
	//   💑
	//   1F491                                             couple with heart
	//   💑    🏻
	//   1F491 1F3FB                                       couple with heart: light skin tone
	//   🧑    🏾         ❤              🧑    🏻
	//   1F9D1 1F3FE 200D 2764 FE0F 200D 1F9D1 1F3FB       couple with heart: person, person, medium-dark skin tone, light skin tone
	//if isEmoji(e, 0x1F491) {
	//	// TODO
	//}

	// "Family" supports settings the four family members' gender (no skintone
	// support):
	//
	//   👪
	//   1F46A                                     family
	//   👨         👩         👦
	//   1F468 200D 1F469 200D 1F466               family: man, woman, boy
	//   👩         👩         👧         👦
	//   1F469 200D 1F469 200D 1F467 200D 1F466    family: woman, woman, girl, boy
	//   👨         👧
	//   1F468 200D 1F467                          family: man, girl
	//   👨         👧         👦
	//   1F468 200D 1F467 200D 1F466               family: man, girl, boy
	//if isEmoji(e, 0x1F46A) {
	//	// TODO
	//}

	e = e.applyGender(mod & (ModPerson | ModMale | ModFemale))
	e = e.applyTone(mod &^ (ModPerson | ModMale | ModFemale))
	return e
}

func (e Emoji) applyGender(g EmojiModifier) Emoji {
	switch {
	// Append male or female sign
	//   1F937 1F3FD                   # 🤷🏽 E4.0 person shrugging: medium skin tone
	//   1F937 1F3FB 200D 2642 FE0F    # 🤷🏻‍♂️ E4.0 man shrugging: light skin tone
	case e.gender == genderSign:
		switch g {
		case ModMale:
			e.Name = strings.Replace(e.Name, "person", "man", 1)
			e.Codepoints = append(e.Codepoints, []rune{0x2642, 0xfe0f}...)
		case ModFemale:
			e.Name = strings.Replace(e.Name, "person", "woman", 1)
			e.Codepoints = append(e.Codepoints, []rune{0x2640, 0xfe0f}...)
		}
	// Replace first "person" with "man" or "woman".
	//   1F9D1 200D 1F692              # 🧑‍🚒 E12.1 firefighter
	//   1F9D1 1F3FB 200D 1F692        # 🧑🏻‍🚒 E12.1 firefighter: light skin tone
	//   1F469 200D 1F692              # 👩‍🚒 E4.0 woman firefighter
	//   1F469 1F3FB 200D 1F692        # 👩🏻‍🚒 E4.0 woman firefighter: light skin tone
	case e.gender == genderRole:
		switch g {
		case ModMale:
			e.Name = "man " + e.Name
			e.Codepoints = append([]rune{0x1f468}, e.Codepoints[1:]...)
		case ModFemale:
			e.Name = "woman " + e.Name
			e.Codepoints = append([]rune{0x1f469}, e.Codepoints[1:]...)
		}
	}
	return e
}

var tonemap = map[EmojiModifier]rune{
	ModNone:        0,
	ModLight:       0x1f3fb,
	ModMediumLight: 0x1f3fc,
	ModMedium:      0x1f3fd,
	ModMediumDark:  0x1f3fe,
	ModDark:        0x1f3ff,
}
var tonenames = map[EmojiModifier]string{
	ModNone:        "",
	ModLight:       "light",
	ModMediumLight: "medium-light",
	ModMedium:      "mediun",
	ModMediumDark:  "medium-dark",
	ModDark:        "dark",
}

// Skintone always comes after the base emoji and doesn't required a ZWJ.
func (e Emoji) applyTone(t EmojiModifier) Emoji {
	if tcp := tonemap[t]; tcp > 0 {
		e.Name = e.Name + ": " + tonenames[t] + " skin tone"
		e.Codepoints = append(append([]rune{e.Codepoints[0]}, tcp), e.Codepoints[1:]...)
		l := len(e.Codepoints) - 1
		if e.Codepoints[l] == 0xfe0f {
			e.Codepoints = e.Codepoints[:l]
		}
	}
	return e
}
