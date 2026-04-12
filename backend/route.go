package backend

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
	"github.com/urfave/cli/v3"

	"github.com/ShotaKitazawa/kube-portal/backend/controller"
	"github.com/ShotaKitazawa/kube-portal/backend/view"
	"github.com/ShotaKitazawa/kube-portal/flag"
	"github.com/ShotaKitazawa/kube-portal/internal/infrastructure/kubernetes"
	"github.com/ShotaKitazawa/kube-portal/internal/middleware"
)

func Run(ctx context.Context, cmd *cli.Command) error {
	// New Instances
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))
	// OIDC configuration
	provider, err := oidc.NewProvider(ctx, cmd.String(flag.OIDCProviderURL.Name))
	if err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}
	// Verify access tokens using the API audience as the expected aud claim
	oidcVerifier := provider.Verifier(&oidc.Config{
		ClientID: cmd.String(flag.OIDCAudience.Name),
	})
	// New Clients
	k8sClient, err := kubernetes.NewClient(logger, cmd.String(flag.KubeConfigPath.Name))
	if err != nil {
		panic(err)
	}
	// New Controllers
	k8sController := controller.NewAPIController(logger, k8sClient,
		&view.JSONForAPI{ShowUntaggedLinks: cmd.Bool(flag.ShowUntaggedLinks.Name)})

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
			AllowOrigins:  []string{os.Getenv("ALLOWED_ORIGIN_URL")},
			AllowMethods:  []string{"*"},
			AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	// Listen
	return http.ListenAndServe(cmd.String(flag.BindAddr.Name), e)
}
