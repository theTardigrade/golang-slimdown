package tokenization

import (
	"strings"
	"unicode"
)

type TokenType uint8

const (
	TokenTypeEmpty TokenType = iota
	TokenTypeDocumentDoctype
	TokenTypeDocumentHTMLBound
	TokenTypeDocumentHeadBound
	TokenTypeDocumentBodyBound
	TokenTypeParagraphBound
	TokenTypeText
	TokenTypeSpace
	TokenTypeTab
	TokenTypeCarriageReturn
	TokenTypeBackslash
	TokenTypeAsterisk
	TokenTypeAsteriskDouble
	TokenTypeAsteriskTriple
	TokenTypeUnderscore
	TokenTypeUnderscoreDouble
	TokenTypeUnderscoreTriple
	TokenTypeEqualsDouble
	TokenTypeBacktick
	TokenTypeExclamation
	TokenTypeParenthesisOpen
	TokenTypeParenthesisClose
	TokenTypeSquareBracketOpen
	TokenTypeSquareBracketClose
	TokenTypeLink
	TokenTypeImage
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeEmpty:
		return "EMP"
	case TokenTypeDocumentDoctype:
		return "DOC_TYP"
	case TokenTypeDocumentHTMLBound:
		return "DOC_HTM"
	case TokenTypeDocumentHeadBound:
		return "DOC_HED"
	case TokenTypeDocumentBodyBound:
		return "DOC_BDY"
	case TokenTypeParagraphBound:
		return "PAR_BND"
	case TokenTypeText:
		return "TXT"
	case TokenTypeSpace:
		return "SPC"
	case TokenTypeTab:
		return "TAB"
	case TokenTypeCarriageReturn:
		return "CAR_RET"
	case TokenTypeBackslash:
		return "BKS"
	case TokenTypeAsterisk:
		return "AST"
	case TokenTypeAsteriskDouble:
		return "AST_DUB"
	case TokenTypeAsteriskTriple:
		return "AST_TRI"
	case TokenTypeUnderscore:
		return "UND"
	case TokenTypeUnderscoreDouble:
		return "UND_DUB"
	case TokenTypeUnderscoreTriple:
		return "UND_TRI"
	case TokenTypeEqualsDouble:
		return "EQU_DUB"
	case TokenTypeBacktick:
		return "BAK_TIK"
	case TokenTypeExclamation:
		return "EXL"
	case TokenTypeParenthesisOpen:
		return "PRN_OPN"
	case TokenTypeParenthesisClose:
		return "PRN_CLS"
	case TokenTypeSquareBracketOpen:
		return "SQU_BRK_OPN"
	case TokenTypeSquareBracketClose:
		return "SQU_BRK_CLS"
	case TokenTypeLink:
		return "LNK"
	case TokenTypeImage:
		return "IMG"
	}

	return "UNK"
}

func (t TokenType) ClassName() string {
	var builder strings.Builder

	for _, r := range t.String() {
		if r == '_' {
			builder.WriteByte('-')
		} else {
			builder.WriteRune(unicode.ToLower(r))
		}
	}

	return builder.String()
}
