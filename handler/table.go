package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	. "github.com/panzerit/htmlkit"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/model"
	"github.com/panzerit/runway/service"
	"github.com/panzerit/runway/template/html"
	"github.com/panzerit/runway/template/page"
)

const DEFAULT_LIMIT = 10

var ErrNoSlice = errors.New("data must be a slice")

type tableHandler struct {
	service service.Service
	logger  *slog.Logger
	appName string
}

func NewTableHandler(s service.Service, l *slog.Logger, appName string) Handler {
	return &tableHandler{
		service: s,
		logger:  l.With("handler", "table"),
		appName: appName,
	}
}

func (h *tableHandler) Register(g *echo.Group) {
	g.GET("/model/:model", h.getTable)
	g.GET("/model/:model/page/:page", h.getTable)
	g.GET("/model/:model/page/:page/limit/:limit", h.getTable)
	g.GET("/model/:model/:id/edit", h.getInlineForm)
	g.DELETE("/model/:model/:id", h.deleteRowHandler)
}

func (h *tableHandler) getTable(c echo.Context) error {
	uc := c.(*data.UserContext)
	model := c.Param("model")

	limit := getParamWithDefault(c, "limit", DEFAULT_LIMIT)
	currentPage := getParamWithDefault(c, "page", 1)

	count, err := h.service.FindRowCount(model)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error counting rows for model %s: %v", model, err))
	}
	if currentPage < 1 || int64(limit*currentPage) > count {
		render404(c, h.appName)
	}

	data, err := h.service.FindRows(model, limit, currentPage)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error fetching rows for model %s: %v", model, err))
		return nil
	}

	t, err := RenderTable(data, model)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Error rendering table for model %s: %v", model, err))
	}

	pagination := getPaginationLinks(model, limit, currentPage, count)

	return render(c, http.StatusOK, page.Table(h.appName, model, &uc.User, string(t), pagination))
}

func (h *tableHandler) getInlineForm(c echo.Context) error {
	m := c.Param("model")
	_ = m

	return c.HTML(http.StatusCreated, "<tr><form><td> NEW </td><td> NEW </td><td> NEW </td><td> NEW </td></form></tr>")
}

func (h *tableHandler) createRowHandler(c echo.Context) error {
	m := c.Param("model")
	f, _ := c.FormParams()

	j, err := json.Marshal(f)
	if err != nil {
		return err
	}

	h.logger.Debug("creating row from JSON data", "model", m, "data", string(j))

	user := &model.User{} // TODO: this should be dynamic based on the model name
	err = json.Unmarshal(j, user)
	fmt.Println("Unmarshalled user:", user)
	fmt.Println("JSON data:", string(j))
	fmt.Println("Error:", err)

	err = h.service.CreateRowFromJSON(m, j)
	if err != nil {
		h.logger.Error("failed to create row", "model", m, "error", err)
		return nil
	}

	return c.HTML(http.StatusCreated, "<tr><td> NEW </td></tr>")
}

func (h *tableHandler) deleteRowHandler(c echo.Context) error {
	m := c.Param("model")
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		fmt.Println("Error parsing ID:", err)
	}

	log.Printf("Deleting row with ID %s from model %s", id, m)
	fmt.Println("Deleting row with ID", id, "from model", m)

	result := h.service.DeleteRow(m, id)
	_ = result
	_ = m

	return c.NoContent(http.StatusOK)
}

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
		tr.AddChild(Th(Text(t.Field(i).Name)))
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

func getPaginationLinks(model string, limit, page int, count int64) []templ.Component {
	lastPage := int(count / int64(limit))
	pages := getPages(page, lastPage)

	links := []templ.Component{}

	linkString := "/admin/model/" + model + "/page/%d"

	if limit != DEFAULT_LIMIT {
		linkString = linkString + "/limit/" + fmt.Sprintf("%d", limit)
	}

	for _, p := range pages {
		linkText := fmt.Sprintf("%d", p)
		switch p {
		case -4:
			links = append(links, html.Link(fmt.Sprintf(linkString, 1), html.WithText("First")))
		case -3:
			links = append(links, html.Link(fmt.Sprintf(linkString, page-1), html.WithText("Prev")))
		case -2:
			links = append(links, html.Link(fmt.Sprintf(linkString, page+1), html.WithText("Next")))
		case -1:
			links = append(links, html.Link(fmt.Sprintf(linkString, lastPage), html.WithText("Last")))
		default:
			links = append(links, html.Link(fmt.Sprintf(linkString, p), html.WithText(linkText)))
		}
	}

	return links
}

func getPages(page, pages int) []int {
	result := []int{}

	if page > 2 {
		result = append(result, -4)
	}

	if page > 1 {
		result = append(result, -3)
	}

	for i := range 5 {
		if page-5+i >= 1 {
			result = append(result, page-5+i)
		}
	}

	result = append(result, page)

	for i := range 5 {
		if page+1+i <= pages {
			result = append(result, page+1+i)
		}
	}

	if page < pages {
		result = append(result, -2)
	}

	if page < pages-1 {
		result = append(result, -1)
	}

	return result
}
