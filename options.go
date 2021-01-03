package slimdown

type Options struct {
	AllowHTML                 bool
	CleanEmptyTags            bool
	DebugPrintTokens          bool
	EnableBackslashTransforms bool
	EnableCodeTags            bool
	EnableDocumentTags        bool
	EnableEmTags              bool
	EnableLinks               bool
	EnableMarkTags            bool
	EnableParagraphTags       bool
	EnableStrongTags          bool
	MaxConsecutiveSpaces      int
	SpacesToTab               int
	TabToSpaces               int
}

var (
	DefaultOptions = Options{
		AllowHTML:                 false,
		DebugPrintTokens:          false,
		EnableBackslashTransforms: false,
		EnableCodeTags:            false,
		EnableDocumentTags:        false,
		EnableEmTags:              true,
		EnableLinks:               false,
		EnableParagraphTags:       true,
		EnableStrongTags:          true,
		SpacesToTab:               0,
		TabToSpaces:               0,
		MaxConsecutiveSpaces:      0,
	}
)
