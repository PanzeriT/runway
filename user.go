package runway

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/panzerit/runway/model"
)

func (a *App) CreateUser(firstName, lastName, email string) uuid.UUID {
	// TODO: Remove since this is just for testing purposes
	user := model.User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err := a.service.CreateRow(&user)
	if err != nil {
		fmt.Println("Error creating user:", err)
	}

	return user.ID
}

func (a *App) DeleteUser(id uuid.UUID) {
	// TODO: Remove since this is just for testing purposes

	err := a.service.DeleteRow("user", id)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		return
	}
}
