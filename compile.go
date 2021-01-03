package slimdown

import (
	"bytes"
	"html"
	"html/template"
	"net/url"
	"strings"
)

func CompileString(input string, options *Options) (output template.HTML, err error) {
	return Compile([]byte(input), options)
}

func Compile(input []byte, options *Options) (output template.HTML, err error) {
	tokens := tokenCollectionNew(input)

	if options == nil {
		options = &Options{}
	}

	err = compileTokenize(options, tokens)
	if err != nil {
		return
	}

	if options.DebugPrintTokens {
		debugPrintTokens(tokens)
	}

	err = compileGenerateHTML(options, tokens)
	if err != nil {
		return
	}

	if options.CleanEmptyTags {
		compileClean(tokens)
	}

	output = tokens.HTML()

	return
}

func compileTokenize(options *Options, tokens *tokenCollection) (err error) {
	backslashTokens := tokenCollectionNew(tokens.Input)

	if options.EnableDocumentTags {
		tokens.PushNewEmpty(tokenTypeDocumentDoctype)

		defer tokens.PushNewEmpty(tokenTypeDocumentHTMLBound)
		tokens.PushNewEmpty(tokenTypeDocumentHTMLBound)

		tokens.PushNewEmpty(tokenTypeDocumentHeadBound)
		tokens.PushNewEmpty(tokenTypeDocumentHeadBound)

		defer tokens.PushNewEmpty(tokenTypeDocumentBodyBound)
		tokens.PushNewEmpty(tokenTypeDocumentBodyBound)
	}

	if options.EnableParagraphTags {
		defer tokens.PushNewEmpty(tokenTypeParagraphBound)
		tokens.PushNewEmpty(tokenTypeParagraphBound)
	}

	for i, b := range tokens.Input {
		switch b {
		case '*':
			match := true

			if t := tokens.Peek(); t != nil {
				switch t.Type {
				case tokenTypeAsteriskDouble:
					t.Type = tokenTypeAsteriskTriple
				case tokenTypeAsterisk:
					t.Type = tokenTypeAsteriskDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				tokens.PushNewSingle(tokenTypeAsterisk, i)
			}
		case '_':
			if t := tokens.Peek(); t != nil && t.Type == tokenTypeUnderscore {
				t.Type = tokenTypeUnderscoreDouble
				t.InputEndIndex++
			} else {
				tokens.PushNewSingle(tokenTypeUnderscore, i)
			}
		case '\\':
			backslashTokens.PushAsIs(
				tokens.PushNewSingle(tokenTypeBackslash, i),
			)
			tokens.PushNewSingle(tokenTypeEmpty, i)
		case '=':
			var handled bool

			if t := tokens.Peek(); t != nil && t.Type == tokenTypeText {
				if l := t.Len(); l == 1 {
					if b := t.Bytes(); b[0] == '=' {
						t.Type = tokenTypeEqualsDouble
						handled = true
					}
				}
			}

			if !handled {
				tokens.PushNewSingle(tokenTypeText, i)
			}
		case '`':
			tokens.PushNewSingle(tokenTypeBacktick, i)
		case '\r':
			tokens.PushNewSingle(tokenTypeCarriageReturn, i)
		case '\n':
			tokens.PushNewSingle(tokenTypeParagraphBound, i)
			tokens.PushNewSingle(tokenTypeParagraphBound, i)
		case '\t':
			if tts := options.TabToSpaces; tts > 0 {
				for j := 0; j < tts; j++ {
					tokens.PushNewSingle(tokenTypeSpace, i)
				}
			} else {
				tokens.PushNewSingle(tokenTypeTab, i)
			}
		case '(':
			tokens.PushNewSingle(tokenTypeParenthesisOpen, i)
		case ')':
			tokens.PushNewSingle(tokenTypeParenthesisClose, i)
		case '[':
			tokens.PushNewSingle(tokenTypeSquareBracketOpen, i)
		case ']':
			tokens.PushNewSingle(tokenTypeSquareBracketClose, i)
		case ' ':
			t := tokens.PushNewSingle(tokenTypeSpace, i)

			if stt := options.SpacesToTab; stt > 0 {
				if stt == 1 {
					t.Type = tokenTypeTab
				} else {
					potentialTypes := make([]tokenType, stt-1)

					for i := range potentialTypes {
						potentialTypes[i] = tokenTypeSpace
					}

					if prevs, foundPrevs := t.PrevNTypes(potentialTypes); foundPrevs {
						t.Type = tokenTypeTab

						for _, p := range prevs.Data {
							p.Type = tokenTypeEmpty
						}
					}
				}
			}
		default:
			if t := tokens.Peek(); t != nil && t.Type == tokenTypeText {
				t.InputEndIndex++
			} else {
				tokens.PushNewSingle(tokenTypeText, i)
			}
		}
	}

	if options.EnableBackslashTransforms && backslashTokens.Len() > 0 {
		if err = compileTokenizeBackslashTransforms(backslashTokens); err != nil {
			return
		}
	}

	return
}

func compileTokenizeBackslashTransforms(tokens *tokenCollection) (err error) {
	for _, t := range tokens.Data {
		var isHandled bool

		if nextEmpty := t.Next; nextEmpty != nil && nextEmpty.Type == tokenTypeEmpty {
			if nextNextText := nextEmpty.Next; nextNextText != nil && nextNextText.Type == tokenTypeText {
				if nextNextText.Len() > 0 {
					switch isHandled = true; tokens.Input[nextNextText.InputStartIndex] {
					case 'n':
						t.Type = tokenTypeParagraphBound
						t.Next.Type = tokenTypeParagraphBound
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
						nextNextText.InputStartIndex++
					}
				}
			}
		}

		if !isHandled {
			err = ErrCompileBackslashTransformUnknown
			return
		}
	}

	return
}

func compileGenerateHTML(options *Options, tokens *tokenCollection) (err error) {
	tokenStack := tokenCollectionNewEmpty()

	for _, t := range tokens.Data {
		if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
			return
		}
	}

	if tokenStack.Len() > 0 {
		for {
			if t := tokenStack.PopAsIs(); t != nil {
				t.Type = tokenTypeText
				compileGenerateHTMLToken(options, t, tokenStack)
			} else {
				break
			}
		}
	}

	return
}

func compileGenerateHTMLToken(options *Options, t *token, tokenStack *tokenCollection) (err error) {
	switch y := t.Type; y {
	case tokenTypeEmpty:
		break
	case tokenTypeSpace:
		t.HTML = []byte{' '}

		if mcs := options.MaxConsecutiveSpaces; mcs > 0 {
			nextSpaceTokens, foundNextSpaceTokens := t.NextUntilEndOfPotentialTypes(y)
			if foundNextSpaceTokens {
				for i, t2 := range nextSpaceTokens.Data {
					if i+1 < mcs {
						continue
					}
					t2.Type = tokenTypeEmpty
				}
			}
		}
	case tokenTypeParagraphBound,
		tokenTypeDocumentBodyBound,
		tokenTypeDocumentHeadBound,
		tokenTypeDocumentHTMLBound:
		err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
	case tokenTypeSquareBracketOpen:
		if !options.EnableLinks {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		textTokens, foundTextTokens := t.NextUntilEndOfPotentialTypes(
			tokenTypeText,
			tokenTypeSpace,
			tokenTypeAsterisk,
			tokenTypeAsteriskDouble,
			tokenTypeUnderscore,
			tokenTypeUnderscoreDouble,
		)
		if !foundTextTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypes([]tokenType{
			tokenTypeSquareBracketClose,
			tokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextUntilEndOfPotentialTypes(
			tokenTypeText,
			tokenTypeAsterisk,
			tokenTypeAsteriskDouble,
			tokenTypeUnderscore,
			tokenTypeUnderscoreDouble,
		)
		if !foundLinkTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		finalToken := linkTokens.Get(-1).Next
		if finalToken == nil || finalToken.Type != tokenTypeParenthesisClose {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		var textBuff, linkBuff bytes.Buffer

		for _, t2 := range textTokens.Data {
			textBuff.Write(t2.Bytes())
		}
		for _, t2 := range linkTokens.Data {
			linkBuff.Write(t2.Bytes())
		}

		textBytes, linkString := textBuff.Bytes(), linkBuff.String()

		var linkURL *url.URL
		linkURL, err = url.Parse(linkString)
		if err != nil {
			return
		}

		linkString = linkURL.String()
		if strings.Contains(linkString, `"`) {
			err = ErrCompileURLCannotContainDoubleQuote
			return
		}

		var buff bytes.Buffer

		buff.WriteString(`<a href="`)
		buff.WriteString(linkString)
		buff.WriteString(`">`)
		buff.Write(textBytes)
		buff.WriteString(`</a>`)

		compileGenerateHTMLTokenHandleBuffer(t, &buff)

		allTokens := tokenCollectionNewEmpty()

		allTokens.PushAsIs(textTokens.Data...)
		allTokens.PushAsIs(midTokens.Data...)
		allTokens.PushAsIs(linkTokens.Data...)
		allTokens.PushAsIs(finalToken)

		for _, t2 := range allTokens.Data {
			t2.Type = tokenTypeEmpty
		}
	case tokenTypeBackslash,
		tokenTypeParenthesisOpen,
		tokenTypeParenthesisClose,
		tokenTypeSquareBracketClose:
		compileGenerateHTMLTokenHandleBytes(t)
	case tokenTypeText:
		compileGenerateHTMLTokenHandleBytes(t)

		if !options.AllowHTML {
			t.HTML = []byte(html.EscapeString(string(t.HTML)))
		}
	case tokenTypeTab:
		t.HTML = []byte{'\t'}
	case tokenTypeCarriageReturn:
		t.HTML = []byte{}

		if _, foundNewline := t.NextUntilEndOfPotentialTypes(tokenTypeParagraphBound); !foundNewline {
			prevs, foundPrevs := t.PrevUntilStartOfPotentialTypes(
				tokenTypeParagraphBound,
			)
			if foundPrevs {
				for _, t2 := range prevs.Data {
					if t2.Type == tokenTypeText {
						t2.Type = tokenTypeEmpty
					}
				}
			}
		}
	case tokenTypeBacktick:
		if !options.EnableCodeTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if next := t.Next; next != nil && next.Type == y {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if prev := t.Prev; prev != nil && prev.Type == y {
			if prevPrev := prev.Prev; prevPrev == nil || prevPrev.Type != y {
				break
			}
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
	case tokenTypeUnderscoreTriple, tokenTypeAsteriskTriple:
		if options.EnableStrongTags && options.EnableEmTags {
			err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
			break
		}

		fallthrough
	case tokenTypeUnderscoreDouble, tokenTypeAsteriskDouble:
		if !options.EnableStrongTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
	case tokenTypeUnderscore, tokenTypeAsterisk:
		if !options.EnableEmTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
	case tokenTypeEqualsDouble:
		if !options.EnableMarkTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack)
	case tokenTypeDocumentDoctype:
		t.HTML = []byte(`<!DOCTYPE html>`)
	default:
		err = ErrCompileTokenTypeUnknown
	}

	return
}

func compileGenerateHTMLTokenHandleBuffer(t *token, b *bytes.Buffer) {
	t.HTML = b.Bytes()
}

func compileGenerateHTMLTokenHandleBytes(t *token) {
	t.HTML = t.Bytes()
}

func compileGenerateHTMLTokenHandleTag(t *token, tokenStack *tokenCollection) (err error) {
	y := t.Type

	if y != tokenTypeBacktick && tokenStack.ContainsType(tokenTypeBacktick) {
		compileGenerateHTMLTokenHandleBytes(t)
		return
	}

	if t2 := tokenStack.Peek(); t2 != nil && t2.Type == y {
		tags, foundTag := tokenTypeTagMap[y]
		tagsLen := len(tags)
		if !foundTag || tagsLen < 1 {
			err = ErrCompileTagsForTokenNotFound
			return
		}

		t.HTML, t2.HTML = []byte{}, []byte{}

		for i := 0; i < tagsLen; i++ {
			t2.HTML = append(t2.HTML, []byte("<"+tags[i]+">")...)
			t.HTML = append(t.HTML, []byte("</"+tags[tagsLen-i-1]+">")...)
		}

		tokenStack.PopAsIs()

		t.Collection.TagPairCleanData = append(t.Collection.TagPairCleanData, [2]*token{t2, t})
	} else {
		tokenStack.PushAsIs(t)
	}

	return
}

func compileClean(tokens *tokenCollection) {
	for _, pair := range tokens.TagPairCleanData {
		startTagToken, endTagToken := pair[0], pair[1]

		shouldClean := true
		prevs, foundPrevs := endTagToken.PrevUntilMeetToken(startTagToken)
		if foundPrevs {
			for _, p := range prevs.Data {
				switch p.Type {
				case tokenTypeText:
					if p.Len() > 0 {
						shouldClean = false
					}
				default:
					break
				}

				if !shouldClean {
					break
				}
			}
		}

		if shouldClean {
			startTagToken.Type = tokenTypeEmpty
			endTagToken.Type = tokenTypeEmpty

			for _, p := range prevs.Data {
				p.Type = tokenTypeEmpty
			}
		}
	}
}
