package controller

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/nebiros/krss/internal/controller/output"
	apiMiddleware "github.com/nebiros/krss/internal/middleware"
	"github.com/pkg/errors"
)

type Base struct {
}

func (m *Base) Session(c echo.Context) (*sessions.Session, error) {
	s, err := session.Get("session", c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s, nil
}

func (m *Base) UserSession(c echo.Context) (output.UserOutput, error) {
	u, isType := c.Get(apiMiddleware.IsLoggedInContextKey).(output.UserOutput)
	if !isType {
		return output.UserOutput{}, errors.WithStack(echo.ErrUnauthorized)
	}

	return u, nil
}
