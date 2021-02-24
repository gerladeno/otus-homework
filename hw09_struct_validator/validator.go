package hw09_struct_validator //nolint:golint,stylecheck
import (
	"errors"
	"fmt"
	"reflect"
)

var ErrNotStruct = errors.New("not a struct")
var ErrUnsupportedType = errors.New("not supported field type")

type ValidationError struct {
	Field string
	Err   error
}

type Tag struct {
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errStr string
	for _, err := range v {
		errStr += fmt.Sprintf("%s: %s", err.Err.Error(), err.Field)
	}
	return errStr
}

func (v *ValidationErrors) add(err ValidationError) {
	*v = append(*v, err)
}

func Validate(v interface{}) error {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var errs ValidationErrors
	validateFields(v, "", &errs)
	return errs
}

func validateFields(v interface{}, currentPath string, errs *ValidationErrors) {
	reflectValue := reflect.ValueOf(v)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)
		var fieldPath string
		if currentPath == "" {
			fieldPath = field.Name
		} else {
			fieldPath = currentPath + "." + field.Name
		}
		if field.Type.Kind() == reflect.Struct {
			validateFields(reflectValue.Field(i).Interface(), fieldPath, errs)
			continue
		}
		if validators, ok := field.Tag.Lookup("validate"); ok {
			err := validate(reflectValue.Field(i), validators)
			if err != nil {
				errs.add(ValidationError{
					Field: fieldPath,
					Err:   err,
				})
			}
		}
	}
}

func validate(field reflect.Value, validators string) error {
	var err error
	switch field.Kind() {
	case reflect.String:
		err = validateString(field, validators)
	case reflect.Slice:
		err = validateSlice(field, validators)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err = validateInt(field, validators)
	default:
		return ErrUnsupportedType
	}
	return err
}
