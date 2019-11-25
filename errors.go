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
