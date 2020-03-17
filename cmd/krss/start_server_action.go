package main

import (
	"fmt"
	"html/template"

	"github.com/go-playground/validator/v10"

	apiMiddleware "github.com/nebiros/krss/internal/middleware"

	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nebiros/krss/internal/router"
	"github.com/urfave/cli/v2"
)

func startServerAction(c *cli.Context) error {
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
	e.Use(middleware.Secure())
	e.Use(middleware.Gzip())

	err := router.ConfigureRoutes(e)
	if err != nil {
		return errors.WithStack(err)
	}

	return e.Start(fmt.Sprintf(":%s", port))
}
