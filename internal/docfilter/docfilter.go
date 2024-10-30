package docfilter

import (
	"github.com/FlutterDizaster/file-server/internal/docfilter/filters"
	"github.com/FlutterDizaster/file-server/internal/models"
)

// DocumentsFilter used to filter metadata by given filters.
//
// If no filters are provided, all metadata will be returned.
// Filters can be added with AddFilter function.
//
// FilterData function returns filtered metadata.
//
// Must be initialized with New function.
type DocumentsFilter struct {
	limit  int
	offset int

	filters []filters.Filter
}

// New creates new DocumentsFilter instance.
//
// limit and offset are used to limit and offset returned metadata.
//
// Returns *DocumentsFilter.
func New(limit, offset int) *DocumentsFilter {
	return &DocumentsFilter{
		limit:  limit,
		offset: offset,
	}
}

// AddFilter adds a new filter to the DocumentsFilter based on the provided key and value.
//
// key specifies the type of filter to apply and must be one of the predefined filter keys.
// value is used to configure the filter and should match the expected format for the specified key.
//
// All filters must be added before calling FilterData.
//
// you can find all available filter keys and their corresponding values in the filters package.
//
// Returns an error if the key is unknown or the value is invalid for the specified filter.
func (f *DocumentsFilter) AddFilter(key, value string) error {
	filter, err := filters.QuerryFilter(key, value)
	if err != nil {
		return err
	}

	f.filters = append(f.filters, filter)

	return nil
}

// FilterData filters the given slice of metadata according to the filters
// set in the DocumentsFilter and returns a new slice of filtered metadata.
//
// The function will return up to `limit` number of metadata items that
// match all the filters. The `offset` is used to skip the first `offset`
// matching metadata.
//
// The method will return an empty slice if no metadata match the filters
// or if the `limit` is set to 0.
//
// The function will not modify the original slice.
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
