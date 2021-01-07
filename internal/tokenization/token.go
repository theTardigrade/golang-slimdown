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

func (t *Token) Next() *Token {
	if t == nil {
		return nil
	}

	t2 := t

	for {
		if t2 = t2.RawNext; t2 == nil {
			return nil
		}

		if t2.Type != TokenTypeEmpty {
			return t2
		}
	}
}

func (t *Token) PrevNCollection(n int) (prevs *TokenCollection, foundAll bool) {
	if t == nil {
		return
	}

	prevs = TokenCollectionNewEmpty()
	t2 := t
	var i int

	for {
		if i >= n {
			break
		}

		if t2 = t2.Prev(); t2 == nil {
			return
		}

		prevs.PushAsIs(t2)

		i++
	}

	if prevs.Len() == n {
		foundAll = true
	}

	return
}

func (t *Token) NextNCollection(n int) (nexts *TokenCollection, foundAll bool) {
	if t == nil {
		return
	}

	nexts = TokenCollectionNewEmpty()
	t2 := t
	var i int

	for {
		if i >= n {
			break
		}

		if t2 = t2.Next(); t2 == nil {
			return
		}

		nexts.PushAsIs(t2)

		i++
	}

	if nexts.Len() == n {
		foundAll = true
	}

	return
}

func (t *Token) PrevNTypesCollection(types []TokenType) (prevs *TokenCollection, foundAll bool) {
	if t == nil {
		return
	}

	prevs = TokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i int

	for {
		if i >= l {
			break
		}

		if t2 = t2.Prev(); t2 == nil || t2.Type != types[i] {
			return
		}

		prevs.PushAsIs(t2)

		i++
	}

	if prevs.Len() == l {
		foundAll = true
	}

	return
}

func (t *Token) NextNTypesCollection(types []TokenType) (nexts *TokenCollection, foundAll bool) {
	if t == nil {
		return
	}

	nexts = TokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i int

	for {
		if i >= l {
			break
		}

		if t2 = t2.Next(); t2 == nil || t2.Type != types[i] {
			return
		}

		nexts.PushAsIs(t2)

		i++
	}

	if nexts.Len() == l {
		foundAll = true
	}

	return
}

func (t *Token) PrevsCollectionUntilMeetToken(stopToken *Token) (prevs *TokenCollection, foundAny bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev(); t2 == nil || t2 == stopToken {
			break
		}

		prevs.PushAsIs(t2)
	}

	if prevs.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) PrevsCollectionUntilStartOfPotentialTypes(potentialTypes ...TokenType) (prevs *TokenCollection, foundAny bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev(); t2 == nil {
			break
		}

		var match bool

		for _, y := range potentialTypes {
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
		foundAny = true
	}

	return
}

func (t *Token) PrevsCollectionUntilEndOfPotentialTypes(potentialTypes ...TokenType) (prevs *TokenCollection, foundAny bool) {
	prevs = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev(); t2 == nil {
			break
		}

		var match bool

		for _, y := range potentialTypes {
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
		foundAny = true
	}

	return
}

func (t *Token) NextsCollectionUntilStartOfPotentialTypes(potentialTypes ...TokenType) (nexts *TokenCollection, foundAny bool) {
	nexts = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Next(); t2 == nil {
			break
		}

		var match bool

		for _, y := range potentialTypes {
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
		foundAny = true
	}

	return
}

func (t *Token) NextsCollectionUntilEndOfPotentialTypes(potentialTypes ...TokenType) (nexts *TokenCollection, foundAny bool) {
	nexts = TokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Next(); t2 == nil {
			break
		}

		var match bool

		for _, y := range potentialTypes {
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
		foundAny = true
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
