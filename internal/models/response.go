package models

import (
	"github.com/mailru/easyjson"
)

//go:generate easyjson -all -omit_empty response.go
type Response struct {
	Error    *ResponseError                `json:"error"`
	Response easyjson.MarshalerUnmarshaler `json:"response"`
	Data     easyjson.MarshalerUnmarshaler `json:"data"`
}

type ResponseUploading struct {
	JSON JSONString `json:"json"`
	File string     `json:"file"`
}

type ResponseFilesList struct {
	Docs []Metadata `json:"docs"`
}

type ResponseError struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func NewErrorResponseJSON(code int, text string) ([]byte, error) {
	resp := &Response{
		Error: &ResponseError{
			Code: code,
			Text: text,
		},
	}

	respData, err := resp.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return respData, nil
}
