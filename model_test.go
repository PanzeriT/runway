package runway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterModel(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expectedErr error
	}{
		{
			"passing an int",
			1,
			ErrUnableToRegisterType,
		},
		{
			"passing a string",
			"string",
			ErrUnableToRegisterType,
		},
		{
			"passing an struct",
			User{},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterModel(tt.input)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestCheckAnnotation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			"valid annotation",
			"hidden",
			nil,
		},
		{
			"unknown annotation",
			"unknown",
			ErrUnknownAnnotation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkAnnotation(tt.input)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
