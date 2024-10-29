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

type Filter interface {
	Apply(data models.Metadata) bool
}

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
