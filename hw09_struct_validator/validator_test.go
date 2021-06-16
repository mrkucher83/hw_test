package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
		// meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	IncorrectTagValue struct {
		Content string `validate:"len1:5"`
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

func TestPassValidate(t *testing.T) {
	tests := []struct {
		in interface{}
	}{
		{
			in: App{"v 2.0"},
		},
		{
			in: Response{404, "Not found"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Nil(t, err)

			_ = tt
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{"v 2.0.1"},
			expectedErr: ErrLength,
		},
		{
			in:          Response{403, "Forbidden"},
			expectedErr: fmt.Errorf("should be equal one of: 200,404,500"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// err := Validate(tt.in)
			// require.Equal(t, tt.expectedErr, err)

			_ = tt
		})
	}
}

func TestIncorrectTagValue(t *testing.T) {
	testStruct := IncorrectTagValue{"qwerty"}
	err := Validate(testStruct)

	require.ErrorIs(t, err, ErrUnknownTag)
}
