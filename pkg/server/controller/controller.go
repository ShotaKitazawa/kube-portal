package controller

import "github.com/labstack/echo/v4"

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c Controller) ListIngressInfo(ctx echo.Context) error {
	// TODO
	return nil
}
