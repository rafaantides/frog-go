package utils

import (
	"fmt"
	"frog-go/internal/core/errors"
	"strconv"
	"time"
)

func ToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func ToDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.ErrEmptyField

	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ToDateUnsafe(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func ToNillableDate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ToDateTimeString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func SafeToNillableDateTimeString(date *time.Time) *string {
	if date == nil || date.IsZero() {
		return nil
	}
	formatted := date.Format("2006-01-02")
	return &formatted
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}
