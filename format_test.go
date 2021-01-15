package main

import (
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	f, err := NewFormat("%(a) %(b l:auto) %(c)", false, false)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f.Line(map[string]string{
			"a": "col a",
			"b": "col b",
			"c": "col c",
		})
	}
}
