package controller

import (
	"net/http"

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

	_, err = ctrl.feedModel.FeedsByUserID(u.UserID)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, "feeds.gohtml", nil)
}

func (ctrl *Feed) AddFeed(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)

	return c.Render(http.StatusOK, "add_feed.gohtml", map[string]interface{}{
		"csrfToken": csrfToken,
	})
}

func (ctrl *Feed) DoAddFeed(c echo.Context) error {
	in := new(input.AddFeedInput)
	if err := c.Bind(in); err != nil {
		return errors.WithStack(err)
	}

	if err := c.Validate(in); err != nil {
		return errors.WithStack(err)
	}

	_, err := ctrl.feedModel.CreateFeed(in.ToCreateFeed())
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Redirect(http.StatusSeeOther, "/feeds")
}
