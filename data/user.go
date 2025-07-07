package data

import "github.com/labstack/echo/v4"

type User struct {
	Name  string
	Email string
}

type UserContext struct {
	echo.Context
	User
}
