package tokenization

type Token struct {
	RawPrev         *Token
	RawNext         *Token
	Collection      *TokenCollection
	Attributes      map[string]string
	HTML            []byte
	InputStartIndex int
	InputEndIndex   int
	Type            TokenType
}

func (t *Token) SimpleCloneForClosingTag() *Token {
	return &Token{
		Collection: t.Collection,
		Type:       t.Type,
	}
}

func (t *Token) Prev() *Token {
	if t == nil {
		return nil
	}

	t2 := t

	for {
		if t2 = t2.RawPrev; t2 == nil {
			return nil
		}

		if t2.Type != TokenTypeEmpty {
			return t2
		}
	}
}

func (t *Token) PrevN(n int) (prevs *TokenCollection, found bool) {
	if t == nil {
		return
	}

	prevs = TokenCollectionNewEmpty()
	t2 := t
	var i, j int

	for {
		if i+j >= n {
			break
		}

		if t2 = t2.RawPrev; t2 == nil {
			return
		}

		if t2.Type == TokenTypeEmpty {
			j++
			continue
		}

		prevs.PushAsIs(t2)

		i++
	}

	if prevs.Len() == n {
		found = true
	}

	return
}

func (t *Token) PrevNTypes(types []TokenType) (prevs *TokenCollection, found bool) {
	if t == nil {
		return
	}

	prevs = TokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i, j int

	for {
		if i+j >= l {
			break
		}

		if t2 = t2.RawPrev; t2 == nil || t2.Type != types[i] {
			return
		}

		if t2.Type == TokenTypeEmpty {
			j++
			continue
		}

		prevs.PushAsIs(t2)

		i++
	}

	if prevs.Len() == l {
		found = true
	}

	return
}

func (t *Token) NextNTypes(types []TokenType) (nexts *TokenCollection, found bool) {
	if t == nil {
		return
	}

	nexts = TokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i, j int

	for {
		if i+j >= l {
			break
		}

		if t2 = t2.RawNext; t2 == nil || t2.Type != types[i] {
			return
		}

		if t2.Type == TokenTypeEmpty {
			j++
			continue
		}

		nexts.PushAsIs(t2)

		i++
	}

	if nexts.Len() == l {
		found = true
	}

	return
}

func (t *Token) PrevUntilMeetToken(stopToken *Token) (prevs *TokenCollection, found bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.RawPrev; t2 == nil || t2 == stopToken {
			break
		}

		if t2.Type != TokenTypeEmpty {
			prevs.PushAsIs(t2)
		}
	}

	if prevs.Len() > 0 {
		found = true
	}

	return
}

func (t *Token) PrevUntilStartOfPotentialTypes(possibleTypes ...TokenType) (prevs *TokenCollection, found bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.RawPrev; t2 == nil {
			break
		}

		if t2.Type == TokenTypeEmpty {
			continue
		}

		var match bool

		for _, y := range possibleTypes {
			if t2.Type == y {
				match = true
				break
			}
		}

		if match {
			break
		}

		prevs.PushAsIs(t2)
	}

	if prevs.Len() > 0 {
		found = true
	}

	return
}

func (t *Token) PrevUntilEndOfPotentialTypes(possibleTypes ...TokenType) (prevs *TokenCollection, found bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.RawPrev; t2 == nil {
			break
		}

		if t2.Type == TokenTypeEmpty {
			continue
		}

		var match bool

		for _, y := range possibleTypes {
			if t2.Type == y {
				match = true
				break
			}
		}

		if !match {
			break
		}

		prevs.PushAsIs(t2)
	}

	if prevs.Len() > 0 {
		found = true
	}

	return
}

func (t *Token) NextUntilStartOfPotentialTypes(possibleTypes ...TokenType) (nexts *TokenCollection, found bool) {
	nexts = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.RawNext; t2 == nil {
			break
		}

		if t2.Type == TokenTypeEmpty {
			continue
		}

		var match bool

		for _, y := range possibleTypes {
			if t2.Type == y {
				match = true
				break
			}
		}

		if match {
			break
		}

		nexts.PushAsIs(t2)
	}

	if nexts.Len() > 0 {
		found = true
	}

	return
}

func (t *Token) NextUntilEndOfPotentialTypes(possibleTypes ...TokenType) (nexts *TokenCollection, found bool) {
	nexts = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.RawNext; t2 == nil {
			break
		}

		if t2.Type == TokenTypeEmpty {
			continue
		}

		var match bool

		for _, y := range possibleTypes {
			if t2.Type == y {
				match = true
				break
			}
		}

		if !match {
			break
		}

		nexts.PushAsIs(t2)
	}

	if nexts.Len() > 0 {
		found = true
	}

	return
}

func (t *Token) Len() int {
	if l := len(t.HTML); l > 0 {
		return l
	}

	e, s := t.InputEndIndex, t.InputStartIndex

	if e <= s {
		return 0
	}

	return e - s
}

func (t *Token) Bytes() []byte {
	if t.Type == TokenTypeEmpty || t.Collection == nil || t.Len() <= 0 {
		return []byte{}
	}

	return t.Collection.Input[t.InputStartIndex:t.InputEndIndex]
}

func (t *Token) String() string {
	return string(t.Bytes())
}
