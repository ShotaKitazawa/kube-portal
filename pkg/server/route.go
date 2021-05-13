package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ShotaKitazawa/kube-portal/pkg/cmd/kubeportal"
	"github.com/ShotaKitazawa/kube-portal/pkg/kubernetes"
	"github.com/ShotaKitazawa/kube-portal/pkg/server/controller"
)

func Run(opts *kubeportal.Opts) error {
	// New
	e := echo.New()
	k8sClient, err := kubernetes.NewClient("") // TODO
	if err != nil {
		panic(err)
	}
	c := controller.New(k8sClient)

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Route
	// TODO client/assets 以下の静的コンテンツをサーブする設定
	api := e.Group("/api")
	api.GET("/list", c.ListIngressInfo)

	// Listen
	return http.ListenAndServe(strings.Join([]string{opts.BindAddr, strconv.Itoa(opts.BindPort)}, ":"), e)
}
