package slimdown

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"sync"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

const (
	debugPrintTokensMinIndent = 4
)

var (
	debugPrintTokensIntMaxLen           int
	debugPrintTokensIndent              int
	debugPrintTokensIndentCalculateOnce sync.Once
)

func debugPrintTokensIndentCalculate() {
	debugPrintTokensIndentCalculateOnce.Do(func() {
		var i int

		for {
			i++

			tt := tokenization.TokenType(i)
			if tt.String() == "UNK" {
				debugPrintTokensIntMaxLen = len(strconv.Itoa(i - 1))
				debugPrintTokensIndent = debugPrintTokensMinIndent + debugPrintTokensIntMaxLen
				break
			}
		}
	})
}

func debugPrintTokens(tokens *tokenization.TokenCollection) {
	debugPrintTokensIndentCalculate()

	var builder strings.Builder

	for i, t := range tokens.Data {
		if i > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(
			fmt.Sprintf("%*[2]d:%[2]s", debugPrintTokensIndent, t.Type),
		)
	}

	fmt.Println(builder.String())
}

func debugPrintOutput(output template.HTML) {
	fmt.Printf("%s\n", output)
}
