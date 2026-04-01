package view

import "github.com/labstack/echo/v5"

type CommonResponse struct {
	Message string `json:"message"`
}

func ResponseJSON(ctx echo.Context, code int, message string) error {
	return ctx.JSON(code, CommonResponse{message})
}
