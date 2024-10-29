package models

import (
	"github.com/google/uuid"
)

//easyjson:json
type Metadatas []Metadata

//go:generate easyjson -all -omit_empty metadata.go
type Metadata struct {
	ID       *uuid.UUID `json:"id"`
	Name     string     `json:"name"`
	File     bool       `json:"file"`
	Public   bool       `json:"public"`
	Mime     string     `json:"mime"`
	Created  string     `json:"created"`
	OwnerID  *uuid.UUID `json:"owner_id"`
	Grant    []string   `json:"grant"`
	JSON     JSONString `json:"json"`
	FileSize int64      `json:"file-size"`
	URL      string     `json:"url"`
}
