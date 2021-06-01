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
		"tabToSpaces",
	} {
		prefix := filePathPrefix + key
		input := assets.Load(prefix + "Input.md")
		output := assets.Load(prefix + "Output.html")

		testCompileStringInput[key] = input
		testCompileStringExpectedOutput[key] = output
	}

}

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
