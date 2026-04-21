package platform

import (
	"strings"
)

func Split(msg string, delim string, maxSize int) []string {
	var splits []string

	msgStr := string(msg)
	paragraphs := strings.SplitAfter(msgStr, delim)

	var chunk string
	for _, para := range paragraphs {
		if len(chunk)+len(para) < maxSize {
			chunk += para
		} else {
			if len(chunk) > 0 {
				splits = append(splits, chunk)
			}
			chunk = para
		}
	}
	// Any leftovers should be accounted for
	if len(chunk) > 0 {
		splits = append(splits, chunk)
	}

	return splits
}

type PreprocessingFormatter func(string) string
type BoldFormatter func(string) string
type ItalicsFormatter func(string) string
type SuperscriptFormatter func(string) string

type FormatType string

const (
	Bold        FormatType = "*"
	Italics     FormatType = "_"
	Superscript FormatType = "^"
	Null        FormatType = "0"
)

type FormatBlock struct {
	Start int
	End   int
	Type  FormatType
}

var formattypes = []string{
	string(Bold),
	string(Italics),
	string(Superscript),
}

func NextFormatBlock(str string, offset int) FormatBlock {
	var block FormatBlock

	first := -1
	var firstF string
	for _, f := range formattypes {
		i := strings.Index(str[offset:], f)
		if i != -1 {
			idx := i + offset
			if first == -1 || idx < first {
				first = idx
				firstF = f
			}
		}
	}

	if first == -1 {
		block.Type = Null
		return block
	}

	block.Start = first

	block.End = strings.Index(str[block.Start+1:], firstF)
	if block.End == -1 {
		block.Type = Null
		return block
	}

	block.End = block.End + block.Start + 1 // Account for starting offset + 2 markup symbols
	block.Type = FormatType(firstF)

	return block
}

func Format(str string, preprocess PreprocessingFormatter, bold BoldFormatter, ita ItalicsFormatter, sup SuperscriptFormatter) string {
	str = preprocess(str)

	var builder strings.Builder
	builder.Grow(len(str))
	pos := 0
	for {
		block := NextFormatBlock(str, pos)
		if block.Type == Null {
			break
		}

		builder.WriteString(str[pos:block.Start]) // Add any text before the formatter
		fmtStr := str[block.Start+1 : block.End]  // Ignore the symbols

		switch block.Type {
		case Bold:
			fmtStr = bold(fmtStr)
		case Italics:
			fmtStr = ita(fmtStr)
		case Superscript:
			fmtStr = sup(fmtStr)
		}

		builder.WriteString(fmtStr)

		pos = block.End + 1
	}

	// Any leftovers
	builder.WriteString(str[pos:])

	return builder.String()
}
