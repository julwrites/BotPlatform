package platform

import (
	"testing"
)

func BenchmarkFormat(b *testing.B) {
	preprocess := func(s string) string { return s }
	bold := func(s string) string { return "*" + s + "*" }
	italics := func(s string) string { return "_" + s + "_" }
	super := func(s string) string { return "^" + s + "^" }

	input := "This is a _test_ message with *bold* and _italics_ and ^superscript^ text. " +
		"Let's add more _italics_ and *bold* to make it longer. " +
		"And some ^1234^ numbers too."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Format(input, preprocess, bold, italics, super)
	}
}
