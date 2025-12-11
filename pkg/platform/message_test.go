package platform

import (
	"fmt"
	"testing"
)

// Message tests
func TestSplit(t *testing.T) {
	chunks := Split("This is a block of text that has been delimited by spaces only", " ", 20)

	if len(chunks) != 4 {
		t.Errorf(fmt.Sprintf("Failed TestSplit multiple chunks scenario, got %d instead of 4", len(chunks)))
	}
	if chunks[0] != "This is a block of " {
		t.Errorf(fmt.Sprintf("Failed TestSplit content check, got '%s'", chunks[0]))
	}

	mono := Split("This is a block of text that has been delimited by spaces and nothing else", " ", 100)

	if len(mono) != 1 {
		t.Errorf("Failed TestSplit single chunk scenario")
	}

	newlines := Split("Line1\nLine2\nLine3", "\n", 100)
	if len(newlines) != 1 {
		t.Errorf("Failed TestSplit newlines scenario, got %d chunks", len(newlines))
	}
	if newlines[0] != "Line1\nLine2\nLine3" {
		t.Errorf("Failed TestSplit newlines content, got %q", newlines[0])
	}
}

func TestNextFormatBlock(t *testing.T) {
	italics := NextFormatBlock(" _Italics_ *Bold* ^1234^ Text", 0)

	if italics.Type != Italics {
		t.Errorf("Failed TestNextFormatBlock italics format blocks scenario")
	}

	bold := NextFormatBlock(" _Italics_ *Bold* ^1234^ Text", 10)

	if bold.Type != Bold {
		t.Errorf("Failed TestNextFormatBlock bold format blocks scenario")
	}

	sup := NextFormatBlock(" _Italics_ *Bold* ^1234^ Text", 17)

	if sup.Type != Superscript {
		t.Errorf("Failed TestNextFormatBlock bold format blocks scenario")
	}

	null := NextFormatBlock("Text only no formatting", 0)

	if null.Type != Null {
		t.Errorf("Failed TestNextFormatBlock no format blocks scenario")
	}
}

func TestFormat(t *testing.T) {
	normal := func(s string) string { return s }
	bold := func(s string) string { return "Bold" }
	italics := func(s string) string { return "Italics" }
	super := func(s string) string { return "Superscript" }

	{
		output := Format("_Italics_ *Bold* ^1234^ Text", normal, bold, italics, super)

		if output != "Italics Bold Superscript Text" {
			t.Errorf(fmt.Sprintf("Failed TestFormat basic test, got %s", output))
		}
	}
	{
		output := Format("_Italics_ *Boldbold**Bold* _Ita_ ^1234^ Text^Super^Normal", normal, bold, italics, super)

		if output != "Italics BoldBold Italics Superscript TextSuperscriptNormal" {
			t.Errorf(fmt.Sprintf("Failed TestFormat compound test, got %s", output))
		}
	}
}
