package model

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
			err := Register(tt.input)

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
			"valid hidden annotation",
			"hidden",
			nil,
		},
		{
			"valid create annotation",
			"create",
			nil,
		},
		{
			"valid hidden and create annotation",
			"hidden,create",
			nil,
		},
		{
			"valid hidden and create annotation (reversed)",
			"create,hidden",
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

func TestCheckCreateAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		model     any
		modelName string
		fields    []Field
	}{
		{
			"check user creation field names",
			User{},
			"user",
			[]Field{
				NewField("Email"),
				NewField("FirstName"),
				NewField("LastName"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.model)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.fields, registeredModels[tt.modelName].CreationFieldNames)
		})
	}
}

func TestCheckHiddenAnnotation(t *testing.T) {
	tests := []struct {
		name      string
		model     any
		modelName string
		fields    []Field
	}{
		{
			"check user hidden field names",
			User{},
			"user",
			[]Field{
				NewField("ID"),
				NewField("CreatedAt"),
				NewField("UpdatedAt"),
				NewField("DeactivatedAt"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.model)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.fields, registeredModels[tt.modelName].HiddenFieldNames)
		})
	}
}

func TestFieldNameToHTMLId(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		doesPanic bool
	}{
		{
			"simple field name",
			"firstName",
			"first-name",
			false,
		},
		{
			"complex field name",
			"firstAndLastName",
			"first-and-last-name",
			false,
		},
		{
			"field name with numbers",
			"user123",
			"user123",
			false,
		},
		{
			"field name with special characters",
			"user-name!",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.doesPanic {
				assert.Panics(t, func() { fieldNameToHTMLId(tt.input) })
				return
			}

			got := fieldNameToHTMLId(tt.input)

			assert.Equal(t, tt.expected, got)
		})
	}
}
