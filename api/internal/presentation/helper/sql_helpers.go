package helper

import (
	"database/sql"
	"time"
)

// NewNullString creates a sql.NullString from a string pointer.
func NewNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// NewNullTime creates a sql.NullTime from a time.Time pointer.
func NewNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// NewNullInt32 creates a sql.NullInt32 from an int32 pointer.
func NewNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

// NewNullFloat64 creates a sql.NullFloat64 from a float64 pointer.
func NewNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}
