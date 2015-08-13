package command

type Error struct {
	Method    string `json:"method,omitempty" xml:"method,omitempty"`
	ClientErr string `json:"error,omitempty"  xml:"error,omitempty"`
	Err       error  `json:"-"                xml:"-"`
}

func (cerr Error) Error() string {
	return cerr.Err.Error()
}

type ParseError struct {
	Method      string   `json:"method,omitempty" xml:"method,omitempty"`
	ClientErr   string   `json:"error,omitempty"  xml:"error,omitempty"`
	RecordNum   int      `json:"record-number"    xml:"record-number"`
	RecordField int      `json:"record-field"     xml:"record-field"`
	Record      []string `json:"parsed-record"    xml:"parsed-record"`
	Err         error    `json:"-"                xml:"-"`
}

func (cpe ParseError) Error() string {
	return cpe.Err.Error()
}

type RunError struct {
	Method    string `json:"method,omitempty" xml:"method,omitempty"`
	ClientErr string `json:"error,omitempty"  xml:"error,omitempty"`
	Err       error  `json:"-"                xml:"-"`
	Cmd       Commander
}

func (cre RunError) Error() string {
	return cre.Err.Error()
}
