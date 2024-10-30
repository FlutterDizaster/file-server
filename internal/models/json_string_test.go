package models

import (
	"testing"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"github.com/stretchr/testify/assert"
)

func TestJSONString_UnmarshalEasyJSON(t *testing.T) {
	type test struct {
		name string
		j    JSONString
		want JSONString
	}
	tests := []test{
		{
			name: "test",
			j:    JSONString("test"),
			want: JSONString("test"),
		},
		{
			name: "empty",
			j:    JSONString(""),
			want: JSONString(""),
		},
		{
			name: "json",
			j:    JSONString("{\"test\":\"test\"}"),
			want: JSONString("{\"test\":\"test\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := jlexer.Lexer{Data: []byte(tt.j)}
			tt.j.UnmarshalEasyJSON(&l)

			assert.Equal(t, tt.want, tt.j)
		})
	}
}
func TestJSONString_MarshalEasyJSON(t *testing.T) {
	type test struct {
		name string
		j    JSONString
		want JSONString
	}
	tests := []test{
		{
			name: "test",
			j:    JSONString("test"),
			want: JSONString("test"),
		},
		{
			name: "empty",
			j:    JSONString(""),
			want: JSONString(""),
		},
		{
			name: "json",
			j:    JSONString("{\"test\":\"test\"}"),
			want: JSONString("{\"test\":\"test\"}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := jwriter.Writer{}
			tt.j.MarshalEasyJSON(&w)

			result := JSONString(w.Buffer.BuildBytes())

			assert.Equal(t, tt.want, result)
		})
	}
}
