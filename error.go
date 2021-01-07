package slimdown

import "errors"

var (
	ErrCompileTokenStackOverflow        = errors.New("token stack overflow") // unused
	ErrCompileTokenTypeUnknown          = errors.New("token type unknown")
	ErrCompileBackslashTransformUnknown = errors.New("backslash transform unknown")
)
