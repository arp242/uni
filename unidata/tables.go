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

// Unicode versions.
const (
	UnicodeLatest = Unicode(iota)
	Unicode1_1
	Unicode2
	Unicode2_1
	Unicode3
	Unicode3_1
	Unicode3_2
	Unicode4
	Unicode4_1
	Unicode5
	Unicode5_1
	Unicode5_2
	Unicode6
	Unicode6_1
	Unicode6_2
	Unicode6_3
	Unicode7
	Unicode8
	Unicode9
	Unicode10
	Unicode11
	Unicode12
	Unicode12_1
	Unicode13
	Unicode14
	Unicode15
	Unicode15_1
	Unicode16
	Unicode17
)

// Unicodes is a list of all Unicode versions since 1.1.
var Unicodes = map[Unicode]struct {
	Name, Released string
}{
	Unicode1_1:    {"1.1", "June, 1993"},
	Unicode2:      {"2.0", "July, 1996"},
	Unicode2_1:    {"2.1", "May, 1998"},
	Unicode3:      {"3.0", "September, 1999"},
	Unicode3_1:    {"3.1", "March, 2001"},
	Unicode3_2:    {"3.2", "March, 2002"},
	Unicode4:      {"4.0", "April, 2003"},
	Unicode4_1:    {"4.1", "March, 2005"},
	Unicode5:      {"5.0", "July, 2006"},
	Unicode5_1:    {"5.1", "March, 2008"},
	Unicode5_2:    {"5.2", "October, 2009"},
	Unicode6:      {"6.0", "October, 2010"},
	Unicode6_1:    {"6.1", "January, 2012"},
	Unicode6_2:    {"6.2", "September, 2012"},
	Unicode6_3:    {"6.3", "September, 2013"},
	Unicode7:      {"7.0", "June, 2014"},
	Unicode8:      {"8.0", "June, 2015"},
	Unicode9:      {"9.0", "June, 2016"},
	Unicode10:     {"10.0", "June, 2017"},
	Unicode11:     {"11.0", "June, 2018"},
	Unicode12:     {"12.0", "March, 2019"},
	Unicode12_1:   {"12.1", "May, 2019"},
	Unicode13:     {"13.0", "March, 2020"},
	Unicode14:     {"14.0", "September, 2021"},
	Unicode15:     {"15.0", "September, 2022"},
	Unicode15_1:   {"15.1", "September, 2023"},
	Unicode16:     {"16.0", "September, 2024"},
	Unicode17:     {"17.0", "September, 2025"},
	UnicodeLatest: {"17.0", "September, 2025"},
}
