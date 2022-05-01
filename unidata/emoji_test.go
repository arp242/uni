package unidata

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestEmoji(t *testing.T) {
	var (
		shrug     = Emoji{Codepoints: []rune("🤷"), Name: "person shrugging", gender: genderSign, skinTones: true}
		handshake = Emoji{Codepoints: []rune("🤝"), Name: "handshake", skinTones: true}
	)
	tests := []struct {
		mod  []EmojiModifier
		in   Emoji
		want Emoji
	}{
		{[]EmojiModifier{ModMale},
			shrug,
			Emoji{Codepoints: []rune("🤷♂\ufe0f")}},
		{[]EmojiModifier{ModFemale},
			shrug,
			Emoji{Codepoints: []rune("🤷♀\ufe0f")}},
		{[]EmojiModifier{ModFemale | ModDark},
			shrug,
			Emoji{Codepoints: []rune("🤷🏿♀")}},

		{[]EmojiModifier{ModDark},
			handshake,
			Emoji{Codepoints: []rune("🤝🏿")}},
		{[]EmojiModifier{ModDark, ModLight},
			handshake,
			Emoji{Codepoints: []rune("🫱🏿‍🫲🏻")}},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := tt.in.With(tt.mod[0], tt.mod[1:]...)
			if !reflect.DeepEqual(have.Codepoints, tt.want.Codepoints) {
				t.Errorf("codepoints wrong\nhave: %-30s %q %s\nwant: %-30s %q",
					strings.Trim(fmt.Sprintf("% X", have.Codepoints), "[]"),
					have.Codepoints, have.String(),
					strings.Trim(fmt.Sprintf("% X", tt.want.Codepoints), "[]"),
					tt.want.Codepoints)
			}
		})
	}
}
