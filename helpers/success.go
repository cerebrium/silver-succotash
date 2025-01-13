package helpers

import (
	"github.com/labstack/echo/v4"
)

// Json the response if it isn't html
func Success(ctx echo.Context, content interface{}) error {
	return ctx.JSON(200, content)
}
