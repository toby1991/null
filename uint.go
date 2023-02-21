package null

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Uint is an nullable uint64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint struct {
	sql.NullString
}

// NewInt creates a new Uint
func NewUint(i uint64, valid bool) Uint {
	return Uint{
		NullString: sql.NullString{
			String: strconv.FormatUint(i, 10),
			Valid:  valid,
		},
	}
}

// IntFrom creates a new Uint that will always be valid.
func UintFrom(i uint64) Uint {
	return NewUint(i, true)
}

// IntFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint64) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	return NewUint(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Uint) ValueOrZero() uint64 {
	if !i.Valid {
		return 0
	}
	parseUint, _ := strconv.ParseUint(i.String, 10, 64)
	return parseUint
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number, string, and null input.
// 0 will not be considered a null Uint.
func (i *Uint) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		i.Valid = false
		return nil
	}

	var _n uint64
	if err := json.Unmarshal(data, &_n); err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			// special case: accept string input
			if typeError.Value != "string" {
				return fmt.Errorf("null: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("null: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return fmt.Errorf("null: couldn't convert string to int: %w", err)
			}
			i.String = strconv.FormatUint(n, 10)
			i.Valid = true
			return nil
		}
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}

	i.String = strconv.FormatUint(_n, 10)
	i.Valid = true
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	n, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't convert string to int: %w", err)
	}
	i.String = strconv.FormatUint(n, 10)
	i.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint is null.
func (i Uint) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}

	return []byte(i.String), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint is null.
func (i Uint) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(i.String), nil
}

// SetValid changes this Uint's value and also sets it to be non-null.
func (i *Uint) SetValid(n uint64) {
	i.String = strconv.FormatUint(n, 10)
	i.Valid = true
}

// Ptr returns a pointer to this Uint's value, or a nil pointer if this Uint is null.
func (i Uint) Ptr() *uint64 {
	if !i.Valid {
		return nil
	}

	parseUint, err := strconv.ParseUint(i.String, 10, 64)
	if err != nil {
		return nil
	}
	return &parseUint
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Uint with a 0 value will not be considered zero.
func (i Uint) IsZero() bool {
	return !i.Valid
}

// Equal returns true if both ints have the same value or are both null.
func (i Uint) Equal(other Uint) bool {
	return i.Valid == other.Valid && (!i.Valid || i.String == other.String)
}
