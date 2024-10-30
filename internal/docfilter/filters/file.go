package filters

import (
	"strconv"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

// FileFilter used to filter metadata by file flag.
type FileFilter struct {
	isFile bool
}

// NewFileFilter creates new FileFilter instance.
//
// value must be 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False..
//
// Returns ErrInvalidFilterValue if value is invalid.
func NewFileFilter(value string) (*FileFilter, error) {
	isFile, err := strconv.ParseBool(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &FileFilter{
		isFile: isFile,
	}, nil
}

// Apply filter to given metadata.
//
// Return true if metadata is valid for given filter, false otherwise.
func (f *FileFilter) Apply(data models.Metadata) bool {
	return f.isFile == data.File
}
