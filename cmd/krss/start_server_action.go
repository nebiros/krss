package main

import (
	"fmt"
	"html/template"

	"github.com/gorilla/securecookie"

	"github.com/go-playground/validator/v10"

	apiMiddleware "github.com/nebiros/krss/internal/middleware"

	"github.com/pkg/errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nebiros/krss/internal/router"
	"github.com/urfave/cli/v2"
)

func startServerAction(c *cli.Context) error {
	cookieAuthKey := securecookie.GenerateRandomKey(64)
	cookieEncryptionKey := securecookie.GenerateRandomKey(32)

	port := c.String("port")

	e := echo.New()

	e.HideBanner = true

	e.Validator = &apiMiddleware.Validator{
		Validator: validator.New(),
	}

	e.Renderer = &apiMiddleware.TemplateRenderer{
		Template: template.Must(template.ParseGlob("../../web/template/*.gohtml")),
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(sessions.NewCookieStore(cookieAuthKey, cookieEncryptionKey)))

	err := router.ConfigureRoutes(e)
	if err != nil {
		return errors.WithStack(err)
	}

	return e.Start(fmt.Sprintf(":%s", port))
}
