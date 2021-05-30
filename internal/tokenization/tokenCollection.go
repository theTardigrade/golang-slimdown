package tokenization

import (
	"html/template"
	"strings"
)

type TokenCollection struct {
	Data             []*Token
	TagPairCleanData [][2]*Token
	Input            []byte
}

func TokenCollectionNewEmpty() *TokenCollection {
	return &TokenCollection{}
}

func TokenCollectionNew(input []byte) *TokenCollection {
	return &TokenCollection{
		Input: input,
	}
}

func (c *TokenCollection) PushNewEmpty(y TokenType) *Token {
	return c.PushNew(y, 0, 0)
}

func (c *TokenCollection) PushNewSingle(y TokenType, InputStartIndex int) *Token {
	return c.PushNew(y, InputStartIndex, InputStartIndex+1)
}

func (c *TokenCollection) PushNew(y TokenType, InputStartIndex int, InputEndIndex int) *Token {
	t := &Token{
		Type:            y,
		InputStartIndex: InputStartIndex,
		InputEndIndex:   InputEndIndex,
	}

	c.Push(t)

	return t
}

func (c *TokenCollection) Push(tokens ...*Token) {
	for _, t := range tokens {
		if t == nil {
			continue
		}

		if prev := c.Peek(); prev != nil {
			t.RawPrev = prev
			prev.RawNext = t
		} else {
			t.RawPrev = nil
		}

		t.RawNext = nil
		t.Collection = c

		c.Data = append(c.Data, t)
	}
}

func (c *TokenCollection) PushAsIs(tokens ...*Token) {
	for _, t := range tokens {
		if t == nil {
			continue
		}

		c.Data = append(c.Data, t)
	}
}

func (c *TokenCollection) Peek() *Token {
	for i := c.Len() - 1; i >= 0; i-- {
		if t := c.Data[i]; t.Type != TokenTypeEmpty {
			return t
		}
	}

	return nil
}

func (c *TokenCollection) Pop() *Token {
	if i := c.Len() - 1; i >= 0 {
		t := c.Data[i]

		if prev := t.RawPrev; prev != nil {
			prev.RawNext = nil
		}

		t.RawPrev = nil
		t.RawNext = nil

		c.Data = c.Data[:i]

		return t
	}

	return nil
}

func (c *TokenCollection) PopAsIs() *Token {
	if i := c.Len() - 1; i >= 0 {
		t := c.Data[i]

		c.Data = c.Data[:i]

		return t
	}

	return nil
}

func (c *TokenCollection) Swap(index1, index2 int) (success bool) {
	if index1 < 0 || index2 < 0 {
		return
	}

	if l := c.Len(); index1 >= l || index2 >= l {
		return
	}

	t1, t2 := c.Data[index1], c.Data[index2]
	t1CachedNext, t1CachedPrev := t1.RawNext, t1.RawPrev

	t1.RawNext, t1.RawPrev = t2.RawNext, t2.RawPrev
	t2.RawNext, t2.RawPrev = t1CachedNext, t1CachedPrev

	c.Data[index2], c.Data[index1] = t1, t2

	return
}

func (c *TokenCollection) ContainsType(y TokenType) bool {
	for t := c.Peek(); t != nil; t = t.RawPrev {
		if t.Type == y {
			return true
		}
	}

	return false
}

func (c *TokenCollection) Get(index int) *Token {
	l := c.Len()

	if index < 0 {
		index += l
	}

	if index >= 0 && index < l {
		return c.Data[index]
	}

	return nil
}

func (c *TokenCollection) Len() int {
	return len(c.Data)
}

func (c *TokenCollection) String() string {
	var builder strings.Builder

	for _, t := range c.Data {
		if b := t.Bytes(); len(b) > 0 {
			builder.Write(b)
		}
	}

	return builder.String()
}

func (c *TokenCollection) HTML() template.HTML {
	var builder strings.Builder

	for _, t := range c.Data {
		if t.Type != TokenTypeEmpty {
			if h := t.HTML; len(h) > 0 {
				builder.Write(h)
			}
		}
	}

	return template.HTML(builder.String())
}

func (c *TokenCollection) SetAllTokenTypesToEmpty() {
	c.SetAllTokenTypes(TokenTypeEmpty)
}

func (c *TokenCollection) SetAllTokenTypes(y TokenType) {
	for _, t := range c.Data {
		t.Type = y
	}
}
