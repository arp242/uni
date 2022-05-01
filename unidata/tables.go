package unidata

// Unicode planes.
const (
	PlaneUnknown    = Plane(iota)
	PlaneBMP        // 0:     Basic Multilingual Plane
	PlaneSMP        // 1:     Supplementary Multilingual Plane
	PlaneSIP        // 2:     Supplementary Ideographic Plane
	PlaneTIP        // 3:     Tertiary Ideographic Plane
	PlaneUnassigned // 4-13:  Unassigned
	PlaneSpecial    // 14:    Supplementary Special-purpose Plane
	PlanePrivate    // 15-16: Supplementary Private Use Area planes
)

// Planes is a list of all Unicode planes.
var Planes = map[Plane]struct {
	Range [2]rune
	Name  string
}{
	PlaneBMP:        {[2]rune{0x00000, 0x00FFFF}, "Basic Multilingual Plane"},
	PlaneSMP:        {[2]rune{0x10000, 0x01FFFF}, "Supplementary Multilingual Plane"},
	PlaneSIP:        {[2]rune{0x20000, 0x02FFFF}, "Supplementary Ideographic Plane"},
	PlaneTIP:        {[2]rune{0x30000, 0x03FFFF}, "Tertiary Ideographic Plane"},
	PlaneUnassigned: {[2]rune{0x40000, 0x0DFFFF}, "Unassigned"},
	PlaneSpecial:    {[2]rune{0xE0000, 0x0EFFFF}, "Supplementary Special-purpose Plane"},
	PlanePrivate:    {[2]rune{0xF0000, 0x10FFFF}, "Supplementary Private Use Area planes"},
}

// Unicode widths.
const (
	WidthAmbiguous = Width(iota) // Ambiguous, A
	WidthFullWidth               // FullWidth, F
	WidthHalfWidth               // Halfwidth, H
	WidthNarrow                  // Narrow, N
	WidthNeutral                 // Neutral (Not East Asian), Na
	WidthWide                    // Wide, W
)

// Widths is a list of all Unicode Widths.
var Widths = map[Width]string{
	WidthAmbiguous: "ambiguous",
	WidthFullWidth: "full",
	WidthHalfWidth: "half",
	WidthNarrow:    "narrow",
	WidthNeutral:   "neutral",
	WidthWide:      "wide",
}
