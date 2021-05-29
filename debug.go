package slimdown

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

func debugPrintTokens(tokens *tokenization.TokenCollection) {
	var builder strings.Builder

	for i, t := range tokens.Data {
		if i > 0 {
			builder.WriteByte(' ')
		}

		builder.WriteString(
			fmt.Sprintf("%[1]s(%[1]d)", t.Type),
		)
	}

	fmt.Println(builder.String())
}

func debugPrintOutput(output template.HTML) {
	fmt.Printf("%s\n", output)
}
