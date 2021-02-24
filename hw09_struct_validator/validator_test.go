package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string
type NumberWrapper int

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string        `validate:"len:5"`
		Numbers []int         `validate:"min:10|max:50"`
		Wrap    NumberWrapper `validate:"max:100|min:25"`
	}

	SomeStruct struct {
		Digit int    `validate:"min:10"`
		Word  string `validate:"regexp:^\\w+$|len:10"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,403,404"`
		Body string `json:"omitempty"`
	}

	Response2 struct {
		Code int64 `validate:"in:200,403,404,500"`
	}

	Request struct {
		Phone string `validate:"regexp:("`
	}
)

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "1",
				Name:   "Xipe-Totec",
				Age:    5,
				Email:  "hello@world",
				Role:   "superuser",
				Phones: []string{"111222333"},
				meta:   json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrLen,
				},
			},
		},
		{
			in: App{
				Version: "ololo",
				Numbers: []int{10, 15, 218},
				Wrap:    12,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Numbers",
					Err:   ErrMax,
				},
				ValidationError{
					Field: "Wrap",
					Err:   ErrMin,
				},
			},
		},
		{
			in: Token{
				Header:    []byte("text/json"),
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 0,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrIn,
				},
			},
		}, {
			in: Response2{
				Code: 404,
			},
			expectedErr: nil,
		},
		{
			in: SomeStruct{
				Digit: 10,
				Word:  "wordword10",
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			assert.Equal(t, &tt.expectedErr, err)
		})
	}
}

func TestValidateInternalProgramError(t *testing.T) {
	tests := []struct {
		in                           interface{}
		expectedInternalProgramError string
	}{
		{
			in:                           666,
			expectedInternalProgramError: "int is not a pointer to struct",
		}, {
			in:                           "string",
			expectedInternalProgramError: "string is not a pointer to struct",
		}, {
			in:                           true,
			expectedInternalProgramError: "bool is not a pointer to struct",
		}, {
			in:                           &reflect.Value{},
			expectedInternalProgramError: "*reflect.Value is not a pointer to struct",
		}, {
			in: Request{
				Phone: "218",
			},
			expectedInternalProgramError: "match error: error parsing regexp: missing closing ): `(`",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			assert.Equal(t, tt.expectedInternalProgramError, err.Error())
		})
	}
}
