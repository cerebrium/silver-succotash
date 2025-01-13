package helpers

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Json the response if it isn't html
func Success(ctx echo.Context, content interface{}) error {
	fmt.Println("The content: ", content)
	return ctx.JSON(200, content)
}
