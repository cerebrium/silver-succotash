package dashboard

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hopdf.com/helpers"
	"hopdf.com/views/dashboard"
)

func DashboardHandler(c echo.Context) error {
	return helpers.Render(c, http.StatusOK, dashboard.Dashboard())
}
