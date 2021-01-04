package slimdown

type tokenTypeDatum struct {
	Tags        []string
	SelfClosing bool
}

var (
	tokenTypeData = map[tokenType]tokenTypeDatum{
		tokenTypeDocumentDoctype:   {Tags: []string{"!DOCTYPE"}, SelfClosing: true},
		tokenTypeDocumentHTMLBound: {Tags: []string{"html"}},
		tokenTypeDocumentHeadBound: {Tags: []string{"head"}},
		tokenTypeDocumentBodyBound: {Tags: []string{"body"}},
		tokenTypeParagraphBound:    {Tags: []string{"p"}},
		tokenTypeEqualsDouble:      {Tags: []string{"mark"}},
		tokenTypeAsterisk:          {Tags: []string{"em"}},
		tokenTypeAsteriskDouble:    {Tags: []string{"strong"}},
		tokenTypeAsteriskTriple:    {Tags: []string{"strong", "em"}},
		tokenTypeUnderscore:        {Tags: []string{"em"}},
		tokenTypeUnderscoreDouble:  {Tags: []string{"strong"}},
		tokenTypeBacktick:          {Tags: []string{"code"}},
		tokenTypeLink:              {Tags: []string{"a"}},
		tokenTypeImage:             {Tags: []string{"img"}, SelfClosing: true},
	}
)
