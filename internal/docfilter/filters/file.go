package filters

import (
	"strconv"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type FileFilter struct {
	isFile bool
}

func NewFileFilter(value string) (*FileFilter, error) {
	isFile, err := strconv.ParseBool(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return &FileFilter{
		isFile: isFile,
	}, nil
}

func (f *FileFilter) Apply(data models.Metadata) bool {
	return f.isFile == data.File
}
