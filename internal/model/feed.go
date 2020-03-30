package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/nebiros/krss/internal/model/entity"
	"github.com/pkg/errors"
)

type FeedInterface interface {
	FeedsByUserID(userID int) (entity.Feeds, error)
	CreateFeed(feed entity.CreateFeed) (int, error)
	CreateFeedWithTx(tx *sqlx.Tx, feed entity.CreateFeed) (int, error)
}

type Feed struct {
	*sqlx.DB
}

func NewFeed(db *sqlx.DB) *Feed {
	return &Feed{
		DB: db,
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
