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
