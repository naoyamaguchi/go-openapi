package openapi

import (
	"fmt"
	"strconv"
)

type requiredError struct {
	RequiredField string
}

func ErrRequired(requiredField string) error {
	return requiredError{
		RequiredField: requiredField,
	}
}

func (e requiredError) Error() string {
	return fmt.Sprintf("%s field is required", strconv.Quote(e.RequiredField))
}

func (e requiredError) Is(target error) bool {
	if target, ok := target.(requiredError); ok {
		return e.RequiredField == target.RequiredField
	}
	return false
}

type unknownKeyError struct {
	Key string
}

func ErrUnknownKey(key string) error {
	return unknownKeyError{
		Key: key,
	}
}

func (e unknownKeyError) Error() string {
	return fmt.Sprintf("unknown key: %s", e.Key)
}

func (e unknownKeyError) Is(target error) bool {
	if target, ok := target.(unknownKeyError); ok {
		return e.Key == target.Key
	}
	return false
}
