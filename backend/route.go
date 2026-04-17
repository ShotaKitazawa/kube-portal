package main

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

	"github.com/ShotaKitazawa/kube-portal/flag"
	"github.com/ShotaKitazawa/kube-portal/internal/api"
	"github.com/ShotaKitazawa/kube-portal/internal/handler"
	"github.com/ShotaKitazawa/kube-portal/internal/infrastructure/kubernetes"
)

func Run(ctx context.Context, cmd *cli.Command) error {
	// New Instances
	e := echo.New()
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	}))

	// New Clients
	k8sClient, err := kubernetes.NewClient(logger, cmd.String(flag.KubeConfigPath.Name))
	if err != nil {
		return fmt.Errorf("failed to initialize k8s client: %w", err)
	}

	disableOIDC := cmd.Bool(flag.DisableOIDC.Name)

	var provider *oidc.Provider
	var oidcVerifier *oidc.IDTokenVerifier
	if !disableOIDC {
		provider, err = oidc.NewProvider(ctx, cmd.String(flag.OIDCProviderURL.Name))
		if err != nil {
			return fmt.Errorf("failed to initialize provider: %w", err)
		}
		oidcVerifier = provider.Verifier(&oidc.Config{
			ClientID: cmd.String(flag.OIDCAudience.Name),
		})
	}

	// Build ogen server
	srv, err := api.NewServer(
		handler.NewHandler(k8sClient, cmd.Bool(flag.ShowUntaggedLinks.Name), disableOIDC, handler.OIDCConfig{
			Issuer:   cmd.String(flag.OIDCProviderURL.Name),
			ClientID: cmd.String(flag.OIDCClientID.Name),
			Audience: cmd.String(flag.OIDCAudience.Name),
		}, provider),
		handler.NewSecurityHandler(oidcVerifier, cmd.String(flag.RoleAttributePath.Name), logger),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize api server: %w", err)
	}

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

	// Routes
	e.Static("/", "./frontend/dist")
	e.Any("/api/*", echo.WrapHandler(http.StripPrefix("/api", srv)))

	// Listen
	return http.ListenAndServe(cmd.String(flag.BindAddr.Name), e)
}
