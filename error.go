package slimdown

import "errors"

var (
	ErrCompileTokenStackOverflow          = errors.New("token stack overflow")
	ErrCompileTokenTypeUnknown            = errors.New("token type unknown")
	ErrCompileBackslashTransformUnknown   = errors.New("backslash transform unknown")
	ErrCompileURLCannotContainDoubleQuote = errors.New("compiled URL cannot contain the double quote symbol")
)
