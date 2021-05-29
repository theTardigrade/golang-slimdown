package slimdown

import (
	"bytes"
	"html"
	"html/template"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

const (
	compileUseConcurrencyTokensMinLen = 1 << 12
)

func CompileString(input string, options *Options) (output template.HTML, err error) {
	return Compile([]byte(input), options)
}

func Compile(input []byte, options *Options) (output template.HTML, err error) {
	tokens := tokenization.TokenCollectionNew(input)

	if options == nil || options == &DefaultOptions {
		options = DefaultOptions.Clone()
	}

	err = compileTokenize(options, tokens)
	if err != nil {
		return
	}

	if options.DebugPrintTokens {
		debugPrintTokens(tokens)
	}

	if options.UseConcurrency && tokens.Len() < compileUseConcurrencyTokensMinLen {
		options = options.Clone()
		options.UseConcurrency = false
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
	hashTokens := tokenization.TokenCollectionNew(tokens.Input)

	if options.EnableDocumentTags {
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentDoctype)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHTMLBound)

		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentHeadBound)

		defer tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
		tokens.PushNewEmpty(tokenization.TokenTypeDocumentBodyBound)
	}

	if options.EnableParagraphs {
		defer tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)
		tokens.PushNewEmpty(tokenization.TokenTypeParagraphBound)
	}

	for i, b := range tokens.Input {
		switch b {
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
				tokens.PushNewSingle(tokenization.TokenTypeHyphen, i)
			}
		case '\\':
			backslashTokens.PushAsIs(
				tokens.PushNewSingle(tokenization.TokenTypeBackslash, i),
			)
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
				hashTokens.PushAsIs(
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
			tokens.PushNewSingle(tokenization.TokenTypeBacktick, i)
		case '!':
			tokens.PushNewSingle(tokenization.TokenTypeExclamation, i)
		case '\r':
			tokens.PushNewSingle(tokenization.TokenTypeCarriageReturn, i)
		case '\n':
			compileTokenizeTransformNewLineBreak(
				tokens.PushNewSingle(tokenization.TokenTypeLineBreak, i),
			)
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

					for j := range potentialTypes {
						potentialTypes[j] = tokenization.TokenTypeSpace
					}

					if prevs, foundPrevs := t.PrevNTypesCollection(potentialTypes); foundPrevs {
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

	if options.EnableHeadings && hashTokens.Len() > 0 {
		if err = compileTokenizeHashHeadings(hashTokens); err != nil {
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

func compileTokenizeHashHeadings(tokens *tokenization.TokenCollection) (err error) {
	for _, t := range tokens.Data {
		prevBound := t.Prev()

		if prevBound == nil || prevBound.Type != tokenization.TokenTypeParagraphBound {
			continue
		}

		nextBound := t.NextOfType(tokenization.TokenTypeParagraphBound)

		if nextBound == nil {
			continue
		}

		nextSpace := t.Next()

		if nextSpace == nil || nextSpace.Type != tokenization.TokenTypeSpace {
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
		default:
			return ErrCompileTokenTypeUnknown
		}

		prevBound.Type = tt
		nextBound.Type = tt

		tt = tokenization.TokenTypeEmpty

		nextSpace.Type = tt
		t.Type = tt
	}

	return nil
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

	if !options.UseConcurrency {
		for _, t := range tokens.Data {
			if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
				return
			}
		}
	} else {
		var wg sync.WaitGroup
		errChan := make(chan error)

		go func(errChan chan<- error) {
			for _, t := range tokens.Data {
				for _, y := range tokenization.TokenTypeCompileGenerateHTMLUseConcurrency {
					if y == t.Type {
						wg.Add(1)

						go func(t *tokenization.Token) {
							wg.Done()

							if err := compileGenerateHTMLToken(options, t, nil); err != nil {
								select {
								case errChan <- err:
								default:
								}
							}
						}(t)

						break
					}
				}
			}
		}(errChan)

		for _, t := range tokens.Data {
			var match bool

			for _, y := range tokenization.TokenTypeCompileGenerateHTMLUseConcurrency {
				if y == t.Type {
					match = true
					break
				}
			}

			if !match {
				if err = compileGenerateHTMLToken(options, t, tokenStack); err != nil {
					return
				}
			}
		}

		wg.Wait()

		select {
		case err = <-errChan:
			return
		default:
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
		tokenization.TokenTypeHeading2Bound:
		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeParagraphBound:
		if !options.EnableParagraphs {
			t.HTML = []byte{'\n'}
			break
		}

		err = compileGenerateHTMLTokenHandleTag(t, tokenStack, options)
	case tokenization.TokenTypeLineBreak:
		if !options.EnableParagraphs {
			t.HTML = []byte{'\n'}
			break
		}

		err = compileGenerateHTMLTokenHandleTagFromSingleToken(t, tokenStack, options)
	case tokenization.TokenTypeExclamation:
		if !options.EnableImages {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		squareBracketOpenToken := t.Next()
		if squareBracketOpenToken == nil || squareBracketOpenToken.Type != tokenization.TokenTypeSquareBracketOpen {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		textTokens, foundTextTokens := squareBracketOpenToken.NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentText...,
		)
		if !foundTextTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypesCollection([]tokenization.TokenType{
			tokenization.TokenTypeSquareBracketClose,
			tokenization.TokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListImageSegmentLink...,
		)
		if !foundLinkTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		finalToken := linkTokens.Get(-1).Next()
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
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkString = linkURL.String()
		if strings.Contains(linkString, `"`) {
			compileGenerateHTMLTokenHandleBytes(t)
			break
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

		textTokens, foundTextTokens := t.NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListLinkSegmentText...,
		)
		if !foundTextTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		midTokens, foundMidTokens := textTokens.Get(-1).NextNTypesCollection([]tokenization.TokenType{
			tokenization.TokenTypeSquareBracketClose,
			tokenization.TokenTypeParenthesisOpen,
		})
		if !foundMidTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkTokens, foundLinkTokens := midTokens.Get(-1).NextsCollectionUntilEndOfPotentialTypes(
			tokenization.TokenTypeListLinkSegmentLink...,
		)
		if !foundLinkTokens {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		finalToken := linkTokens.Get(-1).Next()
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
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		linkString = linkURL.String()
		if strings.Contains(linkString, `"`) {
			compileGenerateHTMLTokenHandleBytes(t)
			break
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
		tokenization.TokenTypeSquareBracketClose,
		tokenization.TokenTypeHash,
		tokenization.TokenTypeHashDouble,
		tokenization.TokenTypeHashTriple,
		tokenization.TokenTypeHashQuadruple,
		tokenization.TokenTypeHashQuintuple,
		tokenization.TokenTypeHashSextuple:
		compileGenerateHTMLTokenHandleBytes(t)
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
	case tokenization.TokenTypeHyphenTriple:
		if !options.EnableHyphenTransforms {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if next := t.Next(); next != nil && next.Type == tokenization.TokenTypeSpace {
			next.Type = tokenization.TokenTypeSpaceHair
		}

		if prev := t.Prev(); prev != nil && prev.Type == tokenization.TokenTypeSpace {
			prev.Type = tokenization.TokenTypeSpaceHair
		}

		t.HTML = []byte("—")
	case tokenization.TokenTypeHyphenDouble:
		if !options.EnableHyphenTransforms {
			compileGenerateHTMLTokenHandleBytes(t)
			break
		}

		if next := t.Next(); next != nil && next.Type == tokenization.TokenTypeSpace {
			next.Type = tokenization.TokenTypeSpaceHair
		}

		if prev := t.Prev(); prev != nil && prev.Type == tokenization.TokenTypeSpace {
			prev.Type = tokenization.TokenTypeSpaceHair
		}

		t.HTML = []byte("–")
	case tokenization.TokenTypeHyphen:
		if options.EnableHyphenTransforms {
			if next := t.Next(); next != nil && next.Type == tokenization.TokenTypeSpace {
				next.Type = tokenization.TokenTypeSpaceHair
			}

			if prev := t.Prev(); prev != nil && prev.Type == tokenization.TokenTypeSpace {
				prev.Type = tokenization.TokenTypeSpaceHair
			}
		}

		compileGenerateHTMLTokenHandleBytes(t)
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

		prevs, foundPrevs := endTagToken.PrevsCollectionUntilMeetToken(startTagToken)
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
