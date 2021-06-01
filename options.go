package slimdown

type Options struct {
	AllowHTML                 bool
	CleanEmptyTags            bool
	DebugPrintOutput          bool
	DebugPrintTokens          bool
	EnableBackslashTransforms bool
	EnableBlockquotes         bool
	EnableCodeTags            bool
	EnableDocumentTags        bool
	EnableEmTags              bool
	EnableHeadings            bool
	EnableHorizontalRules     bool
	EnableHyphenTransforms    bool
	EnableImages              bool
	EnableLinks               bool
	EnableLists               bool
	EnableMarkTags            bool
	EnableParagraphs          bool
	EnableStrongTags          bool
	MaxConsecutiveTabs        int
	MaxConsecutiveSpaces      int
	SpacesToTab               int
	TabToSpaces               int

	isCloned bool
}

var (
	DefaultOptions = Options{
		AllowHTML:                 false,
		CleanEmptyTags:            false,
		DebugPrintOutput:          false,
		DebugPrintTokens:          false,
		EnableBackslashTransforms: false,
		EnableBlockquotes:         false,
		EnableCodeTags:            true,
		EnableDocumentTags:        false,
		EnableEmTags:              true,
		EnableHeadings:            true,
		EnableHorizontalRules:     true,
		EnableHyphenTransforms:    true,
		EnableImages:              true,
		EnableLinks:               true,
		EnableLists:               true,
		EnableParagraphs:          true,
		EnableStrongTags:          true,
		MaxConsecutiveTabs:        0,
		MaxConsecutiveSpaces:      0,
		SpacesToTab:               0,
		TabToSpaces:               0,
	}
)

func (o *Options) clone() *Options {
	if o.isCloned {
		return o
	}

	return &Options{
		AllowHTML:                 o.AllowHTML,
		CleanEmptyTags:            o.CleanEmptyTags,
		DebugPrintOutput:          o.DebugPrintOutput,
		DebugPrintTokens:          o.DebugPrintTokens,
		EnableBackslashTransforms: o.EnableBackslashTransforms,
		EnableBlockquotes:         o.EnableBlockquotes,
		EnableCodeTags:            o.EnableCodeTags,
		EnableDocumentTags:        o.EnableDocumentTags,
		EnableEmTags:              o.EnableEmTags,
		EnableHeadings:            o.EnableHeadings,
		EnableHorizontalRules:     o.EnableHorizontalRules,
		EnableHyphenTransforms:    o.EnableHyphenTransforms,
		EnableImages:              o.EnableImages,
		EnableLinks:               o.EnableLinks,
		EnableLists:               o.EnableLists,
		EnableParagraphs:          o.EnableParagraphs,
		EnableStrongTags:          o.EnableStrongTags,
		MaxConsecutiveTabs:        o.MaxConsecutiveTabs,
		MaxConsecutiveSpaces:      o.MaxConsecutiveSpaces,
		SpacesToTab:               o.SpacesToTab,
		TabToSpaces:               o.TabToSpaces,
		isCloned:                  true,
	}
}
