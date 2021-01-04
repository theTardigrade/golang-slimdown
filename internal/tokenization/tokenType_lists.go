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
)

var (
	TokenTypeListImageSegmentText = TokenTypeListLinkSegmentText[:]
	TokenTypeListImageSegmentLink = TokenTypeListLinkSegmentLink[:]
)
