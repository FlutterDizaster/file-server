package filters

import (
	"slices"

	"github.com/FlutterDizaster/file-server/internal/models"
)

// GrantFilter used to filter metadata by grant.
type GrantFilter struct {
	login string
}

// NewGrantFilter creates new GrantFilter instance.
//
// login must be login of user with grant access.
//
// Returns ErrInvalidFilterValue if login is empty or invalid.
func NewGrantFilter(login string) (*GrantFilter, error) {
	return &GrantFilter{
		login: login,
	}, nil
}

// Apply implements Filter interface.
//
// Returns true if data has grant access to user with login f.login.
func (f *GrantFilter) Apply(data models.Metadata) bool {
	return slices.Contains(data.Grant, f.login)
}
