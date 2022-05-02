package unidata

import "strings"

// We don't want to be too pedantic about names, e.g. "Uppercase_letter" can
// also be found as "uppercaseletter".
//
// TODO: get tid of blockmap and catmap, and replace with a loop (possibly
// keeping a cache).
func nameKey(cat string) string {
	cat = strings.Replace(cat, " ", "", -1)
	cat = strings.Replace(cat, ",", "", -1)
	cat = strings.Replace(cat, "_", "", -1)
	cat = strings.ToLower(cat)
	return cat
}

// namekey → block
var blockmap = func() map[string]Block {
	b := make(map[string]Block)
	for k, v := range Blocks {
		b[nameKey(v.Name)] = k
	}
	return b
}()

var propmap = func() map[string]Property {
	b := make(map[string]Property)
	for k, v := range Properties {
		b[nameKey(v.Name)] = k
	}
	return b
}()

// namekey → Category
var catmap = map[string]Category{
	// Short-hand.
	"Lu": CatUppercaseLetter, "Ll": CatLowercaseLetter, "Lt": CatTitlecaseLetter,
	"LC": CatCasedLetter, "Lm": CatModifierLetter, "Lo": CatOtherLetter,
	"L": CatLetter, "Mn": CatNonspacingMark, "Mc": CatSpacingMark,
	"Me": CatEnclosingMark, "M": CatMark, "Nd": CatDecimalNumber,
	"Nl": CatLetterNumber, "No": CatOtherNumber, "N": CatNumber,
	"Pc": CatConnectorPunctuation, "Pd": CatDashPunctuation, "Ps": CatOpenPunctuation,
	"Pe": CatClosePunctuation, "Pi": CatInitialPunctuation, "Pf": CatFinalPunctuation,
	"Po": CatOtherPunctuation, "P": CatPunctuation, "Sm": CatMathSymbol,
	"Sc": CatCurrencySymbol, "Sk": CatModifierSymbol, "So": CatOtherSymbol,
	"S": CatSymbol, "Zs": CatSpaceSeparator, "Zl": CatLineSeparator,
	"Zp": CatParagraphSeparator, "Z": CatSeparator, "Cc": CatControl,
	"Cf": CatFormat, "Cs": CatSurrogate, "Co": CatPrivateUse,
	"Cn": CatUnassigned, "C": CatOther,

	// Lower-case shorthand.
	"lu": CatUppercaseLetter, "ll": CatLowercaseLetter, "lt": CatTitlecaseLetter,
	"lc": CatCasedLetter, "lm": CatModifierLetter, "lo": CatOtherLetter,
	"l": CatLetter, "mn": CatNonspacingMark, "mc": CatSpacingMark,
	"me": CatEnclosingMark, "m": CatMark, "nd": CatDecimalNumber,
	"nl": CatLetterNumber, "no": CatOtherNumber, "n": CatNumber,
	"pc": CatConnectorPunctuation, "pd": CatDashPunctuation, "ps": CatOpenPunctuation,
	"pe": CatClosePunctuation, "pi": CatInitialPunctuation, "pf": CatFinalPunctuation,
	"po": CatOtherPunctuation, "p": CatPunctuation, "sm": CatMathSymbol,
	"sc": CatCurrencySymbol, "sk": CatModifierSymbol, "so": CatOtherSymbol,
	"s": CatSymbol, "zs": CatSpaceSeparator, "zl": CatLineSeparator,
	"zp": CatParagraphSeparator, "z": CatSeparator, "cc": CatControl,
	"cf": CatFormat, "cs": CatSurrogate, "co": CatPrivateUse,
	"cn": CatUnassigned, "c": CatOther,

	// Full names, underscores.
	"uppercase_letter": CatUppercaseLetter, "lowercase_letter": CatLowercaseLetter,
	"titlecase_letter": CatTitlecaseLetter, "cased_letter": CatCasedLetter,
	"modifier_letter": CatModifierLetter, "other_letter": CatOtherLetter,
	"letter": CatLetter, "nonspacing_mark": CatNonspacingMark,
	"spacing_mark": CatSpacingMark, "enclosing_mark": CatEnclosingMark,
	"mark": CatMark, "decimal_number": CatDecimalNumber,
	"letter_number": CatLetterNumber, "other_number": CatOtherNumber,
	"number": CatNumber, "connector_punctuation": CatConnectorPunctuation,
	"dash_punctuation": CatDashPunctuation, "open_punctuation": CatOpenPunctuation,
	"close_punctuation": CatClosePunctuation, "initial_punctuation": CatInitialPunctuation,
	"final_punctuation": CatFinalPunctuation, "other_punctuation": CatOtherPunctuation,
	"punctuation": CatPunctuation, "math_symbol": CatMathSymbol,
	"currency_symbol": CatCurrencySymbol, "modifier_symbol": CatModifierSymbol,
	"other_symbol": CatOtherSymbol, "symbol": CatSymbol,
	"space_separator": CatSpaceSeparator, "line_separator": CatLineSeparator,
	"paragraph_separator": CatParagraphSeparator, "separator": CatSeparator,
	"control": CatControl, "format": CatFormat,
	"surrogate": CatSurrogate, "private_use": CatPrivateUse,
	"unassigned": CatUnassigned, "other": CatOther,

	// Without underscore.
	"uppercaseletter": CatUppercaseLetter, "lowercaseletter": CatLowercaseLetter,
	"titlecaseletter": CatTitlecaseLetter, "casedletter": CatCasedLetter,
	"modifierletter": CatModifierLetter, "otherletter": CatOtherLetter,
	"nonspacingmark": CatNonspacingMark, "spacingmark": CatSpacingMark,
	"enclosingmark": CatEnclosingMark, "decimalnumber": CatDecimalNumber,
	"letternumber": CatLetterNumber, "othernumber": CatOtherNumber,
	"connectorpunctuation": CatConnectorPunctuation, "dashpunctuation": CatDashPunctuation,
	"openpunctuation": CatOpenPunctuation, "closepunctuation": CatClosePunctuation,
	"initialpunctuation": CatInitialPunctuation, "finalpunctuation": CatFinalPunctuation,
	"otherpunctuation": CatOtherPunctuation, "mathsymbol": CatMathSymbol,
	"currencysymbol": CatCurrencySymbol, "modifiersymbol": CatModifierSymbol,
	"othersymbol": CatOtherSymbol, "spaceseparator": CatSpaceSeparator,
	"lineseparator": CatLineSeparator, "paragraphseparator": CatParagraphSeparator,
	"privateuse": CatPrivateUse,
}
