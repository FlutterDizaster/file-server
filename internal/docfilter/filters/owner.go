package filters

import (
	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

type OwnerFilter struct {
	id uuid.UUID
}

func NewOwnerFilter(value string) (*OwnerFilter, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &OwnerFilter{
		id: id,
	}, nil
}

func (f *OwnerFilter) Apply(meta models.Metadata) bool {
	return f.id == *meta.OwnerID
}
