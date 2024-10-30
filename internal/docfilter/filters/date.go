package filters

import (
	"strings"
	"time"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type dateFilterMode int

const (
	dateFilterModeAfter dateFilterMode = iota
	dateFilterModeBefore
	dateFilterModeEqual
	dateFilterModeBetween
)

// DateFilter used to filter documents by date.
// Must be initialized with NewDateFilter function.
//
// Value formats:
//
//   - "<date>" : after date.
//   - "<date>" : before date.
//   - "=<date>" : equal (exact date).
//   - "<date>~<date>" : between dates.
//
// Date format is "2006-01-02 15:04:05".
type DateFilter struct {
	date    time.Time
	endDate time.Time
	mode    dateFilterMode
}

// NewDateFilter creates new DateFilter instance.
//
// value format:
//
//   - "<date>" : after date.
//   - "<date>" : before date.
//   - "=<date>" : equal (exact date).
//   - "<date>~<date>" : between dates.
//
// Date format is "2006-01-02 15:04:05".
//
// Returns ErrInvalidFilterValue if value is invalid.
func NewDateFilter(value string) (Filter, error) {
	f := &DateFilter{}

	err := f.parseFilterValue(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return f, nil
}

// Apply checks if given metadata date matches filter value.
//
// Returns true if filter matches, false otherwise.
//
// If metadata date is invalid, returns false.
func (f *DateFilter) Apply(data models.Metadata) bool {
	date, err := time.Parse(time.DateTime, data.Created)
	if err != nil {
		return false
	}

	switch f.mode {
	case dateFilterModeAfter:
		return date.After(f.date)
	case dateFilterModeBefore:
		return date.Before(f.date)
	case dateFilterModeEqual:
		return date.Equal(f.date)
	case dateFilterModeBetween:
		return date.After(f.date) && date.Before(f.endDate)
	}
	return false
}

func (f *DateFilter) parseFilterValue(value string) error {
	switch {
	case strings.HasPrefix(value, ">"):
		f.mode = dateFilterModeAfter

	case strings.HasPrefix(value, "<"):
		f.mode = dateFilterModeBefore

	case strings.HasPrefix(value, "="):
		f.mode = dateFilterModeEqual

	case strings.Contains(value, "~"):
		f.mode = dateFilterModeBetween
		//nolint:gocritic // value difinitely contains "~"
		date, err := time.Parse(time.DateTime, value[:strings.Index(value, "~")])
		if err != nil {
			return err
		}
		f.date = date

		f.endDate, err = time.Parse(time.DateTime, value[strings.Index(value, "~")+1:])
		if err != nil {
			return err
		}
		return nil

	default:
		return apperrors.ErrInvalidFilterValue
	}

	date, err := time.Parse(time.DateTime, value[1:])
	if err != nil {
		return err
	}
	f.date = date

	return nil
}
