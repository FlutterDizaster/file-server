package models

//go:generate easyjson -all -omit_empty credentials.go
type Credentials struct {
	Token    string `json:"token"`
	Login    string `json:"login"`
	Password string `json:"pswd"`
}
