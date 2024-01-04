package server

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"github.com/ShotaKitazawa/kube-portal/cmd/kube-portal/flag"
	"github.com/ShotaKitazawa/kube-portal/server/controller"
	"github.com/ShotaKitazawa/kube-portal/server/infrastructure/github"
	"github.com/ShotaKitazawa/kube-portal/server/infrastructure/kubernetes"
	"github.com/ShotaKitazawa/kube-portal/server/utils"
)

func Run(opts *flag.Opts) error {
	// New Instances
	e := echo.New()
	l := logrus.New()
	// New Clients
	k8sClient, err := kubernetes.NewClient(opts.KubeConfigPath)
	if err != nil {
		panic(err)
	}
	githubClient := github.NewGitHubClient(l)
	// New Controllers
	k8sController := controller.NewK8sController(l,
		k8sClient, opts.JwtSecret, opts.ShowUntaggedLinks)
	u, err := url.Parse(fmt.Sprintf("%s/auth/callback", opts.BaseUrl))
	if err != nil {
		panic(err)
	}
	oauthController, err := controller.NewOAuthController(
		l, githubClient,
		opts.GitHubOAuthKey, opts.GitHubOAuthSecret, u.String(),
		opts.JwtSecret, time.Duration(opts.ExpiredHour)*time.Hour,
		utils.Split(opts.GitHubAllowUsers, ",")...,
	)
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
	auth := e.Group("/auth")
	auth.GET("/callback", oauthController.Callback)
	auth.GET("/logout", oauthController.Logout)
	auth.GET("/login", oauthController.Login)

	// Listen
	return http.ListenAndServe(opts.BindAddr, e)
}
