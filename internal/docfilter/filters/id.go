package filters

import (
	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

// IDFilter used to filter metadata by id.
type IDFilter struct {
	id uuid.UUID
}

// NewIDFilter creates new IDFilter instance.
//
// id must be valid uuid.
//
// Returns ErrInvalidFilterValue if id is invalid.
func NewIDFilter(id string) (IDFilter, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return IDFilter{}, apperrors.ErrInvalidFilterValue
	}

	return IDFilter{
		id: parsedID,
	}, nil
}

// Apply implements filters.Filter interface.
//
// It checks if given metadata id equals to filter id.
func (f IDFilter) Apply(meta models.Metadata) bool {
	return f.id == *meta.ID
}
