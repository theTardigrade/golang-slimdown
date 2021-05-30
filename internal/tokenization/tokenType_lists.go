package tokenization

var (
	TokenTypeListLinkSegmentText = []TokenType{
		TokenTypeText,
		TokenTypeSpace,
		TokenTypeTab,
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
		TokenTypeText,
		TokenTypeAsterisk,
		TokenTypeAsteriskDouble,
		TokenTypeUnderscore,
		TokenTypeUnderscoreDouble,
	}
	TokenTypeListLinkSegmentTitle = []TokenType{
		TokenTypeText,
		TokenTypeAsterisk,
		TokenTypeAsteriskDouble,
		TokenTypeUnderscore,
		TokenTypeUnderscoreDouble,
		TokenTypeSpace,
	}
)

var (
	TokenTypeListImageSegmentText = TokenTypeListLinkSegmentText[:]
	TokenTypeListImageSegmentLink = TokenTypeListLinkSegmentLink[:]
)
