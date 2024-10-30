package filters

import (
	"strconv"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

// PublicFilter used to filter metadata by public flag.
type PublicFilter struct {
	isPublic bool
}

// NewPublicFilter creates new PublicFilter instance.
//
// value must be 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// Returns ErrInvalidFilterValue if value is invalid.
func NewPublicFilter(isPublic string) (*PublicFilter, error) {
	b, err := strconv.ParseBool(isPublic)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &PublicFilter{
		isPublic: b,
	}, nil
}

// Apply implements filters.Filter interface.
//
// It returns true if the document is public according to filter and false otherwise.
func (f *PublicFilter) Apply(data models.Metadata) bool {
	return f.isPublic == data.Public
}
