package slimdown

import (
	"fmt"
	"html/template"
)

func ParseString(input string, options *Options) (output template.HTML, err error) {
	return Parse([]byte(input), options)
}

func Parse(input []byte, options *Options) (output template.HTML, err error) {
	tokens := &tokenCollection{}

	if options == nil {
		options = &Options{}
	}

	parseTokenize(input, options, tokens)

	for _, t := range tokens.Data {
		fmt.Print(t.String() + " ")
	}
	fmt.Println()

	err = parseGenerateHTML(input, options, tokens)
	if err != nil {
		return
	}

	output = tokens.HTML()

	return
}

func parseTokenize(input []byte, options *Options, tokens *tokenCollection) {
	backslashTokens := &tokenCollection{}

	defer tokens.PushNewEmpty(tokenTypeNewline)
	tokens.PushNewEmpty(tokenTypeNewline)

	for i, b := range input {
		switch b {
		case '*':
			if t := tokens.Peek(); t != nil && t.Type == tokenTypeAsterisk {
				t.Type = tokenTypeAsteriskDouble
				t.EndIndex++
			} else {
				tokens.PushNewSingle(tokenTypeAsterisk, i)
			}
		case '_':
			if t := tokens.Peek(); t != nil && t.Type == tokenTypeUnderscore {
				t.Type = tokenTypeUnderscoreDouble
				t.EndIndex++
			} else {
				tokens.PushNewSingle(tokenTypeUnderscore, i)
			}
		case '\\':
			backslashTokens.PushAsIs(
				tokens.PushNewSingle(tokenTypeBackslash, i),
			)
		case '`':
			tokens.PushNewSingle(tokenTypeBacktick, i)
		case '\r':
			tokens.PushNewSingle(tokenTypeCarriageReturn, i)
		case '\n':
			tokens.PushNewSingle(tokenTypeNewline, i)
		case '\t':
			tokens.PushNewSingle(tokenTypeTab, i)
		default:
			if t := tokens.Peek(); t != nil && t.Type == tokenTypeText {
				t.EndIndex++
			} else {
				tokens.PushNewSingle(tokenTypeText, i)
			}
		}
	}

	if options.EnableBackslashTransforms && backslashTokens.Len() > 0 {
		parseTokenizeBackslashTransforms(input, backslashTokens)
	}
}

func parseTokenizeBackslashTransforms(input []byte, tokens *tokenCollection) {
	for _, t := range tokens.Data {
		var isHandled bool

		if next := t.Next; next != nil {
			if next.Type == tokenTypeText && next.EndIndex > next.StartIndex {
				switch isHandled = true; input[next.StartIndex] {
				case 'n':
					t.Type = tokenTypeNewline
				case 'r':
					t.Type = tokenTypeCarriageReturn
				case 't':
					t.Type = tokenTypeTab
				case '\\':
					break
				default:
					isHandled = false
				}

				if isHandled {
					next.StartIndex++
				}
			}
		}

		if !isHandled {
			panic("")
		}
	}
}

func parseGenerateHTML(input []byte, options *Options, tokens *tokenCollection) (err error) {
	tokenStack := &tokenCollection{}

	for _, t := range tokens.Data {
		if err = parseGenerateHTMLToken(input, options, t, tokenStack); err != nil {
			return
		}
	}

	if tokenStack.Len() > 0 {
		err = ErrParseTokenStackOverflow
		return
	}

	return
}

func parseGenerateHTMLToken(input []byte, options *Options, t *token, tokenStack *tokenCollection) (err error) {
	switch y := t.Type; y {
	case tokenTypeNewline:
		if options.EnableParagraphs {
			tag := tokenTypeTagMap[y]
			t.HTML = []byte{}

			if t.Prev != nil {
				t.HTML = append(t.HTML, []byte("</"+tag+">")...)
			}

			if t.Next != nil {
				t.HTML = append(t.HTML, []byte("<"+tag+">")...)
			}

			break
		}

		fallthrough
	case tokenTypeText, tokenTypeBackslash:
		parseGenerateHTMLTokenRaw(input, t)
	// case tokenTypeBackslash:
	// 	if next := t.Next; next != nil {
	// 		if next.Type == tokenTypeText && next.EndIndex > next.StartIndex {
	// 			switch input[next.StartIndex] {
	// 			case 'n':
	// 				t.Type = tokenTypeNewline
	// 			case 'r':
	// 				t.Type = tokenTypeCarriageReturn
	// 			default:
	// 				next.StartIndex--
	// 			}

	// 			next.StartIndex++
	// 		}
	// 	}

	// 	if t.Type != tokenTypeBackslash {
	// 		if err = parseGenerateHTMLToken(input, t, tokenStack, useParagraphs); err != nil {
	// 			return
	// 		}
	// 	}
	case tokenTypeTab:
		t.HTML = []byte{'\t'}
	case tokenTypeCarriageReturn:
		t.HTML = []byte{}

		if next := t.Next; next == nil || next.Type != tokenTypeNewline {
			for t2 := t.Prev; t2 != nil && t2.Type != tokenTypeNewline; t2 = t2.Prev {
				if t2.Type == tokenTypeText {
					t2.HTML = []byte{}
				}
			}
		}
	case tokenTypeBacktick:
		if !options.EnableCode {
			parseGenerateHTMLTokenRaw(input, t)
			return
		}

		if next := t.Next; next != nil && next.Type == y {
			parseGenerateHTMLTokenRaw(input, t)
			break
		}

		if prev := t.Prev; prev != nil && prev.Type == y {
			if prevPrev := prev.Prev; prevPrev == nil || prevPrev.Type != y {
				break
			}
		}

		fallthrough
	case tokenTypeAsterisk, tokenTypeAsteriskDouble,
		tokenTypeUnderscore, tokenTypeUnderscoreDouble:
		if y != tokenTypeBacktick && tokenStack.ContainsType(tokenTypeBacktick) {
			t.HTML = input[t.StartIndex:t.EndIndex]
			break
		}

		if t2 := tokenStack.PeekTwin(); t2 != nil && t2.Type == y {
			tag := tokenTypeTagMap[y]
			t2.HTML = []byte("<" + tag + ">")
			t.HTML = []byte("</" + tag + ">")

			tokenStack.Pop()
		} else {
			if _, err = tokenStack.PushTwin(t); err != nil {
				return
			}
		}
	default:
		err = ErrParseTokenTypeUnknown
	}

	return
}

func parseGenerateHTMLTokenRaw(input []byte, t *token) {
	if t.EndIndex > t.StartIndex {
		t.HTML = input[t.StartIndex:t.EndIndex]
	} else {
		t.HTML = []byte{}
	}
}
