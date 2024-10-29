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

type DateFilter struct {
	date    time.Time
	endDate time.Time
	mode    dateFilterMode
}

func NewDateFilter(value string) (Filter, error) {
	f := &DateFilter{}

	err := f.parseFilterValue(value)
	if err != nil {
		return nil, apperrors.ErrInvalidFilterValue
	}

	return f, nil
}

func (f *DateFilter) Apply(data models.Metadata) bool {
	switch f.mode {
	case dateFilterModeAfter:
		return data.Created.After(f.date)
	case dateFilterModeBefore:
		return data.Created.Before(f.date)
	case dateFilterModeEqual:
		return data.Created.Equal(f.date)
	case dateFilterModeBetween:
		return data.Created.After(f.date) && data.Created.Before(f.endDate)
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
