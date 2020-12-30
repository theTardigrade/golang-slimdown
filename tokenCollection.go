package slimdown

import (
	"html/template"
	"strings"
)

type tokenCollection struct {
	Data []*token
}

func (c *tokenCollection) PushNewEmpty(y tokenType) *token {
	return c.PushNew(y, 0, 0)
}

func (c *tokenCollection) PushNewSingle(y tokenType, startIndex int) *token {
	return c.PushNew(y, startIndex, startIndex+1)
}

func (c *tokenCollection) PushNew(y tokenType, startIndex int, endIndex int) *token {
	t := &token{
		Type:       y,
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}

	c.Push(t)

	return t
}

func (c *tokenCollection) Push(t *token) {
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

	c.Data = append(c.Data, t)

	return
}

func (c *tokenCollection) PushAsIs(t *token) {
	if t == nil {
		return
	}

	c.Data = append(c.Data, t)

	return
}

func (c *tokenCollection) PushTwin(t *token) (t2 *token, err error) {
	t2, err = t.CreateTwin()
	if err != nil {
		return
	}

	c.Push(t2)

	return
}

func (c *tokenCollection) Peek() *token {
	if l := c.Len(); l > 0 {
		return c.Data[l-1]
	}

	return nil
}

func (c *tokenCollection) PeekTwin() *token {
	if l := c.Len(); l > 0 {
		return c.Data[l-1].Twin
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

func (c *tokenCollection) ContainsType(y tokenType) bool {
	for t := c.Peek(); t != nil; t = t.Prev {
		if t.Type == y {
			return true
		}
	}

	return false
}

func (c *tokenCollection) Len() int {
	return len(c.Data)
}

func (c *tokenCollection) String() string {
	var builder strings.Builder

	for _, t := range c.Data {
		builder.Write(t.HTML)
	}

	return builder.String()
}

func (c *tokenCollection) HTML() template.HTML {
	return template.HTML(c.String())
}
