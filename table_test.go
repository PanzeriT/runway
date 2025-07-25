package runway

import (
	"errors"
	"testing"
	"time"

	"github.com/panzerit/runway/model"
	"github.com/stretchr/testify/assert"
)

type Dummy struct {
	Field1 int
	Field2 string
}
type User struct {
	ID        int
	FirstName string
	CreateAt  time.Time
}

type UserWithHiddenID struct {
	ID        int `runway:"hidden"`
	FirstName string
	CreateAt  time.Time
}

func TestRenderTable_NoSlice(t *testing.T) {
	tests := []struct {
		name  string
		input any
	}{
		{
			"no slice but int",
			1,
		},
		{
			"no slice but string",
			"string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RenderTable(tt.input, "test")

			assert.NotNil(t, err)

			if !errors.Is(err, ErrNoSlice) {
				t.Errorf("expected ErrNoSlice, got %v", err)
			}
		})
	}
}

func TestRenderTable_EmptySlice(t *testing.T) {
	expected := "<table>No data available</table>"

	got, err := RenderTable([]User{}, "user")

	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func TestRenderTable_SliceOfUsers(t *testing.T) {
	items := []User{
		{ID: 1, FirstName: "John"},
		{ID: 2, FirstName: "Jane"},
	}
	// don't test the end, since there are buttons up for change
	// for now, the data isn't tested either
	expected := []string{
		"<table>",
		"<tr><th>ID</th><th>FirstName</th><th>CreateAt</th>",
		"<tr><td>1</td><td>John</td>",
		"<tr><td>2</td><td>Jane</td>",
		"</table>",
	}

	model.Register(User{})
	got, err := RenderTable(items, "user")

	assert.Nil(t, err)
	for _, e := range expected {
		assert.Contains(t, string(got), e)
	}
}

func TestRenderTable_SliceOfDummies(t *testing.T) {
	items := []Dummy{
		{1, "number1"},
		{2, "number2"},
	}
	expected := []string{
		"<table>",
		"<tr><th>Field1</th><th>Field2</th></tr>",
		"<tr><td>1</td><td>number1</td>",
		"<tr><td>2</td><td>number2</td>",
		"</table>",
	}

	model.Register(Dummy{})
	got, err := RenderTable(items, "dummy")

	assert.Nil(t, err)
	for _, e := range expected {
		assert.Contains(t, string(got), e)
	}
}

func TestRenderTable_SliceOfUsersWithHiddenID(t *testing.T) {
	items := []UserWithHiddenID{
		{ID: 1, FirstName: "John"},
		{ID: 2, FirstName: "Jane"},
	}
	// don't test the end, since there are buttons up for change
	// for now, the data isn't tested either
	expected := []string{
		"<table>",
		"<tr><th>FirstName</th><th>CreateAt</th>",
		"<tr><td>John</td>",
		"<tr><td>Jane</td>",
		"</table>",
	}

	model.Register(UserWithHiddenID{})
	got, err := RenderTable(items, "user_with_hidden_id")

	assert.Nil(t, err)
	for _, e := range expected {
		assert.Contains(t, string(got), e)
	}
}
