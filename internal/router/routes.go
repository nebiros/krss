package router

import (
	"net/http"
	"time"

	"github.com/bluele/gcache"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
	"github.com/nebiros/krss/internal/controller"
	"github.com/nebiros/krss/internal/db"
	apiMiddleware "github.com/nebiros/krss/internal/middleware"
	"github.com/nebiros/krss/internal/model"
	"github.com/pkg/errors"
)

func ConfigureRoutes(e *echo.Echo) error {
	dbClient, err := db.NewClient(db.DefaultClientOptions)
	if err != nil {
		return errors.WithStack(err)
	}

	cacheClient := gcache.New(10).ARC().Build()

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	feedParser := gofeed.NewParser()
	feedParser.Client = httpClient

	userController := controller.NewUser(model.NewUser(dbClient))

	e.GET("/", userController.Login).Name = "user_login"
	e.POST("/", userController.DoLogin).Name = "user_do_login"

	feedController := controller.NewFeed(model.NewFeed(dbClient, feedParser, cacheClient))

	e.GET("/feeds", feedController.Feeds, apiMiddleware.IsLoggedIn())
	e.GET("/feeds/new", feedController.NewFeed, apiMiddleware.IsLoggedIn())
	e.POST("/feeds/new", feedController.DoNewFeed, apiMiddleware.IsLoggedIn())
	e.GET("/feeds/:feed_id", feedController.Show, apiMiddleware.IsLoggedIn())
	e.GET("/feeds/:feed_id/items/:slug", feedController.ShowItem, apiMiddleware.IsLoggedIn())
	e.GET("/feeds/:feed_id/items/:slug/read", feedController.ReadItem, apiMiddleware.IsLoggedIn())

	return nil
}
