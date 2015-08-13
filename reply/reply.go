package reply

import (
	"io"
	"strings"
)

var newLine *strings.Replacer

func init() {
	newLine = strings.NewReplacer("\n", "")
}

type ReplyWriter interface {
	WriteCSV(io.Writer) error
	//WriteJSON(io.Writer) error
	//WriteXML(io.Writer) error
}
