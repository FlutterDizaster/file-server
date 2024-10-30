package filters

import (
	"strings"

	"github.com/FlutterDizaster/file-server/internal/models"
)

// NameFilter used to filter metadata by name.
type NameFilter struct {
	name string
}

// NewNameFilter creates a new NameFilter instance with the given name value.
// The value can have a wildcard prefix or suffix, indicated by '*', to match
// the start or end of a name, respectively.
func NewNameFilter(value string) (*NameFilter, error) {
	return &NameFilter{
		name: value,
	}, nil
}

// Apply implements the Filter interface and returns true if the given
// metadata's name matches the given name value, either directly or with
// wildcard prefix or suffix.
func (f *NameFilter) Apply(data models.Metadata) bool {
	if strings.HasPrefix(f.name, "*") {
		return strings.HasSuffix(data.Name, f.name[1:])
	}

	if strings.HasSuffix(f.name, "*") {
		return strings.HasPrefix(data.Name, f.name[:len(f.name)-1])
	}

	return strings.Contains(data.Name, f.name)
}
