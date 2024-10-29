package filters

import (
	"slices"

	"github.com/FlutterDizaster/file-server/internal/models"
)

type GrantFilter struct {
	login string
}

func NewGrantFilter(login string) (*GrantFilter, error) {
	return &GrantFilter{
		login: login,
	}, nil
}

func (f *GrantFilter) Apply(data models.Metadata) bool {
	return slices.Contains(data.Grant, f.login)
}
