package tokenization

var (
	TokenTypeListLinkSegmentText = []TokenType{
		TokenTypeTextGroup,
		TokenTypeSpaceGroup,
		TokenTypeTabGroup,
		TokenTypeBackslash,
		TokenTypeAsterisk,
		TokenTypeAsteriskDouble,
		TokenTypeAsteriskTriple,
		TokenTypeUnderscore,
		TokenTypeUnderscoreDouble,
		TokenTypeUnderscoreTriple,
		TokenTypeHyphen,
		TokenTypeHyphenDouble,
		TokenTypeHyphenTriple,
		TokenTypeEqualsDouble,
		TokenTypeBacktick,
		TokenTypeExclamation,
		TokenTypeParenthesisOpen,
		TokenTypeParenthesisClose,
	}
	TokenTypeListLinkSegmentLink = []TokenType{
		TokenTypeTextGroup,
		TokenTypeAsterisk,
		TokenTypeAsteriskDouble,
		TokenTypeUnderscore,
		TokenTypeUnderscoreDouble,
	}
	TokenTypeListLinkSegmentTitle = []TokenType{
		TokenTypeTextGroup,
		TokenTypeAsterisk,
		TokenTypeAsteriskDouble,
		TokenTypeUnderscore,
		TokenTypeUnderscoreDouble,
		TokenTypeSpaceGroup,
	}
)

var (
	TokenTypeListImageSegmentText  = TokenTypeListLinkSegmentText[:]
	TokenTypeListImageSegmentLink  = TokenTypeListLinkSegmentLink[:]
	TokenTypeListImageSegmentTitle = TokenTypeListLinkSegmentTitle[:]
)
