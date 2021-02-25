package hw09structvalidator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnknownIntValidator = errors.New("incorrect or unknown validator for integer types")
var ErrInvalidIntValue = errors.New("invalid int value")
var ErrInvalidIntMin = errors.New("int value is lesser than limit")
var ErrInvalidIntMax = errors.New("int value exceeds limit")
var ErrUnknownStringValidator = errors.New("incorrect or unknown validator for strings")
var ErrInvalidStringValue = errors.New("invalid string value")
var ErrInvalidStringRegexp = errors.New("string doesn't match regexp")
var ErrInvalidStringLength = errors.New("string length exceeds the limit")
var ErrInvalidSlice = errors.New("invalid slice")

func validateString(field reflect.Value, validators string) error {
	var valid bool
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "len:"):
			l, err := strconv.Atoi(validator[4:])
			if err != nil {
				return ErrUnknownStringValidator
			}
			if len(field.String()) > l {
				return ErrInvalidStringLength
			}
		case strings.HasPrefix(validator, "in:"):
			for _, val := range strings.Split(validator[3:], ",") {
				if val == field.String() {
					valid = true
					break
				}
			}
			if !valid {
				return ErrInvalidStringValue
			}
		case strings.HasPrefix(validator, "regexp:"):
			re := validator[7:]
			Re := regexp.MustCompile(re)
			if !Re.MatchString(field.String()) {
				return ErrInvalidStringRegexp
			}
		default:
			return ErrUnknownStringValidator
		}
	}
	return nil
}

func validateSlice(field reflect.Value, validators string) error {
	slice := field.Interface()
	var errs []struct {
		i int
		e error
	}
	switch t := slice.(type) {
	case []string:
		for i, elem := range t {
			err := validateString(reflect.ValueOf(elem), validators)
			if err != nil {
				errs = append(errs, struct {
					i int
					e error
				}{i: i, e: err})
			}
		}
	case []int:
		for i, elem := range t {
			err := validateInt(reflect.ValueOf(elem), validators)
			if err != nil {
				errs = append(errs, struct {
					i int
					e error
				}{i: i, e: err})
			}
		}
	default:
		return ErrUnsupportedType
	}
	if errs != nil {
		var combinedError string
		for _, err := range errs {
			combinedError += fmt.Sprintf("%d: %s\n", err.i, err.e.Error())
		}
		return fmt.Errorf("%w\n%s", ErrInvalidSlice, combinedError)
	}
	return nil
}

func validateInt(field reflect.Value, validators string) error {
	var valid bool
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "in:"):
			for _, val := range strings.Split(validator[3:], ",") {
				i, err := strconv.Atoi(val)
				if err != nil {
					return ErrUnknownIntValidator
				}
				if i == int(field.Int()) {
					valid = true
					break
				}
			}
			if !valid {
				return ErrInvalidIntValue
			}
		case strings.HasPrefix(validator, "min:"):
			var (
				min int
				err error
			)
			if min, err = strconv.Atoi(validator[4:]); err != nil {
				return ErrUnknownIntValidator
			}
			if int(field.Int()) < min {
				return ErrInvalidIntMin
			}
		case strings.HasPrefix(validator, "max:"):
			var (
				max int
				err error
			)
			if max, err = strconv.Atoi(validator[4:]); err != nil {
				return ErrUnknownIntValidator
			}
			if int(field.Int()) > max {
				return ErrInvalidIntMax
			}
		default:
			return ErrUnknownIntValidator
		}
	}
	return nil
}
