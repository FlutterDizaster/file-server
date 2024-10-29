package filters

import (
	"strings"

	"github.com/FlutterDizaster/file-server/internal/models"
)

type MimeFilter struct {
	mime string
}

func NewMimeFilter(value string) (*MimeFilter, error) {
	return &MimeFilter{
		mime: value,
	}, nil
}

func (f *MimeFilter) Apply(meta models.Metadata) bool {
	return strings.Contains(meta.Mime, f.mime)
}
