package runway

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/panzerit/runway/model"
)

func (a *App) CreateUser() {
	user := model.User{
		ID:        uuid.New(),
		FirstName: "Thomas",
		LastName:  "Panzeri",
		Email:     "thomas@panzeri.rocks",
	}

	result := a.db.Create(&user)

	fmt.Println(user.ID, user.CreatedAt, result.Error, result.RowsAffected)
}
