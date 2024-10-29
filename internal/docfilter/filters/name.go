package filters

import (
	"strings"

	"github.com/FlutterDizaster/file-server/internal/models"
)

type NameFilter struct {
	name string
}

func NewNameFilter(value string) (*NameFilter, error) {
	return &NameFilter{
		name: value,
	}, nil
}

func (f *NameFilter) Apply(data models.Metadata) bool {
	if strings.HasPrefix(f.name, "*") {
		return strings.HasSuffix(data.Name, f.name[1:])
	}

	if strings.HasSuffix(f.name, "*") {
		return strings.HasPrefix(data.Name, f.name[:len(f.name)-1])
	}

	return strings.Contains(data.Name, f.name)
}
