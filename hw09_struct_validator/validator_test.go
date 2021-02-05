package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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

	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "1",
				Name:   "Xipe-Totec",
				Age:    5438,
				Email:  "hello@world",
				Role:   "admin",
				Phones: []string{"111222333"},
				meta:   json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Name",
					Err:   ErrInvalidString,
				},
			},
		},
		//{
		//	in: App{
		//		Version: "ololo",
		//	},
		//	expectedErr: ValidationErrors{
		//		ValidationError{
		//			Field: "Name",
		//			Err:   ErrInvalidString,
		//		},
		//	},
		//}, {
		//	in: Token{
		//		Header:    nil,
		//		Payload:   nil,
		//		Signature: nil,
		//	},
		//	expectedErr: ValidationErrors{
		//		ValidationError{
		//			Field: "Header",
		//			Err:   ErrInvalidString,
		//		},
		//	},
		//}, {
		//	in: Response{
		//		Code: 0,
		//		Body: "",
		//	},
		//	expectedErr: ValidationErrors{
		//		ValidationError{
		//			Field: "Body",
		//			Err:   ErrInvalidString,
		//		},
		//	},
		//},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
