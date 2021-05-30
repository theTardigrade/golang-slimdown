package slimdown

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
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
		EnableMarkTags:            true,
		EnableParagraphs:          true,
		EnableStrongTags:          true,
		MaxConsecutiveSpaces:      0,
		SpacesToTabs:              0,
		TabsToSpaces:              0,
	}
	testCompileStringInput          []byte
	testCompileStringExpectedOutput []byte
)

func init() {
	var once sync.Once
	once.Do(testInit)

	const filePathPrefix = "testAssets/compileString"
	var err error

	testCompileStringInput, err = os.ReadFile(globalFilepath.Join(filePathPrefix + "Input.md"))
	if err != nil {
		panic(err)
	}

	testCompileStringExpectedOutput, err = os.ReadFile(globalFilepath.Join(filePathPrefix + "Output.html"))
	if err != nil {
		panic(err)
	}
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
