package slimdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/theTardigrade/golang-slimdown/internal/test/assets"
)

var (
	testCompileStringInput          = make(map[string][]byte)
	testCompileStringExpectedOutput = make(map[string][]byte)
	testCompileStringOptions        = make(map[string]*Options)
)

func init() {
	const filePathPrefix = "compileString/"

	for _, key := range []string{
		"blockquotes",
		"spacesToTab",
		"tabToSpaces",
	} {
		prefix := filePathPrefix + key
		input := assets.Load(prefix + "Input.md")
		output := assets.Load(prefix + "Output.html")

		testCompileStringInput[key] = input
		testCompileStringExpectedOutput[key] = output
	}
}

/* blockquotes */

func init() {
	testCompileStringOptions["blockquotes"] = &Options{
		DebugPrintTokens:  true,
		EnableBlockquotes: true,
		EnableHeadings:    true,
		EnableParagraphs:  true,
	}
}

func TestCompileString_blockquotes(t *testing.T) {
	const key = "blockquotes"

	output, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(testCompileStringExpectedOutput[key]), string(output))
}

func BenchmarkCompileString_blockquotes(b *testing.B) {
	const key = "blockquotes"

	for i := 0; i < b.N; i++ {
		_, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
		if err != nil {
			panic(err)
		}
	}
}

/* spacesToTab */

func init() {
	testCompileStringOptions["spacesToTab"] = &Options{
		DebugPrintTokens: true,
		EnableParagraphs: true,
		SpacesToTab:      4,
	}
}

func TestCompileString_spacesToTab(t *testing.T) {
	const key = "spacesToTab"

	output, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(testCompileStringExpectedOutput[key]), string(output))
}

func BenchmarkCompileString_spacesToTab(b *testing.B) {
	const key = "spacesToTab"

	for i := 0; i < b.N; i++ {
		_, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
		if err != nil {
			panic(err)
		}
	}
}

/* tabToSpaces */

func init() {
	testCompileStringOptions["tabToSpaces"] = &Options{
		DebugPrintTokens: true,
		EnableParagraphs: true,
		TabToSpaces:      1,
	}
}

func TestCompileString_tabToSpaces(t *testing.T) {
	const key = "tabToSpaces"

	output, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
	if err != nil {
		panic(err)
	}

	assert.Equal(t, string(testCompileStringExpectedOutput[key]), string(output))
}

func BenchmarkCompileString_tabToSpaces(b *testing.B) {
	const key = "tabToSpaces"

	for i := 0; i < b.N; i++ {
		_, err := Compile(testCompileStringInput[key], testCompileStringOptions[key])
		if err != nil {
			panic(err)
		}
	}
}
