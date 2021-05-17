package server

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/kube-portal/cmd/kubeportal"
	"github.com/ShotaKitazawa/kube-portal/server/controller"
	"github.com/ShotaKitazawa/kube-portal/server/infrastructure/kubernetes"
)

func Run(opts *kubeportal.Opts) error {
	// New
	e := echo.New()
	l := logrus.New()
	k8sClient, err := kubernetes.NewClient(opts.KubeConfigPath)
	if err != nil {
		panic(err)
	}
	k8sController := controller.NewK8sController(k8sClient, opts.JwtSecret)
	u, err := url.Parse(fmt.Sprintf("%s/auth/callback", opts.BaseUrl))
	if err != nil {
		panic(err)
	}
	oauthController, err := controller.NewOAuthController(l, opts.GitHubOAuthKey, opts.GitHubOAuthSecret, u.String(), opts.JwtSecret, opts.GitHubAllowUsers...)
	if err != nil {
		panic(err)
	}

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Route
	e.Static("/", "./client/out")
	api := e.Group("/api")
	api.GET("/list", k8sController.ListIngressInfo, k8sController.JwtAuthMiddleware)
	//auth := e.Group("/auth")
	e.GET("/auth/callback", oauthController.Callback)
	e.GET("/auth/logout", oauthController.Logout)
	e.GET("/auth/login", oauthController.Login)

	// Listen
	return http.ListenAndServe(strings.Join([]string{opts.BindAddr, strconv.Itoa(opts.BindPort)}, ":"), e)
}
