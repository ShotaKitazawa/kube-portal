package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ShotaKitazawa/kube-portal/internal/api"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jmespath/go-jmespath"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type SecurityHandler struct {
	oidcVerifier      *oidc.IDTokenVerifier
	roleAttributePath string
	logger            *slog.Logger
}

var _ api.SecurityHandler = (*SecurityHandler)(nil)

func NewSecurityHandler(v *oidc.IDTokenVerifier, roleAttributePath string, l *slog.Logger) *SecurityHandler {
	return &SecurityHandler{oidcVerifier: v, roleAttributePath: roleAttributePath, logger: l}
}

func (s *SecurityHandler) HandleBearerAuth(ctx context.Context, _ api.OperationName, t api.BearerAuth) (context.Context, error) {
	idToken, err := s.oidcVerifier.Verify(ctx, t.Token)
	if err != nil {
		return ctx, fmt.Errorf("%w: %w", ErrUnauthorized, err)
	}

	var claims any
	if err := idToken.Claims(&claims); err != nil {
		return ctx, fmt.Errorf("failed to parse claims: %w", err)
	}

	res, err := jmespath.Search(s.roleAttributePath, claims)
	if err != nil {
		return ctx, fmt.Errorf("failed to search claims by JMESPath: %w", err)
	}
	typedRes, ok := res.(bool)
	if !ok || !typedRes {
		return ctx, fmt.Errorf("%w: JMESPath result is not true", ErrForbidden)
	}

	return setIsLogin(ctx, true), nil
}
