package hw09structvalidator

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
	var err error
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "len:"):
			err = validateStringLen(field.String(), validator[4:])
		case strings.HasPrefix(validator, "in:"):
			err = validateStringInValues(field.String(), validator[3:])
		case strings.HasPrefix(validator, "regexp:"):
			err = validateStringMatchRe(field.String(), validator[7:])
		default:
			err = ErrUnknownStringValidator
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func validateInt(field reflect.Value, validators string) error {
	var err error
	for _, validator := range strings.Split(validators, "|") {
		switch {
		case strings.HasPrefix(validator, "in:"):
			err = validateIntIn(int(field.Int()), validator[3:])
		case strings.HasPrefix(validator, "min:"):
			err = validateIntMin(int(field.Int()), validator[4:])
		case strings.HasPrefix(validator, "max:"):
			err = validateIntMax(int(field.Int()), validator[4:])
		default:
			err = ErrUnknownIntValidator
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type sliceErr struct {
	i int
	e error
}

func validateSlice(field reflect.Value, validators string) error {
	slice := field.Interface()
	var errs []sliceErr
	switch t := slice.(type) {
	case []string:
		errs = validateSliceOfStrings(t, validators)
	case []int:
		errs = validateSliceOfInt(t, validators)
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

func validateSliceOfStrings(t []string, validators string) []sliceErr {
	var errs []sliceErr
	for i, elem := range t {
		err := validateString(reflect.ValueOf(elem), validators)
		if err != nil {
			errs = append(errs, sliceErr{i: i, e: err})
		}
	}
	return errs
}

func validateSliceOfInt(t []int, validators string) []sliceErr {
	var errs []sliceErr
	for i, elem := range t {
		err := validateInt(reflect.ValueOf(elem), validators)
		if err != nil {
			errs = append(errs, sliceErr{i: i, e: err})
		}
	}
	return errs
}

func validateStringLen(s, validator string) error {
	l, err := strconv.Atoi(validator)
	if err != nil {
		return ErrUnknownStringValidator
	}
	if len(s) > l {
		return ErrInvalidStringLength
	}
	return nil
}

func validateStringInValues(s, validator string) error {
	values := strings.Split(validator, ",")
	var valid bool
	for _, val := range values {
		if val == s {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidStringValue
	}
	return nil
}

func validateStringMatchRe(s, re string) error {
	Re := regexp.MustCompile(re)
	if !Re.MatchString(s) {
		return ErrInvalidStringRegexp
	}
	return nil
}

func validateIntIn(value int, validator string) error {
	var valid bool
	for _, val := range strings.Split(validator, ",") {
		i, err := strconv.Atoi(val)
		if err != nil {
			return ErrUnknownIntValidator
		}
		if i == value {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidIntValue
	}
	return nil
}

func validateIntMin(value int, validator string) error {
	min, err := strconv.Atoi(validator)
	if err != nil {
		return ErrUnknownIntValidator
	}
	if value < min {
		return ErrInvalidIntMin
	}
	return nil
}

func validateIntMax(value int, validator string) error {
	max, err := strconv.Atoi(validator)
	if err != nil {
		return ErrUnknownIntValidator
	}
	if value > max {
		return ErrInvalidIntMax
	}
	return nil
}
