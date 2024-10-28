package models

import (
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

type JSONString string

func (j *JSONString) UnmarshalJSON(b []byte) error {
	*j = JSONString(b)
	return nil
}

func (j JSONString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

func (j JSONString) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawString(string(j))
}

func (j *JSONString) UnmarshalEasyJSON(l *jlexer.Lexer) {
	*j = JSONString(l.Data)
}
