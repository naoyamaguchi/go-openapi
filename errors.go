package openapi

import (
	"fmt"
	"strconv"
)

type RequiredError struct {
	RequiredField string
}

func ErrRequired(requiredField string) error {
	return RequiredError{
		RequiredField: requiredField,
	}
}

func (e RequiredError) Error() string {
	return fmt.Sprintf("%s field is required", strconv.Quote(e.RequiredField))
}

func (e RequiredError) Is(target error) bool {
	if target, ok := target.(RequiredError); ok {
		return e.RequiredField == target.RequiredField
	}
	return false
}

type UnknownKeyError struct {
	Key string
}

func ErrUnknownKey(key string) error {
	return UnknownKeyError{
		Key: key,
	}
}

func (e UnknownKeyError) Error() string {
	return fmt.Sprintf("unknown key: %s", e.Key)
}

func (e UnknownKeyError) Is(target error) bool {
	if target, ok := target.(UnknownKeyError); ok {
		return e.Key == target.Key
	}
	return false
}
