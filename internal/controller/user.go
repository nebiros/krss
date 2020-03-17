package controller

import (
	"net/http"

	"github.com/nebiros/krss/internal/controller/input"

	"github.com/nebiros/krss/internal/model"

	"github.com/labstack/echo/v4"
)

type User struct {
	userModel model.UserInterface
}

func NewUser(userModel model.UserInterface) *User {
	return &User{
		userModel: userModel,
	}
}

func (ctrl *User) Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.gohtml", nil)
}

func (ctrl *User) DoLogin(c echo.Context) error {
	in := new(input.LoginInput)
	if err := c.Bind(in); err != nil {
		return err
	}

	if err := c.Validate(in); err != nil {
		return err
	}

	return nil
}
