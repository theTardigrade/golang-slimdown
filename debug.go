package slimdown

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/theTardigrade/golang-slimdown/internal/tokenization"
)

const (
	debugTokenMinIndent = 4
)

var (
	debugTokenIntMaxLen int
	debugTokenIndent    int
)

func init() {
	var i int

	for {
		i++

		tt := tokenization.TokenType(i)
		if tt.String() == "UNK" {
			debugTokenIntMaxLen = len(strconv.Itoa(i))
			debugTokenIndent = debugTokenMinIndent + debugTokenIntMaxLen
			break
		}
	}
}

func debugPrintTokens(tokens *tokenization.TokenCollection) {
	var builder strings.Builder

	for i, t := range tokens.Data {
		if i > 0 {
			builder.WriteByte('\n')
		}

		builder.WriteString(
			fmt.Sprintf("%*[2]d:%[2]s", debugTokenIndent, t.Type),
		)
	}

	fmt.Println(builder.String())
}

func debugPrintOutput(output template.HTML) {
	fmt.Printf("%s\n", output)
}
