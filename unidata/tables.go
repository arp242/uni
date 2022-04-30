package unidata

const (
	WidthAmbiguous = uint8(iota) // Ambiguous, A
	WidthFullWidth               // FullWidth, F
	WidthHalfWidth               // Halfwidth, H
	WidthNarrow                  // Narrow, N
	WidthNeutral                 // Neutral (Not East Asian), Na
	WidthWide                    // Wide, W
)

var WidthNames = map[uint8]string{
	WidthAmbiguous: "ambiguous",
	WidthFullWidth: "full",
	WidthHalfWidth: "half",
	WidthNarrow:    "narrow",
	WidthNeutral:   "neutral",
	WidthWide:      "wide",
}

// http://www.unicode.org/reports/tr44/#General_Category_Values
const (
	CatUnknown              = Category(iota)
	CatUppercaseLetter      // Lu – an uppercase letter
	CatLowercaseLetter      // Ll – a lowercase letter
	CatTitlecaseLetter      // Lt – a digraphic character, with first part uppercase
	CatCasedLetter          // LC – Lu | Ll | Lt
	CatModifierLetter       // Lm – a modifier letter
	CatOtherLetter          // Lo – other letters, including syllables and ideographs
	CatLetter               // L  – Lu | Ll | Lt | Lm | Lo
	CatNonspacingMark       // Mn – a nonspacing combining mark (zero advance width)
	CatSpacingMark          // Mc – a spacing combining mark (positive advance width)
	CatEnclosingMark        // Me – an enclosing combining mark
	CatMark                 // M  – Mn | Mc | Me
	CatDecimalNumber        // Nd – a decimal digit
	CatLetterNumber         // Nl – a letterlike numeric character
	CatOtherNumber          // No – a numeric character of other type
	CatNumber               // N  – Nd | Nl | No
	CatConnectorPunctuation // Pc – a connecting punctuation mark, like a tie
	CatDashPunctuation      // Pd – a dash or hyphen punctuation mark
	CatOpenPunctuation      // Ps – an opening punctuation mark (of a pair)
	CatClosePunctuation     // Pe – a closing punctuation mark (of a pair)
	CatInitialPunctuation   // Pi – an initial quotation mark
	CatFinalPunctuation     // Pf – a final quotation mark
	CatOtherPunctuation     // Po – a punctuation mark of other type
	CatPunctuation          // P  – Pc | Pd | Ps | Pe | Pi | Pf | Po
	CatMathSymbol           // Sm – a symbol of mathematical use
	CatCurrencySymbol       // Sc – a currency sign
	CatModifierSymbol       // Sk – a non-letterlike modifier symbol
	CatOtherSymbol          // So – a symbol of other type
	CatSymbol               // S  – Sm | Sc | Sk | So
	CatSpaceSeparator       // Zs – a space character (of various non-zero widths)
	CatLineSeparator        // Zl – U+2028 LINE SEPARATOR only
	CatParagraphSeparator   // Zp – U+2029 PARAGRAPH SEPARATOR only
	CatSeparator            // Z  – Zs | Zl | Zp
	CatControl              // Cc – a C0 or C1 control code
	CatFormat               // Cf – a format control character
	CatSurrogate            // Cs – a surrogate code point
	CatPrivateUse           // Co – a private-use character
	CatUnassigned           // Cn – a reserved unassigned code point or a noncharacter
	CatOther                // C  – Cc | Cf | Cs | Co | Cn
)

var Planes = map[string][2]rune{
	"Basic Multilingual Plane":              {0, 0xFFFF},
	"Supplementary Multilingual Plane":      {0x10000, 0x1FFFF},
	"Supplementary Ideographic Plane":       {0x20000, 0x2FFFF},
	"Tertiary Ideographic Plane":            {0x30000, 0x3FFFF},
	"Unassigned":                            {0x40000, 0xDFFFF},
	"Supplementary Special-purpose Plane":   {0xE0000, 0xEFFFF},
	"Supplementary Private Use Area planes": {0xF0000, 0x10FFFF},
}

var Blockmap = make(map[string]string)

func init() {
	for k := range Blocks {
		Blockmap[CanonicalCategory(k)] = k
	}
}

var (
	Catmap = map[string]Category{
		// Short-hand.
		"Lu": CatUppercaseLetter,
		"Ll": CatLowercaseLetter,
		"Lt": CatTitlecaseLetter,
		"LC": CatCasedLetter,
		"Lm": CatModifierLetter,
		"Lo": CatOtherLetter,
		"L":  CatLetter,
		"Mn": CatNonspacingMark,
		"Mc": CatSpacingMark,
		"Me": CatEnclosingMark,
		"M":  CatMark,
		"Nd": CatDecimalNumber,
		"Nl": CatLetterNumber,
		"No": CatOtherNumber,
		"N":  CatNumber,
		"Pc": CatConnectorPunctuation,
		"Pd": CatDashPunctuation,
		"Ps": CatOpenPunctuation,
		"Pe": CatClosePunctuation,
		"Pi": CatInitialPunctuation,
		"Pf": CatFinalPunctuation,
		"Po": CatOtherPunctuation,
		"P":  CatPunctuation,
		"Sm": CatMathSymbol,
		"Sc": CatCurrencySymbol,
		"Sk": CatModifierSymbol,
		"So": CatOtherSymbol,
		"S":  CatSymbol,
		"Zs": CatSpaceSeparator,
		"Zl": CatLineSeparator,
		"Zp": CatParagraphSeparator,
		"Z":  CatSeparator,
		"Cc": CatControl,
		"Cf": CatFormat,
		"Cs": CatSurrogate,
		"Co": CatPrivateUse,
		"Cn": CatUnassigned,
		"C":  CatOther,

		// Lower-case shorthand.
		"lu": CatUppercaseLetter,
		"ll": CatLowercaseLetter,
		"lt": CatTitlecaseLetter,
		"lc": CatCasedLetter,
		"lm": CatModifierLetter,
		"lo": CatOtherLetter,
		"l":  CatLetter,
		"mn": CatNonspacingMark,
		"mc": CatSpacingMark,
		"me": CatEnclosingMark,
		"m":  CatMark,
		"nd": CatDecimalNumber,
		"nl": CatLetterNumber,
		"no": CatOtherNumber,
		"n":  CatNumber,
		"pc": CatConnectorPunctuation,
		"pd": CatDashPunctuation,
		"ps": CatOpenPunctuation,
		"pe": CatClosePunctuation,
		"pi": CatInitialPunctuation,
		"pf": CatFinalPunctuation,
		"po": CatOtherPunctuation,
		"p":  CatPunctuation,
		"sm": CatMathSymbol,
		"sc": CatCurrencySymbol,
		"sk": CatModifierSymbol,
		"so": CatOtherSymbol,
		"s":  CatSymbol,
		"zs": CatSpaceSeparator,
		"zl": CatLineSeparator,
		"zp": CatParagraphSeparator,
		"z":  CatSeparator,
		"cc": CatControl,
		"cf": CatFormat,
		"cs": CatSurrogate,
		"co": CatPrivateUse,
		"cn": CatUnassigned,
		"c":  CatOther,

		// Full names, underscores.
		"uppercase_letter":      CatUppercaseLetter,
		"lowercase_letter":      CatLowercaseLetter,
		"titlecase_letter":      CatTitlecaseLetter,
		"cased_letter":          CatCasedLetter,
		"modifier_letter":       CatModifierLetter,
		"other_letter":          CatOtherLetter,
		"letter":                CatLetter,
		"nonspacing_mark":       CatNonspacingMark,
		"spacing_mark":          CatSpacingMark,
		"enclosing_mark":        CatEnclosingMark,
		"mark":                  CatMark,
		"decimal_number":        CatDecimalNumber,
		"letter_number":         CatLetterNumber,
		"other_number":          CatOtherNumber,
		"number":                CatNumber,
		"connector_punctuation": CatConnectorPunctuation,
		"dash_punctuation":      CatDashPunctuation,
		"open_punctuation":      CatOpenPunctuation,
		"close_punctuation":     CatClosePunctuation,
		"initial_punctuation":   CatInitialPunctuation,
		"final_punctuation":     CatFinalPunctuation,
		"other_punctuation":     CatOtherPunctuation,
		"punctuation":           CatPunctuation,
		"math_symbol":           CatMathSymbol,
		"currency_symbol":       CatCurrencySymbol,
		"modifier_symbol":       CatModifierSymbol,
		"other_symbol":          CatOtherSymbol,
		"symbol":                CatSymbol,
		"space_separator":       CatSpaceSeparator,
		"line_separator":        CatLineSeparator,
		"paragraph_separator":   CatParagraphSeparator,
		"separator":             CatSeparator,
		"control":               CatControl,
		"format":                CatFormat,
		"surrogate":             CatSurrogate,
		"private_use":           CatPrivateUse,
		"unassigned":            CatUnassigned,
		"other":                 CatOther,

		// Without underscore.
		"uppercaseletter":      CatUppercaseLetter,
		"lowercaseletter":      CatLowercaseLetter,
		"titlecaseletter":      CatTitlecaseLetter,
		"casedletter":          CatCasedLetter,
		"modifierletter":       CatModifierLetter,
		"otherletter":          CatOtherLetter,
		"nonspacingmark":       CatNonspacingMark,
		"spacingmark":          CatSpacingMark,
		"enclosingmark":        CatEnclosingMark,
		"decimalnumber":        CatDecimalNumber,
		"letternumber":         CatLetterNumber,
		"othernumber":          CatOtherNumber,
		"connectorpunctuation": CatConnectorPunctuation,
		"dashpunctuation":      CatDashPunctuation,
		"openpunctuation":      CatOpenPunctuation,
		"closepunctuation":     CatClosePunctuation,
		"initialpunctuation":   CatInitialPunctuation,
		"finalpunctuation":     CatFinalPunctuation,
		"otherpunctuation":     CatOtherPunctuation,
		"mathsymbol":           CatMathSymbol,
		"currencysymbol":       CatCurrencySymbol,
		"modifiersymbol":       CatModifierSymbol,
		"othersymbol":          CatOtherSymbol,
		"spaceseparator":       CatSpaceSeparator,
		"lineseparator":        CatLineSeparator,
		"paragraphseparator":   CatParagraphSeparator,
		"privateuse":           CatPrivateUse,
	}

	Catnames = map[Category]string{
		CatUppercaseLetter:      "Uppercase_Letter",
		CatLowercaseLetter:      "Lowercase_Letter",
		CatTitlecaseLetter:      "Titlecase_Letter",
		CatCasedLetter:          "Cased_Letter",
		CatModifierLetter:       "Modifier_Letter",
		CatOtherLetter:          "Other_Letter",
		CatLetter:               "Letter",
		CatNonspacingMark:       "Nonspacing_Mark",
		CatSpacingMark:          "Spacing_Mark",
		CatEnclosingMark:        "Enclosing_Mark",
		CatMark:                 "Mark",
		CatDecimalNumber:        "Decimal_Number",
		CatLetterNumber:         "Letter_Number",
		CatOtherNumber:          "Other_Number",
		CatNumber:               "Number",
		CatConnectorPunctuation: "Connector_Punctuation",
		CatDashPunctuation:      "Dash_Punctuation",
		CatOpenPunctuation:      "Open_Punctuation",
		CatClosePunctuation:     "Close_Punctuation",
		CatInitialPunctuation:   "Initial_Punctuation",
		CatFinalPunctuation:     "Final_Punctuation",
		CatOtherPunctuation:     "Other_Punctuation",
		CatPunctuation:          "Punctuation",
		CatMathSymbol:           "Math_Symbol",
		CatCurrencySymbol:       "Currency_Symbol",
		CatModifierSymbol:       "Modifier_Symbol",
		CatOtherSymbol:          "Other_Symbol",
		CatSymbol:               "Symbol",
		CatSpaceSeparator:       "Space_Separator",
		CatLineSeparator:        "Line_Separator",
		CatParagraphSeparator:   "Paragraph_Separator",
		CatSeparator:            "Separator",
		CatControl:              "Control",
		CatFormat:               "Format",
		CatSurrogate:            "Surrogate",
		CatPrivateUse:           "Private_Use",
		CatUnassigned:           "Unassigned",
		CatOther:                "Other",
	}
)

var (
	ranges = [][]rune{
		{0x3400, 0x4DB5},
		{0x4E00, 0x9FEF},
		{0xAC00, 0xD7A3},
		{0xD800, 0xDB7F},
		{0xDB80, 0xDBFF},
		{0xDC00, 0xDFFF},
		{0xE000, 0xF8FF},
		{0x17000, 0x187F1},
		{0x20000, 0x2A6D6},
		{0x2A700, 0x2B734},
		{0x2B740, 0x2B81D},
		{0x2B820, 0x2CEA1},
		{0x2CEB0, 0x2EBE0},
		{0xF0000, 0xFFFFD},
		{0x100000, 0x10FFFD},
	}

	rangeNames = []string{
		"<CJK Ideograph Extension A>",
		"<CJK Ideograph>",
		"<Hangul Syllable>",
		"<Non Private Use High Surrogate>",
		"<Private Use High Surrogate>",
		"<Low Surrogate>",
		"<Private Use>",
		"<Tangut Ideograph>",
		"<CJK Ideograph Extension B>",
		"<CJK Ideograph Extension C>",
		"<CJK Ideograph Extension D>",
		"<CJK Ideograph Extension E>",
		"<CJK Ideograph Extension F>",
		"<Plane 15 Private Use>",
		"<Plane 16 Private Use>",
	}
)
