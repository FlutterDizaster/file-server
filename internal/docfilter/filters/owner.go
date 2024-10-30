package filters

import (
	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

// OwnerFilter used to filter metadata by owner id.
type OwnerFilter struct {
	id uuid.UUID
}

// NewOwnerFilter creates a new OwnerFilter instance with the given owner ID value.
// The value should be a valid UUID string. If the value is invalid, it returns
// ErrInvalidFilterValue.
func NewOwnerFilter(value string) (*OwnerFilter, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &OwnerFilter{
		id: id,
	}, nil
}

// Apply checks if the owner ID in the given metadata matches the owner ID
// specified in the OwnerFilter. Returns true if they match, otherwise false.
func (f *OwnerFilter) Apply(meta models.Metadata) bool {
	return f.id == *meta.OwnerID
}
