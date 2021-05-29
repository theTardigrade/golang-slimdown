package slimdown

type Options struct {
	AllowHTML                 bool
	CleanEmptyTags            bool
	DebugPrintOutput          bool
	DebugPrintTokens          bool
	EnableBackslashTransforms bool
	EnableCodeTags            bool
	EnableDocumentTags        bool
	EnableEmTags              bool
	EnableHeadings            bool
	EnableHyphenTransforms    bool
	EnableImages              bool
	EnableLinks               bool
	EnableMarkTags            bool
	EnableParagraphs          bool
	EnableStrongTags          bool
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
		EnableCodeTags:            false,
		EnableDocumentTags:        false,
		EnableEmTags:              true,
		EnableHyphenTransforms:    true,
		EnableImages:              false,
		EnableLinks:               true,
		EnableParagraphs:          true,
		EnableStrongTags:          true,
		MaxConsecutiveSpaces:      0,
		SpacesToTab:               0,
		TabToSpaces:               0,
	}
)

func (o *Options) Clone() *Options {
	if o.isCloned {
		return o
	}

	return &Options{
		AllowHTML:                 o.AllowHTML,
		CleanEmptyTags:            o.CleanEmptyTags,
		DebugPrintTokens:          o.DebugPrintTokens,
		EnableBackslashTransforms: o.EnableBackslashTransforms,
		EnableCodeTags:            o.EnableCodeTags,
		EnableDocumentTags:        o.EnableDocumentTags,
		EnableEmTags:              o.EnableEmTags,
		EnableHeadings:            o.EnableHeadings,
		EnableImages:              o.EnableImages,
		EnableLinks:               o.EnableLinks,
		EnableParagraphs:          o.EnableParagraphs,
		EnableStrongTags:          o.EnableStrongTags,
		MaxConsecutiveSpaces:      o.MaxConsecutiveSpaces,
		SpacesToTab:               o.SpacesToTab,
		TabToSpaces:               o.TabToSpaces,
		isCloned:                  true,
	}
}
