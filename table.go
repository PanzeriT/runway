package runway

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	. "github.com/panzerit/htmlkit"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/template/page"
)

// TODO: clean-up
var (
	ErrNoSlice = errors.New("data must be a slice")
	ErrNoData  = errors.New("no data provided")
	ErrNoModel = errors.New("model not found")
)

func (a *App) tableHandler(c echo.Context) error {
	uc := c.(*data.UserContext)

	m := c.Param("model")

	data, err := a.service.FindRows(m, 0, 0) // TODO: implement pagination
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching rows for model %s: %v", m, err))
		return nil
	}

	t, err := RenderTable(data, m)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error rendering table for model %s: %v", m, err))
	}

	return Render(c, http.StatusOK, page.Table(a.name, &uc.User, string(t)))
}

func (a *App) deleteRowHandler(c echo.Context) error {
	m := c.Param("model")
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Error parsing ID:", err)
	}

	log.Printf("Deleting row with ID %s from model %s", id, m)
	fmt.Println("Deleting row with ID", id, "from model", m)

	result := a.service.DeleteRow(m, id)
	_ = result
	_ = m

	return c.NoContent(http.StatusOK)
}

type tagName = string

// TODO: clean-up
const (
	tableTag tagName = "table"
	trTag    tagName = "tr"
	tdTag    tagName = "td"
	thTag    tagName = "th"
)

func RenderTable(data any, model string) ([]byte, error) {
	// check for errors or missing data
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		return []byte{}, ErrNoSlice
	}
	if v.Len() == 0 {
		return Table(Text("No data available")).Render(), nil
	}

	// evalute the slice
	t := v.Index(0).Type()

	// check annotations
	numFields := t.NumField()
	skipFields := make([]int, 0)
	for i := 0; i < numFields; i++ {
		tag := t.Field(i).Tag.Get("runway")
		if strings.Contains(tag, "hidden") {
			skipFields = append(skipFields, i)
		}
	}

	// generate table with header
	table := Table(Class("min-w-full divide-y divide-gray-200 dark:divide-gray-700"))
	tr := Tr(Class("p-8 mb-8 text-md font-normal text-gray-800"))
	for i := 0; i < numFields; i++ {
		if slices.Contains(skipFields, i) {
			continue
		}
		tr.AddChild(Td(Text(t.Field(i).Name)))
	}
	table.AddChild(tr)

	// generate table rows
	for i := 0; i < v.Len(); i++ {
		tr := Tr()

		row := v.Index(i)
		for j := 0; j < numFields; j++ {
			if slices.Contains(skipFields, j) {
				continue
			}
			cell := row.Field(j)
			tr.AddChild(Td(Text(template.HTMLEscapeString(fmt.Sprint(cell.Interface())))))
		}

		tr.AddChild(Td(Text("Test")))
		tr.AddChild(Td(Raw(`EDIT --- <button class="text-red-500" hx-delete="` + model + `/` + row.FieldByName("ID").Interface().(uuid.UUID).String() + `" hx-confirm="Are you sure?" hx-target="closest tr" hx-swap="outerHTML">Delete</button>`)))
		table.AddChild(tr)
	}

	return table.Render(), nil
}
