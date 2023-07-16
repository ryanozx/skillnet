package helpers

import (
	"errors"
	"strconv"
)

var (
	ErrIsNull = errors.New("null value")
)

type NullableUint struct {
	value    uint
	hasValue bool
}

func ParseNullableUint(str string) (*NullableUint, error) {
	output := NullableUint{}
	if str == "" {
		output.hasValue = false
		return &output, nil
	}
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return &output, err
	}
	output.value = uint(val)
	output.hasValue = true
	return &output, nil
}

func (v *NullableUint) IsNull() bool {
	return !v.hasValue
}

func (v *NullableUint) GetValue() (uint, error) {
	if v.IsNull() {
		return 0, ErrIsNull
	}
	return v.value, nil
}
