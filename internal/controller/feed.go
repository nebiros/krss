package controller

import (
	"net/http"
	"strconv"

	"github.com/go-shiori/go-readability"

	apiMiddleware "github.com/nebiros/krss/internal/middleware"
	"github.com/nebiros/krss/internal/model/entity"

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

	return c.Render(http.StatusOK, "feed/feeds", apiMiddleware.IncludeData{
		Title: "feeds",
		Data: struct {
			Feeds entity.Feeds
		}{Feeds: fs},
	})
}

func (ctrl *Feed) NewFeed(c echo.Context) error {
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)

	return c.Render(http.StatusOK, "feed/new", apiMiddleware.IncludeData{
		Title: "new feed",
		Data: struct {
			CSRFToken string
		}{CSRFToken: csrfToken},
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

func (ctrl *Feed) Show(c echo.Context) error {
	feedID, err := strconv.Atoi(c.Param("feed_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := ctrl.feedModel.FeedByFeedID(feedID)
	if err != nil {
		return errors.WithStack(err)
	}

	fp, err := ctrl.feedModel.ParseFeed(f)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, "feed/show", apiMiddleware.IncludeData{
		Title: f.Title,
		Data: struct {
			Feed        entity.Feed
			Description string
			Items       []*gofeed.Item
		}{
			Feed:        f,
			Description: fp.Description,
			Items:       fp.Items,
		},
	})
}

func (ctrl *Feed) ShowItem(c echo.Context) error {
	feedID, err := strconv.Atoi(c.Param("feed_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	slugID := c.Param("slug")

	f, err := ctrl.feedModel.FeedByFeedID(feedID)
	if err != nil {
		return errors.WithStack(err)
	}

	fp, err := ctrl.feedModel.ParseFeed(f)
	if err != nil {
		return errors.WithStack(err)
	}

	item := ctrl.feedModel.FeedItemBySlugIDWithItems(fp.Items, slugID)

	if item == nil {
		return echo.ErrNotFound
	}

	return c.Render(http.StatusOK, "item/show", apiMiddleware.IncludeData{
		Title: item.Title,
		Data: struct {
			Feed entity.Feed
			Item *gofeed.Item
		}{
			Feed: f,
			Item: item,
		},
	})
}

func (ctrl *Feed) ReadItem(c echo.Context) error {
	feedID, err := strconv.Atoi(c.Param("feed_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	slugID := c.Param("slug")

	f, err := ctrl.feedModel.FeedByFeedID(feedID)
	if err != nil {
		return errors.WithStack(err)
	}

	fp, err := ctrl.feedModel.ParseFeed(f)
	if err != nil {
		return errors.WithStack(err)
	}

	item := ctrl.feedModel.FeedItemBySlugIDWithItems(fp.Items, slugID)

	if item == nil {
		return echo.ErrNotFound
	}

	ir, err := ctrl.feedModel.ReadItem(item)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(http.StatusOK, "item/read", apiMiddleware.IncludeData{
		Title: item.Title,
		Data: struct {
			ItemRead readability.Article
		}{
			ItemRead: ir,
		},
	})
}
