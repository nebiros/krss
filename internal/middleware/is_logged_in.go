package middleware

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const IsLoggedInContextKey = "IsLoggedIn"

type IsLoggedInConfig struct {
	Skipper     middleware.Skipper
	SessionName string
}

var (
	DefaultIsLoggedInConfig = IsLoggedInConfig{
		Skipper:     middleware.DefaultSkipper,
		SessionName: "session",
	}
)

func IsLoggedIn() echo.MiddlewareFunc {
	return IsLoggedInWithConfig(DefaultIsLoggedInConfig)
}

func IsLoggedInWithConfig(config IsLoggedInConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultIsLoggedInConfig.Skipper
	}

	if len(config.SessionName) <= 0 {
		config.SessionName = DefaultIsLoggedInConfig.SessionName
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			sess, err := session.Get(config.SessionName, c)
			if err != nil {
				return err
			}

			if u, exists := sess.Values["user"]; !exists {
				return echo.ErrUnauthorized
			} else {
				c.Set(IsLoggedInContextKey, u)
			}

			return next(c)
		}
	}
}
