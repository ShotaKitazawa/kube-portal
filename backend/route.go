package backend

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2"

	"github.com/ShotaKitazawa/kube-portal/backend/controller"
	"github.com/ShotaKitazawa/kube-portal/backend/infrastructure/kubernetes"
	"github.com/ShotaKitazawa/kube-portal/backend/middleware"
	"github.com/ShotaKitazawa/kube-portal/backend/view"
	"github.com/ShotaKitazawa/kube-portal/flag"
)

func Run(ctx context.Context, cmd *cli.Command) error {
	// New Instances
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))
	// OAuth2 configuration
	provider, err := oidc.NewProvider(ctx, cmd.String(flag.OIDCProviderURL.Name))
	if err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}
	serverURL, _ := url.Parse(cmd.String(flag.ServerURL.Name)) // serverURL was already parsed
	oauthConfig := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.Endpoint().AuthURL,
			TokenURL: provider.Endpoint().TokenURL,
		},
		ClientID:     cmd.String(flag.OIDCClientID.Name),
		ClientSecret: cmd.String(flag.OIDCClientSecret.Name),
		RedirectURL:  serverURL.ResolveReference(&url.URL{Path: "/auth/callback"}).String(),
		Scopes:       []string{"openid", "profile", "email"},
	}
	oidcVerifier := provider.Verifier(&oidc.Config{
		ClientID: oauthConfig.ClientID,
	})
	// New Clients
	k8sClient, err := kubernetes.NewClient(logger, cmd.String(flag.KubeConfigPath.Name))
	if err != nil {
		panic(err)
	}
	// New Controllers
	k8sController := controller.NewAPIController(logger, k8sClient,
		&view.JSONForAPI{ShowUntaggedLinks: cmd.Bool(flag.ShowUntaggedLinks.Name)})
	oauthController := controller.NewAuthController(logger, oauthConfig)

	// Middlewares
	e.Use(echomiddleware.Recover())
	e.Use(slogecho.NewWithConfig(
		slog.New(slog.NewJSONHandler(os.Stderr,
			&slog.HandlerOptions{
				AddSource: false,
			})),
		slogecho.Config{
			WithUserAgent:     true,
			WithRequestID:     true,
			WithRequestHeader: true,
		}))
	if cmd.Bool(flag.Development.Name) {
		e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
			AllowOrigins:     []string{os.Getenv("ALLOWED_ORIGIN_URL")},
			AllowMethods:     []string{"*"},
			AllowCredentials: true,
		}))
	}

	// Route
	e.Static("/", "./frontend/out")
	api := e.Group("/api")
	api.Use(middleware.AuthWithConfig(middleware.AuthConfig{
		OIDCVerifier:      oidcVerifier,
		RoleAttributePath: cmd.String(flag.RoleAttributePath.Name),
	}, logger))
	api.GET("/list", k8sController.ListIngressInfo)
	auth := e.Group("/auth")
	auth.GET("/callback", oauthController.Callback)
	auth.GET("/logout", oauthController.Logout)
	auth.GET("/login", oauthController.Login)

	// Listen
	return http.ListenAndServe(cmd.String(flag.BindAddr.Name), e)
}
