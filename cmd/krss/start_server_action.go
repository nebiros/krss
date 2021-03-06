package main

import (
	"fmt"

	"github.com/nebiros/krss/internal/router"

	"github.com/gorilla/securecookie"

	"github.com/go-playground/validator/v10"

	apiMiddleware "github.com/nebiros/krss/internal/middleware"

	"github.com/pkg/errors"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	tr, err := apiMiddleware.NewTemplateRenderer()
	if err != nil {
		return errors.WithStack(err)
	}

	e.Renderer = tr

	csrfFormConfig := middleware.DefaultCSRFConfig
	csrfFormConfig.TokenLookup = "form:csrf"

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(middleware.CSRFWithConfig(csrfFormConfig))
	e.Use(session.Middleware(sessions.NewCookieStore(cookieAuthKey, cookieEncryptionKey)))

	if err := router.ConfigureRoutes(e); err != nil {
		return errors.WithStack(err)
	}

	return e.Start(fmt.Sprintf(":%s", port))
}
