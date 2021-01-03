package slimdown

import (
	"fmt"
	"strings"
)

func debugPrintTokens(tokens *tokenCollection) {
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
