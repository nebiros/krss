package controller

import (
	"net/http"

	"github.com/nebiros/krss/internal/controller/output"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"

	"github.com/labstack/echo/v4/middleware"

	"github.com/pkg/errors"

	"github.com/nebiros/krss/internal/controller/input"

	"github.com/nebiros/krss/internal/model"

	"github.com/labstack/echo/v4"
)

type User struct {
	Base

	userModel model.UserInterface
}

func NewUser(userModel model.UserInterface) *User {
	return &User{
		userModel: userModel,
	}
}

func (ctrl *User) Login(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)

	return c.Render(http.StatusOK, "login.gohtml", map[string]interface{}{
		"csrfToken": csrfToken,
	})
}

func (ctrl *User) DoLogin(c echo.Context) error {
	in := new(input.LoginInput)
	if err := c.Bind(in); err != nil {
		return errors.WithStack(err)
	}

	if err := c.Validate(in); err != nil {
		return errors.WithStack(err)
	}

	u, err := ctrl.userModel.SignIn(in.Email, in.Password)
	if err != nil {
		return errors.WithStack(err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return errors.WithStack(err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["user"] = u.ToUserOutput()

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(http.StatusSeeOther, "/feeds")
}

func (ctrl *User) Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return errors.WithStack(err)
	}

	sess.Values["user"] = output.UserOutput{}
	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
