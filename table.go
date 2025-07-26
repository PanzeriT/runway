package runway

import (
	"encoding/json"
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
	"github.com/panzerit/runway/model"
	"github.com/panzerit/runway/template/html"
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

	return Render(c, http.StatusOK, page.Table(a.name, m, &uc.User, string(t)))
}

func (a *App) showEditFormHandler(c echo.Context) error {
	logger.Debug("showing edit form for model", "model", c.Param("model"))

	return c.HTML(http.StatusCreated, "<tr><form><td> NEW </td><td> NEW </td><td> NEW </td><td> NEW </td></form></tr>")
}

func (a *App) createRowHandler(c echo.Context) error {
	m := c.Param("model")
	f, _ := c.FormParams()

	j, err := json.Marshal(f)
	if err != nil {
		return err
	}

	logger.Debug("creating row from JSON data", "model", m, "data", string(j))

	user := &model.User{} // TODO: this should be dynamic based on the model name
	err = json.Unmarshal(j, user)
	fmt.Println("Unmarshalled user:", user)
	fmt.Println("JSON data:", string(j))
	fmt.Println("Error:", err)

	err = a.service.CreateRowFromJSON(m, j)
	if err != nil {
		logger.Error("failed to create row", "model", m, "error", err)
		return nil
	}

	return c.HTML(http.StatusCreated, "<tr><td> NEW </td></tr>")
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

		editButton := html.LinkButton("#",
			html.WithText("Edit"),
			html.WithHxGet(fmt.Sprintf("/admin/model/%s/%s/edit", model, row.FieldByName("ID").String())),
			html.WithHxTarget("closest tr"),
			html.WithHxSwap("outerHTML"),
		)
		deleteButton := html.LinkButton("#",
			html.WithText("Delete"),
			html.WithClass("bg-red-600 hover:bg-gradient-to-r hover:from-red-600 hover:to-red-800"),
			html.WithHxDelete(fmt.Sprintf("/admin/model/%s/%s/delete", model, row.FieldByName("ID").String())),
			html.WithHxConfirm("Are you sure?"),
			html.WithHxTarget("closest tr"),
			html.WithHxSwap("outerHTML"),
		)

		tr.AddChild(Td(Component(editButton), Component(deleteButton)))
		table.AddChild(tr)
	}

	return table.Render(), nil
}
