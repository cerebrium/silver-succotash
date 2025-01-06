package helpers

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"hopdf.com/localware"
)

func AddLocalware(app *echo.Group, db_ref *sql.DB) {
	app.Use(localware.AddDb(db_ref))

	app.Use(localware.WithHeaderAuthorizationMiddleware)
	// Add user struct to context
	app.Use(localware.AddLocalUser)
}
