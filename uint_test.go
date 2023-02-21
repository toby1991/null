package null

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"
)

var (
	uintJSON       = []byte(`12345`)
	uintStringJSON = []byte(`"12345"`)
	nullUintJSON   = []byte(`{"Uint64":12345,"Valid":true}`)
)

func TestUintFrom(t *testing.T) {
	i := UintFrom(12345)
	assertUint(t, i, "UintFrom()")

	zero := UintFrom(0)
	if !zero.Valid {
		t.Error("UintFrom(0)", "is invalid, but should be valid")
	}
}

func TestUintFromPtr(t *testing.T) {
	n := uint64(12345)
	iptr := &n
	i := UintFromPtr(iptr)
	assertUint(t, i, "UintFromPtr()")

	null := UintFromPtr(nil)
	assertNullUint(t, null, "UintFromPtr(nil)")
}

func TestUnmarshalUint(t *testing.T) {
	var i Uint
	err := json.Unmarshal(uintJSON, &i)
	maybePanic(err)
	assertUint(t, i, "uint json")

	var si Uint
	err = json.Unmarshal(uintStringJSON, &si)
	maybePanic(err)
	assertUint(t, si, "uint string json")

	var ni Uint
	err = json.Unmarshal(nullUintJSON, &ni)
	if err == nil {
		panic("err should not be nill")
	}

	var bi Uint
	err = json.Unmarshal(floatBlankJSON, &bi)
	if err == nil {
		panic("err should not be nill")
	}

	var null Uint
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullUint(t, null, "null json")

	var badType Uint
	err = json.Unmarshal(boolJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullUint(t, badType, "wrong type json")

	var invalid Uint
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullUint(t, invalid, "invalid json")
}

func TestUnmarshalNonUintegerNumber(t *testing.T) {
	var i Uint
	err := json.Unmarshal(floatJSON, &i)
	if err == nil {
		panic("err should be present; non-integer number coerced to uint")
	}
}

func TestUnmarshalUint64Overflow(t *testing.T) {
	uint64Overflow := uint64(math.MaxUint64)

	// Max int64 should decode successfully
	var i Uint
	err := json.Unmarshal([]byte(strconv.FormatUint(uint64Overflow, 10)), &i)
	maybePanic(err)

	// Attempt to overflow
	//uint64Overflow++ // = 0
	//err = json.Unmarshal([]byte(strconv.FormatUint(uint64Overflow, 10)), &i)
	//if err == nil {
	//	panic("err should be present; decoded value overflows int64")
	//}
}

func TestTextUnmarshalUint(t *testing.T) {
	var i Uint
	err := i.UnmarshalText([]byte("12345"))
	maybePanic(err)
	assertUint(t, i, "UnmarshalText() uint")

	var blank Uint
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullUint(t, blank, "UnmarshalText() empty uint")

	var null Uint
	err = null.UnmarshalText([]byte("null"))
	maybePanic(err)
	assertNullUint(t, null, `UnmarshalText() "null"`)

	var invalid Uint
	err = invalid.UnmarshalText([]byte("hello world"))
	if err == nil {
		panic("expected error")
	}
}

func TestMarshalUint(t *testing.T) {
	i := UintFrom(12345)
	data, err := json.Marshal(i)
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty json marshal")

	// invalid values should be encoded as null
	null := NewUint(0, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "null", "null json marshal")
}

func TestMarshalUintText(t *testing.T) {
	i := UintFrom(12345)
	data, err := i.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "12345", "non-empty text marshal")

	// invalid values should be encoded as null
	null := NewUint(0, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "", "null text marshal")
}

func TestUintPointer(t *testing.T) {
	i := UintFrom(12345)
	ptr := i.Ptr()
	if *ptr != 12345 {
		t.Errorf("bad %s uint: %#v ≠ %d\n", "pointer", ptr, 12345)
	}

	null := NewUint(0, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s uint: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestUintIsZero(t *testing.T) {
	i := UintFrom(12345)
	if i.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewUint(0, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewUint(0, true)
	if zero.IsZero() {
		t.Errorf("IsZero() should be false")
	}
}

func TestUintSetValid(t *testing.T) {
	change := NewUint(0, false)
	assertNullUint(t, change, "SetValid()")
	change.SetValid(12345)
	assertUint(t, change, "SetValid()")
}

func TestUintScan(t *testing.T) {
	var i Uint
	err := i.Scan(12345)
	maybePanic(err)
	assertUint(t, i, "scanned uint")

	var null Uint
	err = null.Scan(nil)
	maybePanic(err)
	assertNullUint(t, null, "scanned null")
}

func TestUintValueOrZero(t *testing.T) {
	valid := NewUint(12345, true)
	if valid.ValueOrZero() != 12345 {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewUint(12345, false)
	if invalid.ValueOrZero() != 0 {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestUintEqual(t *testing.T) {
	int1 := NewUint(10, false)
	int2 := NewUint(10, false)
	assertUintEqualIsTrue(t, int1, int2)

	int1 = NewUint(10, false)
	int2 = NewUint(20, false)
	assertUintEqualIsTrue(t, int1, int2)

	int1 = NewUint(10, true)
	int2 = NewUint(10, true)
	assertUintEqualIsTrue(t, int1, int2)

	int1 = NewUint(10, true)
	int2 = NewUint(10, false)
	assertUintEqualIsFalse(t, int1, int2)

	int1 = NewUint(10, false)
	int2 = NewUint(10, true)
	assertUintEqualIsFalse(t, int1, int2)

	int1 = NewUint(10, true)
	int2 = NewUint(20, true)
	assertUintEqualIsFalse(t, int1, int2)
}

func assertUint(t *testing.T, i Uint, from string) {
	if i.String != "12345" {
		t.Errorf("bad %s uint: %s ≠ %s\n", from, i.String, "12345")
	}
	if !i.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullUint(t *testing.T, i Uint, from string) {
	if i.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertUintEqualIsTrue(t *testing.T, a, b Uint) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Uint{%v, Valid:%t} and Uint{%v, Valid:%t} should return true", a.String, a.Valid, b.String, b.Valid)
	}
}

func assertUintEqualIsFalse(t *testing.T, a, b Uint) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Uint{%v, Valid:%t} and Uint{%v, Valid:%t} should return false", a.String, a.Valid, b.String, b.Valid)
	}
}
