package index

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/helpers"
	"hopdf.com/views/index"
)

func IndexHandler(c echo.Context) error {
	return helpers.Render(c, http.StatusOK, index.Index())
}
