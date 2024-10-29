package filters

import (
	"strconv"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type PublicFilter struct {
	isPublic bool
}

func NewPublicFilter(isPublic string) (*PublicFilter, error) {
	b, err := strconv.ParseBool(isPublic)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &PublicFilter{
		isPublic: b,
	}, nil
}

func (f *PublicFilter) Apply(data models.Metadata) bool {
	return f.isPublic == data.Public
}
