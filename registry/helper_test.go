package registry_test

import (
	"testing"

	"github.com/panzerit/runway/apps/root"
	"github.com/panzerit/runway/apps/user"
	"github.com/panzerit/runway/registry"
	"github.com/stretchr/testify/assert"
)

func TestInferName(t *testing.T) {
	tests := []struct {
		name         string
		app          registry.App
		expectedName string
	}{
		{
			"root",
			&root.RootApp{},
			"root",
		},
		{
			"user",
			&user.UserApp{},
			"user",
		},
	}

	for _, tt := range tests {
		t.Run("check app name for "+tt.name, func(t *testing.T) {
			got := registry.InferName(tt.app)
			assert.Equal(t, tt.expectedName, got, "app name was inferred wrongly")
		})
	}
}

func TestInferPath(t *testing.T) {
	tests := []struct {
		name         string
		app          registry.App
		expectedPath string
	}{
		{
			"root",
			&root.RootApp{},
			"",
		},
		{
			"user",
			&user.UserApp{},
			"user",
		},
	}

	for _, tt := range tests {
		t.Run("check app name for "+tt.name, func(t *testing.T) {
			got := registry.InferPath(tt.app)
			assert.Equal(t, tt.expectedPath, got, "app path was inferred wrongly")
		})
	}
}
