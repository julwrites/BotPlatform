package platform

import (
	"sort"
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

func NextFormatBlock(str string, offset int) FormatBlock {
	var candidates []int
	var formattypes = []string{
		string(Bold),
		string(Italics),
		string(Superscript),
	}
	var block FormatBlock

	for _, f := range formattypes {
		i := strings.Index(str[offset:], f)
		if i == -1 {
			continue
		}
		candidates = append(candidates, i+offset)
	}

	if len(candidates) == 0 {
		block.Type = Null
		return block
	}

	sort.Ints(candidates)
	c := string(str[candidates[0]])

	block.Start = candidates[0]

	block.End = strings.Index(str[block.Start+1:], c)
	if block.End == -1 {
		block.Type = Null
		return block
	}

	block.End = block.End + block.Start + 1 // Account for starting offset + 2 markup symbols
	block.Type = FormatType(c)

	return block
}

func Format(str string, preprocess PreprocessingFormatter, bold BoldFormatter, ita ItalicsFormatter, sup SuperscriptFormatter) string {
	var outStr string

	str = preprocess(str)

	pos := 0
	for true {
		block := NextFormatBlock(str, pos)
		if block.Type == Null {
			break
		}

		outStr = outStr + str[pos:block.Start]   // Add any text before the formatter
		fmtStr := str[block.Start+1 : block.End] // Ignore the symbols

		switch block.Type {
		case Bold:
			fmtStr = bold(fmtStr)
			break
		case Italics:
			fmtStr = ita(fmtStr)
			break
		case Superscript:
			fmtStr = sup(fmtStr)
			break
		}

		outStr = outStr + fmtStr

		pos = block.End + 1
	}

	// Any leftovers
	outStr = outStr + str[pos:]

	return outStr
}
