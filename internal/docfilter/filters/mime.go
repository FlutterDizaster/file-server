package filters

import (
	"strings"

	"github.com/FlutterDizaster/file-server/internal/models"
)

// MimeFilter used to filter metadata by mime type.
type MimeFilter struct {
	mime string
}

// NewMimeFilter creates new MimeFilter instance.
//
// value must be mime type to filter by.
//
// Returns ErrInvalidFilterValue if value is invalid.
func NewMimeFilter(value string) (*MimeFilter, error) {
	return &MimeFilter{
		mime: value,
	}, nil
}

// Apply checks if the mime type in the given metadata contains the mime type
// specified in the MimeFilter. Returns true if it contains the mime type,
// otherwise false.
func (f *MimeFilter) Apply(meta models.Metadata) bool {
	return strings.Contains(meta.Mime, f.mime)
}
