package tokenization

type Token struct {
	RawPrev         *Token
	RawNext         *Token
	ListCollection  *TokenListCollection
	Attributes      map[string]string
	HTML            []byte
	InputStartIndex int
	InputEndIndex   int
	Type            TokenType
	Indent          int
}

func tokenNew(y TokenType, inputStartIndex int, inputEndIndex int) *Token {
	return &Token{
		Type:            y,
		InputStartIndex: inputStartIndex,
		InputEndIndex:   inputEndIndex,
	}
}

func (t *Token) SimpleCloneForClosingTag() *Token {
	return &Token{
		ListCollection: t.ListCollection,
		Type:           t.Type,
	}
}

func (t *Token) TypeDatum() *TokenTypeDatum {
	if d, ok := tokenTypeData[t.Type]; ok {
		return &d
	}

	return nil
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

func (t *Token) PrevOfType(y TokenType) *Token {
	if t == nil {
		return nil
	}

	t2 := t

	for {
		if t2 = t2.RawPrev; t2 == nil {
			return nil
		}

		if t2.Type == y {
			return t2
		}
	}
}

func (t *Token) NextOfType(y TokenType) *Token {
	if t == nil {
		return nil
	}

	t2 := t

	for {
		if t2 = t2.RawNext; t2 == nil {
			return nil
		}

		if t2.Type == y {
			return t2
		}
	}
}

func (t *Token) PrevNCollection(n int) (prevs *TokenSliceCollection, foundAll bool) {
	if t == nil {
		return
	}

	prevs = TokenSliceCollectionNew()
	t2 := t
	var i int

	for {
		if i >= n {
			break
		}

		if t2 = t2.Prev(); t2 == nil {
			return
		}

		prevs.Push(t2)

		i++
	}

	if prevs.Len() == n {
		foundAll = true
	}

	return
}

func (t *Token) NextNCollection(n int) (nexts *TokenSliceCollection, foundAll bool) {
	if t == nil {
		return
	}

	nexts = TokenSliceCollectionNew()
	t2 := t
	var i int

	for {
		if i >= n {
			break
		}

		if t2 = t2.Next(); t2 == nil {
			return
		}

		nexts.Push(t2)

		i++
	}

	if nexts.Len() == n {
		foundAll = true
	}

	return
}

func (t *Token) PrevNTypesCollection(types []TokenType) (prevs *TokenSliceCollection, foundAll bool) {
	if t == nil {
		return
	}

	prevs = TokenSliceCollectionNew()
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

		prevs.Push(t2)

		i++
	}

	if prevs.Len() == l {
		foundAll = true
	}

	return
}

func (t *Token) NextNTypesCollection(types []TokenType) (nexts *TokenSliceCollection, foundAll bool) {
	if t == nil {
		return
	}

	nexts = TokenSliceCollectionNew()
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

		nexts.Push(t2)

		i++
	}

	if nexts.Len() == l {
		foundAll = true
	}

	return
}

func (t *Token) PrevsCollectionUntilMeetToken(stopToken *Token) (prevs *TokenSliceCollection, foundAny bool) {
	prevs = TokenSliceCollectionNew()
	t2 := t

	for {
		if t2 = t2.Prev(); t2 == nil || t2 == stopToken {
			break
		}

		prevs.Push(t2)
	}

	if prevs.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) NextsCollectionUntilMeetType(y TokenType) (nexts *TokenSliceCollection, foundAny bool) {
	nexts = TokenSliceCollectionNew()
	t2 := t

	for {
		if t2 = t2.Next(); t2 == nil || t2.Type == y {
			break
		}

		nexts.Push(t2)
	}

	if nexts.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) PrevsCollectionUntilStartOfPotentialTypes(potentialTypes ...TokenType) (prevs *TokenSliceCollection, foundAny bool) {
	prevs = TokenSliceCollectionNew()
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

		prevs.Push(t2)
	}

	if prevs.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) PrevsCollectionUntilEndOfPotentialTypes(potentialTypes ...TokenType) (prevs *TokenSliceCollection, foundAny bool) {
	prevs = TokenSliceCollectionNew()
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

		prevs.Push(t2)
	}

	if prevs.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) NextsCollectionUntilStartOfPotentialTypes(potentialTypes ...TokenType) (nexts *TokenSliceCollection, foundAny bool) {
	nexts = TokenSliceCollectionNew()
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

		nexts.Push(t2)
	}

	if nexts.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) NextsCollectionUntilEndOfPotentialTypes(potentialTypes ...TokenType) (nexts *TokenSliceCollection, foundAny bool) {
	nexts = TokenSliceCollectionNew()
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

		nexts.Push(t2)
	}

	if nexts.Len() > 0 {
		foundAny = true
	}

	return
}

func (t *Token) Len() int {
	if t.Type == TokenTypeEmpty {
		return 0
	}

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
	if t.Type == TokenTypeEmpty || t.ListCollection == nil || t.Len() <= 0 {
		return []byte{}
	}

	return t.ListCollection.Input[t.InputStartIndex:t.InputEndIndex]
}

func (t *Token) String() string {
	return string(t.Bytes())
}
