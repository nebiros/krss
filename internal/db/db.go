package db

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DefaultClientOptions = ClientOptions{
		Filename: "krss.db",
		Cache:    "shared",
		Mode:     "rw",
	}
)

// https://github.com/mattn/go-sqlite3#connection-string
type ClientOptions struct {
	Filename string
	Cache    string
	Mode     string
}

func NewClient(options ClientOptions) (*sqlx.DB, error) {
	dsn, err := url.Parse(fmt.Sprintf("file:%s", options.Filename))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	dsn.Query().Add("cache", options.Cache)
	dsn.Query().Add("mode", options.Mode)

	db, err := sqlx.Open("sqlite3", dsn.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}
