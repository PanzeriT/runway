package runway

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/panzerit/runway/data"
	"github.com/panzerit/runway/model"
	"github.com/panzerit/runway/template/page"
	"gorm.io/gorm"
)

func (a *App) tableHandler(c echo.Context) error {
	uc := c.(*data.UserContext)
	model := c.Param("model")

	columns, d := processSchemas(model, a.db)
	log.Printf("User %s (%s) is accessing the %s table", uc.User.Name, uc.User.Email, a.name)

	return Render(c, http.StatusOK, page.Table(a.name, &uc.User, columns, d))
}

func processSchemas(name string, db *gorm.DB) (data.Columns, [][]string) {
	fmt.Printf("Processing schemas for model: %s\n", name)

	m := registeredModels[name]
	t := reflect.TypeOf(m) // This is just a placeholder for actual schema processing logic
	columns := make(data.Columns, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columns[i] = data.Field{
			Name:       field.Name,
			Type:       field.Type.String(),
			IsReadOnly: false,
			IsHidden:   false,
		}
	}

	var users []model.User

	result := db.Find(&users)

	anySlice := make([][]string, result.RowsAffected)
	for i, user := range users {
		anySlice[i] = make([]string, t.NumField())
		fmt.Printf("Processing user %d: %s\n", i, user.FirstName)

		v := reflect.ValueOf(user)
		for j := 0; j < t.NumField(); j++ {
			field := v.Field(j)
			anySlice[i][j] = field.String()
		}
	}

	_ = anySlice
	fmt.Printf("%#v\n", anySlice)
	return columns, anySlice
}
