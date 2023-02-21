package zero

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Uint is a nullable uint64.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Uint struct {
	sql.NullString
}

// NewUint creates a new Uint
func NewUint(i uint64, valid bool) Uint {
	return Uint{
		NullString: sql.NullString{
			String: strconv.FormatUint(i, 10),
			Valid:  valid,
		},
	}
}

// UintFrom creates a new Uint that will be null if zero.
func UintFrom(i uint64) Uint {
	return NewUint(i, i != 0)
}

// UintFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint64) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	n := NewUint(*i, true)
	return n
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
// It supports number and null input.
// 0 will be considered a null Uint.
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
				return fmt.Errorf("zero: JSON input is invalid type (need int or string): %w", err)
			}
			var str string
			if err := json.Unmarshal(data, &str); err != nil {
				return fmt.Errorf("zero: couldn't unmarshal number string: %w", err)
			}
			n, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return fmt.Errorf("zero: couldn't convert string to int: %w", err)
			}
			i.String = strconv.FormatUint(n, 10)
			i.Valid = n != 0
			return nil
		}
		return fmt.Errorf("zero: couldn't unmarshal JSON: %w", err)
	}

	i.String = strconv.FormatUint(_n, 10)
	i.Valid = _n != 0
	return nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint if the input is a blank, or zero.
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
	i.Valid = n != 0
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Uint is null.
func (i Uint) MarshalJSON() ([]byte, error) {
	parseUint, err := strconv.ParseUint(i.String, 10, 64)
	if err != nil {
		return nil, err
	}
	n := parseUint
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(n, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Uint is null.
func (i Uint) MarshalText() ([]byte, error) {
	parseUint, err := strconv.ParseUint(i.String, 10, 64)
	if err != nil {
		return nil, err
	}
	n := parseUint
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatUint(n, 10)), nil
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

	// todo: may cause ptr error
	parseUint, err := strconv.ParseUint(i.String, 10, 64)
	if err != nil {
		return nil
	}
	return &parseUint
}

// IsZero returns true for null or zero Uints, for future omitempty support (Go 1.4?)
func (i Uint) IsZero() bool {
	n, err := strconv.ParseUint(i.String, 10, 64)
	if err != nil {
		return true
	}
	return !i.Valid || n == 0
}

// Equal returns true if both ints have the same value or are both either null or zero.
func (i Uint) Equal(other Uint) bool {
	return i.ValueOrZero() == other.ValueOrZero()
}
