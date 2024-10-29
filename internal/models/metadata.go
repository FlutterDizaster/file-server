package models

import (
	"time"

	"github.com/google/uuid"
)

//go:generate easyjson -all -omit_empty metadata.go
type Metadata struct {
	ID       *uuid.UUID `json:"id"`
	Name     string     `json:"name"`
	File     bool       `json:"file"`
	Public   bool       `json:"public"`
	Mime     string     `json:"mime"`
	Created  *time.Time `json:"created"`
	OwnerID  *uuid.UUID `json:"owner_id"`
	Grant    []string   `json:"grant"`
	JSON     JSONString `json:"json"`
	FileSize int64      `json:"file-size"`
	URL      string     `json:"url"`
}
