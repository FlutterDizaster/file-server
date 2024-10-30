package models

//go:generate easyjson -all -omit_empty requests.go
type FilesListRequest struct {
	Login  string `json:"login"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
