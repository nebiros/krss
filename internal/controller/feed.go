package controller

import (
	"net/http"

	"github.com/mmcdole/gofeed"

	"github.com/nebiros/krss/internal/controller/input"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"github.com/nebiros/krss/internal/model"
)

type Feed struct {
	Base

	feedModel model.FeedInterface
}

func NewFeed(feedModel model.FeedInterface) *Feed {
	return &Feed{
		feedModel: feedModel,
	}
}

func (ctrl *Feed) Feeds(c echo.Context) error {
	u, err := ctrl.UserSession(c)
	if err != nil {
		return errors.WithStack(err)
	}

	fs, err := ctrl.feedModel.FeedsByUserID(u.UserID)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, "feed/feeds", map[string]interface{}{
		"title": "feeds",
		"feeds": fs,
	})
}

func (ctrl *Feed) NewFeed(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)

	return c.Render(http.StatusOK, "feed/new", map[string]interface{}{
		"title":     "new feed",
		"csrfToken": csrfToken,
	})
}

func (ctrl *Feed) DoNewFeed(c echo.Context) error {
	u, err := ctrl.UserSession(c)
	if err != nil {
		return errors.WithStack(err)
	}

	in := new(input.NewFeedInput)
	if err := c.Bind(in); err != nil {
		return errors.WithStack(err)
	}

	if err := c.Validate(in); err != nil {
		return errors.WithStack(err)
	}

	if len(in.Title) <= 0 {
		fp := gofeed.NewParser()

		f, err := fp.ParseURL(in.URL)
		if err != nil {
			return errors.WithStack(err)
		}

		in.Title = f.Title
	}

	createFeed := in.ToCreateFeed()
	createFeed.UserID = u.UserID

	if _, err := ctrl.feedModel.CreateFeed(createFeed); err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(http.StatusSeeOther, "/feeds")
}
