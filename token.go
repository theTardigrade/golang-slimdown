package slimdown

type token struct {
	Prev, Next      *token
	Collection      *tokenCollection
	Attributes      map[string]string
	HTML            []byte
	InputStartIndex int
	InputEndIndex   int
	Type            tokenType
}

func (t *token) SimpleCloneForClosingTag() *token {
	return &token{
		Type: t.Type,
	}
}

func (t *token) PrevNTypes(types []tokenType) (prevs *tokenCollection, found bool) {
	if t == nil {
		return
	}

	prevs = tokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i, j int

	for {
		if i+j >= l {
			break
		}

		if t2 = t2.Prev; t2 == nil || t2.Type != types[i] {
			return
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) NextNTypes(types []tokenType) (nexts *tokenCollection, found bool) {
	if t == nil {
		return
	}

	nexts = tokenCollectionNewEmpty()
	t2 := t
	l := len(types)
	var i, j int

	for {
		if i+j >= l {
			break
		}

		if t2 = t2.Next; t2 == nil || t2.Type != types[i] {
			return
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) PrevUntilMeetToken(stopToken *token) (prevs *tokenCollection, found bool) {
	prevs = tokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev; t2 == nil || t2 == stopToken {
			break
		}

		if t2.Type != tokenTypeEmpty {
			prevs.PushAsIs(t2)
		}
	}

	if prevs.Len() > 0 {
		found = true
	}

	return
}

func (t *token) PrevUntilStartOfPotentialTypes(possibleTypes ...tokenType) (prevs *tokenCollection, found bool) {
	prevs = tokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev; t2 == nil {
			break
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) PrevUntilEndOfPotentialTypes(possibleTypes ...tokenType) (prevs *tokenCollection, found bool) {
	prevs = tokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Prev; t2 == nil {
			break
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) NextUntilStartOfPotentialTypes(possibleTypes ...tokenType) (nexts *tokenCollection, found bool) {
	nexts = tokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Next; t2 == nil {
			break
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) NextUntilEndOfPotentialTypes(possibleTypes ...tokenType) (nexts *tokenCollection, found bool) {
	nexts = tokenCollectionNewEmpty()
	t2 := t

	for {
		if t2 = t2.Next; t2 == nil {
			break
		}

		if t2.Type == tokenTypeEmpty {
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

func (t *token) Len() int {
	if l := len(t.HTML); l > 0 {
		return l
	}

	e, s := t.InputEndIndex, t.InputStartIndex

	if e <= s {
		return 0
	}

	return e - s
}

func (t *token) Bytes() []byte {
	if t.Type == tokenTypeEmpty || t.Collection == nil || t.Len() <= 0 {
		return []byte{}
	}

	return t.Collection.Input[t.InputStartIndex:t.InputEndIndex]
}

func (t *token) String() string {
	return string(t.Bytes())
}
