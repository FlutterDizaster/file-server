package filters

import (
	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

type IDFilter struct {
	id uuid.UUID
}

func NewIDFilter(id string) (IDFilter, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return IDFilter{}, apperrors.ErrInvalidFilterValue
	}

	return IDFilter{
		id: parsedID,
	}, nil
}

func (f IDFilter) Apply(meta models.Metadata) bool {
	return f.id == *meta.ID
}
