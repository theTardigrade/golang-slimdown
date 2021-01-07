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
