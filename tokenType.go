package slimdown

import (
	"strings"
	"unicode"
)

type tokenType uint8

const (
	tokenTypeEmpty tokenType = iota
	tokenTypeDocumentDoctype
	tokenTypeDocumentHTMLBound
	tokenTypeDocumentHeadBound
	tokenTypeDocumentBodyBound
	tokenTypeParagraphBound
	tokenTypeText
	tokenTypeSpace
	tokenTypeTab
	tokenTypeCarriageReturn
	tokenTypeBackslash
	tokenTypeAsterisk
	tokenTypeAsteriskDouble
	tokenTypeAsteriskTriple
	tokenTypeUnderscore
	tokenTypeUnderscoreDouble
	tokenTypeUnderscoreTriple
	tokenTypeEqualsDouble
	tokenTypeBacktick
	tokenTypeExclamation
	tokenTypeParenthesisOpen
	tokenTypeParenthesisClose
	tokenTypeSquareBracketOpen
	tokenTypeSquareBracketClose
	tokenTypeLink
	tokenTypeImage
)

func (t tokenType) String() string {
	switch t {
	case tokenTypeEmpty:
		return "EMP"
	case tokenTypeDocumentDoctype:
		return "DOC_TYP"
	case tokenTypeDocumentHTMLBound:
		return "DOC_HTM"
	case tokenTypeDocumentHeadBound:
		return "DOC_HED"
	case tokenTypeDocumentBodyBound:
		return "DOC_BDY"
	case tokenTypeParagraphBound:
		return "PAR_BND"
	case tokenTypeText:
		return "TXT"
	case tokenTypeSpace:
		return "SPC"
	case tokenTypeTab:
		return "TAB"
	case tokenTypeCarriageReturn:
		return "CAR_RET"
	case tokenTypeBackslash:
		return "BKS"
	case tokenTypeAsterisk:
		return "AST"
	case tokenTypeAsteriskDouble:
		return "AST_DUB"
	case tokenTypeAsteriskTriple:
		return "AST_TRI"
	case tokenTypeUnderscore:
		return "UND"
	case tokenTypeUnderscoreDouble:
		return "UND_DUB"
	case tokenTypeUnderscoreTriple:
		return "UND_TRI"
	case tokenTypeEqualsDouble:
		return "EQU_DUB"
	case tokenTypeBacktick:
		return "BAK_TIK"
	case tokenTypeExclamation:
		return "EXL"
	case tokenTypeParenthesisOpen:
		return "PRN_OPN"
	case tokenTypeParenthesisClose:
		return "PRN_CLS"
	case tokenTypeSquareBracketOpen:
		return "SQU_BRK_OPN"
	case tokenTypeSquareBracketClose:
		return "SQU_BRK_CLS"
	case tokenTypeLink:
		return "LNK"
	case tokenTypeImage:
		return "IMG"
	}

	panic(ErrTokenTypeStringNotFound)
}

func (t tokenType) ClassName() string {
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
