package slimdown

import (
	"bytes"
	"html"
	"html/template"
	"net/url"
	"sort"

	"github.com/theTardigrade/golang-slimdown/internal/debug"
	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

func CompileStringDefault(input string) (output template.HTML, err error) {
	return CompileString(input, nil)
}

func CompileDefault(input []byte) (output template.HTML, err error) {
	return Compile(input, nil)
}

func CompileString(input string, options *Options) (output template.HTML, err error) {
	return Compile([]byte(input), options)
}

func Compile(input []byte, options *Options) (output template.HTML, err error) {
	tokens := tokenization.TokenCollectionNew(input)

	if options == nil || options == &DefaultOptions {
		options = DefaultOptions.clone()
	}

	err = compileTokenize(options, tokens)
	if err != nil {
		return
	}

	if options.DebugPrintTokens {
		debug.PrintTokens(tokens)
	}

	err = compileGenerateHTML(options, tokens)
	if err != nil {
		return
	}

	if options.CleanEmptyTags {
		compileClean(tokens)
	}

	output = tokens.HTML()

	if options.DebugPrintOutput {
		debug.PrintOutput(output)
	}

	return
}

func compileTokenize(options *Options, tokens *tokenization.TokenCollection) (err error) {
	var backslashTokens *tokenization.TokenCollection
	if options.EnableBackslashTransforms {
		backslashTokens = tokenization.TokenCollectionNew(tokens.Input)
	}

	var hyphenTokens *tokenization.TokenCollection
	if options.EnableHyphenTransforms {
		hyphenTokens = tokenization.TokenCollectionNewEmpty()
	}

	var headingTokens *tokenization.TokenCollection
	if options.EnableHeadings {
		headingTokens = tokenization.TokenCollectionNewEmpty()
	}

	var linkTokens *tokenization.TokenCollection
	if options.EnableLinks {
		linkTokens = tokenization.TokenCollectionNewEmpty()
	}

	var imageTokens *tokenization.TokenCollection
	if options.EnableImages {
		imageTokens = tokenization.TokenCollectionNewEmpty()
	}

	if options.EnableDocumentTags {
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentDoctype)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)

		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
	}

	defer tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)
	tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)

	for i, b := range tokens.Input {
		switch b {
		// TODO: add em and en dashes
		case 138: // SPA_HAR
			var match bool

			if t := tokens.Peek(); t != nil && t.Type == tokenization.TokenTypeText {
				if l := t.Len(); l >= 2 {
					b1 := tokens.Input[t.InputEndIndex-1]
					b2 := tokens.Input[t.InputEndIndex-2]
					if b1 == 128 && b2 == 226 {
						t.InputEndIndex -= 2
						if l == 2 {
							t.Type = tokenization.TokenTypeEmpty
						}
						match = true
					}
				}

				if match {
					tokens.PushNewSingle(tokenization.TokenTypeSpaceHair, i)
				} else {
					t.InputEndIndex++
					match = true
				}
			}

			if !match {
				tokens.PushNewSingle(tokenization.TokenTypeText, i)
			}
		case '*':
			var match bool

			if t := tokens.Peek(); t != nil {
				switch match = true; t.Type {
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
			var match bool

			if t := tokens.Peek(); t != nil {
				switch match = true; t.Type {
				case tokenization.TokenTypeUnderscoreDouble:
					t.Type = tokenization.TokenTypeUnderscoreTriple
				case tokenization.TokenTypeUnderscore:
					t.Type = tokenization.TokenTypeUnderscoreDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				tokens.PushNewSingle(tokenization.TokenTypeUnderscore, i)
			}
		case '-':
			var match bool

			if t := tokens.Peek(); t != nil {
				switch match = true; t.Type {
				case tokenization.TokenTypeHyphenDouble:
					t.Type = tokenization.TokenTypeHyphenTriple
				case tokenization.TokenTypeHyphen:
					t.Type = tokenization.TokenTypeHyphenDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				hyphenTokens.PushAsIs(
					tokens.PushNewSingle(tokenization.TokenTypeHyphen, i),
				)
			}
		case '\\':
			t := tokens.PushNewSingle(tokenization.TokenTypeBackslash, i)

			if options.EnableBackslashTransforms {
				backslashTokens.PushAsIs(t)
			}
		case '#':
			var match bool

			if t := tokens.Peek(); t != nil {
				switch match = true; t.Type {
				case tokenization.TokenTypeHashQuintuple:
					t.Type = tokenization.TokenTypeHashSextuple
				case tokenization.TokenTypeHashQuadruple:
					t.Type = tokenization.TokenTypeHashQuintuple
				case tokenization.TokenTypeHashTriple:
					t.Type = tokenization.TokenTypeHashQuadruple
				case tokenization.TokenTypeHashDouble:
					t.Type = tokenization.TokenTypeHashTriple
				case tokenization.TokenTypeHash:
					t.Type = tokenization.TokenTypeHashDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				headingTokens.PushAsIs(
					tokens.PushNewSingle(tokenization.TokenTypeHash, i),
				)
			}
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
			var match bool

			if t := tokens.Peek(); t != nil {
				switch match = true; t.Type {
				case tokenization.TokenTypeBacktick:
					t.Type = tokenization.TokenTypeBacktickDouble
				default:
					match = false
				}

				if match {
					t.InputEndIndex++
				}
			}

			if !match {
				tokens.PushNewSingle(tokenization.TokenTypeBacktick, i)
			}
		case '!':
			t := tokens.PushNewSingle(tokenization.TokenTypeExclamation, i)

			if options.EnableImages {
				imageTokens.PushAsIs(t)
			}
		case '\r':
			tokens.PushNewSingle(tokenization.TokenTypeCarriageReturn, i)
		case '\n':
			compileTokenizeTransformNewLineBreak(
				tokens.PushNewSingle(tokenization.TokenTypeLineBreak, i),
			)
		case '\t':
			if tts := options.TabsToSpaces; tts > 0 {
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
			t := tokens.PushNewSingle(tokenization.TokenTypeSquareBracketOpen, i)

			if options.EnableLinks {
				linkTokens.PushAsIs(t)
			}
		case ']':
			tokens.PushNewSingle(tokenization.TokenTypeSquareBracketClose, i)
		case '<':
			t := tokens.PushNewSingle(tokenization.TokenTypeAngleBracketOpen, i)

			if options.EnableLinks {
				linkTokens.PushAsIs(t)
			}
		case '>':
			t := tokens.PushNewSingle(tokenization.TokenTypeAngleBracketClose, i)

			if options.EnableHeadings {
				headingTokens.PushAsIs(t)
			}
		case ' ':
			t := tokens.PushNewSingle(tokenization.TokenTypeSpace, i)

			if stt := options.SpacesToTabs; stt > 0 {
				if stt == 1 {
					t.Type = tokenization.TokenTypeTab
				} else {
					potentialTypes := make([]tokenization.TokenType, stt-1)

					for j := range potentialTypes {
						potentialTypes[j] = tokenization.TokenTypeSpace
					}

					if prevs, foundPrevs := t.PrevNTypesCollection(potentialTypes); foundPrevs {
						t.Type = tokenization.TokenTypeTab

						prevs.SetAllTokenTypesToEmpty()
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

	if options.EnableImages && imageTokens.Len() > 0 {
		if err = compileTokenizeImages(imageTokens); err != nil {
			return
		}
	}

	if options.EnableLinks && linkTokens.Len() > 0 {
		if err = compileTokenizeLinks(linkTokens); err != nil {
			return
		}
	}

	if options.EnableHeadings && headingTokens.Len() > 0 {
		if err = compileTokenizeHeadings(headingTokens); err != nil {
			return
		}
	}

	if options.EnableHyphenTransforms && hyphenTokens.Len() > 0 {
		if err = compileTokenizeHyphenTransforms(hyphenTokens); err != nil {
			return
		}
	}

	if options.EnableBackslashTransforms && backslashTokens.Len() > 0 {
		if err = compileTokenizeBackslashTransforms(backslashTokens); err != nil {
			return
		}
	}

	return
}

func compileTokenizeTransformNewLineBreak(t *tokenization.Token) {
	if prev := t.Prev(); prev != nil {
		if prev.Type == tokenization.TokenTypeCarriageReturn {
			prev.Type = tokenization.TokenTypeEmpty
			prev = t.Prev()
		}

		if prev != nil {
			if prev.Type == tokenization.TokenTypeLineBreak {
				prev.Type = tokenization.TokenTypeParagraphBound
				t.Type = tokenization.TokenTypeParagraphBound
			}
		}
	}
}

func compileTokenizeImages(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		if t.Type != tokenization.TokenTypeExclamation {
			continue
		}

		squareBracketOpenToken := t.Next()
		if squareBracketOpenToken == nil || squareBracketOpenToken.Type != tokenization.TokenTypeSquareBracketOpen {
			continue
		}

		textTokens, foundTextTokens := squareBracketOpenToken.NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentText...,
		)
		if !foundTextTokens {
			continue
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypesCollection([]tokenization.TokenType{
			tokenization.TokenTypeSquareBracketClose,
			tokenization.TokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			continue
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentLink...,
		)
		if !foundLinkTokens {
			continue
		}

		lastLinkToken := linkTokens.Get(-1)
		var finalToken *tokenization.Token
		var titleTokens *tokenization.TokenCollection

		spaceTokens, foundSpaceTokens := lastLinkToken.NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeSpace,
		)
		if foundSpaceTokens {
			var foundTitleTokens bool
			titleTokens, foundTitleTokens = spaceTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
				tokenization.TokenTypeListImageSegmentTitle...,
			)
			if !foundTitleTokens {
				continue
			}

			finalToken = titleTokens.Get(-1).Next()
		} else {
			finalToken = lastLinkToken.Next()
		}

		if finalToken == nil || finalToken.Type != tokenization.TokenTypeParenthesisClose {
			continue
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
			continue
		}

		linkString = linkURL.String()
		textString = html.EscapeString(textString)

		t.Type = tokenization.TokenTypeImage
		t.Attributes = map[string]string{
			"alt": textString,
			"src": linkString,
		}

		if foundSpaceTokens {
			var titleBuff bytes.Buffer

			for _, t2 := range titleTokens.Data {
				titleBuff.Write(t2.Bytes())
			}

			titleString := titleBuff.String()

			t.Attributes["title"] = titleString

			spaceTokens.SetAllTokenTypesToEmpty()
			titleTokens.SetAllTokenTypesToEmpty()
		}

		textTokens.SetAllTokenTypesToEmpty()
		midTokens.SetAllTokenTypesToEmpty()
		linkTokens.SetAllTokenTypesToEmpty()

		squareBracketOpenToken.Type = tokenization.TokenTypeImage

		finalToken.Type = tokenization.TokenTypeEmpty
	}

	return
}

func compileTokenizeLinks(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		var textTokens, midTokens, linkTokens, spaceTokens, titleTokens *tokenization.TokenCollection
		var foundTextTokens, foundMidTokens, foundLinkTokens, foundSpaceTokens bool
		var finalToken *tokenization.Token
		var expectedFinalTokenType tokenization.TokenType

		switch t.Type {
		case tokenization.TokenTypeSquareBracketOpen:
			textTokens, foundTextTokens = t.NextsCollectionUntilEndOfPotentialTypes(
				tokenization.TokenTypeListLinkSegmentText...,
			)
			if !foundTextTokens {
				continue
			}

			midTokens, foundMidTokens = textTokens.Get(-1).NextNTypesCollection([]tokenization.TokenType{
				tokenization.TokenTypeSquareBracketClose,
				tokenization.TokenTypeParenthesisOpen,
			})
			if !foundMidTokens {
				continue
			}

			linkTokens, foundLinkTokens = midTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
				tokenization.TokenTypeListLinkSegmentLink...,
			)
			if !foundLinkTokens {
				continue
			}

			expectedFinalTokenType = tokenization.TokenTypeParenthesisClose
		case tokenization.TokenTypeAngleBracketOpen:
			linkTokens, foundLinkTokens = t.NextsCollectionUntilEndOfPotentialTypes(
				tokenization.TokenTypeListLinkSegmentLink...,
			)
			if !foundLinkTokens {
				continue
			}

			expectedFinalTokenType = tokenization.TokenTypeAngleBracketClose
		default:
			continue
		}

		{
			lastLinkToken := linkTokens.Get(-1)
			spaceTokens, foundSpaceTokens = lastLinkToken.NextsCollectionUntilEndOfPotentialTypes(
				tokenization.TokenTypeSpace,
			)
			if foundSpaceTokens {
				var foundTitleTokens bool
				titleTokens, foundTitleTokens = spaceTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
					tokenization.TokenTypeListLinkSegmentTitle...,
				)
				if !foundTitleTokens {
					continue
				}

				finalToken = titleTokens.Get(-1).Next()
			} else {
				finalToken = lastLinkToken.Next()
			}

			if finalToken == nil || finalToken.Type != expectedFinalTokenType {
				continue
			}
		}

		var linkBuff bytes.Buffer

		for _, t2 := range linkTokens.Data {
			linkBuff.Write(t2.Bytes())
		}

		linkString := linkBuff.String()

		linkURL, err2 := url.Parse(linkString)
		if err2 != nil {
			continue
		}
		linkString = linkURL.String()

		t.Type = tokenization.TokenTypeLink
		t.Attributes = map[string]string{"href": linkString}

		if foundSpaceTokens {
			var titleBuff bytes.Buffer

			for _, t2 := range titleTokens.Data {
				titleBuff.Write(t2.Bytes())
			}

			titleString := titleBuff.String()

			t.Attributes["title"] = titleString

			spaceTokens.SetAllTokenTypesToEmpty()
			titleTokens.SetAllTokenTypesToEmpty()
		}

		if foundTextTokens {
			textTokens.SetAllTokenTypesToEmpty()
		}

		if foundMidTokens {
			midTokens.SetAllTokenTypesToEmpty()
		}

		finalToken.Type = tokenization.TokenTypeLink
	}

	return
}

func compileTokenizeHeadings(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		prevBound := t.Prev()
		if prevBound == nil || prevBound.Type != tokenization.TokenTypeParagraphBound {
			continue
		}

		nextSpace := t.Next()
		if nextSpace == nil || nextSpace.Type != tokenization.TokenTypeSpace {
			continue
		}

		nextBound := t.NextOfType(tokenization.TokenTypeParagraphBound)
		if nextBound == nil {
			continue
		}

		var tt tokenization.TokenType

		switch t.Type {
		case tokenization.TokenTypeHash:
			tt = tokenization.TokenTypeHeading1Bound
		case tokenization.TokenTypeHashDouble:
			tt = tokenization.TokenTypeHeading2Bound
		case tokenization.TokenTypeHashTriple:
			tt = tokenization.TokenTypeHeading3Bound
		case tokenization.TokenTypeHashQuadruple:
			tt = tokenization.TokenTypeHeading4Bound
		case tokenization.TokenTypeHashQuintuple:
			tt = tokenization.TokenTypeHeading5Bound
		case tokenization.TokenTypeHashSextuple:
			tt = tokenization.TokenTypeHeading6Bound
		case tokenization.TokenTypeAngleBracketClose:
			tt = tokenization.TokenTypeBlockquoteBound
		default:
			continue
		}

		prevBound.Type = tt
		nextBound.Type = tt

		tt = tokenization.TokenTypeEmpty

		nextSpace.Type = tt
		t.Type = tt
	}

	return
}

func compileTokenizeHyphenTransforms(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		if next := t.Next(); next != nil && next.Type == tokenization.TokenTypeSpace {
			next.Type = tokenization.TokenTypeSpaceHair
		}

		if prev := t.Prev(); prev != nil && prev.Type == tokenization.TokenTypeSpace {
			prev.Type = tokenization.TokenTypeSpaceHair
		}

		switch t.Type {
		case tokenization.TokenTypeHyphenTriple:
			t.Type = tokenization.TokenTypeDashEm
		case tokenization.TokenTypeHyphenDouble:
			t.Type = tokenization.TokenTypeDashEn
		}
	}

	return
}

func compileTokenizeBackslashTransforms(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		var isHandled bool

		if nextText := t.Next(); nextText != nil && nextText.Type == tokenization.TokenTypeText {
			if nextText.Len() > 0 {
				switch isHandled = true; tokens.Input[nextText.InputStartIndex] {
				case 'n':
					t.Type = tokenization.TokenTypeLineBreak
					compileTokenizeTransformNewLineBreak(t)
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
					if nextText.InputStartIndex++; nextText.Len() == 0 {
						nextText.Type = tokenization.TokenTypeEmpty
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
				if err = compileGenerateHTMLToken(options, t, nil); err != nil {
					return
				}
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
	case tokenization.TokenTypeText:
		compileGenerateHTMLTokenHandleBytes(t)

		if !options.AllowHTML {
			t.HTML = []byte(html.EscapeString(string(t.HTML)))
		}
	case tokenization.TokenTypeSpace:
		t.HTML = []byte{' '}

		if mcs := options.MaxConsecutiveSpaces; mcs > 0 {
			nextSpaceTokens, foundNextSpaceTokens := t.NextsCollectionUntilEndOfPotentialTypes(y)
			if foundNextSpaceTokens {
				for i, t2 := range nextSpaceTokens.Data {
					if i+1 < mcs {
						continue
					}
					t2.Type = tokenization.TokenTypeEmpty
				}
			}
		}
	case tokenization.TokenTypeDocumentBodyBound,
		tokenization.TokenTypeDocumentHeadBound,
		tokenization.TokenTypeDocumentHTMLBound,
		tokenization.TokenTypeHeading1Bound,
		tokenization.TokenTypeHeading2Bound,
		tokenization.TokenTypeHeading3Bound,
		tokenization.TokenTypeHeading4Bound,
		tokenization.TokenTypeHeading5Bound,
		tokenization.TokenTypeHeading6Bound,
		tokenization.TokenTypeBlockquoteBound:
		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeParagraphBound:
		if options.EnableHorizontalRules {
			if next := t.Next(); next != nil {
				if next.Type == tokenization.TokenTypeAsteriskTriple ||
					next.Type == tokenization.TokenTypeUnderscoreTriple {
					if nextNext := next.Next(); nextNext != nil {
						if nextNext.Type == tokenization.TokenTypeParagraphBound {
							t.Type = tokenization.TokenTypeEmpty
							next.Type = tokenization.TokenTypeHorizontalRule
							nextNext.Type = tokenization.TokenTypeEmpty
							break
						}
					}
				}
			}
		}

		if !options.EnableParagraphs {
			if prev := t.Prev(); prev != nil {
				if next := t.Next(); next != nil {
					t.HTML = []byte{'\n'}
					break
				}
			}
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeLineBreak:
		if !options.EnableParagraphs {
			t.HTML = []byte{'\n'}
			break
		}

		err = compileGenerateHTMLTokenHandleTagFromSingleToken(t, tokenStack, options)
	case tokenization.TokenTypeBackslash,
		tokenization.TokenTypeParenthesisOpen,
		tokenization.TokenTypeParenthesisClose,
		tokenization.TokenTypeSquareBracketOpen,
		tokenization.TokenTypeSquareBracketClose,
		tokenization.TokenTypeExclamation,
		tokenization.TokenTypeHash,
		tokenization.TokenTypeHashDouble,
		tokenization.TokenTypeHashTriple,
		tokenization.TokenTypeHashQuadruple,
		tokenization.TokenTypeHashQuintuple,
		tokenization.TokenTypeHashSextuple,
		tokenization.TokenTypeHyphen,
		tokenization.TokenTypeHyphenDouble,
		tokenization.TokenTypeHyphenTriple:
		compileGenerateHTMLTokenHandleBytes(t)
	case tokenization.TokenTypeAngleBracketOpen:
		t.HTML = []byte{'&', 'l', 't', ';'}
	case tokenization.TokenTypeAngleBracketClose:
		t.HTML = []byte{'&', 'g', 't', ';'}
	case tokenization.TokenTypeDashEm:
		t.HTML = []byte("—")
	case tokenization.TokenTypeDashEn:
		t.HTML = []byte("–")
	case tokenization.TokenTypeTab:
		t.HTML = []byte{'\t'}
	case tokenization.TokenTypeCarriageReturn:
		t.HTML = []byte{}

		if _, foundNewline := t.NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeParagraphBound,
			tokenization.TokenTypeLineBreak,
		); !foundNewline {
			prevs, foundPrevs := t.PrevsCollectionUntilStartOfPotentialTypes(
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
	case tokenization.TokenTypeBacktickDouble:
		var insideCode bool

		if options.EnableCodeTags {
			d := tokenStack.Data
			for i := len(d) - 1; i >= 0; i-- {
				if t2 := d[i]; t2.Type == tokenization.TokenTypeBacktick {
					insideCode = true
					break
				}
			}
		}

		t.HTML = []byte{'`'}

		if !insideCode {
			t.HTML = append(t.HTML, '`')
		}
	case tokenization.TokenTypeBacktick:
		if !options.EnableCodeTags {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if next := t.Next(); next != nil && next.Type == y {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if prev := t.Prev(); prev != nil && prev.Type == y {
			if prevPrev := prev.Prev(); prevPrev == nil || prevPrev.Type != y {
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
	case tokenization.TokenTypeSpaceHair:
		t.HTML = []byte(" ")
	case tokenization.TokenTypeHorizontalRule:
		if !options.EnableHorizontalRules {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		err = compileGenerateHTMLTokenHandleTagFromSingleToken(t, tokenStack, options)
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

	datum := t.TypeDatum()

	if t2 := tokenStack.Peek(); t2 != nil && t2.Type == y {
		if t2.Attributes == nil {
			t2.Attributes = make(map[string]string)
		}
		attributes := t2.Attributes

		var tags []string
		var tagsLen int

		if datum == nil {
			tags = []string{"span"}
			tagsLen = 1
			attributes["class"] = t2.Type.ClassName()
		} else {
			tags = append(tags, datum.Tags...)
			tagsLen = len(tags)
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
					a := attributesSlice[j]
					if k := a[0]; k != "" {
						buff2.WriteByte(' ')
						buff2.WriteString(a[0])
						if v := a[1]; v != "" {
							buff2.WriteByte('=')
							buff2.WriteByte('"')
							buff2.WriteString(v)
							buff2.WriteByte('"')
						}
					}
				}
				buff2.WriteByte('>')

				if !datum.SelfClosing {
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

		if datum := startTagToken.TypeDatum(); datum != nil && datum.SelfClosing {
			continue
		}

		shouldClean := true

		prevs, foundPrevs := endTagToken.PrevsCollectionUntilMeetToken(startTagToken)
		if foundPrevs {
			for _, p := range prevs.Data {
				switch p.Type {
				case tokenization.TokenTypeText:
					if p.Len() > 0 {
						shouldClean = false
					}
				}

				if !shouldClean {
					break
				}
			}
		}

		if shouldClean {
			tt := tokenization.TokenTypeEmpty

			startTagToken.Type = tt
			endTagToken.Type = tt

			prevs.SetAllTokenTypes(tt)
		}
	}
}
