package slimdown

var (
	tokenTypeListLinkSegmentText = []tokenType{
		tokenTypeText,
		tokenTypeSpace,
		tokenTypeTab,
		tokenTypeBackslash,
		tokenTypeAsterisk,
		tokenTypeAsteriskDouble,
		tokenTypeAsteriskTriple,
		tokenTypeUnderscore,
		tokenTypeUnderscoreDouble,
		tokenTypeUnderscoreTriple,
		tokenTypeEqualsDouble,
		tokenTypeBacktick,
		tokenTypeExclamation,
		tokenTypeParenthesisOpen,
		tokenTypeParenthesisClose,
	}
	tokenTypeListLinkSegmentLink = []tokenType{
		tokenTypeText,
		tokenTypeAsterisk,
		tokenTypeAsteriskDouble,
		tokenTypeUnderscore,
		tokenTypeUnderscoreDouble,
	}
)

var (
	tokenTypeListImageSegmentText = tokenTypeListLinkSegmentText[:]
	tokenTypeListImageSegmentLink = tokenTypeListLinkSegmentLink[:]
)
