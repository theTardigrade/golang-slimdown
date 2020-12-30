package slimdown

import "errors"

var (
	ErrTokenAlreadyTwinned     = errors.New("token cannot have more than one twin")
	ErrParseTokenStackOverflow = errors.New("token stack overflow")
	ErrParseTokenTypeUnknown   = errors.New("token type unknown")
)
