package platform

import (
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	normal := func(s string) string { return s }
	bold := func(s string) string { return "*" + s + "*" }
	italics := func(s string) string { return "_" + s + "_" }
	super := func(s string) string { return "^" + s + "^" }

	input := "_Italics_ *Bold* ^1234^ Text _Italics_ *Bold* ^1234^ Text _Italics_ *Bold* ^1234^ Text"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(input, normal, bold, italics, super)
	}
}
