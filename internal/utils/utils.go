package utils

import (
	"fmt"
	"frog-go/internal/core/errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ToUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

func ToUUID(str string) (uuid.UUID, error) {
	if str == "" {
		return uuid.UUID{}, errors.ErrEmptyField
	}

	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, err
	}

	return parsedUUID, nil
}

func ToUUIDSlice(strs []string) []uuid.UUID {
	var result []uuid.UUID
	for _, s := range strs {
		if id, err := uuid.Parse(s); err == nil {
			result = append(result, id)
		}
	}
	return result
}

func ToNillableUUID(str string) (*uuid.UUID, error) {
	if str == "" {
		return nil, nil
	}

	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return nil, err
	}

	return &parsedUUID, nil
}

func ToDateTime(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, errors.ErrEmptyField

	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ToDateTimeUnsafe(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func ToNillableDateTime(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ToDateTimeString(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func ToNillableDateTimeString(date *time.Time) *string {
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
