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
	TokenTypeBlockquoteBound
	TokenTypeHeading1Bound
	TokenTypeHeading2Bound
	TokenTypeHeading3Bound
	TokenTypeHeading4Bound
	TokenTypeHeading5Bound
	TokenTypeHeading6Bound
	TokenTypeLineBreak
	TokenTypeTextGroup
	TokenTypeSpaceGroup
	TokenTypeSpaceHair
	TokenTypeTabGroup
	TokenTypeCarriageReturn
	TokenTypeHorizontalRule
	TokenTypeBackslash
	TokenTypeAsterisk
	TokenTypeAsteriskDouble
	TokenTypeAsteriskTriple
	TokenTypeUnderscore
	TokenTypeUnderscoreDouble
	TokenTypeUnderscoreTriple
	TokenTypeHash
	TokenTypeHashDouble
	TokenTypeHashTriple
	TokenTypeHashQuadruple
	TokenTypeHashQuintuple
	TokenTypeHashSextuple
	TokenTypeHyphen
	TokenTypeHyphenDouble
	TokenTypeHyphenTriple
	TokenTypeDashEm
	TokenTypeDashEn
	TokenTypeEqualsDouble
	TokenTypeBacktick
	TokenTypeBacktickDouble
	TokenTypeExclamation
	TokenTypeParenthesisOpen
	TokenTypeParenthesisClose
	TokenTypeSquareBracketOpen
	TokenTypeSquareBracketClose
	TokenTypeAngleBracketOpen
	TokenTypeAngleBracketClose
	TokenTypeLinkBound
	TokenTypeImageBound
	TokenTypeUnorderedListBound
	TokenTypeListItemBound
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
	case TokenTypeBlockquoteBound:
		return "QUO_BND"
	case TokenTypeHeading1Bound:
		return "HD1_BND"
	case TokenTypeHeading2Bound:
		return "HD2_BND"
	case TokenTypeHeading3Bound:
		return "HD3_BND"
	case TokenTypeHeading4Bound:
		return "HD4_BND"
	case TokenTypeHeading5Bound:
		return "HD5_BND"
	case TokenTypeHeading6Bound:
		return "HD6_BND"
	case TokenTypeLineBreak:
		return "LBK"
	case TokenTypeTextGroup:
		return "TXT"
	case TokenTypeSpaceGroup:
		return "SPC"
	case TokenTypeSpaceHair:
		return "SPC_HAR"
	case TokenTypeTabGroup:
		return "TAB"
	case TokenTypeCarriageReturn:
		return "CRT"
	case TokenTypeHorizontalRule:
		return "HRL"
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
	case TokenTypeHash:
		return "HSH"
	case TokenTypeHashDouble:
		return "HSH_DUB"
	case TokenTypeHashTriple:
		return "HSH_TRI"
	case TokenTypeHashQuadruple:
		return "HSH_QUA"
	case TokenTypeHashQuintuple:
		return "HSH_QUI"
	case TokenTypeHashSextuple:
		return "HSH_SXT"
	case TokenTypeHyphen:
		return "HYP"
	case TokenTypeHyphenDouble:
		return "HYP_DUB"
	case TokenTypeHyphenTriple:
		return "HYP_TRI"
	case TokenTypeDashEm:
		return "DSH_MMM"
	case TokenTypeDashEn:
		return "DSH_NNN"
	case TokenTypeEqualsDouble:
		return "EQU_DUB"
	case TokenTypeBacktick:
		return "BTK"
	case TokenTypeBacktickDouble:
		return "BTK_DUB"
	case TokenTypeExclamation:
		return "EXL"
	case TokenTypeParenthesisOpen:
		return "PRN_OPN"
	case TokenTypeParenthesisClose:
		return "PRN_CLS"
	case TokenTypeSquareBracketOpen:
		return "SBK_OPN"
	case TokenTypeSquareBracketClose:
		return "SBK_CLS"
	case TokenTypeAngleBracketOpen:
		return "ABK_OPN"
	case TokenTypeAngleBracketClose:
		return "ABK_CLS"
	case TokenTypeLinkBound:
		return "LNK_BND"
	case TokenTypeImageBound:
		return "IMG_BND"
	case TokenTypeUnorderedListBound:
		return "UND_LST_BND"
	case TokenTypeListItemBound:
		return "LST_ITM_BND"
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
