package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nebiros/krss/internal/controller"
	"github.com/nebiros/krss/internal/db"
	apiMiddleware "github.com/nebiros/krss/internal/middleware"
	"github.com/nebiros/krss/internal/model"
	"github.com/pkg/errors"
)

func ConfigureRoutes(e *echo.Echo) error {
	csrfFormConfig := middleware.DefaultCSRFConfig
	csrfFormConfig.TokenLookup = "form:csrf"

	dbClient, err := db.NewClient(db.DefaultClientOptions)
	if err != nil {
		return errors.WithStack(err)
	}

	userController := controller.NewUser(model.NewUser(dbClient))

	e.GET("/", userController.Login, middleware.CSRF())
	e.POST("/", userController.DoLogin, middleware.CSRFWithConfig(csrfFormConfig))

	feedController := controller.NewFeed(model.NewFeed(dbClient))

	e.GET("/feeds", feedController.Feeds, apiMiddleware.IsLoggedIn())
	e.GET("/feeds/new", feedController.NewFeed, apiMiddleware.IsLoggedIn(), middleware.CSRF())
	e.POST("/feeds/new", feedController.DoNewFeed, apiMiddleware.IsLoggedIn(), middleware.CSRFWithConfig(csrfFormConfig))

	return nil
}
