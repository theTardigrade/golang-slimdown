package slimdown

import "strconv"

type tokenType uint8

const (
	tokenTypeText tokenType = iota
	tokenTypeCarriageReturn
	tokenTypeTab
	tokenTypeNewline
	tokenTypeBackslash
	tokenTypeAsterisk
	tokenTypeAsteriskDouble
	tokenTypeUnderscore
	tokenTypeUnderscoreDouble
	tokenTypeBacktick
)

var (
	tokenTypeTagMap = map[tokenType]string{
		tokenTypeNewline:          "p",
		tokenTypeAsterisk:         "em",
		tokenTypeAsteriskDouble:   "strong",
		tokenTypeUnderscore:       "em",
		tokenTypeUnderscoreDouble: "strong",
		tokenTypeBacktick:         "code",
	}
)

type token struct {
	Type       tokenType
	StartIndex int
	EndIndex   int
	HTML       []byte
	Prev, Next *token
	Twin       *token
}

func (t *token) String() string {
	return strconv.FormatUint(uint64(t.Type), 10)
}

func (t *token) CreateTwin() (t2 *token, err error) {
	if t.Twin != nil {
		err = ErrTokenAlreadyTwinned
		return
	}

	t2 = &token{
		Type:       t.Type,
		StartIndex: t.StartIndex,
		EndIndex:   t.EndIndex,
	}

	copy(t2.HTML, t.HTML)

	t2.Twin = t
	t.Twin = t2

	return
}
