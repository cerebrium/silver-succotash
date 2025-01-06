package notFound

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/helpers"
	"hopdf.com/views/notFound"
)

func NotFoundHandler(c echo.Context) error {
	return helpers.Render(c, http.StatusOK, notFound.NotFound())
}
