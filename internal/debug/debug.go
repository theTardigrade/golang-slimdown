package debug

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

const (
	printTokensMinIndent = 4
)

var (
	printTokensIntMaxLen           int
	printTokensIndent              int
	printTokensIndentCalculateOnce sync.Once
)

func printTokensIndentCalculate() {
	printTokensIndentCalculateOnce.Do(func() {
		var i int

		for increment := 1; i >= 0; increment = int(math.Ceil(float64(increment) * 1.25)) {
			i += increment

			if tt := tokenization.TokenType(i); tt.String() == "UNK" {
				for j, l := i-1, i-increment; j >= l; j-- {
					if tt = tokenization.TokenType(j); tt.String() != "UNK" {
						i = -j
						break
					}
				}
			}
		}

		i *= -1
		printTokensIntMaxLen = len(strconv.Itoa(i))
		printTokensIndent = printTokensMinIndent + printTokensIntMaxLen
	})
}

func PrintTokens(tokens *tokenization.TokenCollection) {
	printTokensIndentCalculate()

	var builder strings.Builder

	for i, t := range tokens.Data {
		if i > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(
			fmt.Sprintf("%*[2]d:%[2]s", printTokensIndent, t.Type),
		)
	}

	fmt.Println(builder.String())
}

func PrintOutput(output template.HTML) {
	fmt.Printf("%s\n", output)
}