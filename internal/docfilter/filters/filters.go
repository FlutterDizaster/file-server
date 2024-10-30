package filters

import (
	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type FilterKey string

const (
	FilterKeyOwner  FilterKey = "owner"
	FilterKeyFile   FilterKey = "file"
	FilterKeyName   FilterKey = "name"
	FilterKeyMime   FilterKey = "mime"
	FilterKeyPublic FilterKey = "public"
	FilterKeyDate   FilterKey = "created"
	FilterKeyGrant  FilterKey = "grant"
	FilterKeyID     FilterKey = "id"
)

// Filter used to filter metadata.
type Filter interface {
	Apply(data models.Metadata) bool
}

// QuerryFilter creates new Filter instance based on given key and value.
//
// key:
//   - "owner" : filter by owner id.
//   - "file" : filter by file status.
//   - "name" : filter by name.
//   - "mime" : filter by mime type.
//   - "public" : filter by public status.
//   - "created" : filter by creation date.
//   - "grant" : filter by user login.
//   - "id" : filter by document id.
//
// value format depends on filter type.
//
// Returns ErrUnknownFilter if key is unknown.
func QuerryFilter(key, value string) (Filter, error) {
	switch FilterKey(key) {
	case FilterKeyOwner:
		return NewOwnerFilter(value)
	case FilterKeyFile:
		return NewFileFilter(value)
	case FilterKeyName:
		return NewNameFilter(value)
	case FilterKeyMime:
		return NewMimeFilter(value)
	case FilterKeyPublic:
		return NewPublicFilter(value)
	case FilterKeyDate:
		return NewDateFilter(value)
	case FilterKeyGrant:
		return NewGrantFilter(value)
	case FilterKeyID:
		return NewIDFilter(value)
	default:
		return nil, apperrors.ErrUnknownFilter
	}
}
