package tokenization

import (
	"html/template"
	"strings"
)

type TokenListCollection struct {
	HeadToken        *Token
	TailToken        *Token
	len              int
	TagPairCleanData [][2]*Token
	Input            []byte
}

func TokenListCollectionNew(input []byte) *TokenListCollection {
	return &TokenListCollection{
		Input: input,
	}
}

func (c *TokenListCollection) PushNewEmpty(y TokenType) *Token {
	return c.InsertNewEmptyAfter(nil, y)
}

func (c *TokenListCollection) InsertNewEmptyAfter(referenceToken *Token, y TokenType) *Token {
	return c.InsertNewAfter(referenceToken, y, 0, 0)
}

func (c *TokenListCollection) InsertNewEmptyBefore(referenceToken *Token, y TokenType) *Token {
	return c.InsertNewBefore(referenceToken, y, 0, 0)
}

func (c *TokenListCollection) PushNewSingle(y TokenType, inputStartIndex int) *Token {
	return c.InsertNewSingleAfter(nil, y, inputStartIndex)
}

func (c *TokenListCollection) InsertNewSingleAfter(referenceToken *Token, y TokenType, inputStartIndex int) *Token {
	return c.InsertNewAfter(referenceToken, y, inputStartIndex, inputStartIndex+1)
}

func (c *TokenListCollection) InsertNewSingleBefore(referenceToken *Token, y TokenType, inputStartIndex int) *Token {
	return c.InsertNewBefore(referenceToken, y, inputStartIndex, inputStartIndex+1)
}

func (c *TokenListCollection) PushNew(y TokenType, inputStartIndex int, inputEndIndex int) *Token {
	t := tokenNew(y, inputStartIndex, inputEndIndex)

	c.Push(t)

	return t
}

func (c *TokenListCollection) InsertNewAfter(referenceToken *Token, y TokenType, inputStartIndex int, inputEndIndex int) *Token {
	t := tokenNew(y, inputStartIndex, inputEndIndex)

	c.InsertAfter(referenceToken, t)

	return t
}

func (c *TokenListCollection) InsertNewBefore(referenceToken *Token, y TokenType, inputStartIndex int, inputEndIndex int) *Token {
	t := tokenNew(y, inputStartIndex, inputEndIndex)

	c.InsertBefore(referenceToken, t)

	return t
}

func (c *TokenListCollection) Push(tokens ...*Token) {
	for _, t := range tokens {
		if t == nil {
			continue
		}

		if c.TailToken == nil {
			c.HeadToken = t
			c.TailToken = t
			t.RawPrev = nil
			t.RawNext = nil
		} else {
			tail := c.TailToken

			tail.RawNext = t
			c.TailToken = t
			t.RawPrev = tail
			t.RawNext = nil
		}

		if t.ListCollection != c {
			c.len++
			t.ListCollection = c
		}
	}
}

func (c *TokenListCollection) InsertAfter(referenceToken *Token, tokens ...*Token) {
	if referenceToken == nil || referenceToken == c.TailToken {
		c.Push(tokens...)
		return
	}

	for _, t := range tokens {
		if t == nil {
			continue
		}

		n := referenceToken.RawNext

		if n != nil {
			n.RawPrev = t
		}

		referenceToken.RawNext = t

		t.RawPrev = referenceToken
		t.RawNext = n

		if t.ListCollection != c {
			c.len++
			t.ListCollection = c
		}
	}
}

func (c *TokenListCollection) InsertBefore(referenceToken *Token, tokens ...*Token) {
	if referenceToken == nil {
		c.Push(tokens...)
		return
	}

	for _, t := range tokens {
		if t == nil {
			continue
		}

		p := referenceToken.RawPrev

		if p != nil {
			p.RawNext = t
		}

		referenceToken.RawPrev = t

		t.RawNext = referenceToken
		t.RawPrev = p

		if t.ListCollection != c {
			c.len++
			t.ListCollection = c
		}
	}
}

// func (c *TokenListCollection) PushAsIs(tokens ...*Token) {
// 	for _, t := range tokens {
// 		if t == nil {
// 			continue
// 		}

// 		c.Data = append(c.Data, t)
// 	}
// }

func (c *TokenListCollection) Peek() *Token {
	// for i := c.Len() - 1; i >= 0; i-- {
	// 	if t := c.Data[i]; t.Type != TokenTypeEmpty {
	// 		return t
	// 	}
	// }

	// if i := c.Len() - 1; i >= 0 {
	// 	return c.Data[i]
	// }

	return c.TailToken
}

func (c *TokenListCollection) Pop() *Token {
	if c.TailToken != nil {
		t := c.TailToken
		p := t.RawPrev

		c.TailToken = p

		if p != nil {
			p.RawNext = nil
			t.RawPrev = nil
		}

		c.len--

		return t
	}

	return nil
}

// func (c *TokenListCollection) Swap(index1, index2 int) (success bool) {
// 	if index1 < 0 || index2 < 0 {
// 		return
// 	}

// 	if l := c.Len(); index1 >= l || index2 >= l {
// 		return
// 	}

// 	t1, t2 := c.Data[index1], c.Data[index2]
// 	t1CachedNext, t1CachedPrev := t1.RawNext, t1.RawPrev

// 	t1.RawNext, t1.RawPrev = t2.RawNext, t2.RawPrev
// 	t2.RawNext, t2.RawPrev = t1CachedNext, t1CachedPrev

// 	c.Data[index2], c.Data[index1] = t1, t2

// 	return
// }

func (c *TokenListCollection) Len() int {
	return c.len
}

func (c *TokenListCollection) String() string {
	var builder strings.Builder

	for t := c.HeadToken; t != nil; t = t.RawNext {
		if b := t.Bytes(); len(b) > 0 {
			builder.Write(b)
		}
	}

	return builder.String()
}

func (c *TokenListCollection) HTML() template.HTML {
	var builder strings.Builder

	for t := c.HeadToken; t != nil; t = t.RawNext {
		if t.Type != TokenTypeEmpty {
			if h := t.HTML; len(h) > 0 {
				builder.Write(h)
			}
		}
	}

	return template.HTML(builder.String())
}
