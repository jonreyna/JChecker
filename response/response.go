package response

import (
	"io"
	"strings"
)

var newLine *strings.Replacer

func init() {
	newLine = strings.NewReplacer("\n", "")
}

type Response interface {
	WriteCSV(io.Writer)
}
