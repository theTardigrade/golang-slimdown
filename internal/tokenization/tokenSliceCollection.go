package tokenization

type TokenSliceCollection struct {
	Tokens []*Token
}

func TokenSliceCollectionNew() *TokenSliceCollection {
	return &TokenSliceCollection{}
}

func (c *TokenSliceCollection) Push(tokens ...*Token) {
	for _, t := range tokens {
		if t == nil {
			continue
		}

		c.Tokens = append(c.Tokens, t)
	}
}

func (c *TokenSliceCollection) Pop() *Token {
	if i := c.Len() - 1; i >= 0 {
		t := c.Tokens[i]

		c.Tokens = c.Tokens[:i]

		return t
	}

	return nil
}

func (c *TokenSliceCollection) Peek() *Token {
	// for i := c.Len() - 1; i >= 0; i-- {
	// 	if t := c.Tokens[i]; t.Type != TokenTypeEmpty {
	// 		return t
	// 	}
	// }

	if i := c.Len() - 1; i >= 0 {
		return c.Tokens[i]
	}

	return nil
}

func (t *TokenSliceCollection) Len() int {
	return len(t.Tokens)
}

func (c *TokenSliceCollection) SetAllTokenTypesToEmpty() {
	c.SetAllTokenTypes(TokenTypeEmpty)
}

func (c *TokenSliceCollection) SetAllTokenTypes(y TokenType) {
	for _, t := range c.Tokens {
		t.Type = y
	}
}

func (c *TokenSliceCollection) Get(index int) *Token {
	l := c.Len()

	if index < 0 {
		index += l
	}

	if index >= 0 && index < l {
		return c.Tokens[index]
	}

	return nil
}

func (c *TokenSliceCollection) ContainsType(y TokenType) bool {
	for _, t := range c.Tokens {
		if t.Type == y {
			return true
		}
	}

	return false
}
