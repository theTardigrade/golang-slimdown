package tokenization

type TokenTypeDatum struct {
	Tags        []string
	SelfClosing bool
}

var (
	TokenTypeData = map[TokenType]TokenTypeDatum{
		TokenTypeDocumentDoctype:   {Tags: []string{"!DOCTYPE"}, SelfClosing: true},
		TokenTypeDocumentHTMLBound: {Tags: []string{"html"}},
		TokenTypeDocumentHeadBound: {Tags: []string{"head"}},
		TokenTypeDocumentBodyBound: {Tags: []string{"body"}},
		TokenTypeParagraphBound:    {Tags: []string{"p"}},
		TokenTypeHeading1Bound:     {Tags: []string{"h1"}},
		TokenTypeHeading2Bound:     {Tags: []string{"h2"}},
		TokenTypeHeading3Bound:     {Tags: []string{"h3"}},
		TokenTypeHeading4Bound:     {Tags: []string{"h4"}},
		TokenTypeHeading5Bound:     {Tags: []string{"h5"}},
		TokenTypeHeading6Bound:     {Tags: []string{"h6"}},
		TokenTypeLineBreak:         {Tags: []string{"br"}, SelfClosing: true},
		TokenTypeEqualsDouble:      {Tags: []string{"mark"}},
		TokenTypeAsterisk:          {Tags: []string{"em"}},
		TokenTypeAsteriskDouble:    {Tags: []string{"strong"}},
		TokenTypeAsteriskTriple:    {Tags: []string{"strong", "em"}},
		TokenTypeUnderscore:        {Tags: []string{"em"}},
		TokenTypeUnderscoreDouble:  {Tags: []string{"strong"}},
		TokenTypeUnderscoreTriple:  {Tags: []string{"strong", "em"}},
		TokenTypeBacktick:          {Tags: []string{"code"}},
		TokenTypeLink:              {Tags: []string{"a"}},
		TokenTypeImage:             {Tags: []string{"img"}, SelfClosing: true},
	}
)
