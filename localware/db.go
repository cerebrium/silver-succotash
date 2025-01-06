package localware

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type DbContext struct {
	echo.Context
	Db *sql.DB
}

// This is a wrapper around the middleware function
// that allows me to pass an extra argument to
// the handle function below that actually extends
// the context struct.
func AddDb(db_ref *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return handleDbAddition(next, db_ref)
	}
}

// Due to the wrapper above, we now can access both
// the next and the db_ref. Add the db_ref to the
// context so it is available on all requests anywhere.
func handleDbAddition(next echo.HandlerFunc, db_ref *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		newCtx := &DbContext{c, db_ref}

		return next(newCtx)
	}
}
