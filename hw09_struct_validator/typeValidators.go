package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var ErrUnknownIntValidator = errors.New("incorrect or unknown validator for integer types")
var ErrInvalidInt = errors.New("invalid int value")
var ErrUnknownStringValidator = errors.New("incorrect or unknown validator for strings")
var ErrInvalidString = errors.New("invalid string value")
var ErrUnknownSliceValidator = errors.New("incorrect or unknown validator for slices")
var ErrInvalidSlice = errors.New("invalid slice")

func validateString(fieldT reflect.StructField, fieldV reflect.Value) error {
	var valid bool
	validators := fieldT.Tag.Get("validate")
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "len:"):
			l, err := strconv.Atoi(validator[4:])
			if err != nil {
				return ErrUnknownStringValidator
			}
			if len(fieldV.String()) > l {
				return ErrInvalidString
			}
		case strings.HasPrefix(validator, "in:"):
			for _, val := range strings.Split(validator[3:], ",") {
				if val == fieldV.String() {
					valid = true
					break
				}
			}
			if !valid {
				return ErrInvalidString
			}
		case strings.HasPrefix(validator, "regexp:"):
			re := validator[7:]
			Re := regexp.MustCompile(re)
			if !Re.MatchString(fieldV.String()) {
				return ErrInvalidString
			}
		default:
			return ErrUnknownStringValidator
		}
	}
	return nil
}

func validateSlice(fieldT reflect.StructField, fieldV reflect.Value) error {
	validators := fieldT.Tag.Get("validate")
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "len:"):
			l, err := strconv.Atoi(validator[4:])
			if err != nil {
				return ErrUnknownSliceValidator
			}
			if fieldV.Len() > l {
				return ErrInvalidSlice
			}
		default:
			return ErrUnknownSliceValidator
		}
	}
	return nil
}

func validateInt(fieldT reflect.StructField, fieldV reflect.Value) error {
	var valid bool
	validators := fieldT.Tag.Get("validate")
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "in:"):
			for _, val := range strings.Split(validator[3:], ",") {
				i, err := strconv.Atoi(val)
				if err != nil {
					return ErrUnknownIntValidator
				}
				if i == int(fieldV.Int()) {
					valid = true
					break
				}
			}
			if !valid {
				return ErrInvalidInt
			}
		case strings.HasPrefix(validator, "min:"):
			var (
				min int
				err error
			)
			if min, err = strconv.Atoi(validator[4:]); err != nil {
				return ErrUnknownIntValidator
			}
			if int(fieldV.Int()) < min {
				return ErrInvalidInt
			}
		case strings.HasPrefix(validator, "max:"):
			var (
				max int
				err error
			)
			if max, err = strconv.Atoi(validator[4:]); err != nil {
				return ErrUnknownIntValidator
			}
			if int(fieldV.Int()) > max {
				return ErrInvalidInt
			}
		default:
			return ErrUnknownIntValidator
		}
	}
	return nil
}
