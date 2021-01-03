package slimdown

import (
	"html/template"
	"strings"
)

type tokenCollection struct {
	Data             []*token
	TagPairCleanData [][2]*token
	Input            []byte
}

func tokenCollectionNewEmpty() *tokenCollection {
	return &tokenCollection{}
}

func tokenCollectionNew(input []byte) *tokenCollection {
	return &tokenCollection{
		Input: input,
	}
}

func (c *tokenCollection) PushNewEmpty(y tokenType) *token {
	return c.PushNew(y, 0, 0)
}

func (c *tokenCollection) PushNewSingle(y tokenType, InputStartIndex int) *token {
	return c.PushNew(y, InputStartIndex, InputStartIndex+1)
}

func (c *tokenCollection) PushNew(y tokenType, InputStartIndex int, InputEndIndex int) *token {
	t := &token{
		Type:            y,
		InputStartIndex: InputStartIndex,
		InputEndIndex:   InputEndIndex,
		Collection:      c,
	}

	c.Push(t)

	return t
}

func (c *tokenCollection) Push(tokens ...*token) {
	for _, t := range tokens {
		if t == nil {
			return
		}

		if prev := c.Peek(); prev != nil {
			t.Prev = prev
			prev.Next = t
		} else {
			t.Prev = nil
		}

		t.Next = nil
		t.Collection = c

		c.Data = append(c.Data, t)
	}

	return
}

func (c *tokenCollection) PushAsIs(tokens ...*token) {
	for _, t := range tokens {
		if t == nil {
			return
		}

		c.Data = append(c.Data, t)
	}

	return
}

func (c *tokenCollection) Peek() *token {
	if l := c.Len(); l > 0 {
		return c.Data[l-1]
	}

	return nil
}

func (c *tokenCollection) Pop() *token {
	if i := c.Len() - 1; i >= 0 {
		t := c.Data[i]

		if prev := t.Prev; prev != nil {
			prev.Next = nil
		}

		t.Prev = nil
		t.Next = nil

		c.Data = c.Data[:i]

		return t
	}

	return nil
}

func (c *tokenCollection) PopAsIs() *token {
	if i := c.Len() - 1; i >= 0 {
		t := c.Data[i]

		c.Data = c.Data[:i]

		return t
	}

	return nil
}

func (c *tokenCollection) Swap(index1, index2 int) (success bool) {
	if index1 < 0 || index2 < 0 {
		return
	}

	if l := c.Len(); index1 >= l || index2 >= l {
		return
	}

	t1, t2 := c.Data[index1], c.Data[index2]
	t1CachedNext, t1CachedPrev := t1.Next, t1.Prev

	t1.Next, t1.Prev = t2.Next, t2.Prev
	t2.Next, t2.Prev = t1CachedNext, t1CachedPrev

	c.Data[index2], c.Data[index1] = t1, t2

	return
}

func (c *tokenCollection) ContainsType(y tokenType) bool {
	for t := c.Peek(); t != nil; t = t.Prev {
		if t.Type == y {
			return true
		}
	}

	return false
}

func (c *tokenCollection) Get(index int) *token {
	l := c.Len()

	if index < 0 {
		index += l
	}

	if index >= 0 && index < l {
		return c.Data[index]
	}

	return nil
}

func (c *tokenCollection) Len() int {
	return len(c.Data)
}

func (c *tokenCollection) String() string {
	var builder strings.Builder

	for _, t := range c.Data {
		if b := t.Bytes(); len(b) > 0 {
			builder.Write(b)
		}
	}

	return builder.String()
}

func (c *tokenCollection) HTML() template.HTML {
	var builder strings.Builder

	for _, t := range c.Data {
		if t.Type != tokenTypeEmpty {
			if h := t.HTML; len(h) > 0 {
				builder.Write(h)
			}
		}
	}

	return template.HTML(builder.String())
}
