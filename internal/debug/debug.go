package debug

import (
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

const (
	printTokensMinIndent = 4
)

var (
	printTokensIntMaxLen int
	printTokensIndent    int
)

func init() {
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
}

func PrintTokens(tokens *tokenization.TokenListCollection) {
	var builder strings.Builder
	var i int

	for t := tokens.HeadToken; t != nil; t = t.RawNext {
		if i++; i > 1 {
			builder.WriteByte('\n')
		}

		builder.WriteString(
			fmt.Sprintf("%*[2]d:%[2]s:%*[3]d", printTokensIndent, t.Type, t.Indent),
		)
	}

	fmt.Println(builder.String())
}

func PrintOutput(output template.HTML) {
	fmt.Printf("%s\n", output)
}
