package slimdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theTardigrade/golang-slimdown/internal/test/assets"
)

var (
	testCompileStringOptions = &Options{
		AllowHTML:                 false,
		CleanEmptyTags:            false,
		DebugPrintOutput:          false,
		DebugPrintTokens:          true,
		EnableBackslashTransforms: true,
		EnableBlockquotes:         true,
		EnableCodeTags:            true,
		EnableDocumentTags:        false,
		EnableEmTags:              true,
		EnableHeadings:            true,
		EnableHorizontalRules:     false,
		EnableHyphenTransforms:    true,
		EnableImages:              true,
		EnableLinks:               true,
		EnableLists:               true,
		EnableMarkTags:            true,
		EnableParagraphs:          true,
		EnableStrongTags:          true,
		MaxConsecutiveSpaces:      2,
		MaxConsecutiveTabs:        2,
		SpacesToTab:               5,
		TabToSpaces:               0,
	}
	testCompileStringInput          []byte
	testCompileStringExpectedOutput []byte
)

func init() {
	const filePathPrefix = "compileString/"

	testCompileStringInput = assets.Load(filePathPrefix + "input.md")
	testCompileStringExpectedOutput = assets.Load(filePathPrefix + "output.html")
}

func TestCompileString(t *testing.T) {
	output, err := Compile(testCompileStringInput, testCompileStringOptions)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(testCompileStringExpectedOutput), string(output))
}

func BenchmarkCompileString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Compile(testCompileStringInput, testCompileStringOptions)
		if err != nil {
			panic(err)
		}
	}
}
