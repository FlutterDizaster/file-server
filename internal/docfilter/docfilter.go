package docfilter

import (
	"github.com/FlutterDizaster/file-server/internal/docfilter/filters"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type DocumentsFilter struct {
	limit  int
	offset int

	filters []filters.Filter
}

func New(limit, offset int) *DocumentsFilter {
	return &DocumentsFilter{
		limit:  limit,
		offset: offset,
	}
}

func (f *DocumentsFilter) AddFilter(key, value string) error {
	filter, err := filters.QuerryFilter(key, value)
	if err != nil {
		return err
	}

	f.filters = append(f.filters, filter)

	return nil
}

func (f DocumentsFilter) FilterData(data []models.Metadata) []models.Metadata {
	result := make([]models.Metadata, 0, f.limit)

	for i := 0; len(result) < f.limit && i < len(data); i++ {
		pass := true

		for _, filter := range f.filters {
			if !filter.Apply(data[i]) {
				pass = false
				break
			}
		}

		if pass && i >= f.offset {
			result = append(result, data[i])
		}
	}

	return result
}
