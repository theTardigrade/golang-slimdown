package slimdown

import (
	"bytes"
	"html"
	"html/template"
	"net/url"
	"sort"
	"strings"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

func CompileString(input string, options *Options) (output template.HTML, err error) {
	return Compile([]byte(input), options)
}

func Compile(input []byte, options *Options) (output template.HTML, err error) {
	tokens := tokenization.TokenCollectionNew(input)

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

func compileTokenize(options *Options, tokens *tokenization.TokenCollection) (err error) {
	backslashTokens := tokenization.TokenCollectionNew(tokens.Input)

	if options.EnableDocumentTags {
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentDoctype)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)

		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
	}

	if options.EnableParagraphTags {
		defer tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)
		tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)
	}

	for i, b := range tokens.Input {
		switch b {
		case '*':
			match := true

			if t := tokens.Peek(); t != nil {
				switch t.Type {
				case tokenization.TokenTypeAsteriskDouble:
					t.Type = tokenization.TokenTypeAsteriskTriple
				case tokenization.TokenTypeAsterisk:
					t.Type = tokenization.TokenTypeAsteriskDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				tokens.PushNewSingle(tokenization.TokenTypeAsterisk, i)
			}
		case '_':
			if t := tokens.Peek(); t != nil && t.Type == tokenization.TokenTypeUnderscore {
				t.Type = tokenization.TokenTypeUnderscoreDouble
				t.InputEndIndex++
			} else {
				tokens.PushNewSingle(tokenization.TokenTypeUnderscore, i)
			}
		case '\\':
			backslashTokens.PushAsIs(
				tokens.PushNewSingle(tokenization.TokenTypeBackslash, i),
			)
			tokens.PushNewSingle(tokenization.TokenTypeEmpty, i)
		case '=':
			var handled bool

			if t := tokens.Peek(); t != nil && t.Type == tokenization.TokenTypeText {
				if l := t.Len(); l == 1 {
					if b := t.Bytes(); b[0] == '=' {
						t.Type = tokenization.TokenTypeEqualsDouble
						handled = true
					}
				}
			}

			if !handled {
				tokens.PushNewSingle(tokenization.TokenTypeText, i)
			}
		case '`':
			tokens.PushNewSingle(tokenization.TokenTypeBacktick, i)
		case '!':
			tokens.PushNewSingle(tokenization.TokenTypeExclamation, i)
		case '\r':
			tokens.PushNewSingle(tokenization.TokenTypeCarriageReturn, i)
		case '\n':
			tokens.PushNewSingle(tokenization.TokenTypeParagraphBound, i)
			tokens.PushNewSingle(tokenization.TokenTypeParagraphBound, i)
		case '\t':
			if tts := options.TabToSpaces; tts > 0 {
				for j := 0; j < tts; j++ {
					tokens.PushNewSingle(tokenization.TokenTypeSpace, i)
				}
			} else {
				tokens.PushNewSingle(tokenization.TokenTypeTab, i)
			}
		case '(':
			tokens.PushNewSingle(tokenization.TokenTypeParenthesisOpen, i)
		case ')':
			tokens.PushNewSingle(tokenization.TokenTypeParenthesisClose, i)
		case '[':
			tokens.PushNewSingle(tokenization.TokenTypeSquareBracketOpen, i)
		case ']':
			tokens.PushNewSingle(tokenization.TokenTypeSquareBracketClose, i)
		case ' ':
			t := tokens.PushNewSingle(tokenization.TokenTypeSpace, i)

			if stt := options.SpacesToTab; stt > 0 {
				if stt == 1 {
					t.Type = tokenization.TokenTypeTab
				} else {
					potentialTypes := make([]tokenization.TokenType, stt-1)

					for i := range potentialTypes {
						potentialTypes[i] = tokenization.TokenTypeSpace
					}

					if prevs, foundPrevs := t.PrevNTypes(potentialTypes); foundPrevs {
						t.Type = tokenization.TokenTypeTab

						for _, p := range prevs.Data {
							p.Type = tokenization.TokenTypeEmpty
						}
					}
				}
			}
		default:
			if t := tokens.Peek(); t != nil && t.Type == tokenization.TokenTypeText {
				t.InputEndIndex++
			} else {
				tokens.PushNewSingle(tokenization.TokenTypeText, i)
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

func compileTokenizeBackslashTransforms(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		var isHandled bool

		if nextEmpty := t.Next; nextEmpty != nil && nextEmpty.Type == tokenization.TokenTypeEmpty {
			if nextNextText := nextEmpty.Next; nextNextText != nil && nextNextText.Type == tokenization.TokenTypeText {
				if nextNextText.Len() > 0 {
					switch isHandled = true; tokens.Input[nextNextText.InputStartIndex] {
					case 'n':
						t.Type = tokenization.TokenTypeParagraphBound
						t.Next.Type = tokenization.TokenTypeParagraphBound
					case 'r':
						t.Type = tokenization.TokenTypeCarriageReturn
					case 't':
						t.Type = tokenization.TokenTypeTab
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

func compileGenerateHTML(options *Options, tokens *tokenization.TokenCollection) (err error) {
	tokenStack := tokenization.TokenCollectionNewEmpty()

	for _, t := range tokens.Data {
		if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
			return
		}
	}

	if tokenStack.Len() > 0 {
		for {
			if t := tokenStack.PopAsIs(); t != nil {
				t.Type = tokenization.TokenTypeText
				compileGenerateHTMLToken(options, t, tokenStack)
			} else {
				break
			}
		}
	}

	return
}

func compileGenerateHTMLToken(options *Options, t *tokenization.Token, tokenStack *tokenization.TokenCollection) (err error) {
	switch y := t.Type; y {
	case tokenization.TokenTypeEmpty:
		break
	case tokenization.TokenTypeSpace:
		t.HTML = []byte{' '}

		if mcs := options.MaxConsecutiveSpaces; mcs > 0 {
			nextSpaceTokens, foundNextSpaceTokens := t.NextUntilEndOfPotentialTypes(y)
			if foundNextSpaceTokens {
				for i, t2 := range nextSpaceTokens.Data {
					if i+1 < mcs {
						continue
					}
					t2.Type = tokenization.TokenTypeEmpty
				}
			}
		}
	case tokenization.TokenTypeParagraphBound,
		tokenization.TokenTypeDocumentBodyBound,
		tokenization.TokenTypeDocumentHeadBound,
		tokenization.TokenTypeDocumentHTMLBound:
		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeExclamation:
		if !options.EnableImages {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		squareBracketOpenToken := t.Next
		if squareBracketOpenToken == nil || squareBracketOpenToken.Type != tokenization.TokenTypeSquareBracketOpen {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		textTokens, foundTextTokens := squareBracketOpenToken.NextUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentText...,
		)
		if !foundTextTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypes([]tokenization.TokenType{
			tokenization.TokenTypeSquareBracketClose,
			tokenization.TokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentLink...,
		)
		if !foundLinkTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		finalToken := linkTokens.Get(-1).Next
		if finalToken == nil || finalToken.Type != tokenization.TokenTypeParenthesisClose {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		var linkBuff, textBuff bytes.Buffer

		for _, t2 := range linkTokens.Data {
			linkBuff.Write(t2.Bytes())
		}
		for _, t2 := range textTokens.Data {
			textBuff.Write(t2.Bytes())
		}

		linkString := linkBuff.String()
		textString := textBuff.String()

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

		t.Type = tokenization.TokenTypeImage
		t.Attributes = map[string]string{
			"alt": textString,
			"src": linkString,
		}

		squareBracketOpenToken.Type = tokenization.TokenTypeImage

		for _, t2 := range textTokens.Data {
			t2.Type = tokenization.TokenTypeEmpty
		}
		for _, t2 := range midTokens.Data {
			t2.Type = tokenization.TokenTypeEmpty
		}
		for _, t2 := range linkTokens.Data {
			t2.Type = tokenization.TokenTypeEmpty
		}

		finalToken.Type = tokenization.TokenTypeEmpty

		if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
			return
		}
	case tokenization.TokenTypeSquareBracketOpen:
		if !options.EnableLinks {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		textTokens, foundTextTokens := t.NextUntilEndOfPotentialTypes(
			tokenization.TokenTypeListLinkSegmentText...,
		)
		if !foundTextTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypes([]tokenization.TokenType{
			tokenization.TokenTypeSquareBracketClose,
			tokenization.TokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextUntilEndOfPotentialTypes(
			tokenization.TokenTypeListLinkSegmentLink...,
		)
		if !foundLinkTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		finalToken := linkTokens.Get(-1).Next
		if finalToken == nil || finalToken.Type != tokenization.TokenTypeParenthesisClose {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		var linkBuff bytes.Buffer

		for _, t2 := range linkTokens.Data {
			linkBuff.Write(t2.Bytes())
		}

		linkString := linkBuff.String()

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

		t.Type = tokenization.TokenTypeLink
		t.Attributes = map[string]string{"href": linkString}

		midTokens.Data[0].Type = tokenization.TokenTypeLink
		midTokens.Data[1].Type = tokenization.TokenTypeEmpty

		for _, t2 := range linkTokens.Data {
			t2.Type = tokenization.TokenTypeEmpty
		}

		finalToken.Type = tokenization.TokenTypeEmpty

		if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
			return
		}
	case tokenization.TokenTypeBackslash,
		tokenization.TokenTypeParenthesisOpen,
		tokenization.TokenTypeParenthesisClose,
		tokenization.TokenTypeSquareBracketClose:
		compileGenerateHTMLTokenHandleBytes(t)
	case tokenization.TokenTypeText:
		compileGenerateHTMLTokenHandleBytes(t)

		if !options.AllowHTML {
			t.HTML = []byte(html.EscapeString(string(t.HTML)))
		}
	case tokenization.TokenTypeTab:
		t.HTML = []byte{'\t'}
	case tokenization.TokenTypeCarriageReturn:
		t.HTML = []byte{}

		if _, foundNewline := t.NextUntilEndOfPotentialTypes(tokenization.TokenTypeParagraphBound); !foundNewline {
			prevs, foundPrevs := t.PrevUntilStartOfPotentialTypes(
				tokenization.TokenTypeParagraphBound,
			)
			if foundPrevs {
				for _, t2 := range prevs.Data {
					if t2.Type == tokenization.TokenTypeText {
						t2.Type = tokenization.TokenTypeEmpty
					}
				}
			}
		}
	case tokenization.TokenTypeBacktick:
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

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeUnderscoreTriple, tokenization.TokenTypeAsteriskTriple:
		if options.EnableStrongTags && options.EnableEmTags {
			err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
			break
		}

		fallthrough
	case tokenization.TokenTypeUnderscoreDouble, tokenization.TokenTypeAsteriskDouble:
		if !options.EnableStrongTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeUnderscore, tokenization.TokenTypeAsterisk:
		if !options.EnableEmTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeEqualsDouble:
		if !options.EnableMarkTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeLink:
		if !options.EnableLinks {
			return
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeImage:
		if !options.EnableImages {
			return
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeDocumentDoctype:
		t.Attributes = map[string]string{
			"html": "",
		}

		err = compileGenerateHTMLTokenHandleTagFromSingleToken(t, tokenStack, options)
	default:
		err = ErrCompileTokenTypeUnknown
	}

	return
}

func compileGenerateHTMLTokenHandleBytes(t *tokenization.Token) {
	t.HTML = t.Bytes()
}

func compileGenerateHTMLTokenHandleTagFromSingleToken(t *tokenization.Token, tokenStack *tokenization.TokenCollection, options *Options) (err error) {
	if err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options); err != nil {
		return
	}

	err = compileGenerateHTMLTokenHandleTag(t.SimpleCloneForClosingTag(), tokenStack, options)

	return
}

func compileGenerateHTMLTokenHandleTag(t *tokenization.Token, tokenStack *tokenization.TokenCollection, options *Options) (err error) {
	y := t.Type

	if y != tokenization.TokenTypeBacktick && tokenStack.ContainsType(tokenization.TokenTypeBacktick) {
		compileGenerateHTMLTokenHandleBytes(t)
		return
	}

	if t2 := tokenStack.Peek(); t2 != nil && t2.Type == y {
		if t2.Attributes == nil {
			t2.Attributes = make(map[string]string)
		}
		attributes := t2.Attributes

		data, foundData := tokenization.TokenTypeData[y]
		tags := data.Tags
		tagsLen := len(data.Tags)
		if !foundData || tagsLen < 1 {
			tags = []string{"span"}
			tagsLen = 1

			attributes["class"] = t2.Type.ClassName()
		}

		attributesSliceLen := len(attributes)
		var attributesSlice [][2]string

		if attributesSliceLen > 0 {
			attributesKeys := make([]string, 0, attributesSliceLen)
			for k, _ := range attributes {
				if k != "" {
					attributesKeys = append(attributesKeys, k)
				} else {
					attributesSliceLen--
				}
			}

			if attributesSliceLen > 1 {
				sort.Strings(attributesKeys)
			}

			attributesSlice = make([][2]string, attributesSliceLen)
			for i, k := range attributesKeys {
				attributesSlice[i] = [2]string{
					k,
					attributes[k],
				}
			}
		}

		{
			var buff, buff2 bytes.Buffer

			for i := 0; i < tagsLen; i++ {
				buff2.WriteByte('<')
				buff2.WriteString(tags[i])
				for j := 0; j < attributesSliceLen; j++ {
					buff2.WriteByte(' ')
					buff2.WriteString(attributesSlice[j][0])
					if v := attributesSlice[j][1]; v != "" {
						buff2.WriteByte('=')
						buff2.WriteByte('"')
						buff2.WriteString(attributesSlice[j][1])
						buff2.WriteByte('"')
					}
				}
				buff2.WriteByte('>')

				if !data.SelfClosing {
					buff.WriteByte('<')
					buff.WriteByte('/')
					buff.WriteString(tags[tagsLen-i-1])
					buff.WriteByte('>')
				}
			}

			t2.HTML = buff2.Bytes()
			t.HTML = buff.Bytes()
		}

		if options.CleanEmptyTags {
			t.Collection.TagPairCleanData = append(t.Collection.TagPairCleanData, [2]*tokenization.Token{t2, t})
		}

		tokenStack.PopAsIs()
	} else {
		tokenStack.PushAsIs(t)
	}

	return
}

func compileClean(tokens *tokenization.TokenCollection) {
	for _, pair := range tokens.TagPairCleanData {
		startTagToken, endTagToken := pair[0], pair[1]
		shouldClean := true

		prevs, foundPrevs := endTagToken.PrevUntilMeetToken(startTagToken)
		if foundPrevs {
			for _, p := range prevs.Data {
				switch p.Type {
				case tokenization.TokenTypeText:
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
			startTagToken.Type = tokenization.TokenTypeEmpty
			endTagToken.Type = tokenization.TokenTypeEmpty

			for _, p := range prevs.Data {
				p.Type = tokenization.TokenTypeEmpty
			}
		}
	}
}
