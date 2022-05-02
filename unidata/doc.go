//go:generate ./gen/gen.zsh

// Package unidata contains information about the Unicode database.
//
// This is an alternative to Go's unicode package in stdlib and provides some
// additional data. In particular, it knows a bit more about codepoints and
// knows about emojis.
//
// The downside is that this package is a bit slower, uses more memory, and
// increases the binary size by about 2M. It should still be plenty fast enough
// for most use cases.
//
// This is updated to Unicode 14.0 (September 2021).
//
// NOTE: be careful in mixing this package and the stdlib unicode package; it
// usually takes a while before the tables in there are updated, and may result
// in inconsistent results. For example, unicode.IsPrint() will return false
// even for printable characters if it doesn't know about them. As of Go 1.18,
// unicode is using Unicode 13.0. The unicode/utf8 and unicode/utf16 are fine,
// as they just deal with the byte encodings.
package unidata
