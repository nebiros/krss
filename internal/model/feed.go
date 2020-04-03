package model

import (
	"fmt"
	"time"

	"github.com/bluele/gcache"
	"github.com/jmoiron/sqlx"
	"github.com/mmcdole/gofeed"
	"github.com/nebiros/krss/internal/model/entity"
	"github.com/pkg/errors"
)

type FeedInterface interface {
	FeedsByUserID(userID int) (entity.Feeds, error)
	CreateFeed(feed entity.CreateFeed) (int, error)
	CreateFeedWithTx(tx *sqlx.Tx, feed entity.CreateFeed) (int, error)
	FeedByFeedID(feedID int) (entity.Feed, error)
	ParseFeed(feed entity.Feed) (*gofeed.Feed, error)
}

type Feed struct {
	*sqlx.DB

	feedParser  *gofeed.Parser
	cacheClient gcache.Cache
}

func NewFeed(dbClient *sqlx.DB, feedParser *gofeed.Parser, cacheClient gcache.Cache) *Feed {
	return &Feed{
		DB:          dbClient,
		feedParser:  feedParser,
		cacheClient: cacheClient,
	}
}

func (m *Feed) FeedsByUserID(userID int) (entity.Feeds, error) {
	if userID <= 0 {
		return entity.Feeds{}, errors.WithStack(&ErrEmptyArgument{
			Name:  "userID",
			Value: userID,
		})
	}

	q := `select
		feeds.feed_id as "feed.feed_id",
		feeds.user_id as "feed.user_id",
		feeds.title as "feed.title",
		feeds.url as "feed.url"
		from
		feeds
		inner join
		users
		on
		feeds.user_id = users.user_id 
		where
		users.user_id = ?`

	stmt, err := m.Preparex(m.Rebind(q))
	if err != nil {
		return entity.Feeds{}, errors.WithStack(err)
	}

	defer stmt.Close()

	var fs entity.Feeds

	if err := stmt.Select(&fs, userID); err != nil {
		return entity.Feeds{}, errors.WithStack(err)
	}

	return fs, nil
}

func (m *Feed) CreateFeed(feed entity.CreateFeed) (int, error) {
	tx, err := m.Beginx()
	if err != nil {
		return -1, errors.WithStack(err)
	}

	defer tx.Rollback()

	feedID, err := m.CreateFeedWithTx(tx, feed)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	if err := tx.Commit(); err != nil {
		return -1, errors.WithStack(err)
	}

	return feedID, nil
}

func (m *Feed) CreateFeedWithTx(tx *sqlx.Tx, feed entity.CreateFeed) (int, error) {
	q := `insert into feeds (user_id,
		title,
		url)
	values (:feed.user_id,
		:feed.title,
		:feed.url)`

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(feed)
	if err != nil {
		return -1, errors.WithStack(err)
	}

	feedID, err := res.LastInsertId()
	if err != nil {
		return -1, errors.WithStack(err)
	}

	return int(feedID), nil
}

func (m *Feed) FeedByFeedID(feedID int) (entity.Feed, error) {
	if feedID <= 0 {
		return entity.Feed{}, errors.WithStack(&ErrEmptyArgument{
			Name:  "feedID",
			Value: feedID,
		})
	}

	q := `select
		feeds.feed_id as "feed.feed_id",
		feeds.user_id as "feed.user_id",
		feeds.title as "feed.title",
		feeds.url as "feed.url"
		from
		feeds 
		where
		feeds.feed_id = ?`

	stmt, err := m.Preparex(m.Rebind(q))
	if err != nil {
		return entity.Feed{}, errors.WithStack(err)
	}

	defer stmt.Close()

	var f entity.Feed

	if err := stmt.Get(&f, feedID); err != nil {
		return entity.Feed{}, errors.WithStack(err)
	}

	return f, nil
}

func (m *Feed) ParseFeed(feed entity.Feed) (*gofeed.Feed, error) {
	if (entity.Feed{}) == feed {
		return nil, errors.WithStack(&ErrEmptyArgument{
			Name:  "feed",
			Value: feed,
		})
	}

	key := fmt.Sprintf("%d:%d", feed.UserID, feed.FeedID)

	if m.cacheClient.Has(key) {
		f, err := m.cacheClient.Get(key)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return f.(*gofeed.Feed), nil
	}

	fp, err := m.feedParser.ParseURL(feed.URL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	expiration, err := time.ParseDuration("30m")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := m.cacheClient.SetWithExpire(key, fp, expiration); err != nil {
		return nil, errors.WithStack(err)
	}

	return fp, nil
}
