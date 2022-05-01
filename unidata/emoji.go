package unidata

// Emoji is an emoji sequence.
type (
	Emoji struct {
		Codepoints      []rune   // Codepoints in this emoji
		Name            string   // Emoji name
		Group, Subgroup int      // Group and subgroup
		CLDR            []string // CLDR names
		SkinTones       bool     // Supports skintones?
		Genders         int      // Supports setting gender?
	}

	// TODO: use this.
	// EmojiGroup        uint16
	// EmojiSubGroup     uint16
	// EmojiGenderType   uint8
	// EmojiSkintoneType uint8

	// TODO: move the code to apply gender and skin tones from uni to here;
	// maybe something like:
	// e.Apply(unidata.Male | unidata.SkinDark)
	//   returns []rune
	//
	// e.Apply(unidata.Male | unidata.SkinDark | unidata.VariantText)
	// e.Apply(unidata.Male | unidata.SkinDark | unidata.VariantGraphic)
)

// Emoji genders types
const (
	GenderNone = 0
	GenderSign = 1
	GenderRole = 2
)

func (e Emoji) GroupName() string    { return EmojiGroups[e.Group] }
func (e Emoji) SubgroupName() string { return EmojiSubgroups[e.GroupName()][e.Subgroup] }

func (e Emoji) String() string {
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
