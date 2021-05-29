package slimdown

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

func TestCompileString(t *testing.T) {
	input, err := os.ReadFile(globalFilepath.Join("testAssets/input.md"))
	if err != nil {
		panic(err)
	}

	expectedOutput, err := os.ReadFile(globalFilepath.Join("testAssets/output.html"))
	if err != nil {
		panic(err)
	}

	output, err := Compile(input, &Options{
		AllowHTML:                 false,
		CleanEmptyTags:            false,
		DebugPrintTokens:          true,
		EnableBackslashTransforms: true,
		EnableCodeTags:            true,
		EnableDocumentTags:        false,
		EnableEmTags:              true,
		EnableHyphenTransforms:    true,
		EnableImages:              true,
		EnableLinks:               true,
		EnableMarkTags:            true,
		EnableParagraphs:          true,
		EnableStrongTags:          true,
		MaxConsecutiveSpaces:      0,
		SpacesToTab:               0,
		TabToSpaces:               0,
		UseConcurrency:            false,
	})
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(output), string(expectedOutput))
}
