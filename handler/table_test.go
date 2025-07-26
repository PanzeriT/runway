package handler

import (
	"fmt"
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
	LastName  string
	CreatedAt time.Time
}

type UserWithHiddenID struct {
	ID        int `runway:"hidden"`
	LastName  string
	CreatedAt time.Time
}

func TestRenderTable_EmptySlice(t *testing.T) {
	expected := "<table>No data available</table>"

	got, err := RenderTable([]User{}, "user")

	assert.Nil(t, err)
	assert.Equal(t, expected, string(got))
}

func TestRenderTable_SliceOfUsers(t *testing.T) {
	items := []User{
		{
			ID:       1,
			LastName: "John",
		},
		{
			ID:       2,
			LastName: "Jane",
		},
	}
	// don't test the end, since there are buttons up for change
	// for now, the data isn't tested either
	expected := []string{
		// "<tr><th>ID</th><th>LastName</th><th>CreatedAt</th>",
		"<table class=",
		"<tr class=",
		"<th>ID</th><th>LastName</th><th>CreatedAt</th>",
		"<td>1</td><td>John</td><td>",
		"<td>2</td><td>Jane</td><td>",
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
		{
			1,
			"number1",
		},
		{
			2,
			"number2",
		},
	}
	expected := []string{
		"<table class=",
		"<tr class=",
		"<th>Field1</th><th>Field2</th></tr>",
		"<td>1</td><td>number1</td>",
		"<td>2</td><td>number2</td>",
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
		{
			ID:       1,
			LastName: "John",
		},
		{
			ID:       2,
			LastName: "Jane",
		},
	}
	// don't test the end, since there are buttons up for change
	// for now, the data isn't tested either
	expected := []string{
		"<table class=",
		"<tr class=",
		"<th>LastName</th>",
		"<td>John</td>",
		"<td>Jane</td>",
		"</table>",
	}
	notExpected := []string{
		"<td>1</td>",
		"<td>2</td>",
	}

	model.Register(UserWithHiddenID{})
	got, err := RenderTable(items, "user_with_hidden_id")

	assert.Nil(t, err)
	for _, e := range expected {
		assert.Contains(t, string(got), e)
	}

	for _, ne := range notExpected {
		assert.NotContains(t, string(got), ne)
	}
}

func TestGetPages(t *testing.T) {
	tests := []struct {
		page     int
		pages    int
		expected []int
	}{
		{
			1,
			1,
			[]int{1},
		},
		{
			1,
			2,
			[]int{1, 2, -2},
		},
		{
			1,
			10,
			[]int{1, 2, 3, 4, 5, 6, -2, -1},
		},
		{
			2,
			10,
			[]int{-3, 1, 2, 3, 4, 5, 6, 7, -2, -1},
		},
		{
			3,
			10,
			[]int{-4, -3, 1, 2, 3, 4, 5, 6, 7, 8, -2, -1},
		},
		{
			8,
			10,
			[]int{-4, -3, 3, 4, 5, 6, 7, 8, 9, 10, -2, -1},
		},
		{
			9,
			10,
			[]int{-4, -3, 4, 5, 6, 7, 8, 9, 10, -2},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("page %d of %d", tt.page, tt.pages), func(t *testing.T) {
			result := getPages(tt.page, tt.pages)

			assert.Equal(t, tt.expected, result)
		})
	}
}
