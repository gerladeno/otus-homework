package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	Embedded struct {
		Value1    int `validate:"min:18"`
		Value2    int `validate:"max:50"`
		Value3    int `validate:"max:50|min:18"`
		Value4    int `validate:"max:50|min:18"`
		Structure App
	}

	Invalid struct {

	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	var nilErr ValidationErrors
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{5, ErrNotStruct},
		{App{"gavno"}, nilErr},
		{[]struct{}{}, ErrNotStruct},
		{Response{200, "ok"}, nilErr},
		{Response{73, "ok"}, ValidationErrors{ValidationError{"Code", ErrInvalidInt}}},
		{Token{}, nilErr},
		{User{
			ID:     "GogaMagoga",
			Name:   "JOPA",
			Age:    33,
			Email:  "gerladeno@gmail.com",
			Role:   "admin",
			Phones: []string{"tel1", "tel2"},
			meta:   nil,
		}, nilErr},
		{User{
			ID:     "longLongLongLongLongLongLongLongLongLongLongLongLongName",
			Name:   "",
			Age:    7,
			Email:  "zhopa",
			Role:   "slave",
			Phones: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
			meta:   nil,
		}, ValidationErrors{
			ValidationError{"ID", ErrInvalidString},
			ValidationError{"Age", ErrInvalidInt},
			ValidationError{"Email", ErrInvalidString},
			ValidationError{"Role", ErrInvalidString},
			ValidationError{"Phones", ErrInvalidSlice},
		}},
		{Embedded{
			Value1:    -3,
			Value2:    80,
			Value3:    40,
			Value4:    55,
			Structure: App{Version: "longVersionName"},
		}, ValidationErrors{
			ValidationError{"Value1", ErrInvalidInt},
			ValidationError{"Value2", ErrInvalidInt},
			ValidationError{"Value4", ErrInvalidInt},
			ValidationError{"Structure.Version", ErrInvalidString},
		}},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			result := Validate(test.in)
			require.Equal(t, result, test.expectedErr)
		})
	}
}
